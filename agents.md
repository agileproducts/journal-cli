# Instructions for agents.

The project should be structured as a monorepo with a folder for each small program.

We will write this in Go. Use red-green TDD.

Tools should be installed via the Go bin so they will run from anywhere in a shell command.

We need a naming convention, e.g. `jcli-upload`, `jcli-check`. 

## Chaining programs together

Every tool should accept its primary input (typically a manuscript path) either as `argv[1]` or, if no argument is given, as a single line read from stdin. Argument input always takes priority over stdin. This lets tools be chained the standard Unix way, e.g.:

```
jcli-upload manuscript.md | jcli-validate
```

as well as run standalone with an explicit argument, or via `xargs`/command substitution when preferred:

```
jcli-validate manuscript.md
jcli-upload manuscript.md | xargs jcli-validate
```

Implement this with a small `ResolvePath(args []string, stdin io.Reader) (string, error)` helper in each program, tested with red-green TDD, that errors clearly if neither an argument nor stdin input is provided.

## Sharing code between programs

If logic needs to be shared between two or more programs (e.g. `resolvepath`), factor it into its own top-level directory with its own `go.mod` and no `main.go` — it's a library module, not an installable tool. Wire it into a program's module with a local `replace` directive so nothing needs to be published:

```
cd some-program
go mod edit -require=jcli-<lib>@v0.0.0 -replace=jcli-<lib>=../<lib>
go mod tidy
```

The `Makefile` distinguishes these automatically: any top-level directory with a `go.mod` is included in `test`/`vet`/`fmt`, but only those that also have a `main.go` are treated as installable `PROGRAMS` for `build`/`install`.

## A description of the peer review process
Each step will represnt a program. We will build this as we go. 

0. Hello - let's start with a non-tool that just shows us the pattern.
1. Upload - a tool that lets you select and 'upload' a manuscript.
2. Validate - checks the uploaded file meets basic criteria.
