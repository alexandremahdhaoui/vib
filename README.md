# vib

Vib (pronounced "vibe") allows users to intuitively manage their bash environment across all their platforms. The name
"vib" comes from the contraction of `vi ~/.bash_profile`.

## Repository structure

| Filepath        | Object     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                        |
|-----------------|------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `expression.go` | Expression | An expression is an object that will be templated into a bash expression, such as a bash function, an alias, an environment variable declaration or any other kind of valid bash statement.<br/>The way an expression is templated into any of these types is defined by a Resolver.<br/>3 built-in expression kinds  with pre-defined resolver exist: Function, Alias & Environment.<br/>A 4th kind is Custom which allows users to extend `vib`. |
| `set.go`        | Set        | A Set is a reusable collection of expressions.                                                                                                                                                                                                                                                                                                                                                                                                     |
| `profile.go`    | Profile    | A profile is a collection of set that defines a user's profile.                                                                                                                                                                                                                                                                                                                                                                                    |
| `resolver.go`   | Resolver   | A resolver is an object used to template an expression a valid . 4 resolver types exists: Exec, Fmt, Plain & Gotemplate.                                                                                                                                                                                                                                                                                                                           |
| `render.go`     | N/A        | This file defines the internal for templating the configuration into bash.                                                                                                                                                                                                                                                                                                                                                                         |

## vib's commands

| Command | Description                                                                                                                                                                                                                        |
|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Get     | Returns the Resource Definition of the specified Kind. Works like `kubectl get`, meaning that we can get the full yaml definition or just a list of defined kinds. We could also get the templated results of the specified kinds. |
| Create  | Create a new instance of the specified Kind                                                                                                                                                                                        |
| Edit    | Edit a specific Resource Definition.                                                                                                                                                                                               |
| Delete  | Deletes a Resource Definition from `vib`.                                                                                                                                                                                          |
| Apply   | Apply resource definitions or a list of files to `vib`.                                                                                                                                                                            |


## Try it

```shell
for x in examples/*; do go run ./... --debug apply -f $x; done
go run ./... --debug render profile profile-0
```