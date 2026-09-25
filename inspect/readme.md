# JCLI Inspect

Displays a submission record as formatted JSON, including its stored manuscript
metadata and all validation attempts. This is a read-only reporting command.

*Given* An uploaded submission, with or without validation attempts
*When* I supply its ID as an argument or a single stdin line to `jcli-inspect`
*Then* It should print the saved record as JSON on stdout and exit 0
*And* It should leave the record unchanged, even if validation previously failed

*Given* A missing, invalid, or corrupt submission record
*When* I run `jcli-inspect`
*Then* It should explain the error on stderr and exit 1

Arguments take precedence over stdin. Example:

```sh
jcli-upload manuscript.md | jcli-validate | jcli-inspect
```

Inspect emits JSON rather than an ID, so it normally ends an action pipeline.
See the root README for storage configuration and saving an ID to inspect failures.
