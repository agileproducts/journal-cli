# JCLI Validate

This program accepts a submission ID and checks its stored manuscript snapshot:

* It has a title (a top-level Markdown heading, e.g. `# Title`)
* It has an abstract section (e.g. `## Abstract`)
* It has a references section (e.g. `## References`)

*Given* That I have a manuscript file with a title, an abstract section, and a references section
*When* I supply its submission ID as an argument or a single stdin line to `jcli-validate`
*Then* It should save a validation attempt with the revision, timestamp, and checks
*And* It should print a passing checklist to stderr
*And* It should print only the submission ID to stdout for a successor program

*Given* That I have a manuscript file missing one or more of a title, an abstract section, or a references section
*When* I supply its submission ID to `jcli-validate`
*Then* It should save the failed validation attempt
*And* It should print ticks and crosses on stderr showing which checks passed and failed
*And* It should exit 1 without emitting an ID

*Given* A missing or invalid submission, unreadable snapshot, or failed result save
*When* I run `jcli-validate`
*Then* It should explain the error on stderr and exit 1 without emitting an ID

Arguments take precedence over stdin. Repeated validation appends attempts. The
original source file is no longer needed. Inspect failures later with
`jcli-inspect "$submission_id"`. Checklist marks are colored on character-device
stderr unless `NO_COLOR` is set to a nonempty value.
