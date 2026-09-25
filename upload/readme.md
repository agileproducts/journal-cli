# JCLI Upload

This program snapshots a local Markdown manuscript into the submission store.
See the root README for `JCLI_HOME` and default store locations.

*Given* That I have a manuscript file on my computer
*And* The file extension of that is '.md'
*When* I supply the file as an argument or a single stdin line to `jcli-upload`
*Then* It should store a copy and create a new revision-1 submission
*And* It should print only the new submission ID to stdout for a successor program
*And* Later changes to the original file should not affect the stored manuscript

*Given* That I have a manuscript file on my computer
*And* The file extension is not '.md'
*When* I supply the file as an argument to `jcli-upload`
*Then* It should exit with an error and explain that at present only markdown files are supported

*Given* A missing or non-regular manuscript file, or a store write failure
*When* I run `jcli-upload`
*Then* It should explain the error on stderr and exit 1 without printing an ID

Each upload creates a new submission. An explicit argument takes precedence over
stdin. Example: `jcli-upload "my manuscript.md" | jcli-validate`.
