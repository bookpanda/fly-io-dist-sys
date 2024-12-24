package main

import (
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	al := make(map[string][]string)
	numbers := []int{}
	numBuffer := []int{}
	lastBroadcast := time.Now()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := body["message"]
		switch messageType := message.(type) {
		case float64: // receive messages to other nodes
			num := int(messageType)
			numBuffer = append(numBuffer, num)
			numbers = append(numbers, num)
		case []interface{}: // only receive messages from other nodes
			for _, m := range messageType {
				num := int(m.(float64))
				numbers = append(numbers, num)
			}
		}

		_, ok := body["receiver"]
		if !ok && (len(numBuffer) > 40 || time.Since(lastBroadcast) > 200*time.Millisecond) {
			broadcastMessage(al, &numBuffer, n)
			numBuffer = []int{}
			lastBroadcast = time.Now()
		}

		body["type"] = "broadcast_ok"
		delete(body, "message")

		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["messages"] = numbers

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "topology_ok"
		topology := (body["topology"])
		mapTopology(topology, al)

		delete(body, "topology")

		return n.Reply(msg, body)
	})

	go func() {
		for range time.Tick(2 * time.Second) {
			broadcastMessage(al, &numBuffer, n)
		}
	}()

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
