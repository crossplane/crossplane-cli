/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pack

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/crossplane/crossplane/apis/workload/v1alpha1"
)

var namespaceBytes = []byte(`
apiVersion: v1
kind: Namespace
metadata:
  name: test
  annotations:
    cool: annotation
  labels:
    cool: label
`)

func TestRunPack(t *testing.T) {
	namespaceReader, _ := os.Open("./testdata/namespace.yaml")
	blankReader, _ := os.Open("./testdata/blank.yaml")

	u := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(namespaceBytes, u); err != nil {
		t.Error(err)
	}

	type args struct {
		prefix    string
		content   io.Reader
		resLabels map[string]string
	}
	type want struct {
		out []v1alpha1.KubernetesApplicationResourceTemplate
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulBlank": {
			args: args{
				prefix:    "test",
				content:   blankReader,
				resLabels: map[string]string{},
			},
			want: want{
				out: nil,
				err: nil,
			},
		},
		"SuccessfulNamespace": {
			args: args{
				prefix:    "test",
				content:   namespaceReader,
				resLabels: map[string]string{"test": "test"},
			},
			want: want{
				out: []v1alpha1.KubernetesApplicationResourceTemplate{
					{
						ObjectMeta: v1.ObjectMeta{
							Name:   "test-test-namespace",
							Labels: map[string]string{"test": "test"},
						},
						Spec: v1alpha1.KubernetesApplicationResourceSpec{
							Template: u,
						},
					},
				},
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := ReadResources(tc.args.prefix, tc.args.content, tc.args.resLabels)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("ReadResources(...): -want error, +got error:\n%s", diff)
			}

			if diff := cmp.Diff(tc.want.out, out); diff != "" {
				t.Errorf("ReadResources(...) Output: -want, +got:\n%s", diff)
			}
		})
	}
}
