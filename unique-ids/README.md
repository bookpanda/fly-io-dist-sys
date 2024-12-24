```bash
go install .

./maelstrom test -w unique-ids --bin ~/go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

./maelstrom test -w unique-ids --bin ~/go/bin/unique-ids --time-limit 3 --rate 1000 --node-count 3 --availability total --nemesis partition
```