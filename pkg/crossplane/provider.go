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
	fieldsProviderCredSecretRef          = append(fieldsSpec, "credentialsSecretRef")
	fieldsProviderCredSecretRefName      = append(fieldsProviderCredSecretRef, "name")
	fieldsProviderCredSecretRefNamespace = append(fieldsProviderCredSecretRef, "namespace")
)

type Provider struct {
	instance *unstructured.Unstructured
}

func NewProvider(u *unstructured.Unstructured) *Provider {
	return &Provider{instance: u}
}

func (o *Provider) GetStatus() string {
	return "N/A"
}

func (o *Provider) GetAge() string {
	return GetAge(o.instance)
}

func (o *Provider) GetObjectDetails() ObjectDetails {
	if o.instance == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *Provider) IsReady() bool {
	return true
}

func (o *Provider) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	u := &unstructured.Unstructured{}
	n := getNestedString(obj, fieldsProviderCredSecretRefName...)
	ns := getNestedString(obj, fieldsProviderCredSecretRefNamespace...)
	if n != "" {
		u.SetName(n)
		u.SetAPIVersion("v1")
		u.SetKind("Secret")
		// For backward compatibility, i.e. namespaced Providers didn't set namespace for the secret.
		if ns != "" {
			u.SetNamespace(ns)
		} else {
			u.SetNamespace(o.instance.GetNamespace())
		}
		related = append(related, u)
	}

	return related, nil
}
