## Overview
Apache Kafka is used for building real-time data pipelines and streaming apps. It is horizontally scalable, fault-tolerant, wicked fast, and runs in production in thousands of companies

The KUDO Kafka operator creates, configures and manages Apache Kafka clusters running on Kubernetes.

- Secure cluster through TLS encryption, Kerberos authentication and Kafka AuthZ
- Prometheus metrics right out of the box with example of Grafana dashboards
- Kerberos support
- Graceful rolling updates for any cluster configuration changes
- Graceful rolling upgrades when upgrading the operator version
- External access through LB/Nodeports
- Mirror-maker integration

## Support Level
- Mixed workload tested with 5 brokers, 4096Mib and 2000m each
  - 5Million msgs/sec with avg message size of 60 bytes
- Base tech support

## License
[Apache License 2.0](https://github.com/kudobuilder/operators/blob/master/LICENSE)
