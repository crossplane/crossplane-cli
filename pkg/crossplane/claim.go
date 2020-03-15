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
	fieldsClaimResource = append(fieldsSpec, "resourceRef")
)

type Claim struct {
	instance *unstructured.Unstructured
}

func NewClaim(u *unstructured.Unstructured) *Claim {
	return &Claim{instance: u}
}

func (o *Claim) GetStatus() string {
	return getResourceStatus(o.instance)
}

func (o *Claim) GetAge() string {
	return GetAge(o.instance)
}

func (o *Claim) GetObjectDetails() ObjectDetails {
	if o.instance == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *Claim) IsReady() bool {
	return o.GetStatus() == resourceBindingPhaseBound
}

func (o *Claim) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	// Get resource reference
	u, err := getObjRef(obj, fieldsClaimResource)
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
		// TODO(hasan): Forward compatible hack for claim -> portableClass, currently apiversion, kind and ns missing
		//  hence we need to manually fill them. This limitation will be removed with
		//  https://github.com/crossplane/crossplane/blob/master/design/one-pager-simple-class-selection.md
		if u.GetAPIVersion() == "" {
			u.SetAPIVersion(o.instance.GetAPIVersion())
			u.SetKind(o.instance.GetKind() + "Class")
			u.SetNamespace(o.instance.GetNamespace())
		}
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
		u.SetNamespace(o.instance.GetNamespace())
		if err != nil {
			return related, err
		}
		related = append(related, u)
	}

	return related, nil
}
