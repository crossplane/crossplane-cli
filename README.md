# Crossplane Stacks CLI

## Installation

Get the commands on the PATH so `kubectl` can pick them up:
```
curl -o /usr/local/bin/kubectl-crossplane-stack-build https://raw.githubusercontent.com/suskin/crossplane-stack-cli/first-cli/bin/kubectl-crossplane-stack-build -s >/dev/null
curl -o /usr/local/bin/kubectl-crossplane-stack-init https://raw.githubusercontent.com/suskin/crossplane-stack-cli/first-cli/bin/kubectl-crossplane-stack-init -s >/dev/null
curl -o /usr/local/bin/kubectl-crossplane-stack-publish https://raw.githubusercontent.com/suskin/crossplane-stack-cli/first-cli/bin/kubectl-crossplane-stack-publish -s >/dev/null
chmod +x /usr/local/bin/kubectl-crossplane-stack-*
```

## Usage

```
kubectl crossplane stack init 'myname/mysubname'
kubectl crossplane stack build
kubectl crossplane stack publish
```



