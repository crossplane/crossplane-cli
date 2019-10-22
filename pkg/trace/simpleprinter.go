package trace

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/crossplaneio/crossplane-cli/pkg/crossplane"

	"github.com/fatih/color"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
)

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
	err = p.printDetails(nodes)
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

	_, err = fmt.Fprintln(p.tabWriter, "KIND\tNAME\tNAMESPACE\tSTATUS\tAGE\t")
	if err != nil {
		return err
	}
	for _, n := range nodes {
		o := n.instance
		c, err := crossplane.ObjectFromUnstructured(o)
		if err != nil {
			return err
		}
		if c == nil {
			// This is not a known crossplane object (e.g. secret) so no related obj.
			s := "N/A"
			if n.state == NodeStateMissing {
				s = "<missing>"
			}
			_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t%v\t%v\t\n", o.GetKind(), o.GetName(), o.GetNamespace(), s, crossplane.GetAge(o))
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t%v\t%v\t\n", o.GetKind(), o.GetName(), o.GetNamespace(), c.GetStatus(), c.GetAge())
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
func (p *SimplePrinter) printDetails(nodes []*Node) error {
	titleF := color.New(color.Bold).Add(color.Underline)
	_, err := titleF.Println("DETAILS")
	if err != nil {
		return err
	}
	fmt.Fprintln(p.tabWriter, "")

	allDetails := ""
	for _, n := range nodes {
		o := n.instance
		c, err := crossplane.ObjectFromUnstructured(o)
		if err != nil {
			return err
		}
		if c == nil {
			continue
		}
		// TODO(hasan): How to print details should be responsibility of simple printer and not crossplane objects.
		//  Needs refactoring accoringly.
		d := c.GetDetails()
		if d != "" {
			d += "\n---\n\n"
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
