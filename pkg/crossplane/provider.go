package crossplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (o *Provider) GetDetails() string {
	return ""
}

func (o *Provider) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	// TODO(hasan): credentialsSecretRef?
	return nil, nil
}
