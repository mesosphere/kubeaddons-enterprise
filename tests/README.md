# Addon Tests

kubeaddons-enterprise use KUTTL for testing.

### Tests structure

Addons are added to the `addons` directory and their respective tests are added to the `tests` directory.

Example:
```
├── addons
│   ├── cassandra
│   │   └── 0.x
│   │       ├── cassandra-2.yaml
│   │       └── cassandra.yaml
└── tests
    ├── cassandra
        └── cassandra-install
            ├── 00-assert.yaml
            ├── 00-install.yaml
            ├── 01-assert.yaml
            └── 01-update.yaml
```
In `tests/cassandra/cassandra-install/00-install.yaml` it should refer to the addon file we are going to test. In this case the latest version is `0.x/cassandra-2.yaml` 
