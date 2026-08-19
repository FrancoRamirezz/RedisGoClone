# RedisGoClone
A decided to dive deep into system design and I thought Redis would be great introduction for this

# Overview

Redis is best described as a In Memory Data Structure Store. Which means it stores everything in memory and executes our commands in a single thread, which entails the performance of Redis. Note for Go since were using go routines we can offset this, by using Sync.Lock in Go Mutex format. Redis Single thread ensures that communication never needs locks and does not cause bottlenecks. 


#Important part: Redis is a key-value store: Here are some data structures 
  Strings
  Hashes (objects/dictionaries)
  Lists
  Sets

<img width="698" height="497" alt="Screenshot 2026-08-19 at 1 32 08 PM" src="https://github.com/user-attachments/assets/5eee0d40-1ccf-4960-8cb3-ffea96b2950f" />

# Features Overview:
Concurrent TCP Server: Listens for client connections on port 6379 the same as Redis

Resp: Speaks a simple wire protocol and their commands support SET key Value, GET key value, DEL key

Expiration TTL: Redis must handle any issues regarding keys. As mentioned before, Redis is key-value store 

Data Persistence: Redis has a durability issues, but does provide solutions: Snapshotting (RDB-like): Periodically writes the entire dataset to disk in a snapshot file, similar to Redis RDB dumps.
– Append-Only Log (AOF): Logs every write command to an append-only file, enabling recovery by replaying commands.


