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
	fieldsManagedClaim = append(fieldsSpec, "claimRef")
)

type Managed struct {
	instance *unstructured.Unstructured
}

func NewManaged(u *unstructured.Unstructured) *Managed {
	return &Managed{instance: u}
}

func (o *Managed) GetStatus() string {
	return getResourceStatus(o.instance)
}

func (o *Managed) GetAge() string {
	return GetAge(o.instance)
}

func (o *Managed) GetObjectDetails() ObjectDetails {
	if o.instance == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *Managed) IsReady() bool {
	return o.GetStatus() == resourceBindingPhaseBound
}

func (o *Managed) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	// Get claim reference
	u, err := getObjRef(obj, fieldsManagedClaim)
	if err != nil {
		return related, err
	}
	if u != nil {
		related = append(related, u)
	}

	// Get class reference
	u, err = getObjRef(obj, fieldsResourceClass)
	if err != nil {
		return related, err
	}
	if u != nil {
		related = append(related, u)
	}

	// Get write to secret reference
	u, err = getObjRef(obj, fieldsWriteConnSecret)
	if err != nil {
		return related, err
	}
	if u != nil {
		u.SetAPIVersion("v1")
		u.SetKind("Secret")
		// For backward compatibility with namespaced managed resources, if namespaced, search secret in the
		// same namespace as managed resource otherwise get namespace from spec.writeConnectionSecretsToNamespace
		if u.GetNamespace() == "" {
			u.SetNamespace(o.instance.GetNamespace())
		}

		related = append(related, u)
	}

	// TODO(hasan): add provider as a reference here.

	return related, nil
}
