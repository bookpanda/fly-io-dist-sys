```bash
go install .

./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 2 --concurrency 2n --time-limit 2 --rate 1000 --consistency-models read-committed --availability total –-nemesis partition

./maelstrom test -w txn-rw-register --bin ~/go/bin/totally-available-transactions --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-committed --availability total –-nemesis partition

./maelstrom serve
```