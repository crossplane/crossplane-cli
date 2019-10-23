package trace

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/crossplaneio/crossplane-cli/pkg/crossplane"
	"github.com/fatih/color"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '

	signStateReady    = "\u2714" // Heavy Check Mark
	signStateNotReady = "!"
	signStateMissing  = "\u2715" //	Multiplication X
)

var detailsTemplate = `{{- if or .AdditionalStatusColumns .RemoteStatus .Conditions }}
{{.Kind}}: {{.Name}}

{{ if .AdditionalStatusColumns -}}
{{ range $i, $a := .AdditionalStatusColumns }}
{{- $a.name | printf "%s	" -}}
{{ end }}
{{ range $i, $a := .AdditionalStatusColumns }}
{{- $a.value | printf "%s	" -}}
{{ end }}
{{ end }}

{{- if .Conditions }}
Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
{{ range $i, $c := .Conditions }}
{{- $c.type }}	{{ $c.status }}	{{ $c.lastTransitionTime }}	{{ $c.reason }}	{{ $c.message }}	
{{ end }}
{{- end }}

{{- if .RemoteStatus }}
Remote Status
{{ .RemoteStatus -}}
{{- end }}
{{ end -}}
`

type SimplePrinter struct {
	tabWriter *tabwriter.Writer
}

func NewSimplePrinter() *SimplePrinter {
	t := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)
	return &SimplePrinter{tabWriter: t}
}

func (p *SimplePrinter) Print(nodes []*Node) error {
	err := p.printOverview(nodes)
	if err != nil {
		return err
	}
	err = p.printAllDetails(nodes)
	if err != nil {
		return err
	}
	return nil
}
func (p *SimplePrinter) printOverview(nodes []*Node) error {
	titleF := color.New(color.Bold).Add(color.Underline)
	_, err := titleF.Println("OVERVIEW")
	if err != nil {
		return err
	}
	fmt.Fprintln(p.tabWriter, "")

	_, err = fmt.Fprintln(p.tabWriter, "STATE\tKIND\tNAME\tNAMESPACE\tSTATUS\tAGE\t")
	if err != nil {
		return err
	}
	for _, n := range nodes {
		o := n.Instance
		stateSign := signStateReady
		status := "N/A"
		if n.State == NodeStateMissing {
			status = "<missing>"
			stateSign = signStateMissing
		} else if n.State == NodeStateNotReady {
			stateSign = signStateNotReady
		}

		c := crossplane.ObjectFromUnstructured(o)
		if c == nil {
			// This is not a known crossplane object (e.g. secret) so no related obj.
			_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t\n", stateSign, o.GetKind(), o.GetName(), o.GetNamespace(), status, crossplane.GetAge(o))
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t\n", stateSign, o.GetKind(), o.GetName(), o.GetNamespace(), c.GetStatus(), c.GetAge())
			if err != nil {
				return err
			}
		}
	}
	fmt.Fprintln(p.tabWriter, "")
	err = p.tabWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (p *SimplePrinter) printAllDetails(nodes []*Node) error {
	titleF := color.New(color.Bold).Add(color.Underline)
	_, err := titleF.Println("DETAILS")
	if err != nil {
		return err
	}
	fmt.Fprintln(p.tabWriter, "")

	allDetails := ""
	for _, n := range nodes {
		d := getDetailsText(n)
		if d != "" {
			d += "---\n"
		}
		allDetails += d
	}
	fmt.Fprintln(p.tabWriter, strings.Trim(strings.TrimSpace(allDetails), "-"))
	err = p.tabWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

func getDetailsText(node *Node) string {
	if node == nil {
		return "<error: node to trace is nil>"
	}

	o := node.Instance
	c := crossplane.ObjectFromUnstructured(o)
	if c == nil {
		return ""
	}
	d := c.GetObjectDetails()
	tmpl, err := template.New("details").Parse(detailsTemplate)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	dBuf := new(bytes.Buffer)
	err = tmpl.Execute(dBuf, d)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	return dBuf.String()
}
