package crossplane

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/duration"
)

var (
	fieldsSpec                        = []string{"spec"}
	fieldsStatus                      = []string{"status"}
	fieldsStatusState                 = append(fieldsStatus, "state")
	fieldsResourceClass               = append(fieldsSpec, "classRef")
	fieldsWriteConnSecret             = append(fieldsSpec, "writeConnectionSecretToRef")
	fieldsConditionedStatus           = append(fieldsStatus, "conditionedStatus")
	fieldsConditionedStatusConditions = append(fieldsConditionedStatus, "conditions")
	fieldsStatusClusterRef            = append(fieldsStatus, "clusterRef")
	fieldsStatusClusterRefName        = append(fieldsStatusClusterRef, "name")

	conditionKeyType               = "type"
	conditionKeyStatus             = "status"
	conditionKeyLastTransitionTime = "lastTransitionTime"
	conditionKeyReason             = "reason"
	conditionKeyMessage            = "message"

	resourceDetailsTemplate = `%v: %v

State: %v
State Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

func GetAge(u *unstructured.Unstructured) string {
	ts := u.GetCreationTimestamp()
	if ts.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(ts.Time))
}

func getResourceStatus(u *unstructured.Unstructured) string {
	return getNestedString(u.Object, "status", "bindingPhase")
}

func getResourceDetails(u *unstructured.Unstructured) string {
	d := fmt.Sprintf(resourceDetailsTemplate, u.GetKind(), u.GetName(), getResourceStatus(u))
	cs, f, err := unstructured.NestedSlice(u.Object, "status", "conditions")
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

func getObjRef(obj map[string]interface{}, path []string) (*unstructured.Unstructured, error) {
	a, aFound, err := unstructured.NestedString(obj, append(path, "apiVersion")...)
	if err != nil {
		return nil, err
	}
	k, kFound, err := unstructured.NestedString(obj, append(path, "kind")...)
	if err != nil {
		return nil, err
	}
	n, nFound, err := unstructured.NestedString(obj, append(path, "name")...)
	if err != nil {
		return nil, err
	}
	ns, nsFound, err := unstructured.NestedString(obj, append(path, "namespace")...)
	if err != nil {
		return nil, err
	}

	if !aFound && !kFound && !nFound && !nsFound {
		return nil, errors.New("Failed to find a reference!")
	}

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	u.SetAPIVersion(a)
	u.SetKind(k)
	u.SetName(n)
	u.SetNamespace(ns)

	return u, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func normalizedGroupKind(gvk schema.GroupVersionKind) string {
	return strings.ToLower(fmt.Sprintf("%s.%s", gvk.Kind, gvk.Group))
}

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, found, err := unstructured.NestedString(obj, fields...)
	if err != nil {
		return "<unknown>"
	}
	if !found {
		return " "
	}
	return val
}

func getNestedLabelSelector(obj map[string]interface{}, fields ...string) string {
	selMap, found, err := unstructured.NestedMap(obj, fields...)

	if err != nil {
		return "<unknown>"
	}
	if !found {
		return ""
	}
	selector := ""
	for k, v := range selMap {
		if selector != "" {
			selector += ","
		}
		val := v.(string)
		selector += fmt.Sprintf("%s=%s", k, val)
	}
	return selector
}

func getNestedInt64(obj map[string]interface{}, fields ...string) int64 {
	val, found, err := unstructured.NestedInt64(obj, fields...)
	if !found || err != nil {
		return -1
	}
	return val
}

func getKindsFromGroupKinds(allGks ...[]string) []string {
	allKinds := make([]string, 0)
	for _, gks := range allGks {
		for _, gk := range gks {
			gkp := schema.ParseGroupKind(gk)
			allKinds = append(allKinds, gkp.Kind)
		}
	}
	return allKinds
}
