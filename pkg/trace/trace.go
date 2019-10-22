package trace

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type NodeState string

const (
	NodeStateMissing NodeState = "Missing"
)

type Node struct {
	instance *unstructured.Unstructured
	gvr      schema.GroupVersionResource
	id       string
	related  []*Node
	state    NodeState
}

func NewNode(res schema.GroupVersionResource, instance *unstructured.Unstructured) *Node {
	return &Node{
		gvr:      res,
		instance: instance,
		related:  nil,
		state:    "",
	}
}

func (n *Node) GetId() string {
	return GetNodeIdFor(n.gvr, n.instance)
}

func (n *Node) IsFetched() bool {
	return n.instance.GetUID() != ""
}

func GetNodeIdFor(res schema.GroupVersionResource, i *unstructured.Unstructured) string {
	return strings.ToLower(fmt.Sprintf("%s-%s-%s", res.String(), i.GetNamespace(), i.GetName()))
}

type GraphBuilder interface {
	BuildGraph(string, string, string) (*Node, []*unstructured.Unstructured, error)
}

type Printer interface {
	Print([]*Node) error
}
