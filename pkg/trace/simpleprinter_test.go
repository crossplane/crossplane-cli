package trace

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSimplePrinter_getDetailsText(t *testing.T) {
	type args struct {
		node *Node
	}
	type want struct {
		result string
	}
	cases := map[string]struct {
		args
		want
	}{
		"NilObject": {
			args: args{node: nil},
			want: want{
				result: "<error: node to trace is nil>",
			},
		},
		"Simple": {
			args: args{&Node{
				instance: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "workload.crossplane.io/v1alpha1",
						"kind":       "KubernetesApplicationResource",
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
				},
				gvr:     schema.GroupVersionResource{},
				id:      "",
				related: nil,
				state:   "",
			}},
			want: want{
				result: `
KubernetesApplicationResource: test

STATE	TEMPLATE-KIND	TEMPLATE-NAME	
Scheduled	Deployment	wordpress	

Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
Synced	True	2019-10-21T12:50:17Z	Successfully reconciled managed resource	Test message of Synced	
Ready	False	2019-10-21T12:54:35Z	Managed resource is available for use	Test message of Ready	

Remote Status
loadBalancer:
  ingress:
  - ip: 1.2.3.4

`,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gotResult := getDetailsText(tc.args.node)
			if diff := cmp.Diff(tc.want.result, gotResult); diff != "" {
				t.Errorf("getDetailsText(...): -want result, +got result: %s", diff)
			}
		})
	}
}
