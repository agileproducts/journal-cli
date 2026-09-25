# Journal CLI

The aim of this project is to make a model of the submission and peer review system for an academic journal modelled as a CLI. It takes inpsiration from the UNIX philosophy. It consusts of a set of small, simple programs that do one thing well and can be chained together to achieve great things.

## Building, testing, and installing

Requires Go 1.26.5 or later and Make. The local store supports macOS and Linux.
Each program and shared library lives in its own directory with its own Go module.
The root `Makefile` discovers all top-level directories containing a `go.mod` for
testing, vetting, and formatting. Modules also containing `main.go` are programs
and are built and installed:

```sh
make test     # run every module's test suite (including shared libraries)
make install  # install every jcli-* binary to GOBIN, or $GOPATH/bin by default
make build    # build every binary into ./bin (leaves $GOPATH/bin untouched)
make vet      # go vet every module
make fmt      # report files that gofmt would reformat
make clean    # remove ./bin
make list     # show discovered modules and installable programs
```

Run `make install` again after changing code, including shared libraries. Ensure
your Go bin directory is on `PATH`. `make build` lets you use `./bin/jcli-upload`,
`./bin/jcli-validate`, and `./bin/jcli-inspect` directly during development.

## Persistent, composable workflows

`jcli-upload` accepts a local `.md` file path, stores a snapshot, and prints a new
submission ID (`sub-` followed by 32 hex characters). **Validate now accepts a
submission ID, rather than a local file path.** It saves its findings and prints
the same ID on success. `jcli-inspect` displays the record as formatted JSON.

All three accept their input as the first argument, or as one line from stdin
when no argument is given. Arguments take priority. Action commands reserve
stdout for IDs and stderr for feedback; inspect uses stdout for JSON.

```sh
set -o pipefail  # bash/zsh: report failure from any stage
jcli-upload "my manuscript.md" | jcli-validate | jcli-inspect

# Save the ID to resume later, including after failed validation.
id=$(jcli-upload "my manuscript.md")
jcli-validate "$id"
jcli-inspect "$id"
```
Validation checks for an H1 title and Abstract and References H2 headings (case
insensitive section names; `##Abstract` is also accepted). These are structural
presence checks, not checks of section content. A failed validation saves its
results, prints the failed checks to stderr, exits 1, and emits no ID. Inspect the
saved ID separately after a failure. Repeated validation appends an attempt;
inspect succeeds even if those attempts failed. Set `NO_COLOR=1` for plain output.

### Where state lives

Store location, in priority order:

1. `JCLI_HOME`, if set (relative paths resolve against the working directory).
2. `$XDG_DATA_HOME/jcli`, if `XDG_DATA_HOME` is an absolute path.
3. `~/.local/share/jcli`.

For an isolated experiment, export the setting so every pipeline stage uses it:

```sh
export JCLI_HOME="$HOME/.local/share/jcli-demo"
jcli-upload manuscript.md | jcli-validate | jcli-inspect
```

```text
<store>/submissions/sub-<32 hex characters>/
├── manuscript.md
├── submission.json
└── .lock              # created on the first validation
```

The JSON record contains a schema version, creation timestamp, original filename,
manuscript revision, and validation attempts with timestamps and individual
checks. The tools record facts rather than enforcing one fixed workflow order.

Uploads create new revision-1 submissions. Editing or deleting the original file
does not affect the snapshot. Re-uploading creates a new ID; replacing revisions
is a future increment. Let the tools manage the stored files: record updates are
atomic and serialized with per-submission locks, but direct edits bypass that
coordination. Use this store on a local filesystem.

See [PLAN.md](PLAN.md) for the persistence implementation plan and decisions.
