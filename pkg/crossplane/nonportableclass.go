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
	u *unstructured.Unstructured
}

func NewNonPortableClass(u *unstructured.Unstructured) *NonPortableClass {
	return &NonPortableClass{u: u}
}

func (o *NonPortableClass) GetAge() string {
	return GetAge(o.u)
}

func (o *NonPortableClass) GetStatus() string {
	return "N/A"
}

func (o *NonPortableClass) IsReady() bool {
	return true
}

func (o *NonPortableClass) GetObjectDetails() ObjectDetails {
	if o.u == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.u)
}

func (o *NonPortableClass) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u.Object
	u, err := getObjRef(obj, fieldsNPCProvider)
	if err != nil {
		return related, err
	}

	// TODO(hasan): Hack to set apiVersion for Provider until full object reference is available.
	if u.GetAPIVersion() == "" {
		oApiVersion := o.u.GetAPIVersion()
		s := strings.Split(oApiVersion, ".")
		a := strings.Join(s[1:], ".")

		u.SetAPIVersion(a)
	}
	if u.GetKind() == "" {
		u.SetKind("Provider")
	}
	related = append(related, u)

	return related, nil
}
