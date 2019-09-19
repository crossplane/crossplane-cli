#!/usr/bin/env bash
#
# Sanity test for the CLI. Not meant to be extremely robust
#
# Makes a bunch of assumptions about local setup:
# * Prereqs installed
# * Crossplane up and running and accessible via kubectl

set -e
set -x

TEST_DIR="${1:-test}"

mkdir -p "${TEST_DIR}"
cd "${TEST_DIR}"

GO111MODULE=on kubebuilder init --domain helloworld.stacks.crossplane.io
yes | GO111MODULE=on kubebuilder create api --group samples --version v1alpha1 --kind HelloWorld

kubectl crossplane stack init --cluster 'crossplane-examples/hello-world'

sed -i'.bak' -e 's/\/\/ your logic here/r.Log.V(0).Info("Hello World!", "instance", req.NamespacedName)/' controllers/helloworld_controller.go
rm -f controllers/helloworld_controller.go.bak

GO111MODULE=on make manager manifests
kubectl crossplane stack build local-build
kubectl crossplane stack install --cluster 'crossplane-examples/hello-world' 'crossplane-examples-hello-world' localhost:5000

finished='false'
set +e
for i in $( seq 1 10 ); do
  echo "Attempt ${i} to create resource . . ." >&2
  kubectl apply -f config/samples/*.yaml

  if [[ $? -ne 0 ]]; then
    echo "Create failed. Waiting before retry. . ." >&2
    sleep 2
    continue
  else
    finished='true'
    echo "Create SUCCESS" >&2
    break
  fi
done
set -e

if [[ "${finished}" == 'false' ]]; then
  echo "Error: Couldn't create resource after retries!" >&2
  exit 1
fi

pod_name="$( kubectl get pods -A | grep hello-world | grep -v Completed | awk '{ print $2 }' )"
kubectl logs "${pod_name}" | grep 'Hello World'
