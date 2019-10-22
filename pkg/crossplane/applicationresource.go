package crossplane

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fieldsAppResSecrets      = append(fieldsSpec, "secrets")
	fieldsAppResTemplate     = append(fieldsSpec, "template")
	fieldsAppResTemplateKind = append(fieldsAppResTemplate, "kind")
	fieldsAppResTemplateName = append(fieldsAppResTemplate, "metadata", "name")
	fieldsAppResStatusRemote = append(fieldsStatus, "remote")

	applicationResourceDetailsTemplate = `%v

NAME	TEMPLATE-KIND	TEMPLATE-NAME	CLUSTER	STATUS
%v	%v	%v	%v	%v	

Remote State
%v

State Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

type ApplicationResource struct {
	u *unstructured.Unstructured
}

func NewApplicationResource(u *unstructured.Unstructured) *ApplicationResource {
	return &ApplicationResource{u: u}
}

func (o *ApplicationResource) GetStatus() string {
	return getNestedString(o.u.Object, fieldsStatusState...)
}

func (o *ApplicationResource) GetAge() string {
	return GetAge(o.u)
}

func (o *ApplicationResource) GetObjectDetails() ObjectDetails {
	u := o.u
	if u == nil {
		return ObjectDetails{}
	}
	od := getObjectDetails(o.u)

	apcs := make([]map[string]string, 0)
	apcs = append(apcs, getColumn("NAME", u.GetName()))
	apcs = append(apcs, getColumn("TEMPLATE-KIND", getNestedString(o.u.Object, fieldsAppResTemplateKind...)))
	apcs = append(apcs, getColumn("TEMPLATE-NAME", getNestedString(o.u.Object, fieldsAppResTemplateName...)))
	apcs = append(apcs, getColumn("CLUSTER", getNestedString(o.u.Object, fieldsStatusClusterRefName...)))
	apcs = append(apcs, getColumn("STATUS", getNestedString(o.u.Object, fieldsStatusState...)))

	od.AdditionalPrinterColumns = apcs

	od.RemoteStatus = o.getRemoteStatus()

	return od
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (o *ApplicationResource) GetDetails() string {
	remoteStatus := o.getRemoteStatus()

	d := fmt.Sprintf(applicationResourceDetailsTemplate, o.u.GetKind(),
		o.u.GetName(), getNestedString(o.u.Object, fieldsAppResTemplateKind...),
		getNestedString(o.u.Object, fieldsAppResTemplateName...),
		getNestedString(o.u.Object, fieldsStatusClusterRefName...),
		getNestedString(o.u.Object, fieldsStatusState...), remoteStatus)

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
		getNestedString(cMap, conditionKeyType)

		d = d + fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t\n",
			getNestedString(cMap, conditionKeyType),
			getNestedString(cMap, conditionKeyStatus),
			getNestedString(cMap, conditionKeyLastTransitionTime),
			getNestedString(cMap, conditionKeyReason),
			getNestedString(cMap, conditionKeyMessage))
	}
	return d
}

func (o *ApplicationResource) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get resource reference
	u, err := getObjRef(obj.Object, fieldsStatusClusterRef)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	secrets, f, err := unstructured.NestedSlice(obj.Object, fieldsAppResSecrets...)
	if err != nil {
		return related, err
	}
	if f {
		for _, val := range secrets {
			s, ok := val.(map[string]interface{})
			if !ok {
				return related, errors.New("failed to get secret reference in KubernetesApplicationResource: " + obj.GetName())
			}
			u, err := getObjRef(s, []string{})
			if err != nil {
				return related, err
			}
			u.SetAPIVersion("v1")
			u.SetKind("Secret")
			u.SetNamespace(obj.GetNamespace())
			related = append(related, u)
		}
	}

	return related, nil
}

func (o *ApplicationResource) getRemoteStatus() string {
	rs, f, err := unstructured.NestedFieldNoCopy(o.u.Object, fieldsAppResStatusRemote...)
	if err != nil {
		// failed to get conditions
		return fmt.Sprintf("<error: %v>", err)
	}
	if !f {
		return "<error: not found>"
	}

	b, err := yaml.Marshal(rs)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	return string(b)
}
