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
