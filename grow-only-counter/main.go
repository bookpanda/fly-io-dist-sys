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
			oldVal, err := kv.ReadInt(ctx, "key")
			if err != nil {
				kv.CompareAndSwap(ctx, "key", 0, 0, true)
				oldVal = 0
			}

			// check if there has been a concurrent write samlam (oldVal no longer same as current val)
			err = kv.CompareAndSwap(ctx, "key", oldVal, oldVal+val, false)
			if err == nil {
				break
			}
		}

		// _, ok := body["receiver"].(bool)
		// if !ok {
		// 	for _, neighbor := range node.NodeIDs() {
		// 		if neighbor != node.ID() {
		// 			go broadcast(node, neighbor, val)
		// 		}
		// 	}
		// }

		body["type"] = "add_ok"
		delete(body, "delta")

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		replicas := node.NodeIDs()
		values := make([]int, len(replicas))

		for i, replica := range replicas {
			values[i] = broadcast(node, replica)
		}

		maxValue := values[0]
		for _, v := range values {
			if v > maxValue {
				maxValue = v
			}
		}

		// ctx := context.Background()
		// var val int
		// val, err := kv.ReadInt(ctx, "key")
		// if err != nil {
		// 	val = 0
		// }

		// for {
		// 	var currentVal int
		// 	currentVal, err := kv.ReadInt(ctx, "key")
		// 	if err != nil {
		// 		kv.CompareAndSwap(ctx, "key", 0, 0, true)
		// 		currentVal = 0
		// 	}

		// 	// check if there has been a concurrent write samlam (currentVal no longer same as current val)
		// 	err = kv.CompareAndSwap(ctx, "key", currentVal, currentVal+0, false)
		// 	if err == nil {
		// 		val = currentVal
		// 		break
		// 	}
		// }

		body["type"] = "read_ok"
		body["value"] = maxValue

		return node.Reply(msg, body)
	})

	node.Handle("read_self", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		var val int
		val, err := kv.ReadInt(ctx, "key")
		if err != nil {
			val = 0
		}

		body["type"] = "read_self_ok"
		body["value"] = val

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
