package trace

import (
	"fmt"
	"io"
	"os"

	"github.com/emicklei/dot"
)

// GraphPrinter is able to print graph definition in Graphviz dot format.
type GraphPrinter struct {
	writer io.Writer
}

// NewGraphPrinter returns a GraphPrinter
func NewGraphPrinter() *GraphPrinter {
	return &GraphPrinter{writer: os.Stdout}
}

// Print prints graph definition in Graphviz dot format for a given root node.
func (p *GraphPrinter) Print(nodes []*Node) error {
	g := dot.NewGraph(dot.Undirected)
	for _, n := range nodes {
		relateds := n.Relateds
		for _, r := range relateds {
			t := g.Node(r.GetId())
			t.Label(getNodeLabel(r))
			t.Attr("penwidth", "2")
			if r.State == NodeStateMissing {
				t.Attr("color", "red")
				t.Attr("style", "dotted")
			} else if r.State == NodeStateNotReady {
				t.Attr("color", "orange")
			}
			f := g.Node(n.GetId())
			f.Label(getNodeLabel(n))
			f.Attr("penwidth", "2")
			if n.State == NodeStateMissing {
				f.Attr("color", "red")
				f.Attr("style", "dotted")
			} else if n.State == NodeStateNotReady {
				f.Attr("color", "orange")
			}
			g.Edge(f, t)
		}
	}
	fmt.Fprintln(p.writer, g.String())
	return nil
}

func getNodeLabel(n *Node) string {
	u := n.Instance
	labelKind := u.GetKind()
	labelName := u.GetName()
	if len(labelName) > 24 {
		labelName = labelName[:12] + "..." + labelName[len(labelName)-12:]
	}
	return fmt.Sprintf("%s\n%s", labelKind, labelName)
}
