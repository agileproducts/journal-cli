# Persistent submissions: implementation plan

## Goal

Keep small, composable commands while making submissions persist between runs.
Upload accepts a local manuscript path; subsequent tools accept a submission ID.
Arguments take precedence over a single line on stdin. Successful action commands
emit only the ID on stdout; diagnostics go to stderr. Inspect emits JSON on stdout.

## Design decisions

- Add a shared `submission` Go module using local `replace` directives.
- Store data in `$JCLI_HOME`, or `$XDG_DATA_HOME/jcli`, or
  `~/.local/share/jcli` by default. Resolve the root once per command.
- Each upload creates a random `sub-<hex>` ID and a new revision-1 submission,
  containing a snapshot `manuscript.md` and a versioned `submission.json` record.
- Record facts, not a mandatory workflow status. Validation attempts record
  their manuscript revision, timestamp, and named check results, including failures.
  Repeated validation appends an attempt, preserving earlier results.
- Publish new submissions by renaming a fully prepared staging directory. Save
  record updates using a synced temporary file and atomic rename. Serialize
  read-modify-write updates with a per-submission Unix advisory file lock,
  released automatically when the process exits. Target macOS and Linux.
- Reject invalid IDs, unreadable/non-regular source files, malformed records, and
  results for a stale manuscript revision. Readers see complete records.
- Revision replacement and remote storage are future increments. This increment
  keeps the stored snapshot immutable through the CLI; re-upload creates a new ID.

## Implementation sequence

1. [x] Shared store: write failing tests for snapshots, reload, root selection,
   invalid inputs, saved failures, stale results, and concurrent updates; implement.
2. [x] Upload: test and implement snapshot creation and ID output, including stdin
   paths, spaces, argument precedence, and errors without success output.
3. [x] Validate: retain structural checks, load the stored snapshot, save all
   results before reporting failure, and emit the ID only on success. Test this
   through a testable command entry point.
4. [x] Inspect: add a read-only command accepting an ID via argument/stdin and
   printing the submission record as formatted JSON. Test existing and missing IDs.
5. [x] Update README, command specs, and agent conventions. Document storage,
   reinstalling, changed input/output contracts, and resumable workflows.
6. [x] Run all module tests, race checks for the store, vet, formatting, builds,
   and CLI smoke tests using an isolated store. Check installation into a temporary
   Go bin; the user can run `make install` for their normal Go bin.

## Acceptance example

```sh
export JCLI_HOME="$HOME/.local/share/jcli-demo"
jcli-upload manuscript.md | jcli-validate | jcli-inspect

# Or stop and resume:
id=$(jcli-upload manuscript.md)
jcli-validate "$id"
jcli-inspect "$id"  # also works after failed validation
```

Tests use temporary directories and do not touch the user's normal submission store.

## Completion notes

- Confirmed failing tests before implementing the store and command behavior.
- All module tests, `make vet`, `make fmt`, `make build`, and `git diff --check`
  pass. `go test -race ./...` passes in `submission/`.
- Installed via `make install` with a temporary `GOBIN` and exercised direct pipes,
  resumption after deleting the source, inspection after failed validation, missing
  inputs, and 20 concurrent validation processes (all attempts preserved).
- Fixed shared input errors (empty explicit arguments and read failures) and
  long-paragraph scanning in validate with red-green regression tests.
- Normal Go-bin installation remains the user's `make install` step.
