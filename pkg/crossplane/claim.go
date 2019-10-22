package crossplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsClaimResource = append(fieldsSpec, "resourceRef")
)

type Claim struct {
	u *unstructured.Unstructured
}

func NewClaim(u *unstructured.Unstructured) *Claim {
	return &Claim{u: u}
}

func (o *Claim) GetStatus() string {
	return getResourceStatus(o.u)
}

func (o *Claim) GetAge() string {
	return GetAge(o.u)
}

func (o *Claim) GetObjectDetails() ObjectDetails {
	if o.u == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.u)
}

func (o *Claim) GetDetails() string {
	return getResourceDetails(o.u)
}

func (o *Claim) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u.Object

	// Get resource reference
	u, err := getObjRef(obj, fieldsClaimResource)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	// Get class reference
	u, err = getObjRef(obj, fieldsResourceClass)
	if err != nil {
		return related, err
	}
	// TODO(hasan): Hack for claim -> portableClass, currently apiversion, kind and ns missing
	//  hence we need to manually fill them. This limitation will be removed with
	//  https://github.com/crossplaneio/crossplane/blob/master/design/one-pager-simple-class-selection.md
	if u.GetAPIVersion() == "" {
		u.SetAPIVersion(o.u.GetAPIVersion())
	}
	if u.GetKind() == "" {
		u.SetKind(o.u.GetKind() + "Class")
	}
	if u.GetNamespace() == "" {
		u.SetNamespace(o.u.GetNamespace())
	}

	related = append(related, u)

	// Get write to secret reference
	u, err = getObjRef(obj, fieldsWriteConnSecret)
	u.SetAPIVersion("v1")
	u.SetKind("Secret")
	u.SetNamespace(o.u.GetNamespace())
	if err != nil {
		return related, err
	}
	related = append(related, u)

	return related, nil
}
