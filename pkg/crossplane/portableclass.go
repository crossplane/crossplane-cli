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

package crossplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsPCClass = []string{"classRef"}
)

type PortableClass struct {
	instance *unstructured.Unstructured
}

func NewPortableClass(u *unstructured.Unstructured) *PortableClass {
	return &PortableClass{instance: u}
}
func (o *PortableClass) GetAge() string {
	return GetAge(o.instance)
}

func (o *PortableClass) GetStatus() string {
	return "N/A"
}

func (o *PortableClass) IsReady() bool {
	return true
}

func (o *PortableClass) GetObjectDetails() ObjectDetails {
	if o.instance == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *PortableClass) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	// Get class reference
	u, err := getObjRef(obj, fieldsPCClass)
	if err != nil {
		return related, err
	}
	if u != nil {
		related = append(related, u)
	}

	return related, nil
}
