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
- External access to Spark UI through LB/Nodeports
- Security features: RPC Auth with Encryption, TLS support, Kerberos
- Additional features for Spark and Spark History Server integration with popular data stores, such as Amazon S3 and HDFS

Learn how to take full advantage of KUDO Spark in [docs](https://github.com/kudobuilder/operators/tree/master/repository/spark) available in [operators](https://github.com/kudobuilder/operators) repository

## Dependencies

Spark addon requires previous installation of:
- Prometheus operator

## Support Level
- Mixed workload tested with:
  - 50 operators
  - 1000 Spark jobs with 1000 drivers and 1000 executor pods running concurrently
  - TeraSort benchmark successfully completed and validated  

## License
[Apache License 2.0](https://github.com/kudobuilder/operators/blob/master/LICENSE)