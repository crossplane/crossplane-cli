package crossplane

import (
	"k8s.io/apimachinery/pkg/api/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsAppDesiredRes          = append(fieldsStatus, "desiredResources")
	fieldsAppSubmittedRes        = append(fieldsStatus, "submittedResources")
	fieldsAppResourceSelector    = append(fieldsSpec, "resourceSelector")
	fieldsAppResourceMatchLabels = append(fieldsAppResourceSelector, "matchLabels")

	stateAppSubmitted = "Submitted"
)

type Application struct {
	instance *unstructured.Unstructured
}

func NewApplication(u *unstructured.Unstructured) *Application {
	return &Application{instance: u}
}

func (o *Application) GetStatus() string {
	return getNestedString(o.instance.Object, fieldsStatusState...)
}

func (o *Application) GetAge() string {
	return GetAge(o.instance)
}

func (o *Application) GetObjectDetails() ObjectDetails {
	u := o.instance
	if u == nil {
		return ObjectDetails{}
	}
	return getObjectDetails(o.instance)
}

func (o *Application) IsReady() bool {
	return o.GetStatus() == stateAppSubmitted
}

func (o *Application) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance.Object

	// Get resource reference
	u, err := getObjRef(obj, fieldsStatusClusterRef)
	if err != nil {
		return related, err
	}
	if u != nil {
		related = append(related, u)
	}

	// Get related resources with resourceSelector
	namespacedResourceKinds := getKindsFromGroupKinds(groupKindsClaim, groupKindsApplicationResource)

	for _, k := range namespacedResourceKinds {
		uArr, err := filterByLabel(metav1.GroupVersionKind{
			Kind: k,
		}, o.instance.GetNamespace(), getNestedLabelSelector(obj, fieldsAppResourceMatchLabels...))
		// Ignore NoMatchError since all resources/kinds may not be available on the API,
		// e.g. ignore if AWS stack is not installed when working GCP only.
		if err != nil && !meta.IsNoMatchError(err) {
			return related, err
		}

		for _, u := range uArr {
			related = append(related, u.DeepCopy())
		}
	}

	return related, nil
}
