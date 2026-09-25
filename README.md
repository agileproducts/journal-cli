# Journal CLI

The aim of this project is to make a model of the submission and peer review system for an academic journal modelled as a CLI. It takes inpsiration from the UNIX philosophy. It consusts of a set of small, simple programs that do one thing well and can be chained together to achieve great things.

## Building, testing, and installing

Each program lives in its own directory with its own Go module. The root `Makefile` auto-discovers every program (any top-level directory containing a `go.mod`) and runs commands across all of them:

```sh
make test     # run every program's test suite
make install  # go install every jcli-* binary to $GOPATH/bin, so they're usable anywhere
make build    # build every binary into ./bin (leaves $GOPATH/bin untouched)
make vet      # go vet every program
make fmt      # report files that gofmt would reformat
make clean    # remove ./bin
make list     # show which program directories were discovered
```

Once installed, the tools can be chained together the Unix way, either with a direct pipe or via `xargs`/command substitution:

```sh
jcli-upload manuscript.md | jcli-validate
jcli-upload manuscript.md | xargs jcli-validate
jcli-validate "$(jcli-upload manuscript.md)"
```


