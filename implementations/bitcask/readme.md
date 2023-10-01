### Implementing a fast Key-Value store using log-structured hash tables.
- stores log of key-value pairs in a CSV or binary or json format
- An in-memory Hashmap to store key and position of the last key-value pair stored in a log file

Key benefits:
- Fast write and reads
- simple design and easy to implement

Drawback:
- All keys must fit in the RAM, since keys are stored in in-memory hastable

Implementation:
- A bitcask instance is a directory and only one process will open that bitcask for writing at a given time. (database server)
- At any moment, only one file is "active" in that directory for writing. 
- once this reach a size threshold, it will be closed and a new active file will be created.
- this closed file is immutable, only used for reading but never for writes.
- Active file is written by appending

    | crc* | timestamp | key_size   | value_size | key | value |
    |---| --------- | ---------- | ---------- | --- | ----- |
    | crc | timestamp | ksz        | value_size | key         | value         |
    | crc | timestamp | ksz        | value_size | key    |  value         |
    
    *crc is optional at this moment
- After append completes, an in-memory hastable called a "keydir" is updated with file, offset and size of most recently written entry for that key
    ```
    key -> { file_id, value_sz, value_pos, tstamp }
    key -> { file_id, value_sz, value_pos, tstamp }
    key -> { file_id, value_sz, value_pos, tstamp }
    ```
- To read the value, simply lookup the keydir and read the data using file_id, position and size that are returned from that lookup. 





API:

|                 |    |
|-------------- | -------------- |
| bitcask.Open(DirectoryName, Opts) : BitCaskHandle \| error   | Open a new or existing Bitcask datastore     |
| bitcask.get(BitcaskHandle, key) : not_found \| error \| any  | Retrieve a value by a key from datastore     |
| bitcask.put(BitcaskHandle, key, value) : OK \| error \| any  | store a value and a key in datastore     |
| bitcask.delete(BitcaskHandle, key) : OK \| error \| any  | Delete a key from the bitcask datastore |



Further reading
- Locking shema for keydir

---

### Ref: 
- https://riak.com/assets/bitcask-intro.pdf
- https://en.wikipedia.org/wiki/Cyclic_redundancy_check
- https://hemantkgupta.medium.com/insights-from-paper-bitcask-a-log-structured-hash-table-for-fast-key-value-data-6fdddd9e6681
- https://healeycodes.com/implementing-bitcask-a-log-structured-hash-table


