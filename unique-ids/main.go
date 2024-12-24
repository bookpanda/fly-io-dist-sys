package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	log.SetOutput(os.Stderr)

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
		id |= seq << 10 // the node wont accept the first 10 digits e.g. 7277266422153085000 last 3 is always 0
		log64Bits(id)

		body["type"] = "generate_ok"
		body["id"] = id

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func log64Bits(num int64) {
	log.Printf("%d\n", num)
	log.Printf("%064b\n", num)
	sign := num >> 63
	log.Printf("sign: %d\n", sign)
	timestamp := num >> 22
	log.Printf("timestamp: %064b\n", timestamp)
	node := num >> 20 & 3
	log.Printf("node: %064b\n", node)
	seq := num & 0xFFFFF
	log.Printf("seq: %064b\n", seq)
}
