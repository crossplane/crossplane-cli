INSTALL_DIR ?= /usr/local/bin

default:
	@echo Please specify a target to run!
.PHONY: default

install:
	ln -si $(abspath bin/kubectl-crossplane-stack-*) $(INSTALL_DIR)/
.PHONY: install

uninstall:
	rm $(INSTALL_DIR)/kubectl-crossplane-stack-*
.PHONY: uninstall

integration-test:
	mkdir -p test
	# The local bin is first in the PATH so that it will still be used,
	# even if there is anything installed elsewhere on the path
	PATH=$(abspath bin):$(PATH) $(abspath .)/scripts/integration-test.sh \
			 $(abspath test)
	rm -r test
.PHONY: integration-test
