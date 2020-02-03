## Overview
Apache Spark is a fast and general-purpose cluster computing system for big data. 
It provides high-level APIs in Scala, Java, Python, and R, and an optimized engine that supports general computation graphs for data analysis. 
It also supports a rich set of higher-level tools including: Spark SQL for SQL and DataFrames, MLlib for machine learning, 
GraphX for graph processing, and Spark Streaming for stream processing.

The KUDO Spark Operator creates, configures, and manages instances of Spark Operator running on Kubernetes.

- Easily create and manage Spark workloads on Kubernetes using `SparkApplication` resource
- Spark History Server support
- Metrics collection into Prometheus via ServiceMonitors  
- Monitor Spark jobs using pre-configured Grafana dashboards
- High Availability deployment mode
- Integration with external batch scheduler (Volcano)
- Graceful rolling updates for any cluster configuration changes
- Graceful rolling upgrades when upgrading the operator version
- External access through LB/Nodeports

## Support Level
- Mixed workload tested with:
  - 50 operators
  - 1000 Spark jobs with 1000 drivers and 1000 executor pods running concurrently
  - TeraSort benchmark successfully completed and validated  

## License
[Apache License 2.0](https://github.com/kudobuilder/operators/blob/master/LICENSE)