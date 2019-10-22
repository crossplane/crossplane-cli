#!/usr/bin/env bash

set -e

if [[ -z "${PREFIX}" ]]; then
  PREFIX=/usr/local
fi

if [[ -z "${RELEASE}" ]]; then
  RELEASE=master
fi

PLATFORM="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
  PLATFORM="darwin"
fi

set -x

if [[ "${RELEASE}" == "master" ]]; then
  echo "trace subcommand will not be available from master, RELEASE must be set to a released version."
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-build https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-build >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-init https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-init >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-publish https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-publish >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-uninstall https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-uninstall >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-generate_install https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-generate_install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-stack-list https://raw.githubusercontent.com/crossplaneio/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-stack-list >/dev/null
else
  curl -sL https://github.com/crossplaneio/crossplane-cli/releases/download/"${RELEASE}"/crossplane-cli-"${RELEASE}"-"${PLATFORM}".tar.gz | tar xz --strip 1 -C "${PREFIX}"/bin
fi
chmod +x "${PREFIX}"/bin/kubectl-crossplane-stack-*
