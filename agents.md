# Instructions for agents.

The project should be structured as a monorepo with a folder for each small program.

We will write this in Go. Use red-green TDD.

Tools should be installed via the Go bin so they will run from anywhere in a shell command.

We need a naming convention, e.g. `jcli-upload`, `jcli-check`. 

## A description of the peer review process
Each step will represnt a program. We will build this as we go. 

0. Hello - let's start with a non-tool that just shows us the pattern.
1. Upload - a tool that lets you select and 'upload' a manuscript.
2. Validate - checks the uploaded file meets basic criteria.
