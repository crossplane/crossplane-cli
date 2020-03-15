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
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsNPCProvider = []string{"specTemplate", "providerRef"}
)

type NonPortableClass struct {
	instance *unstructured.Unstructured
}

func NewNonPortableClass(u *unstructured.Unstructured) *NonPortableClass {
	return &NonPortableClass{instance: u}
}

func (o *NonPortableClass) GetAge() string {
	return GetAge(o.instance)
}

func (o *NonPortableClass) GetStatus() string {
	return "N/A"
}

func (o *NonPortableClass) IsReady() bool {
	return true
}

func (o *NonPortableClass) GetObjectDetails() ObjectDetails {
	if o.instance == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *NonPortableClass) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	// Provider ref
	u, err := getObjRef(obj, fieldsNPCProvider)
	if err != nil {
		return related, err
	}
	if u != nil {
		// TODO(hasan): Hack to find Provider until full object reference is available.
		if u.GetAPIVersion() == "" {
			gv, err := schema.ParseGroupVersion(o.instance.GetAPIVersion())
			if err != nil {
				return related, err
			}
			s := strings.Split(gv.Group, ".")
			g := strings.Join(s[1:], ".")
			k := u.GetKind()
			if k == "" {
				k = "Provider"
			}
			u.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   g,
				Version: "",
				Kind:    k,
			})
		}
		related = append(related, u)
	}

	return related, nil
}
