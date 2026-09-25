# Every top-level Go module in the monorepo (programs and shared libraries).
MODULES := $(shell find . -mindepth 1 -maxdepth 1 -type d -exec test -e '{}/go.mod' \; -print | sed 's|^\./||' | sort)

# Just the installable/buildable programs: modules that also have a main.go.
PROGRAMS := $(shell for d in $(MODULES); do test -e "$$d/main.go" && echo $$d; done)

.PHONY: install test build vet fmt clean list

## Run every module's test suite (programs and shared libraries).
test:
	@for dir in $(MODULES); do \
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

## Vet every module.
vet:
	@for dir in $(MODULES); do \
		echo "==> go vet ./$$dir/..."; \
		(cd $$dir && go vet ./...) || exit 1; \
	done

## Report any files that gofmt would reformat.
fmt:
	@for dir in $(MODULES); do \
		(cd $$dir && gofmt -l .); \
	done

## Remove local build artefacts (does not touch installed $GOPATH/bin binaries).
clean:
	rm -rf bin

## List the discovered modules and, of those, which are installable programs.
list:
	@echo "modules:  $(MODULES)"
	@echo "programs: $(PROGRAMS)"
