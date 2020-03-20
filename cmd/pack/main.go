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

package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/spf13/pflag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/crossplane/crossplane-cli/pkg/pack"
	"github.com/crossplane/crossplane/apis/workload/v1alpha1"
)

type packOptions struct {
	help      bool
	useFile   string
	name      string
	namespace string
}

var pOpts = &packOptions{}

func init() {
	pflag.BoolVarP(&pOpts.help, "help", "h", false, "Shows this help message")
	pflag.StringVarP(&pOpts.useFile, "file", "f", "", "Specifies file path to read from")
	pflag.StringVar(&pOpts.name, "name", fmt.Sprintf("app-%s", fmt.Sprint(rand.Intn(50000))), "Name")
	pflag.StringVarP(&pOpts.namespace, "namespace", "n", "default", "Namespace")
	pflag.Parse()
}

func main() {
	data, err := runPack(pOpts)
	if err != nil {
		panic(err)
	}

	if data != nil {
		// TODO(hasheddan): allow writing output to path using flag
		fmt.Print(string(data))
	}

	os.Exit(0)
}

func runPack(opts *packOptions) ([]byte, error) {
	if opts.help {
		printHelp()
		return nil, nil
	}

	var err error
	var reader io.Reader
	reader = os.Stdin

	if opts.useFile != "" {
		reader, err = os.Open(opts.useFile)
		if err != nil {
			return nil, err
		}
	}

	resLabels := map[string]string{"crossplane-pack": opts.name}
	resources, err := pack.ReadResources(opts.name, reader, resLabels)
	if err != nil {
		return nil, err
	}
	kapp := &v1alpha1.KubernetesApplication{
		ObjectMeta: v1.ObjectMeta{
			Name:      opts.name,
			Namespace: opts.namespace,
		},
		Spec: v1alpha1.KubernetesApplicationSpec{
			ResourceTemplates: resources,
		},
	}
	kapp.SetGroupVersionKind(v1alpha1.KubernetesApplicationGroupVersionKind)
	kapp.Spec.ResourceSelector = &v1.LabelSelector{
		MatchLabels: resLabels,
	}

	// TODO(hasheddan): allow target selectors to be provided as flag
	kapp.Spec.TargetSelector = &v1.LabelSelector{}
	data, err := yaml.Marshal(kapp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `
		"Package resources into a KubernetesApplication.
	`)
	pflag.PrintDefaults()
}
