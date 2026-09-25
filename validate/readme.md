# JCLI Validate

This program checks that an uploaded manuscript file meets basic structural criteria before it can proceed through peer review. When supplied with a manuscript file in a valid format, it checks that:

* It has a title (a top-level Markdown heading, e.g. `# Title`)
* It has an abstract section (e.g. `## Abstract`)
* It has a references section (e.g. `## References`)

*Given* That I have a manuscript file with a title, an abstract section, and a references section
*When* I supply the file as an argument to `jcli-validate`
*Then* It should print a checklist showing that all checks passed
*And* It should echo the file/path so that it could be read by a successor program

*Given* That I have a manuscript file missing one or more of a title, an abstract section, or a references section
*When* I supply the file as an argument to `jcli-validate`
*Then* It should print a checklist of ticks and crosses showing which checks passed and which failed
*And* It should exit with an error explaining that the manuscript failed validation
