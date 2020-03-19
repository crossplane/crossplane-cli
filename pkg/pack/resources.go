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
	"fmt"
	"io"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/crossplane/crossplane/apis/workload/v1alpha1"
)

// ReadResources reads Kubernetes objects and packages them into a slice of
// KubernetesApplicationResourceTemplates.
func ReadResources(prefix string, content io.Reader, resLabels map[string]string) ([]v1alpha1.KubernetesApplicationResourceTemplate, error) {
	var result []v1alpha1.KubernetesApplicationResourceTemplate
	d := kyaml.NewYAMLOrJSONDecoder(content, 4096)
	for {
		obj := &unstructured.Unstructured{}
		if err := d.Decode(obj); err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}
		// Ignore empty objects
		if obj.GetName() == "" {
			continue
		}
		kart := v1alpha1.KubernetesApplicationResourceTemplate{
			ObjectMeta: v1.ObjectMeta{
				// Ensure no collisions among templates in KubernetesApplications
				Name:   fmt.Sprintf("%s-%s-%s", prefix, strings.ReplaceAll(obj.GetName(), ":", "-"), strings.ToLower(obj.GetKind())),
				Labels: resLabels,
			},
			Spec: v1alpha1.KubernetesApplicationResourceSpec{
				Template: obj,
			},
		}
		result = append(result, kart)
	}
	return result, nil
}
