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
		val := int(body["value"].(float64))

		oldVal, err := kv.Read(ctx, "key")
		if err != nil {
			return err
		}
		oldNum := oldVal.(int)

		kv.Write(ctx, "key", oldNum+val)

		body["type"] = "add_ok"

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		val, err := kv.Read(ctx, "key")
		if err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["value"] = val

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
