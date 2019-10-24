package crossplane

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestApplication_GetObjectDetails(t *testing.T) {
	type args struct {
		u *unstructured.Unstructured
	}
	type want struct {
		result ObjectDetails
	}
	cases := map[string]struct {
		args
		want
	}{
		"NilObject": {
			args: args{u: nil},
			want: want{
				result: ObjectDetails{},
			},
		},
		"Simple": {
			args: args{u: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind": "KubernetesApplication",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "unittest",
					},
					"status": map[string]interface{}{
						"state":              "Submitted",
						"desiredResources":   int64(3),
						"submittedResources": int64(2),
						"clusterRef": map[string]interface{}{
							"name": "test-cluster",
						},
						"conditionedStatus": map[string]interface{}{
							"conditions": []interface{}{
								map[string]interface{}{
									"type":               "Synced",
									"status":             "True",
									"lastTransitionTime": "2019-10-21T12:50:17Z",
									"reason":             "Successfully reconciled managed resource",
									"message":            "Test message of Synced",
								},
								map[string]interface{}{
									"type":               "Ready",
									"status":             "False",
									"lastTransitionTime": "2019-10-21T12:54:35Z",
									"reason":             "Managed resource is available for use",
									"message":            "Test message of Ready",
								},
							},
						},
					},
				},
			}},
			want: want{
				result: ObjectDetails{
					Kind:      "KubernetesApplication",
					Name:      "test",
					Namespace: "unittest",
					AdditionalStatusColumns: []map[string]string{
						getColumn("DESIREDRESOURCES", "3"),
						getColumn("STATE", "Submitted"),
						getColumn("SUBMITTEDRESOURCES", "2"),
					},
					Conditions: []map[string]string{
						{
							"type":               "Synced",
							"status":             "True",
							"lastTransitionTime": "2019-10-21T12:50:17Z",
							"reason":             "Successfully reconciled managed resource",
							"message":            "Test message of Synced",
						},
						{
							"type":               "Ready",
							"status":             "False",
							"lastTransitionTime": "2019-10-21T12:54:35Z",
							"reason":             "Managed resource is available for use",
							"message":            "Test message of Ready",
						},
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := &Application{
				instance: tc.args.u,
			}
			gotResult := o.GetObjectDetails()
			if diff := cmp.Diff(tc.want.result, gotResult,
				cmp.Options{cmpopts.SortSlices(func(i, j map[string]string) bool { return i["name"] < j["name"] })}); diff != "" {
				t.Errorf("GetObjectDetails(...): -want result, +got result: %s", diff)
			}
		})
	}
}
