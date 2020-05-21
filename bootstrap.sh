#!/usr/bin/env bash
set -e

if [[ -z "${RELEASE}" ]]; then
  RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/crossplane/crossplane-cli/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
fi

PLATFORM="linux"
if [[ "${OSTYPE}" == "darwin"* ]]; then
  PLATFORM="darwin"
fi

curl -sL https://github.com/crossplane/crossplane-cli/releases/download/"${RELEASE}"/crossplane-cli_"${RELEASE}"_"${PLATFORM}"_amd64.tar.gz \
  | tar -xz -v --strip 1 -C /usr/local/bin 2>&1 \
  | sed 's/x /✓ /g'

chmod +x /usr/local/bin/kubectl-crossplane*

printf "👍 Crossplane CLI installed successfully!"
printf "\n\nHave a nice day! 👋\n"
