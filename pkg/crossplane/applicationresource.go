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

	stateAppResSubmitted = "Submitted"
)

type ApplicationResource struct {
	instance *unstructured.Unstructured
}

func NewApplicationResource(u *unstructured.Unstructured) *ApplicationResource {
	return &ApplicationResource{instance: u}
}

func (o *ApplicationResource) GetStatus() string {
	return getNestedString(o.instance.Object, fieldsStatusState...)
}

func (o *ApplicationResource) GetAge() string {
	return GetAge(o.instance)
}

func (o *ApplicationResource) GetObjectDetails() ObjectDetails {
	u := o.instance
	if u == nil {
		return ObjectDetails{}
	}
	od := getObjectDetails(o.instance)

	od.AdditionalStatusColumns = append(od.AdditionalStatusColumns, getColumn("TEMPLATE-KIND", getNestedString(o.instance.Object, fieldsAppResTemplateKind...)))
	od.AdditionalStatusColumns = append(od.AdditionalStatusColumns, getColumn("TEMPLATE-NAME", getNestedString(o.instance.Object, fieldsAppResTemplateName...)))

	od.RemoteStatus = o.getRemoteStatus()

	return od
}

func (o *ApplicationResource) IsReady() bool {
	return o.GetStatus() == stateAppResSubmitted
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

func (o *ApplicationResource) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.instance

	// Get resource reference
	u, err := getObjRef(obj.Object, fieldsStatusClusterRef)
	if err != nil {
		return related, err
	}
	if u != nil {
		related = append(related, u)
	}

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
			if u != nil {
				u.SetAPIVersion("v1")
				u.SetKind("Secret")
				u.SetNamespace(obj.GetNamespace())
				related = append(related, u)
			}
		}
	}

	return related, nil
}

func (o *ApplicationResource) getRemoteStatus() string {
	rs, f, err := unstructured.NestedFieldNoCopy(o.instance.Object, fieldsAppResStatusRemote...)
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
