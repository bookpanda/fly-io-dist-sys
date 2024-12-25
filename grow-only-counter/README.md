```bash
go install .

./maelstrom test -w g-counter --bin ~/go/bin/grow-only-counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition

./maelstrom test -w g-counter --bin ~/go/bin/grow-only-counter --node-count 3 --rate 100 --time-limit 2 --nemesis partition
```