# Notes on "Designing Data-Intensive Application" by Martin Klepmann


These are my notes on this data intense book.


## Chapter 5: Replication

Keeping a copy of the same data on multiple machines that are connected via network.

Why:
- To reduce access latency (keep data geographically close to users)
- Increase availability     (allow the system to continue working even if some part fails)
- Increase read throughput  ( via read replicas)

Assumption:
- Dataset is so small that each machine can hold a copy of entire dataset.

Popular algorithms to replicate changes between nodes:
- Single leader
- Multi leader
- leaderless replication

### Leaders and Followers:
- Each node that stores a copy of the database is called a _replica_
- _leader-based-replication_ aka _active/passive_ or _master-slave_ replication
- Writes happen on a repilca designate as the _leader_ (aka _master_ or _primary_) 
- The other replics are _followers_ (aka _read replicas_, _slaves, secondaries or hot standbys_)
- leader sends the data change to all of its followers as part of _replication log_ or _change stream_
- Each follower takes the log from leader and updates its local copy of the database accordingly, in the same order as they were processed by leader.
- Clients can read either from leader or from any of the read replicas. 
- This mode of replication is built-in in postgres ref: https://www.postgresql.org/docs/14/high-availability.html
    Mysql, Oracle, NoSQL (MongoDb, rethinkDB and Espresso). Distributed message brokes such as Kafka and RabbitMq also use it.

### Synchronous VS Asynchronous Replication:

Does the replication happen synchronously or asynchronously.
When a client writes something to leader, in Synchronous the client gets a confirmation from leader only when the data is replicated it all its followers.
whereas in Asynchronous replication, leader doesn't wait till the data is replicated to send a confirmation to client.

| replication  | Advantage   | Disadvantage   |
|-------------- | -------------- | -------------- |
| Synchronous    | followers have up-to-date copy     | if followers fail then writes cannot be processed on leader. Cause of network fault etc. |
| Asynchronous    | leader can process writes even if all of its followers are failed      | If a leader fails, then any data that is not yet replicated are lost.     |

*Semi-Synchronous* : If synchronous follower isn't available then one of the other asynchronous follower is made syncrhronous.

ref: https://www.postgresql.org/docs/14/warm-standby.html#SYNCHRONOUS-REPLICATION


