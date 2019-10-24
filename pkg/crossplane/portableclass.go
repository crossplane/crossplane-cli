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
