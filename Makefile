INSTALL_DIR ?= /usr/local/bin

default:
	@echo Please specify a target to run!
.PHONY: default

build:
	GO111MODULE=on go build -o bin/kubectl-crossplane-trace cmd/trace/main.go
	GO111MODULE=on go build -o bin/kubectl-crossplane-pack cmd/pack/main.go
.PHONY: build

test:
	GO111MODULE=on go test ./...
.PHONY: test

clean:
	rm -f bin/kubectl-crossplane-trace
	rm -f bin/kubectl-crossplane-pack
.PHONY: clean

install:
	cp $(abspath bin/kubectl-crossplane*) $(INSTALL_DIR)/
.PHONY: install

integration-test:
	mkdir -p test
	# The local bin is first in the PATH so that it will still be used,
	# even if there is anything installed elsewhere on the path
	PATH=$(abspath bin):$(PATH) $(abspath .)/scripts/integration-test.sh \
			 $(abspath test)
	rm -r test
.PHONY: integration-test
