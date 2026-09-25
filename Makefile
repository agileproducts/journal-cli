# Discover every program in the monorepo (any top-level dir with its own go.mod).
PROGRAMS := $(shell find . -mindepth 1 -maxdepth 1 -type d -exec test -e '{}/go.mod' \; -print | sed 's|^\./||' | sort)

.PHONY: install test build vet fmt clean list

## Run each program's test suite.
test:
	@for dir in $(PROGRAMS); do \
		echo "==> go test ./$$dir/..."; \
		(cd $$dir && go test ./...) || exit 1; \
	done

## go install every program's binary into $GOPATH/bin (or GOBIN), so each
## jcli-* tool is available anywhere in the shell.
install:
	@for dir in $(PROGRAMS); do \
		echo "==> go install ./$$dir"; \
		(cd $$dir && go install .) || exit 1; \
	done

## Build every program into ./bin, without touching $GOPATH/bin.
build:
	@mkdir -p bin
	@for dir in $(PROGRAMS); do \
		echo "==> go build ./$$dir -> bin/jcli-$$dir"; \
		(cd $$dir && go build -o ../bin/jcli-$$dir .) || exit 1; \
	done

## Vet every program.
vet:
	@for dir in $(PROGRAMS); do \
		echo "==> go vet ./$$dir/..."; \
		(cd $$dir && go vet ./...) || exit 1; \
	done

## Report any files that gofmt would reformat.
fmt:
	@for dir in $(PROGRAMS); do \
		(cd $$dir && gofmt -l .); \
	done

## Remove local build artefacts (does not touch installed $GOPATH/bin binaries).
clean:
	rm -rf bin

## List the programs this Makefile has discovered.
list:
	@echo $(PROGRAMS)
