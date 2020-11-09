KUTTL_VERSION=0.6.1
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

.PHONY: kind-test
kind-test: kubeaddons-tests
	make -f kubeaddons-tests/Makefile kind-test

.PHONY: clean
clean:
ifneq (,$(wildcard kubeaddons-tests/Makefile))
	make -f kubeaddons-tests/Makefile clean
endif
	rm -rf kubeaddons-tests

.PHONY: dispatch-test
dispatch-test: 
	mkdir -p bin/
	curl -Lo bin/kubectl-kuttl_$(KUTTL_VERSION) https://github.com/kudobuilder/kuttl/releases/download/v$(KUTTL_VERSION)/kubectl-kuttl_$(KUTTL_VERSION)_$(OS)_$(MACHINE)
	chmod +x bin/kubectl-kuttl_$(KUTTL_VERSION)
	ln -sf ./kubectl-kuttl_$(KUTTL_VERSION) bin/kubectl-kuttl
	KUBEADDONS_TESTS_KUBECONFIG=/workspace/src-git/kubeconfig
	git clone https://github.com/mesosphere/kubeaddons-tests.git --branch master --single-branch
	kubeaddons-tests/run-tests.sh
