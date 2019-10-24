package crossplane

import (
	"strings"

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
	u, err := getObjRef(obj, fieldsNPCProvider)
	if err != nil {
		return related, err
	}

	if u != nil {
		// TODO(hasan): Hack to set apiVersion for Provider until full object reference is available.
		if u.GetAPIVersion() == "" {
			oApiVersion := o.instance.GetAPIVersion()
			s := strings.Split(oApiVersion, ".")
			a := strings.Join(s[1:], ".")

			u.SetAPIVersion(a)
		}
		if u.GetKind() == "" {
			u.SetKind("Provider")
		}
		related = append(related, u)
	}

	return related, nil
}
