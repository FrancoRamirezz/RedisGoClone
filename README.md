# RedisGoClone
A decided to dive deep into system design and I thought Redis would be great introduction for this

# Overview

Redis is best described as a In Memory Data Structure Store. Which means it stores everything in memory(RAM) and executes our commands in a single thread, which uses an event loop and relies on I/O Multiplexing Mechanism to deal with file descriptors(think of network sockets). This is done on OS system calls. Here are some examples:
  epoll on Linux
  kqueue on macOS and FreeBSD

  <img width="607" height="243" alt="Screenshot 2026-08-25 at 3 30 58 PM" src="https://github.com/user-attachments/assets/45a92fdf-bd3b-4b72-8ec4-83794d41f3c7" />


#Important part: Redis is a key-value store and every object in Redis is stored in a string key: Here are some data structures 
  Strings
  Hashes (objects/dictionaries)
  Lists
  Sets

<img width="698" height="497" alt="Screenshot 2026-08-19 at 1 32 08 PM" src="https://github.com/user-attachments/assets/5eee0d40-1ccf-4960-8cb3-ffea96b2950f" />

# Features Overview:
Concurrent TCP Server: Listens for client connections on port 6379 the same as Redis.

Concurrency Design: Each connection is processed by a handler func and we will use a goroutine for this. Which means many clients can issue commands without waiting for others. Note: Redis is single threaded and uses an architecture called an event loop, that can handle all incoming requests and outgoing responses. The event loop will make sure to check for new client connections, incoming data, or any completed tasks. Single threads can help deal concurrency issues like Locks and Corruption of shared memory. An event loop runs a single thread that runs a infinite loop. 
Note for Go since were using go routines we can offset this, by using Sync.Lock in Go Mutex format with Go’s concurrency we can increase throughput on multi-core systems(CPU). 
 

Resp: Speaks a simple wire protocol and their commands support SET key Value, GET key value, DEL key

Expiration TTL: Redis must handle any issues regarding keys. As mentioned before, Redis is key-value store 

Data Persistence: Redis has a durability issues, but does provide solutions: Snapshotting (RDB-like) Periodically writes the entire dataset to disk(most other db use this) in a snapshot file, similar to Redis RDB dumps.
– Append-Only Log (AOF): Logs every write command to an append-only file, enabling recovery by replaying commands.


