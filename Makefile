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
