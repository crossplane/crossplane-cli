package crossplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsProviderCredSecretRef     = append(fieldsSpec, "credentialsSecretRef")
	fieldsProviderCredSecretRefName = append(fieldsProviderCredSecretRef, "name")
)

type Provider struct {
	u *unstructured.Unstructured
}

func NewProvider(u *unstructured.Unstructured) *Provider {
	return &Provider{u: u}
}

func (o *Provider) GetStatus() string {
	return "N/A"
}

func (o *Provider) GetAge() string {
	return GetAge(o.u)
}

func (o *Provider) GetObjectDetails() ObjectDetails {
	if o.u == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.u)
}

func (o *Provider) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u.Object

	u := &unstructured.Unstructured{}
	n := getNestedString(obj, fieldsProviderCredSecretRefName...)
	if n != "" {
		u.SetName(n)
		u.SetAPIVersion("v1")
		u.SetKind("Secret")
		u.SetNamespace(o.u.GetNamespace())
		related = append(related, u)
	}

	return related, nil
}
