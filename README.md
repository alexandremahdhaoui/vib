# vib

Vib allows users to intuitively manage and share their shell environments
with their teams, organization and across any platforms.

## Table of Contents

- [Concepts](#concepts)
- [Repository Structure](#repository-structure)
- [Getting started](#getting-started)
  - [Installation (script)](#installation-script)
  - [Manual installation](#manual-installation)
- [Configure vib](#configure-vib)
  - [Your first ExpressionSet](#your-first-expressionset)
  - [Create more ExpressionSets](#create-more-expressionsets)
  - [Edit your Profile](#edit-your-profile)
- [vib's commands](#vib-s-commands)
- [See Also](#see-also)

## Concepts

`vib` is built around three core concepts which are defined as Kubernetes-style resources. For more details on the API definitions, see the [`pkg/apis/v1alpha1`](./pkg/apis/v1alpha1/README.md) package documentation.

*   **ExpressionSet**: A set of expressions that can be rendered into a desired output. An `ExpressionSet` is a collection of key-value pairs or arbitrary keys that are processed by a `Resolver`.
*   **Resolver**: A special resource that transforms an `ExpressionSet` into a specific output format. For example, the built-in `alias` resolver takes key-value pairs and formats them as `alias key='value'`. `vib` comes with several built-in resolvers, and you can create your own.
*   **Profile**: A resource that references one or more `ExpressionSet`s to create a complete shell environment. Profiles are the top-level resource that you will typically render to configure your shell.

By combining these three concepts, you can create a modular and reusable shell configuration that can be easily shared and customized.

## Repository Structure

This repository is organized into several packages. Here's a brief overview:

*   [`cmd/vib`](./cmd/vib/README.md): The main entrypoint for the `vib` command-line tool.
*   [`pkg/apis/v1alpha1`](./pkg/apis/v1alpha1/README.md): Contains the API definitions for the `vib` custom resources.
*   `internal/`: Contains the internal implementation of `vib`.
    *   [`internal/adapter/codec`](./internal/adapter/codec/README.md): Provides codecs for encoding and decoding `vib` resources.
    *   [`internal/adapter/formatter`](./internal/adapter/formatter/README.md): Provides formatters for `vib` resources.
    *   [`internal/service`](./internal/service/README.md): Contains the `APIServer` implementation.
    *   [`internal/types`](./internal/types/README.md): Defines the core types and interfaces.
    *   [`internal/util`](./internal/util/README.md): Provides utility functions.

## Getting started

### Installation (script)

```bash
VERSION=v1.0.0
URL="https://raw.githubusercontent.com/alexandremahdhaoui/vib/refs/tags/${VERSION}/cmd/vib-installer/vib-installer.sh"
curl -sfL "${URL}" | sh -s "${VERSION}"
```

### Manual installation

#### Install the binary

```bash
go install github.com/alexandremahdhaoui/vib/cmd/vib@v1.0.0
vib --help
```

If `vib --help` returns `vib: command not found`, ensure the Go bin directory is in your path:

```bash
export GOPATH="${GOPATH:-$(go env GOPATH)}"
export GOBIN="${GOBIN:-${GOPATH}/bin}"
export PATH="${GOBIN}:${PATH}"
```

#### Set up vib

Create a new profile named after your hostname:

```bash
VIB_PROFILE="$(hostname)"
vib create profile "${VIB_PROFILE}"
```

Ensure your profile is sourced when opening a new shell:

```bash
cat <<EOF | tee -a "${HOME}/.${SHELL}rc"
. <(vib render profile "$(hostname)")
EOF
```

## Configure vib

During installation, you created a profile named from your machine's hostname.
In this section, you will learn to create `ExpressionSet` to declare bash functions
and how to set this up in your profile.

NB: Please note you can create your resource in different namespaces.
Multiple namespaces can be useful when you share `vib` configuration with other
people such as your team or organization.

NB: `vib` case be used in conjunction with git-based configuration file managers
such as [chezmoi](https://github.com/twpayne/chezmoi).

### Your first ExpressionSet

#### What are ExpressionSet

An `ExpressionSet` is set of expressions that can be rendered into a desired output
and referenced in a profile.

The `ExpressionSet` uses a resolver to transform a list of arbirtray keys or key-value
pairs into your desired pairs.

#### What are Resolvers?

`vib` is highly customizable, hence everything in `vib` is a resource.
Resolvers are special kind of resources that renders an `ExpressionSet`.

By default, 5 resolvers are created in the `vib-system` namespace:

- alias
- environment-exported
- environment
- function
- plain

Run the following command to see how they're defined.

```bash
vib get -n vib-system resolver
```

You can create and use your own resolvers or even modify existing ones if you want.

#### Finally, create an ExpressionSet

Create a new `ExpressionSet` that will declare and export a few environment variables.

```bash
cat <<'EOF' | vib apply -f -
apiVersion: vib.amahdha.com/v1alpha1
kind: ExpressionSet
metadata:
  name: env
spec:
  keyValues:
    - GOPATH: $(go env GOPATH)
    - GOBIN: ${GOPATH}/bin
    - PATH: ${PATH}:${GOBIN}
  resolverRef:
    name: environment-exported
    namespace: vib-system
EOF

vib get expressionset env
```

Quickly test your new `ExpressionSet` by rendering it.

```bash
vib render expressionset env
```

Feel free to customize it, and don't forget to test it.

```bash
vib edit expressionset env
vib render expressionset env
```

Move to the next section once your happy with the results.

### Create more ExpressionSets

In this section we will create an `ExpressionSet` to define a few aliases by
using the `alias` resolver.

```bash
cat <<'EOF' | vib apply -f -
apiVersion: vib.amahdha.com/v1alpha1
kind: ExpressionSet
metadata:
  name: alias
spec:
  keyValues:
    - grep: grep --color=always
    - less: less -N
    - ls: ls --color=always
    - ll: ls -laF
  resolverRef:
    name: alias
    namespace: vib-system
EOF
```

### Edit your Profile

Reference your newly created `ExpressionSet`s in your Profile:

```bash
vib edit profile "$(hostname)"
```

Or run:

```bash
cat <<EOF | vib apply -f -
apiVersion: vib.amahdha.com/v1alpha1
kind: Profile
metadata:
  name: $(hostname)
spec:
  refs:
    - name: env
    - name: alias
EOF
```

## vib's commands

The `vib` tool provides several commands for managing resources. For more details on the command-line interface, see the [`cmd/vib`](./cmd/vib/README.md) package documentation.

| Command | Description |
|---------|-------------|
| Apply   | Applies resources from stdin or a file. |
| Create  | Creates a new resource. |
| Delete  | Deletes a resource. |
| Edit    | Edit a resource. |
| Get     | Get a set of resource by name or list all resources in a namespace. |
| Render  | Renders the specified resource. |

## See Also

- [Documentation Conventions](./docs/doc-convention.md)
- [`cmd/vib` README](./cmd/vib/README.md)
- [`pkg/apis/v1alpha1` README](./pkg/apis/v1alpha1/README.md)
