package main

import (
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
		// if !visited[stNode] {
		for _, neighbor := range node {
			if !visited[neighbor] {
				n.Send(neighbor, map[string]any{"type": "send", "message": num})
			}
		}
		// queue = append(queue, node)
		// visited[node] = true
	}

	// for len(queue) > 0 {
	// 	current := queue[0]
	// 	queue = queue[1:]

	// 	n.Send(current, map[string]any{"type": "send", "message": num})

	// 	for _, node := range al[current] {
	// 		if !visited[node] {
	// 			// visited[node] = true
	// 			queue = append(queue, node)
	// 		}
	// 	}
	// }
	// for _, v := range al[start] {
	// 	if !visited[v] {
	// 		bfs_(al, v, num, visited)
	// 	}
	// }

	// for k, v := range al {
	// 	if visited[k] {
	// 		continue
	// 	}

	// }
}
