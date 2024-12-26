package main

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	node.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		val := int(body["delta"].(float64))

		for {
			var oldVal int
			oldVal, err := kv.ReadInt(ctx, node.ID())
			if err != nil {
				kv.CompareAndSwap(ctx, node.ID(), 0, 0, true)
				oldVal = 0
			}

			// check if there has been a concurrent write samlam (oldVal no longer same as current val)
			err = kv.CompareAndSwap(ctx, node.ID(), oldVal, oldVal+val, false)
			if err == nil {
				val += oldVal
				break
			}
		}

		body["type"] = "add_ok"
		delete(body, "delta")

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		replicas := node.NodeIDs()

		ctx := context.Background()
		sum := 0
		for _, replica := range replicas {
			var val int
			val, err := kv.ReadInt(ctx, replica)
			if err != nil {
				val = 0
			}
			sum += val
		}

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["value"] = sum

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
