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
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-runtime/pkg/test"
)

func TestRunPack(t *testing.T) {
	crossplaneOut, _ := ioutil.ReadFile("./testdata/crossplane_out.yaml")

	type args struct {
		options *packOptions
	}
	type want struct {
		out []byte
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulCrossplane": {
			args: args{
				options: &packOptions{
					useFile:   "./testdata/crossplane.yaml",
					name:      "crossplane-install",
					namespace: "crossplane-system",
				},
			},
			want: want{
				out: crossplaneOut,
			},
		},
		"FailureBadPath": {
			args: args{
				options: &packOptions{
					useFile: "./does/not/exist.yaml",
				},
			},
			want: want{
				err: &os.PathError{
					Op:   "open",
					Path: "./does/not/exist.yaml",
					Err:  errors.New("no such file or directory"),
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := runPack(tc.args.options)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("runPack(...): -want error, +got error:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.out, out); diff != "" {
				t.Errorf("runPack(...) Output: -want, +got:\n%s", diff)
			}
		})
	}
}
