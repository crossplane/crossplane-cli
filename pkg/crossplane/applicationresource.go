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
	return getNestedString(o.u.Object, "status", "state")
}

func (o *ApplicationResource) GetAge() string {
	return GetAge(o.u)
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
		o.u.GetName(), getNestedString(o.u.Object, "spec", "template", "kind"),
		getNestedString(o.u.Object, "spec", "template", "metadata", "name"),
		getNestedString(o.u.Object, "status", "clusterRef", "name"),
		getNestedString(o.u.Object, "status", "state"), remoteStatus)

	cs, f, err := unstructured.NestedSlice(o.u.Object, "status", "conditionedStatus", "conditions")
	if err != nil || !f {
		// failed to get conditions
		return d
	}
	for _, c := range cs {
		cMap := c.(map[string]interface{})
		if cMap == nil {
			d = d + "<error>"
			continue
		}
		getNestedString(cMap, "type")

		d = d + fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t\n",
			getNestedString(cMap, "type"),
			getNestedString(cMap, "status"),
			getNestedString(cMap, "lastTransitionTime"),
			getNestedString(cMap, "reason"),
			getNestedString(cMap, "message"))
	}
	return d
}

func (o *ApplicationResource) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get resource reference
	u, err := getObjRef(obj.Object, applicationClusterRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	secrets, f, err := unstructured.NestedSlice(obj.Object, "spec", "secrets")
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
	rs, f, err := unstructured.NestedFieldNoCopy(o.u.Object, "status", "remote")
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
