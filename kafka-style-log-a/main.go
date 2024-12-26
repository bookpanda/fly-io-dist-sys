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

	node.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx := context.Background()
		key := body["key"].(string)
		val := int(body["msg"].(float64))
		var offset int

		for {
			var oldVal any
			oldVal, err := kv.Read(ctx, key)
			if err != nil {
				kv.CompareAndSwap(ctx, key, 0, 0, true)
				oldVal = []int{}
			}

			currentLogs, ok := oldVal.([]int)
			if !ok {
				currentLogs = []int{}
			}

			// check if there has been a concurrent write samlam (oldVal no longer same as current val)
			newLogs := append(currentLogs, val)
			err = kv.CompareAndSwap(ctx, key, oldVal, newLogs, false)
			if err == nil {
				offset = len(newLogs) - 1
				break
			}
		}

		body["type"] = "send_ok"
		body["offset"] = offset
		delete(body, "key")
		delete(body, "msg")

		return node.Reply(msg, body)
	})

	node.Handle("poll", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		offsets, ok := body["offsets"].(map[string]float64)
		if !ok {
			return nil
		}

		msgs := make(map[string][]int)
		for key, offset := range offsets {
			ctx := context.Background()
			var rawLogs any
			rawLogs, err := kv.Read(ctx, key)
			if err != nil {
				rawLogs = []int{}
			}

			logs, ok := rawLogs.([]int)
			if !ok {
				logs = []int{}
			}

			for i := int(offset); i < len(logs); i++ {
				msgs[key] = append(msgs[key], logs[i])
			}
		}

		body["type"] = "poll_ok"
		body["msg"] = msgs
		delete(body, "offset")

		return node.Reply(msg, body)
	})

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		offsets, ok := body["offsets"].(map[string]float64)
		if !ok {
			return nil
		}

		for key, offset := range offsets {
			ctx := context.Background()
			for {
				var oldVal any
				oldVal, err := kv.Read(ctx, key+"_commit")
				if err != nil {
					kv.CompareAndSwap(ctx, key+"_commit", 0, 0, true)
					oldVal = 0
				}

				currentOffset, ok := oldVal.(int)
				if !ok {
					currentOffset = 0
				}

				if int(offset) <= currentOffset {
					continue
				}

				err = kv.CompareAndSwap(ctx, key+"_commit", oldVal, int(offset), false)
				if err == nil {
					break
				}
			}
		}

		body["type"] = "commit_offsets_ok"
		delete(body, "offsets")

		return node.Reply(msg, body)
	})

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		keys, ok := body["keys"].([]string)
		if !ok {
			return nil
		}

		offsets := make(map[string]int)
		for _, key := range keys {
			ctx := context.Background()
			offset, err := kv.ReadInt(ctx, key+"_commit")
			if err != nil {
				offset = 0
			}

			offsets[key] = offset
		}

		body["type"] = "list_committed_offsets_ok"
		body["offsets"] = offsets
		delete(body, "keys")

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
