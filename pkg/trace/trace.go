/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trace

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type NodeState string

const (
	NodeStateUnknown  NodeState = ""
	NodeStateMissing  NodeState = "Missing"
	NodeStateReady    NodeState = "Ready"
	NodeStateNotReady NodeState = "NotReady"
)

type Node struct {
	Instance *unstructured.Unstructured
	GVR      schema.GroupVersionResource
	Relateds []*Node
	State    NodeState
}

func NewNode(res schema.GroupVersionResource, instance *unstructured.Unstructured) *Node {
	return &Node{
		GVR:      res,
		Instance: instance,
		Relateds: nil,
		State:    "",
	}
}

func (n *Node) GetId() string {
	return GetNodeIdFor(n.GVR, n.Instance)
}

func (n *Node) IsFetched() bool {
	return n.Instance.GetUID() != "" || n.State != NodeStateUnknown
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
