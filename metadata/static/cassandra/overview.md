## Overview
The Apache Cassandra database is the right choice when you need scalability and high availability without compromising performance. Linear scalability and proven fault-tolerance on commodity hardware or cloud infrastructure make it the perfect platform for mission-critical data. Cassandra's support for replicating across multiple datacenters is best-in-class, providing lower latency for your users and the peace of mind of knowing that you can survive regional outages.

The KUDO Cassandra operator creates, configures and manages Apache Cassandra clusters running on Kubernetes.

Learn how to take full advantage of KUDO Cassandra in [docs](https://github.com/kudobuilder/operators/tree/master/repository/cassandra/3.11) available in [operators](https://github.com/kudobuilder/operators) repository

## Support Level
- Mixed workload tested
  - 4 DCs * 33 nodes (132 nodes total) 684K writes/sec with RF 3 in each DC (total 2.7M writes / sec) 
  - 1 DCs * 30 nodes (30 nodes total) 1.2M writes/sec with RF 1  
- Base tech support

## License
[Apache License 2.0](https://github.com/kudobuilder/operators/blob/master/LICENSE)
