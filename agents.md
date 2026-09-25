# Instructions for agents.

The project should be structured as a monorepo with a folder for each small program.

We will write this in Go. Use red-green TDD.

Tools should be installed via the Go bin so they will run from anywhere in a shell command.

We need a naming convention, e.g. `jcli-upload`, `jcli-check`. 

## Chaining programs together

Every input-taking tool should accept its primary input either as `argv[1]` or, if no argument is given, as a single line read from stdin. Upload takes a manuscript path; subsequent tools take a submission ID. Argument input always takes priority over stdin. This lets tools be chained the standard Unix way, e.g.:

```
jcli-upload manuscript.md | jcli-validate
```

as well as run standalone with an explicit argument, or via `xargs`/command substitution when preferred:

```
jcli-validate "$submission_id"
jcli-upload manuscript.md | xargs jcli-validate
```

Use the shared `resolvepath.Resolve(args []string, stdin io.Reader) (string, error)` helper with the full argument vector, including the program name. It errors clearly for missing input, empty explicit arguments, and read failures. Test changes with red-green TDD.

Successful action commands print only the submission ID to stdout. Diagnostics
and checklists go to stderr. On failure, exit nonzero without emitting an ID.
Read-only inspection commands may emit data (e.g. JSON) on stdout instead.

## Persistent submission state

Use the shared `submission` library for storage. Resolve its root with
`submission.DefaultRoot()` (`JCLI_HOME`, absolute `XDG_DATA_HOME` plus `/jcli`,
then `~/.local/share/jcli`). Tests must use temporary directories.

Upload stores a snapshot and creates a new ID. Subsequent commands load that ID
and record facts rather than assuming a fixed workflow order. Results must name
the manuscript revision they apply to. Persist failed validation attempts too.
Use locked read-modify-write operations and atomic record replacement for updates;
future state-changing operations must use the same per-submission lock convention.
The current store targets local filesystems on macOS/Linux.

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

Supporting tool: Inspect - displays a persisted submission record as JSON.
