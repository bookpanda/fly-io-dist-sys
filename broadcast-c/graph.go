package main

import (
	"encoding/json"
	"log"

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
				n.RPC(neighbor, map[string]any{"type": "send", "message": num}, func(msg maelstrom.Message) error {
					var body map[string]any
					if err := json.Unmarshal(msg.Body, &body); err != nil {
						return err
					}

					messageType := body["type"]
					if messageType == "send_ok" {
						log.Printf("Message sent to %s", neighbor)
					}

					return nil
				})
				visited[neighbor] = true
			}
		}
	}
}
