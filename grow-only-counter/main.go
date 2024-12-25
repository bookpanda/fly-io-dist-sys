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
	sum := 0

	node.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		val := int(body["delta"].(float64))

		var oldVal int
		oldVal, err := kv.ReadInt(ctx, node.ID())
		if err != nil {
			kv.CompareAndSwap(ctx, node.ID(), 0, 0, true)
			oldVal = 0
		}

		kv.Write(ctx, node.ID(), oldVal+val)
		sum += val

		_, ok := body["receiver"].(bool)
		if !ok {
			for _, neighbor := range node.NodeIDs() {
				if neighbor != node.ID() {
					broadcast(node, neighbor, val)
				}
			}
		}

		body["type"] = "add_ok"
		delete(body, "delta")

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		var val int
		val, err := kv.ReadInt(ctx, node.ID())
		if err != nil {
			val = 0
		}

		body["type"] = "read_ok"
		body["value"] = val
		// body["value"] = sum

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
