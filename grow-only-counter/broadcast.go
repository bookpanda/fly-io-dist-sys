package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// no need: overengineerd
func broadcast(node *maelstrom.Node, neighbor string) int {
	successCh := make(chan bool, 1)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		var currentVal int

		go func() {
			node.RPC(neighbor, map[string]any{"type": "read_self"}, func(msg maelstrom.Message) error {
				var body map[string]any
				if err := json.Unmarshal(msg.Body, &body); err != nil {
					successCh <- false
					return err
				}

				messageType := body["type"]
				if messageType == "read_self_ok" {
					log.Printf("Broadcast read_self to %s", neighbor)
					successCh <- true
					currentVal = int(body["value"].(float64))
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
			log.Printf("Error broadcasting read self to %s: %v", neighbor, lastErr)
		case success := <-successCh: // Wait for the RPC to complete
			if success {
				return currentVal
			}

			log.Printf("Retrying broadcast read self to %s", neighbor)
			time.Sleep(1 * time.Second)
		}
	}
}
