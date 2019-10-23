package crossplane

import (
	"fmt"
	"strconv"
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
	fieldsStatusConditions            = append(fieldsStatus, "conditions")
	fieldsStatusClusterRef            = append(fieldsStatus, "clusterRef")
	fieldsStatusClusterRefName        = append(fieldsStatusClusterRef, "name")
	fieldsStatusBindingPhase          = append(fieldsStatus, "bindingPhase")

	conditionKeyType               = "type"
	conditionKeyStatus             = "status"
	conditionKeyLastTransitionTime = "lastTransitionTime"
	conditionKeyReason             = "reason"
	conditionKeyMessage            = "message"

	keyColumnName  = "name"
	keyColumnValue = "value"

	resourceBindingPhaseBound = "Bound"
)

func GetAge(u *unstructured.Unstructured) string {
	ts := u.GetCreationTimestamp()
	if ts.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(ts.Time))
}

func getResourceStatus(u *unstructured.Unstructured) string {
	return getNestedString(u.Object, fieldsStatusBindingPhase...)
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
		// ignore if there is no referred object
		return nil, nil
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
		return ""
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

func getColumn(name, value string) map[string]string {
	c := make(map[string]string)
	c[keyColumnName] = name
	c[keyColumnValue] = value
	return c
}

func getObjectDetails(u *unstructured.Unstructured) ObjectDetails {
	if u == nil {
		return ObjectDetails{}
	}
	od := ObjectDetails{
		Kind:      u.GetKind(),
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
	}

	conditions := make([]map[string]string, 0)
	cs, f, err := unstructured.NestedSlice(u.Object, fieldsStatusConditions...)
	if err == nil && f {
		for _, c := range cs {
			condition := make(map[string]string)
			cMap := c.(map[string]interface{})
			if cMap == nil {
				condition[conditionKeyMessage] = "<error: condition status is not a map>"
				continue
			}

			condition[conditionKeyType] = getNestedString(cMap, conditionKeyType)
			condition[conditionKeyStatus] = getNestedString(cMap, conditionKeyStatus)
			condition[conditionKeyLastTransitionTime] = getNestedString(cMap, conditionKeyLastTransitionTime)
			condition[conditionKeyReason] = getNestedString(cMap, conditionKeyReason)
			condition[conditionKeyMessage] = getNestedString(cMap, conditionKeyMessage)

			conditions = append(conditions, condition)
		}
	}
	csc, f, err := unstructured.NestedSlice(u.Object, fieldsConditionedStatusConditions...)
	if err == nil && f {
		for _, c := range csc {
			condition := make(map[string]string)
			cMap := c.(map[string]interface{})
			if cMap == nil {
				condition[conditionKeyMessage] = "<error: condition status is not a map>"
				continue
			}

			condition[conditionKeyType] = getNestedString(cMap, conditionKeyType)
			condition[conditionKeyStatus] = getNestedString(cMap, conditionKeyStatus)
			condition[conditionKeyLastTransitionTime] = getNestedString(cMap, conditionKeyLastTransitionTime)
			condition[conditionKeyReason] = getNestedString(cMap, conditionKeyReason)
			condition[conditionKeyMessage] = getNestedString(cMap, conditionKeyMessage)

			conditions = append(conditions, condition)
		}
	}
	if len(conditions) > 0 {
		od.Conditions = conditions
	}

	asc := make([]map[string]string, 0)
	st, f, err := unstructured.NestedMap(u.Object, fieldsStatus...)
	if err == nil && f {
		for k, _ := range st {
			val, f, err := unstructured.NestedString(st, k)
			if err == nil && f {
				info := make(map[string]string)
				info["name"] = strings.ToUpper(k)
				info["value"] = val
				asc = append(asc, info)
			}
			ival, f, err := unstructured.NestedInt64(st, k)
			if err == nil && f {
				info := make(map[string]string)
				info["name"] = strings.ToUpper(k)
				info["value"] = strconv.Itoa(int(ival))
				asc = append(asc, info)
			}
		}
	}
	if len(asc) > 0 {
		od.AdditionalStatusColumns = asc
	}
	return od
}
