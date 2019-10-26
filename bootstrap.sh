#!/usr/bin/env bash

set -e

if [[ -z "${PREFIX}" ]]; then
  PREFIX=/usr/local
fi

if [[ -z "${RELEASE}" ]]; then
  RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/crossplaneio/crossplane-cli/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
fi

PLATFORM="linux"
if [[ "${OSTYPE}" == "darwin"* ]]; then
  PLATFORM="darwin"
fi

set -x

if [[ "${RELEASE}" == "master" ]]; then
  echo "NOTICE: the trace command is not available from master. RELEASE must be set to a released version (such as v0.2.0). See https://github.com/crossplaneio/crossplane-cli/releases for the full list of releases." >&2
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-build https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-build >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-init https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-init >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-publish https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-publish >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-uninstall https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-uninstall >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-generate_install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-generate_install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-list https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-list >/dev/null
else
  curl -sL https://github.com/crossplaneio/crossplane-cli/releases/download/"${RELEASE}"/crossplane-cli_"${RELEASE}"_"${PLATFORM}"_amd64.tar.gz | tar -xz -v --strip 1 -C "${PREFIX}"/bin
fi

chmod +x "${PREFIX}"/bin/kubectl-crossplane-*
