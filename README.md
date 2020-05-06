# Crossplane CLI

## Installation

Here's the one-liner to install latest released version:

```
curl -sL https://raw.githubusercontent.com/crossplane/crossplane-cli/master/bootstrap.sh | bash
```

The behavior is customizable via environment variables:

```
RELEASE=v0.2.0
PREFIX=${HOME}
curl -sL https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bootstrap.sh | env PREFIX=${PREFIX} RELEASE=${RELEASE} bash
```

You can get the latest, bleeding-edge versions of the `package` commands
by setting `RELEASE` as `master`.

```
RELEASE=master && curl -sL https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bootstrap.sh | RELEASE=${RELEASE} bash
```

But please note, the `trace` command will not be installed in this case,
because it's a built binary and isn't set up to be easily downloaded
from `master`. To use `trace` from `master`, we recommend building and
installing from source.

### Installing from source

If you have the source repository checked out, installing is simple as
long as you have `golang` installed:

```
make build
make install
```

## Uninstallation

If you followed the installation process above, you can remove
everything with:

```
rm /usr/local/bin/kubectl-crossplane*
```

Or, if you customized the installation prefix:

```
PREFIX=/thing
rm "${PREFIX}"/bin/kubectl-crossplane*
```

### Uninstalling from source

If you have the source repository checked out:

```
make uninstall
```

## Release

