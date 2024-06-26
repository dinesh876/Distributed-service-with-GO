WAL   Package Path       : internal/log
GRPC  Package Path       : api/v1
GRPC Server Package Path : internal/Server

WAL:
----

Record  : the data stored in our log
Store   : the file we store records in
Index   : the file we store index entries in
Segment : the abstraction that ties a store and an index together
Log     : the abstraction that ties all the Segment together


Secure Your Service
-------------------
1. Encrypt data in flight to protect against man in the middle attacks
2. Authenticate to identify clients
3. Authorize to determine the persmissions of the identified clients

TLS one way authentication        - only Authenticate the Server
TLS Mutual Authenticate (two way) - two way authentication , in which both client and server validate the other's communication
                                    is more commonly used in machine to machine communication


Metrics: (tool : prometheus)
--------
1. Counter     - Track the number of times an event happened.such as the number of requests that failed
2. Histograms  - show you a distribution of your data.ex: measuring the percentiles of your request duration and sizes
3. Gauges      - Track the current value of something. you can replace that value entirely.useful for saturation type metrics. like a host's disk

Google golden Signals 
---------------------
1. Latency - the time it takes your service to process requests.
2. Traffic - the amount of demand on your service.For a typical web service, this could be requests processed per second.
3. Error   - your service's request failure rate.internal server errors are particularly important
4. Saturation - a measure of your service's capacity. For example if service persists data to disk,at your current ingress rate will you run out 
                of hard drive space spoon? if you have an in-memory store,how much memory is your service using compared to the memory available?

structured logs: (tool: Elasticsearch)
----------------
- Logs describe events in your system.Logs should help to troubleshoot.
- A structured log is a set of name and value ordered pairs encoded in consistent schema and format that's easily read by programs

Traces: (tool: jaegar)
-------
Traces capture request lifecycles and let you track requests as they flow through your system.

service Discovery:(serf,zookeeper & consel)
------------------
- Service discovery is a process of figuring out how to connect to a services.

- Ex: Webservice discovering and connecting to its database


Replication:
------------
- Store multiple copies of the log data when we have multiple servers in a cluster.
- Replication makes our service more resillient to failures.
- Ex : if a node's disks fails and we can't recover its data, Replication can save our
  butts because it ensures that there's  a copy saved on another disk


Coordinate your services with consensus
---------------------------------------
- Consensus algorithm are tools used to get distributed services to agreee on shared state even in the face of failures.
- Need to put the servers in leader and follower relationships where the follower replicate the leader's data

Algorithm: Raft:
---------------
- Raft is a distributed consensus algorithm designed to be easily understood and implemented.
- It's the consensus algorithm behind services like Etcd, the key-value store that backs kubernetes, consul

