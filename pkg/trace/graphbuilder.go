package trace

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/crossplaneio/crossplane-cli/pkg/crossplane"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeGraphBuilder struct {
	client     dynamic.Interface
	restMapper meta.RESTMapper
	nodes      map[string]*Node
}

func NewKubeGraphBuilder(client dynamic.Interface, restMapper meta.RESTMapper) *KubeGraphBuilder {
	return &KubeGraphBuilder{
		client:     client,
		restMapper: restMapper,
		nodes:      map[string]*Node{},
	}
}

func (g *KubeGraphBuilder) BuildGraph(name, namespace, groupRes string) (root *Node, traversed []*Node, err error) {
	queue := list.New()

	traversed = make([]*Node, 0)

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	gr := schema.ParseGroupResource(groupRes)
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gr.Group,
		Version: "",
		Kind:    gr.Resource,
	})
	u.SetName(name)
	u.SetNamespace(namespace)

	root, err = g.addNodeIfNotExist(u)
	if err != nil {
		return nil, nil, err
	}

	err = g.fetchObj(root)
	if err != nil {
		return nil, nil, err
	}
	if root.State == NodeStateMissing {
		return root, nil, errors.New(
			fmt.Sprintf("Object to trace is not found: \"%s\" \"%s\" in namespace \"%s\"", groupRes, name, namespace))
	}
	c := crossplane.ObjectFromUnstructured(root.Instance)
	if c == nil || c.IsReady() {
		// This is not a known crossplane object (e.g. secret) or a ready crossplane object
		root.State = NodeStateReady
	} else {
		root.State = NodeStateNotReady
	}

	// TODO(hasan): figure out if visited can be enough without traversed.
	visited := map[string]bool{}
	traversed = append(traversed, root)
	visited[root.GetId()] = true
	queue.PushBack(root)

	for queue.Len() > 0 {
		qnode := queue.Front()
		node := qnode.Value.(*Node)
		// Skip if object is missing
		if node.State == NodeStateMissing {
			queue.Remove(qnode)
			continue
		}
		err = g.findRelated(node)
		if err != nil {
			return nil, nil, err
		}

		for _, n := range node.Relateds {
			if !n.IsFetched() {
				err := g.fetchObj(n)
				if kerrors.IsNotFound(err) {
					n.State = NodeStateMissing
				} else if err != nil {
					return nil, nil, err
				}
			}
			if n.State == NodeStateUnknown {
				c := crossplane.ObjectFromUnstructured(n.Instance)
				if c == nil || c.IsReady() {
					// This is not a known crossplane object (e.g. secret) or a ready crossplane object
					n.State = NodeStateReady
				} else {
					n.State = NodeStateNotReady
				}
			}
			nid := n.GetId()
			if !visited[nid] {
				traversed = append(traversed, n)
				visited[nid] = true
				queue.PushBack(n)
			}
		}
		queue.Remove(qnode)
	}
	return
}

func (g *KubeGraphBuilder) fetchObj(n *Node) error {
	if n.IsFetched() {
		return nil
	}
	gvr := n.GVR
	u := n.Instance

	u, err := g.client.Resource(gvr).Namespace(u.GetNamespace()).Get(u.GetName(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	n.Instance = u
	return nil
}

func (g *KubeGraphBuilder) filterByLabel(gvk metav1.GroupVersionKind, namespace, selector string) ([]unstructured.Unstructured, error) {
	res, err := g.restMapper.ResourceFor(schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: gvk.Kind})
	if err != nil {
		return nil, err
	}

	list, err := g.client.Resource(res).Namespace(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (g *KubeGraphBuilder) findRelated(n *Node) error {
	n.Relateds = make([]*Node, 0)

	c := crossplane.ObjectFromUnstructured(n.Instance)
	if c == nil {
		// This is not a known crossplane object (e.g. secret) so no related obj.
		return nil
	}
	objs, err := c.GetRelated(g.filterByLabel)
	if err != nil {
		return err
	}
	for _, o := range objs {
		r, err := g.addNodeIfNotExist(o)
		if err != nil {
			return err
		}
		n.Relateds = append(n.Relateds, r)
	}
	return nil
}

func (g *KubeGraphBuilder) addNodeIfNotExist(u *unstructured.Unstructured) (*Node, error) {
	var n *Node
	gvr, err := g.restMapper.ResourceFor(schema.GroupVersionResource{Group: u.GroupVersionKind().Group, Version: u.GroupVersionKind().Version, Resource: u.GetKind()})
	if err != nil {
		return nil, err
	}
	id := GetNodeIdFor(gvr, u)
	if e, ok := g.nodes[id]; ok {
		n = e
	} else {
		n = NewNode(gvr, u)
		g.nodes[id] = n
	}
	return n, nil
}
