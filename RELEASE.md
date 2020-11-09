# Release

The kubeaddons-enterprise repository provides the default addons in the Kommander Catalog.
Changes to `master` branch without soaking should be avoided, as breaking `master` branch would affect all the users using Kommander Catalog.


## Schedule

Schedule will only be triggered if there are any changes to be released. 
In case there are any changes, `master` should be updated twice monthly, on the second and forth Wednesdays, or as needed to address CVEs.


## Development 

All development should be done in feature branches created of from the latest `dev` branch. 
Feature branches should never interact directly with `master` or `staging` branches.

Pull Requests of feature branches should be opened against the `dev` branch. 

## Deploying `staging` addons with Kommander

This is best done on the soak cluster running on AWS, its Kommander already has a `kudo-staging` project.
First, we need to update the `AddonRepository` of the `kudo-staging` project to point to the stating branch. For that, generate a cluster token in soak Konvoy UI to be able to use `kubectl`. There should be a `kudo-stating-%something%` namespace containing the `kubeaddons-enterprise` addonrepository. Edit this and change its `spec.ref` to the staging branch, e.g.:

```
$ kubectl edit -n kudo-staging-vmg4d-j44g9 addonrepositories.kubeaddons.mesosphere.io kubeaddons-enterprise

apiVersion: v1
items:
- apiVersion: kubeaddons.mesosphere.io/v1beta2
  kind: AddonRepository
  metadata:
    creationTimestamp: "2020-09-23T12:57:50Z"
    generation: 1
    name: kubeaddons-enterprise
    namespace: kudo-staging-vmg4d-j44g9
    resourceVersion: "19255350"
    selfLink: /apis/kubeaddons.mesosphere.io/v1beta2/namespaces/kudo-staging-vmg4d-j44g9/addonrepositories/kubeaddons-enterprise
    uid: b225c2df-7faa-4d29-81cc-ebb110600c1b
  spec:
    options:
      credentialsRef: {}
    priority: "1"
    ref: my-staging-branch
    url: https://github.com/mesosphere/kubeaddons-enterprise
  status:
    ready: true
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
```

Once that's done, create a cluster on AWS using the konvoy CLI. Get that cluster's kubeconfig, attach the cluster to soak's Kommander, then add the cluster to the `kudo-staging` project.

You should now be able to view the staging addons in the project's catalog and deploy them on the attached cluster. Verify that the updated addons deploy successfully by adding them in Kommander and checking with `kubectl` that the respective instance resources are created on the attached cluster.

## Release Process

### Staging Update (Second and Forth Thursday)

On the second and forth _**Thursday**_, the `staging` should be updated. By merging `dev` into `staging` branch.

Verify in AWS Soak cluster, the workspace of `kudo-testing` 

- Verify the catalog in `kudo-staging` is accessible
- Verify all addons in `kudo-staging` project are running
- Uninstall any addons installed in the project `kudo-staging` 
- Install `Zookeeper`, `Kafka` and `Cassandra` in `kudo-staging`
- Check the plan status of all three instances through kudo-cli


### Master Update (Second and Forth Wednesday)

On the second and forth _**Wednesday**_:

- Open the PR from `staging` to `master`
- Paste the PR link in #sig-ksphere-catalog and provide soak cluster links to help the reviewers
