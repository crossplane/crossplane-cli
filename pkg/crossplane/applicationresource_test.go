package crossplane

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestApplicationResource_GetObjectDetails(t *testing.T) {
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
					"kind": "KubernetesApplicationResource",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "unittest",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"kind": "Deployment",
							"metadata": map[string]interface{}{
								"name": "wordpress",
							},
						},
					},
					"status": map[string]interface{}{
						"state": "Scheduled",
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
						"remote": map[string]interface{}{
							"loadBalancer": map[string]interface{}{
								"ingress": []interface{}{
									map[string]interface{}{
										"ip": "1.2.3.4",
									},
								},
							},
						},
					},
				},
			}},
			want: want{
				result: ObjectDetails{
					Kind:      "KubernetesApplicationResource",
					Name:      "test",
					Namespace: "unittest",
					AdditionalPrinterColumns: []map[string]string{
						getColumn("NAME", "test"),
						getColumn("TEMPLATE-KIND", "Deployment"),
						getColumn("TEMPLATE-NAME", "wordpress"),
						getColumn("CLUSTER", "test-cluster"),
						getColumn("STATUS", "Scheduled"),
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
					RemoteStatus: `loadBalancer:
  ingress:
  - ip: 1.2.3.4
`,
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := &ApplicationResource{
				u: tc.args.u,
			}
			gotResult := o.GetObjectDetails()
			if diff := cmp.Diff(tc.want.result, gotResult); diff != "" {
				t.Errorf("GetObjectDetails(...): -want result, +got result: %s", diff)
			}
		})
	}
}
