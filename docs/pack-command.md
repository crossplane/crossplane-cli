# Pack Command

The Crossplane `pack` command is used to package objects defined in Kubernetes
manifests into a Crossplane `KubernetesApplication` object by embedding the
objects as `KubernetesApplicationResourceTemplates`. This allows a grouping of
resources to be scheduled from a Crossplane Kubernetes control cluster to a
remote Kubernetes cluster that it has access to.

## Usage

After installation, the `pack` command can be invoked by running `kubectl
crossplane pack`. The command is commonly used by piping the output from a
previous command to it. For example, templating a [Helm](https://helm.sh/)
chart:

```
helm template crossplane -n crossplane-system crossplane/crossplane-alpha | kubectl crossplane pack -
```

The command can also be invoked with the path to an existing manifest file:

```
kubectl crossplane pack -f ./deployment/mycoolconfig.yaml
```
