# vib

Vib (pronounced "vibe") allows users to intuitively manage their bash environment across all their platforms. The name
"vib" comes from the contraction of `vi ~/.bash_profile`.

## Install vib

```shell
go install github.com/alexandremahdhaoui/vib/cmd/vib@v0.0.4
vib --help
```

## vib's commands

| Command | Description                                                                                                                                                                                                                        |
|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Get     | Returns the Resource Definition of the specified Kind. Works like `kubectl get`, meaning that we can get the full yaml definition or just a list of defined kinds. We could also get the templated results of the specified kinds. |
| Create  | Create a new instance of the specified Kind                                                                                                                                                                                        |
| Edit    | Edit a specific Resource Definition.                                                                                                                                                                                               |
| Delete  | Deletes a Resource Definition from `vib`.                                                                                                                                                                                          |
| Apply   | Apply resource definitions or a list of files to `vib`.                                                                                                                                                                            |
