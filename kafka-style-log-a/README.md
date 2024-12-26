```bash
go install .

./maelstrom test -w kafka --bin ~/go/bin/kafka-style-log --node-count 1 --concurrency 2n --time-limit 2 --rate 1000

./maelstrom test -w kafka --bin ~/go/bin/kafka-style-log --node-count 1 --concurrency 2n --time-limit 20 --rate 1000

./maelstrom serve
```