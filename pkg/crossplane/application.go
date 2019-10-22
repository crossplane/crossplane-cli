package crossplane

import (
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/api/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsAppDesiredRes          = append(fieldsStatus, "desiredResources")
	fieldsAppSubmittedRes        = append(fieldsStatus, "submittedResources")
	fieldsAppResourceSelector    = append(fieldsSpec, "resourceSelector")
	fieldsAppResourceMatchLabels = append(fieldsAppResourceSelector, "matchLabels")

	applicationDetailsTemplate = `%v

NAME	CLUSTER	STATUS	DESIRED	SUBMITTED
%v	%v	%v	%v	%v	

State Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

type Application struct {
	u *unstructured.Unstructured
}

func NewApplication(u *unstructured.Unstructured) *Application {
	return &Application{u: u}
}

func (o *Application) GetStatus() string {
	return getNestedString(o.u.Object, fieldsStatusState...)
}

func (o *Application) GetAge() string {
	return GetAge(o.u)
}

func (o *Application) GetObjectDetails() ObjectDetails {
	u := o.u
	if u == nil {
		return ObjectDetails{}
	}
	od := getObjectDetails(o.u)

	apcs := make([]map[string]string, 0)
	apcs = append(apcs, getColumn("NAME", u.GetName()))
	apcs = append(apcs, getColumn("CLUSTER", getNestedString(o.u.Object, fieldsStatusClusterRefName...)))
	apcs = append(apcs, getColumn("STATUS", o.GetStatus()))
	apcs = append(apcs, getColumn("DESIRED", strconv.Itoa(int(getNestedInt64(o.u.Object, fieldsAppDesiredRes...)))))
	apcs = append(apcs, getColumn("SUBMITTED", strconv.Itoa(int(getNestedInt64(o.u.Object, fieldsAppSubmittedRes...)))))

	od.AdditionalPrinterColumns = apcs

	return od
}

func (o *Application) GetDetails() string {
	d := fmt.Sprintf(applicationDetailsTemplate, o.u.GetKind(),
		o.u.GetName(), getNestedString(o.u.Object, fieldsStatusClusterRefName...),
		o.GetStatus(), getNestedInt64(o.u.Object, fieldsAppDesiredRes...),
		getNestedInt64(o.u.Object, fieldsAppSubmittedRes...))

	cs, f, err := unstructured.NestedSlice(o.u.Object, fieldsConditionedStatusConditions...)
	if err != nil || !f {
		// failed to get conditions
		return d
	}
	for _, c := range cs {
		cMap := c.(map[string]interface{})
		if cMap == nil {
			d = d + "<error: condition status is not a map>"
			continue
		}

		d = d + fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t\n",
			getNestedString(cMap, conditionKeyType),
			getNestedString(cMap, conditionKeyStatus),
			getNestedString(cMap, conditionKeyLastTransitionTime),
			getNestedString(cMap, conditionKeyReason),
			getNestedString(cMap, conditionKeyMessage))
	}
	return d
}

func (o *Application) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u.Object

	// Get resource reference
	u, err := getObjRef(obj, fieldsStatusClusterRef)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	// Get related resources with resourceSelector
	resourceKinds := getKindsFromGroupKinds(groupKindsClaim, groupKindsManaged, groupKindsApplicationResource)

	for _, k := range resourceKinds {
		uArr, err := filterByLabel(metav1.GroupVersionKind{
			Kind: k,
		}, o.u.GetNamespace(), getNestedLabelSelector(obj, fieldsAppResourceMatchLabels...))
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
