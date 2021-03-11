
KUTTL_VERSION=0.8.0
KIND_VERSION=0.8.1
KUBERNETES_VERSION ?= 1.17.5
KUBECONFIG?="kubeconfig"

OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
MACHINE=$(shell uname -m)
KIND_MACHINE=$(shell uname -m)
ifeq "$(KIND_MACHINE)" "x86_64"
  KIND_MACHINE=amd64
endif

export PATH := $(shell pwd)/bin/:$(PATH)

ARTIFACTS=dist

kubeaddons-tests:
	git clone --depth 1 https://github.com/mesosphere/kubeaddons-tests.git --branch master --single-branch
	ls -a
	@if [ -f kubeconfig ]; then\
		cp kubeconfig kubeaddons-tests/kubeconfig && ls -a ./kubeaddons-tests; \
	fi

.PHONY: kind-test
kind-test: kubeaddons-tests
ifeq (,$(wildcard kubeconfig))
	$(MAKE) -C kubeaddons-tests kind-test
else
	$(MAKE) -C kubeaddons-tests bin/kubectl-kuttl
	$(MAKE) -o kubeconfig -C kubeaddons-tests kind-test
endif


.PHONY: clean
clean:
ifneq (,$(wildcard kubeaddons-tests/Makefile))
	make -C kubeaddons-tests clean
endif
	rm -rf kubeaddons-tests
