```bash
go install .


./maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

./maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 5 --time-limit 2 --rate 10 --nemesis partition
```