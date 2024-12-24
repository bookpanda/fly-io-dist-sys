package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	prevTimestamp := int64(time.Now().UnixMilli())
	var seq int64 = 0

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		var id int64
		timestamp := int64(time.Now().UnixMilli())
		id |= timestamp << 22

		re := regexp.MustCompile("[0-9]+")
		match := re.FindString(msg.Dest)
		node, err := strconv.ParseInt(match, 10, 64)
		if err != nil {
			log.Fatalf("Error converting string to int64: %v", err)
			return err
		}

		nodeBit0 := node & 1
		nodeBit1 := (node >> 1) & 1
		id |= nodeBit0 << 20
		id |= nodeBit1 << 21

		if timestamp == prevTimestamp {
			seq++
		} else {
			seq = 0
			prevTimestamp = timestamp
		}
		id |= seq

		body["type"] = "generate_ok"
		body["id"] = id

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
