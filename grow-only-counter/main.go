package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

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

		body["type"] = "add_ok"
		delete(body, "delta")

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		replicas := node.NodeIDs()
		ch := make(chan int, len(replicas))
		var wg sync.WaitGroup

		for _, replica := range replicas {
			wg.Add(1)
			go func(replica string) {
				defer wg.Done()
				value := broadcast(node, replica)
				ch <- value
			}(replica)
		}

		values := make([]int, len(replicas))
		for i := 0; i < len(replicas); i++ {
			val := <-ch
			values[i] = val
		}

		// Wait for remaining Goroutines to finish (if needed)
		go func() {
			wg.Wait()
			close(ch) // Ensure channel is closed once all Goroutines complete
		}()

		maxValue := values[0]
		for _, v := range values {
			if v > maxValue {
				maxValue = v
			}
		}

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

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
