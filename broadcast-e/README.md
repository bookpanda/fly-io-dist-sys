```bash
go install .

./maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100

./maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100 --nemesis partition

./maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 25 --time-limit 2 --rate 100 --latency 100
```