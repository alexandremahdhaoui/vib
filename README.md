# vib

Vib (pronounced "vibe") allows users to intuitively manage their bash environment across all their platforms. The name
"vib" comes from the contraction of `vi ~/.bash_profile`.

## Getting started

### Install vib

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

### Set up vib

```bash
cat <<EOF | tee -a "${HOME}/.${SHELL}rc"
. <(vib render profile "${__profile}")
EOF
```

### Example configuration

## vib's commands

| Command | Description                                                                                                                                                                                                                        |
|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Get     | Returns the Resource Definition of the specified Kind. Works like `kubectl get`, meaning that we can get the full yaml definition or just a list of defined kinds. We could also get the templated results of the specified kinds. |
| Create  | Create a new instance of the specified Kind                                                                                                                                                                                        |
| Edit    | Edit a specific Resource Definition.                                                                                                                                                                                               |
| Delete  | Deletes a Resource Definition from `vib`.                                                                                                                                                                                          |
| Apply   | Apply resource definitions or a list of files to `vib`.                                                                                                                                                                            |
