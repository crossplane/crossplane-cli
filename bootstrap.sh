#!/usr/bin/env bash

set -e

if [[ -z "${PREFIX}" ]]; then
  PREFIX=/usr/local
fi

if [[ -z "${RELEASE}" ]]; then
  RELEASE=master
fi

set -x

curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-build https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-build >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-init https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-init >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-publish https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-publish >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-install >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-uninstall https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-uninstall >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-generate_install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-generate_install >/dev/null
curl -s -o "${PREFIX}"/bin/kubectl-crossplane-stack-list https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-list >/dev/null
chmod +x "${PREFIX}"/bin/kubectl-crossplane-stack-*
