package trace

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/crossplaneio/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubeGraphBuilder_fetchObj(t *testing.T) {
	scheme := runtime.NewScheme()

	type args struct {
		client           dynamic.Interface
		instanceToCreate *unstructured.Unstructured
		n                *Node
	}
	type want struct {
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"NotExist": {
			args: args{
				client:           fake.NewSimpleDynamicClient(scheme),
				instanceToCreate: nil,
				n: &Node{
					Instance: &unstructured.Unstructured{Object: map[string]interface{}{
						"kind":       "KubernetesCluster",
						"apiVersion": "compute.crossplane.io/v1alpha1",
						"metadata": map[string]interface{}{
							"name":      "test",
							"namespace": "testnamespace",
						},
					}},
					GVR: schema.GroupVersionResource{
						Group:    "compute.crossplane.io",
						Version:  "v1alpha1",
						Resource: "kubernetescluster",
					},
					Relateds: nil,
					State:    "",
				},
			},
			want: want{
				err: errors.NewNotFound(schema.ParseGroupResource("kubernetescluster.compute.crossplane.io"), "test"),
			},
		},
		"ExistsNamespaced": {
			args: args{
				client: fake.NewSimpleDynamicClient(scheme),
				instanceToCreate: &unstructured.Unstructured{Object: map[string]interface{}{
					"kind":       "KubernetesCluster",
					"apiVersion": "compute.crossplane.io/v1alpha1",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "testnamespace",
					},
					"spec": map[string]interface{}{
						"numNodes": int64(1),
					},
					"status": map[string]interface{}{
						"State": "ready",
					},
				}},
				n: &Node{
					Instance: &unstructured.Unstructured{Object: map[string]interface{}{
						"kind":       "KubernetesCluster",
						"apiVersion": "compute.crossplane.io/v1alpha1",
						"metadata": map[string]interface{}{
							"name":      "test",
							"namespace": "testnamespace",
						},
					}},
					GVR: schema.GroupVersionResource{
						Group:    "compute.crossplane.io",
						Version:  "v1alpha1",
						Resource: "kubernetescluster",
					},
					Relateds: nil,
					State:    "",
				},
			},
			want: want{
				err: nil,
			},
		},
		"ExistsClusterScoped": {
			args: args{
				client: fake.NewSimpleDynamicClient(scheme),
				instanceToCreate: &unstructured.Unstructured{Object: map[string]interface{}{
					"kind":       "StorageClass",
					"apiVersion": "storage.k8s.io/v1",
					"metadata": map[string]interface{}{
						"name": "standard",
					},
					"reclaimPolicy": "Delete",
				}},
				n: &Node{
					Instance: &unstructured.Unstructured{Object: map[string]interface{}{
						"kind":       "StorageClass",
						"apiVersion": "storage.k8s.io/v1",
						"metadata": map[string]interface{}{
							"name": "standard",
						},
					}},
					GVR: schema.GroupVersionResource{
						Group:    "storage.k8s.io",
						Version:  "v1",
						Resource: "storageclass",
					},
					Relateds: nil,
					State:    "",
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fc := tc.client
			g := &KubeGraphBuilder{
				client: fc,
			}
			i := tc.instanceToCreate
			var createdInstance *unstructured.Unstructured
			if i != nil {
				var err error
				if i.GetNamespace() == "" {
					createdInstance, err = fc.Resource(tc.n.GVR).Create(i, metav1.CreateOptions{})
				} else {
					createdInstance, err = fc.Resource(tc.n.GVR).Namespace(i.GetNamespace()).Create(i, metav1.CreateOptions{})
				}

				if err != nil {
					t.Fatalf("failed to prepare fake env, could not create Instance: %v", err)
				}
			}
			gotErr := g.fetchObj(tc.args.n)

			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("g.fetchObj(...): -want error, +got error: %s", diff)
			}
			if tc.err != nil {
				return
			}
			gotResult := tc.args.n.Instance
			if diff := cmp.Diff(createdInstance, gotResult); diff != "" {
				t.Errorf("g.fetchObj(...): -want result, +got result: %s", diff)
			}
		})
	}
}

func TestKubeGraphBuilder_BuildGraph(t *testing.T) {
	scheme := runtime.NewScheme()

	gvrK8S := schema.GroupVersionResource{
		Group:    "compute.crossplane.io",
		Version:  "v1alpha1",
		Resource: "kubernetesclusters",
	}

	computeGV := schema.GroupVersion{
		Group:   "compute.crossplane.io",
		Version: "v1alpha1",
	}

	rm := meta.NewDefaultRESTMapper([]schema.GroupVersion{
		computeGV,
	})
	rm.Add(computeGV.WithKind("kubernetescluster"), meta.RESTScopeNamespace)

	type args struct {
		client           dynamic.Interface
		restMapper       meta.RESTMapper
		instanceToCreate *unstructured.Unstructured
		name             string
		namespace        string
		resource         string
		gvr              schema.GroupVersionResource
	}
	type want struct {
		traversed []*Node
		err       error
	}
	cases := map[string]struct {
		args
		want
	}{
		"NotExist": {
			args: args{
				client:           fake.NewSimpleDynamicClient(scheme),
				restMapper:       rm,
				instanceToCreate: nil,
				name:             "test",
				namespace:        "testnamespace",
				resource:         "kubernetescluster",
			},
			want: want{
				err: errors.NewNotFound(schema.ParseGroupResource("kubernetesclusters.compute.crossplane.io"), "test"),
			},
		},
		"ExistsNamespaced": {
			args: args{
				client:           fake.NewSimpleDynamicClient(scheme),
				restMapper:       rm,
				instanceToCreate: getTestInstanceK8SCluster(),
				name:             "test",
				namespace:        "testnamespace",
				resource:         "KubernetesCluster",
				gvr:              gvrK8S,
			},
			want: want{
				err: nil,
				traversed: []*Node{
					{
						Instance: getTestInstanceK8SCluster(),
						GVR:      gvrK8S,
						Relateds: []*Node{},
						State:    NodeStateNotReady,
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fc := tc.client
			g := NewKubeGraphBuilder(fc, tc.restMapper)

			i := tc.instanceToCreate
			if i != nil {
				_, err := fc.Resource(tc.gvr).Namespace(i.GetNamespace()).Create(i, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to prepare test environment, could not create Instance: %v", err)
				}
			}
			_, trvs, gotErr := g.BuildGraph(tc.name, tc.namespace, tc.resource)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("g.BuildGraph(...): -want error, +got error: %s", diff)
			}
			if tc.err != nil {
				return
			}
			gotResult := trvs
			if diff := cmp.Diff(tc.traversed, gotResult); diff != "" {
				t.Errorf("g.fetchObj(...): -want result, +got result: %s", diff)
			}
		})
	}
}

func getTestInstanceK8SCluster() *unstructured.Unstructured {
	i := getTestInstance()
	i.SetAPIVersion("compute.crossplane.io/v1alpha1")
	i.SetKind("KubernetesCluster")
	return i
}

func getTestInstance() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      "test",
				"namespace": "testnamespace",
			},
			"spec": map[string]interface{}{
				"numNodes": int64(1),
			},
			"status": map[string]interface{}{
				"State": "ready",
			},
		},
	}
}
