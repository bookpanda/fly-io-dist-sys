```bash
go get github.com/jepsen-io/maelstrom/demo/go
go install .

brew install openjdk graphviz gnuplot

./maelstrom test -w echo --bin ~/go/bin/echo --node-count 1 --time-limit 10
```