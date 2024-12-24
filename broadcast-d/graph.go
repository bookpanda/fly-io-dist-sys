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

func broadcastMessage(al map[string][]string, num int, start string, n *maelstrom.Node) {
	var visited = make(map[string]bool)
	visited[start] = true

	for _, node := range al {
		for _, neighbor := range node {
			if !visited[neighbor] {
				go broadcastToNode(n, neighbor, num)
				visited[neighbor] = true
			}
		}
	}
}

func broadcastToNode(n *maelstrom.Node, neighbor string, num int) {
	successCh := make(chan bool, 1)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			lastErr := ctx.Err()
			log.Printf("Error broadcasting message %d to %s: %v", num, neighbor, lastErr)
			time.Sleep(1 * time.Second)
		default:
			n.RPC(neighbor, map[string]any{"type": "broadcast", "message": num, "receiver": true}, func(msg maelstrom.Message) error {
				var body map[string]any
				if err := json.Unmarshal(msg.Body, &body); err != nil {
					successCh <- false
					return err
				}

				messageType := body["type"]
				if messageType == "broadcast_ok" {
					log.Printf("Broadcast message %d to %s", num, neighbor)
					successCh <- true
				} else {
					successCh <- false
				}

				return nil
			})

			// Wait for the RPC to complete
			success := <-successCh
			if success {
				return
			}

			log.Printf("Retrying broadcast message %d to %s", num, neighbor)
			time.Sleep(1 * time.Second)
		}

	}
}
