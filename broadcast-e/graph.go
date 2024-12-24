package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func mapTopology(topo any, al map[string][]string) {
	topology, ok := topo.(map[string]interface{})
	if !ok {
		log.Fatalf("Expected map[string]interface{}, got %T", topo)
		return
	}

	for k, v := range topology {
		arr, ok := v.([]interface{})
		if !ok {
			log.Fatalf("Expected []interface{}, got %T for key %s", v, k)
			return
		}

		var stringArr []string
		for _, s := range arr {
			str, ok := s.(string)
			if !ok {
				log.Fatalf("Expected string, got %T in array for key %s", s, k)
				return
			}
			stringArr = append(stringArr, str)
		}

		al[k] = stringArr
	}
}

// broadcast arrays of nums every 2 sec or array reach 3 numbers first?
func broadcastMessage(al map[string][]string, numBuffer *[]int, n *maelstrom.Node) {
	if len(*numBuffer) == 0 {
		return
	}

	var visited = make(map[string]bool)
	visited[n.ID()] = true

	for _, node := range al {
		for _, neighbor := range node {
			if !visited[neighbor] {
				go broadcastToNode(n, neighbor, *numBuffer)
				visited[neighbor] = true
			}
		}
	}
}

func broadcastToNode(n *maelstrom.Node, neighbor string, numBuffer []int) {
	successCh := make(chan bool, 1)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go func() {
			n.RPC(neighbor, map[string]any{"type": "broadcast", "message": numBuffer, "receiver": true}, func(msg maelstrom.Message) error {
				var body map[string]any
				if err := json.Unmarshal(msg.Body, &body); err != nil {
					successCh <- false
					return err
				}

				messageType := body["type"]
				if messageType == "broadcast_ok" {
					log.Printf("Broadcast message %v to %s", numBuffer, neighbor)
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
			log.Printf("Error broadcasting message %v to %s: %v", numBuffer, neighbor, lastErr)
		case success := <-successCh: // Wait for the RPC to complete
			if success {
				return
			}

			log.Printf("Retrying broadcast message %v to %s", numBuffer, neighbor)
			time.Sleep(1 * time.Second)
		}
	}
}
