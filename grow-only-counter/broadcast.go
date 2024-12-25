package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func broadcast(node *maelstrom.Node, neighbor string, delta int) {
	successCh := make(chan bool, 1)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		go func() {
			node.RPC(neighbor, map[string]any{"type": "add", "delta": delta, "receiver": true}, func(msg maelstrom.Message) error {
				var body map[string]any
				if err := json.Unmarshal(msg.Body, &body); err != nil {
					successCh <- false
					return err
				}

				messageType := body["type"]
				if messageType == "add_ok" {
					log.Printf("Broadcast add %d to %s", delta, neighbor)
					successCh <- true
				} else {
					successCh <- false
				}

				return nil
			})
		}()

		select {
		case <-ctx.Done(): // timeout -> cancel, does not care successCh. if time.After() is used instead,
			// have to be careful with stale go routines
			lastErr := ctx.Err()
			log.Printf("Error broadcasting add %d to %s: %v", delta, neighbor, lastErr)
		case success := <-successCh: // Wait for the RPC to complete
			if success {
				return
			}

			log.Printf("Retrying broadcast add %d to %s", delta, neighbor)
			time.Sleep(1 * time.Second)
		}
	}
}
