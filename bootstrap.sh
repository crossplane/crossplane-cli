#!/usr/bin/env bash

set -e

if [[ -z "${PREFIX}" ]]; then
  PREFIX=/usr/local
fi

if [[ -z "${RELEASE}" ]]; then
  RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/crossplane/crossplane-cli/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
fi

PLATFORM="linux"
if [[ "${OSTYPE}" == "darwin"* ]]; then
  PLATFORM="darwin"
fi

set -x

if [[ "${RELEASE}" == "master" || "${RELEASE}" == release-0.1 ]]; then
  set +x
  echo "NOTICE: the trace and pack commands are not available from master. RELEASE must be set to a released version (such as v0.2.0). See https://github.com/crossplane/crossplane-cli/releases for the full list of releases." >&2
  set -x
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-krew-wrapper https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-krew-wrapper >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-build https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-build >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-init https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-init >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-publish https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-publish >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-install https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-uninstall https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-uninstall >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-generate_install https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-generate_install >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-package-list https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-package-list >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-registry https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-registry >/dev/null
  curl -sL -o "${PREFIX}"/bin/kubectl-crossplane-registry-login https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bin/kubectl-crossplane-registry-login >/dev/null
else
  curl -sL https://github.com/crossplane/crossplane-cli/releases/download/"${RELEASE}"/crossplane-cli_"${RELEASE}"_"${PLATFORM}"_amd64.tar.gz | tar -xz -v --strip 1 -C "${PREFIX}"/bin
fi

chmod +x "${PREFIX}"/bin/kubectl-crossplane*
