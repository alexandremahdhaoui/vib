# Examples

This directory contains examples of how to use `vib`.

## ExpressionSet

`ExpressionSet`s are used to define a set of expressions that can be rendered into a desired output.

- [`vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.alias.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.alias.yaml): Defines a set of common aliases.
- [`vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.env.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.env.yaml): Defines a set of environment variables.
- [`vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.git.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.git.yaml): Defines a set of aliases for Git.
- [`vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.kubectl-aliases.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.kubectl-aliases.yaml): Defines a set of aliases for `kubectl`.
- [`vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.kubectl.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.expressionset.kubectl.yaml): Defines a set of functions for `kubectl`.

## Profile

A `Profile` is used to reference a set of `ExpressionSet`s.

- [`vib.alexandre.mahdhaoui.com_v1alpha1.profile.myprofile.yaml`](./vib.alexandre.mahdhaoui.com_v1alpha1.profile.myprofile.yaml): An example profile that references the `env`, `alias`, and `kubectl` `ExpressionSet`s.
