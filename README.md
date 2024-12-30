# fly.io Distributed Systems Challenge
[Challenge link](https://fly.io/dist-sys/1/)

## 1. Echo
intro to Maelstrom

## 2. Unique ID Generation
I used twitter's 64-bit Snowflake ID to generate unique IDs.
- bit 0-19: sequence number for same millisecond (have to skip first 10 bits, for some reason malestrom leaves last 3 decimal places as 0)
- bit 20-21: node bit
- bit 22-62 (41 bits): timestamp
- bit 63: sign bit

## 3. Broadcast
- go routines for broadcast to other nodes
- batching of broadcast messages to reduce message per operation (by number of message and time since last message)
- go routines for flushing broadcasts incase of no new broadcast messages and it has reached buffer threshold
> If anyone can solve 3e problem with median latency < 1 sec and maximum latency < 2 secs, please feel free to share your solution. I'd love to learn the solution.


## 4. Grow-Only Counter
- write: one key per node
- read: read from all nodes and return sum value

## 5. Kafka-style Log
- hold keys for commits and values
- write: append to log
- read: read from log
- lots of parsing

## 6. Totally-available Transactions
- write to buffer before committing to KV store
- reads see only committed values aka past transactions
- after all reads in that transaction are done, commit the write buffer to KV store