To release a new version of `crossplane-cli`, follow the instructions in
[`.goreleaser.yml`](https://github.com/crossplane/crossplane-cli/blob/master/.goreleaser.yml).

## Usage

### Package commands

```
kubectl crossplane package init 'myname/mysubname'
kubectl crossplane package build
kubectl crossplane package publish
kubectl crossplane package install 'myname/mysubname'
kubectl crossplane package generate-install 'myname/mysubname' | kubectl apply --namespace mynamespace -f -
kubectl crossplane package list
kubectl crossplane package uninstall 'myname-mysubname'
```

### Registry commands

```
kubectl crossplane registry login 'registry.upbound.io'
```

### Pack command

The `pack` command wraps Kubernetes manifests in a `KubernetesApplication` so
that they can easily be deployed to remote clusters.

Examples:

```
# Wrap helm chart output in KubernetesApplication
helm template crossplane -n crossplane-system crossplane/crossplane-alpha | kubectl crossplane pack -

# Wrap manifests in file in KubernetesApplication
kubectl crossplane pack -f ./deployment/mycoolconfig.yaml
```

### Trace command

Trace command aims to ease debugging and troubleshooting process by providing a holistic view for a particular object.
It finds the relevant objects for a given one and provides detailed information.

```
kubectl crossplane trace TYPE[.GROUP] NAME [-n| --namespace NAMESPACE]
```

Examples:

```
# Trace a KubernetesApplication
kubectl crossplane trace KubernetesApplication wordpress-app-83f04457-0b1b-4532-9691-f55cf6c0da6e -n app-project1-dev

# Trace a MySQLInstance
kubectl crossplane trace MySQLInstance wordpress-mysql-83f04457-0b1b-4532-9691-f55cf6c0da6e -n app-project1-dev

# Graph output, which can be visualized with graphviz as follows.
kubectl crossplane trace KubernetesApplication wordpress-app-83f04457-0b1b-4532-9691-f55cf6c0da6e -n app-project1-dev -o dot | dot -Tpng > /tmp/output.png
```

For more information, see [the trace command documentation](docs/trace-command.md).

# Quick Start: Packages

This guide will show you the basics for using the Packages CLI to create
and develop a package with a basic controller that responds to a certain
type of object.

This guide uses [kubebuilder version 2][kubebuilder quick start], but
other approaches can be used too.

Excluding installing prerequisites, it'll take about 5 minutes!

## Install prerequisites

The workflow for working with a package assumes that you have a few
things installed:

* [go](https://golang.org/doc/install)
* [docker](https://www.docker.com/products/docker-desktop)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [kubebuilder, version 2][kubebuilder quick start]
* [GNU make](https://www.gnu.org/software/make/)
* bash
* [crossplane](https://github.com/crossplane/crossplane),
  or your `kubectl` should be set up to talk to a crossplane
  control cluster

`make` and `bash` are probably already installed if you're using a
Unix-like environment or MacOS.

### Installing Kubebuilder

Instructions to install kubebuilder come from the [kubebuilder quick
start][kubebuilder quick start], but are reproduced here for
convenience.

```
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -sL https://go.kubebuilder.io/dl/2.0.0/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.0.0_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

## Install CLI

Installing the `kubectl` plugins is pretty simple.

Copy the plugins to somewhere on your `PATH`. If you have
`/usr/local/bin` on your `PATH`, you can do it like this:

```
RELEASE=master && curl -sL https://raw.githubusercontent.com/crossplane/crossplane-cli/"${RELEASE}"/bootstrap.sh | RELEASE=${RELEASE} bash
```

## Init project folder

First you'll need a folder to work with. Initializing it probably looks
something like the following:

```
mkdir helloworld
cd helloworld
git init
```

## Init Kubebuilder project

Once the project is initialized, we'll set it up as a kubebuilder
project.

### Create a kubebuilder v2 project

If your project is within the `GOPATH`, you can do:

```
GO111MODULE=on kubebuilder init --domain helloworld.stacks.crossplane.io
```

If your project is not within the `GOPATH`, you need to do an extra step first. See the
[kubebuilder quick start guide][kubebuilder quick start] for more details.

The instructions in this guide will use `GO111MODULE=on` in front of all
of the `kubebuilder` commands, but it could also be set once at the
beginning with `export GO111MODULE=on`, or set in other ways.

### Create an API

Next, create an API using `kubebuilder`:

```
yes y | GO111MODULE=on kubebuilder create api --group samples --version v1alpha1 --kind HelloWorld
```

## Initialize the package project

Once the project is initialized as a kubebuilder project and the plugins
are installed, we can initialize the project as a package project using
the `crossplane package init` command.

From within the project directory:

```
kubectl crossplane package init --cluster 'crossplane-examples/hello-world'
```

## Set up the Hello World

Now we'll add some custom functionality to our controller.

Add a "Hello World" log line to the kubebuilder controller, which should
be in the `controllers/` directory:

```
	// your logic here
	r.Log.V(0).Info("Hello World!", "instance", req.NamespacedName)
```

## Build kubebuilder stuff

First, we'll want to build the binaries and yamls in the regular way.

```
GO111MODULE=on make manager manifests
```

## Validate the package locally

During development, we'll often be building and testing our package
locally. The local validation section shows us how we would do that.

These local scenarios will attempt to run a local Docker registry in your Kubernetes cluster to help
with publishing your Package image for local testing scenarios.  This local validation flow has been
tested and verified with the following software:

* Docker Desktop `2.2.0.0` with Kubernetes `v1.15.5` enabled
* Crossplane `v0.7.0`
  [installed](https://crossplane.io/docs/master/install-crossplane.html#installation) into Docker
  Desktop's Kubernetes cluster

### Build package and publish locally

Once the (kubebuilder) project is built in the normal way, we want to
bundle it all into a package and publish the image locally:

```
kubectl crossplane package build local-build
```

### Install package locally

When the package is built, the next step is to install it into our
Crossplane:

```
kubectl crossplane package install --cluster 'crossplane-examples/hello-world' 'crossplane-examples-hello-world' localhost:5000
```

This can also be done using the sample local package install that the
`init` command generates, but it's a good habit to use the `install`
command.

### Create an object for the package to manage

Once the package is installed into our Crossplane, we can use one of the
sample objects to create an object that the package will respond to.

```
kubectl apply -f config/samples/*.yaml
```

### Check the package's output

Once the object has been created, we expect the controller to do
something in response! In our case, we expect it to log a message.

We can check whether the controller is running and logging messages
using a `kubectl get` to identify the pod running the controller,
followed by a `kubectl logs` to read the logs of that pod:

```
$ kubectl get pods -A | grep hello-world | grep -v Completed
default             crossplane-examples-hello-world-65d5c59976-vzppd   1/1     Running     0          91s
$ kubectl logs crossplane-examples-hello-world-65d5c59976-vzppd | grep 'Hello World'
2019-08-28T23:37:51.795Z	INFO	controllers.HelloWorld	Hello World!	{"instance": "default/helloworld-sample"}
```

That's it! We've finished writing, building, and locally validating a
package!


### Remove the package

When we're done with the package and want to remove it and all its
resources, we can `uninstall` it by name:

```
kubectl crossplane package uninstall --cluster 'crossplane-examples-hello-world'
```

## How to build for external publishing

After we finish developing a package locally, we may want to publish it to
an external registry. This section shows the commands to do that.

### Build package

To build the package, we use the `build` subcommand once the binaries and
generated yamls have been built:

```
kubectl crossplane package build
```

### Run publish

Once the package is built, we can use the `publish` subcommand to publish
it to the registry:

```
# You may need to log into dockerhub or the docker registry that the
# image is being pushed to. If that's the case, run:
# $ docker login
kubectl crossplane package publish
```

### Install

Installing the package can be done with a sample which was generated for
us by the `init` that we ran earlier, but it's a little easier to use
the `install` command:

```
kubectl crossplane package install --cluster 'crossplane-examples/hello-world'
```

### Uninstall

Installing the package can be done with some sample yaml which was
generated for us by the `init` that we ran earlier, but it's a
little easier to use the `uninstall` command:

```
kubectl crossplane package uninstall --cluster 'crossplane-examples-hello-world'
```

Note that `uninstall` uses the package's name (which has no `/` characters),
while the `install` uses the image name (which uses `/`).

# Recipes

### Use a different image name when building and testing a package locally

```
PACKAGE_IMG=myprefix/myothername kubectl crossplane package build
PACKAGE_IMG=myprefix/myothername kubectl crossplane package publish
```

### Build locally with an overridden install.yaml

We can specify a different build target to run:

```
kubectl crossplane package build local-build
```
See the `config/package/overrides` directory for details about where the
overrides live. See the `package.Makefile` for details about how
`local-build` works.

### Setup RBAC

We can setup extra permissions to grant access to the resources which
are not part of our package. This can be done by specifying the permissions in the package definition YAML..
```
# Human readable title of application.
title: Sample Wordpress App
...
# RBAC Roles will be generated permitting this package to use all verbs on all
# resources in the groups listed below.
permissionScope: Namespaced
dependsOn:
- crd: "kubernetesclusters.compute.crossplane.io/v1alpha1"
- crd: "mysqlinstances.database.crossplane.io/v1alpha1"
- crd: "kubernetesapplications.workload.crossplane.io/v1alpha1"
...
```
For a detailed example, see [here](https://github.com/crossplane/app-wordpress/blob/37443b45f40b73958bfd804f98a2a93ba50d8590/.registry/app.yaml#L60).

[kubebuilder quick start]: https://book.kubebuilder.io/quick-start.html
