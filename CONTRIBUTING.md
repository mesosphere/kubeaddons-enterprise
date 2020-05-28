# Contributing

The [Kubeaddons Contribution Guidelines](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md) apply here in addition to anything else provided in this section.

## Adding a new version

When adding a new version for an addon, make sure that tests are installing the correct version. 

For example for jenkins with have this test step:
```
apiVersion: kudo.dev/v1beta1
kind: TestStep
commands:
  - command: kubectl apply -f ../../../addons/jenkins/1.x/jenkins.yaml
    namespaced: true
```

if `jenkins` is bumped to `jenkins-2.yaml` or `2.x/jenkins.yaml`, we will need to bump also in the tests.

Once we can make sure that we always test the latest version of addons, this requirement will be removed.

You can read more around testing structure in [testing documentation](./tests/README.md).
