# Contributing

The [Kubeaddons Contribution Guidelines](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md) apply here in addition to anything else provided in this section.

## Adding a new version

When adding a new version for an addon, make sure that tests are installing the correct version. 

For example for jenkins with have this test step:
```
apiVersion: kudo.dev/v1beta1
kind: TestStep
commands:
  - command: kubectl apply -f ../../1.x/jenkins.yaml
    namespaced: true
```

if `jenkins` is bumped to `jenkins-2.yaml` or `2.x/jenkins.yaml`, we will need to bump also in the tests.

Once we can make sure that we always test the latest version of addons, this requirement will be removed.

## Adding tests

The tests for this repository live in [kubeaddons-enterprise-tests](https://github.com/mesosphere/kubeaddons-enterprise-tests) repo.
Any new addon that doesn't provide tests will make the CI fail. Please read [adding-tests](https://github.com/mesosphere/kubeaddons-enterprise-tests#adding-tests) documentation in [kubeaddons-enterprise-tests](https://github.com/mesosphere/kubeaddons-enterprise-tests) repo.
