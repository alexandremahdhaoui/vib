# vib

Vib allows users to intuitively manage and share their shell environments
with their teams, organization and across any platforms.

## Getting started

### Installation

The easiest way to install `vib` is to use the `make install` command. This will download and run the installer script with the latest version.

```bash
make install
```

If your system's hostname contains invalid characters for a resource name (e.g., dots), the installer will prompt you to enter a valid profile name.

Alternatively, you can run the installer script directly:

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
export GOBIN="${GOBIN:-${GOBIN:-${GOPATH}/bin}}"
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

### Testing the Installer

To verify the installer script, you can run the following command:

```bash
make test-install
```

This will run a test script that mocks the `hostname`, `go`, and `vib` commands to ensure the installer behaves as expected.

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
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: ExpressionSet
metadata:
  name: env
spec:
  keyValues:
    - GOPATH: $(go env GOPATH)
    - GOBIN: ${GOBIN}/bin
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
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
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
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
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

| Command | Description |
|---------|-------------|
| Apply   | Applies resources from stdin or a file. |
| Create  | Creates a new resource. |
| Delete  | Deletes a resource. |
| Edit    | Edit a resource. |
| Get     | Get a set of resource by name or list all resources in a namespace. |
| Render  | Renders the specified resource. |
