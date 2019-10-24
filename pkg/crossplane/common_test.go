package crossplane

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_getObjectDetails(t *testing.T) {
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
		"NoStatus": {
			args: args{u: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind": "KubernetesApplication",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "unittest",
					},
				},
			}},
			want: want{
				result: ObjectDetails{
					Kind:      "KubernetesApplication",
					Name:      "test",
					Namespace: "unittest",
				},
			},
		},
		"AllCommonFields": {
			args: args{u: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind": "KubernetesApplication",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "unittest",
					},
					"status": map[string]interface{}{
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
			gotResult := getObjectDetails(tc.args.u)
			if diff := cmp.Diff(tc.want.result, gotResult); diff != "" {
				t.Errorf("getObjectDetails(...): -want result, +got result: %s", diff)
			}
		})
	}
}
