# vib

Vib allows users to intuitively manage and share their shell environments
with their teams, organization and across any platforms.

## Getting started

### Installation (script)

```bash
VERSION=v1.0.0
URL="https://raw.githubusercontent.com/alexandremahdhaoui/vib/refs/tags/${VERSION}/cmd/vib-installer/vib-installer.sh"
# TODO: fix script: hostname can be an invalid resource name:
# -> regex the hostname, if invalid set profile name to default
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
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
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

| Command | Description                                                                                                                                                                                                                        |
|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Get     | Returns the Resource Definition of the specified Kind. Works like `kubectl get`, meaning that we can get the full yaml definition or just a list of defined kinds. We could also get the templated results of the specified kinds. |
| Create  | Create a new resource of the specified Kind                                                                                                                                                                                        |
| Edit    | Edit a specific Resource Definition.                                                                                                                                                                                               |
| Delete  | Deletes a Resource Definition from `vib`.                                                                                                                                                                                          |
| Apply   | Apply resource definitions or a list of files to `vib`.                                                                                                                                                                            |
