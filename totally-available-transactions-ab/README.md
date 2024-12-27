```bash
go install .

# a
./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 1 --time-limit 2 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total

./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total

# b
./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted

./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition

./maelstrom serve
```