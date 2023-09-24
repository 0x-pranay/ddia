Implementing a fast Key-Value store using log-structured hash tables.
- stores log of key-value pairs in a CSV or binary or json format
- An in-memory Hashmap to store key and position of the last key-value pair stored in a log file


Ref: 
- https://riak.com/assets/bitcask-intro.pdf
- https://en.wikipedia.org/wiki/Cyclic_redundancy_check
- https://hemantkgupta.medium.com/insights-from-paper-bitcask-a-log-structured-hash-table-for-fast-key-value-data-6fdddd9e6681
