package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	node.Handle("txn", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		rawTxn := body["txn"]
		txnArr, ok := rawTxn.([]interface{})
		if !ok {
			return nil
		}

		result := []interface{}{}
		for _, txn := range txnArr {
			transaction, ok := txn.([]interface{})
			if !ok {
				return nil
			}

			operation, ok := transaction[0].(string)
			if !ok {
				return nil
			}

			ctx := context.Background()
			switch operation {
			case "r":
				key := int(transaction[1].(float64))
				val, err := kv.ReadInt(ctx, strconv.Itoa(key))
				if err != nil {
					result = append(result, []interface{}{"r", key, nil})
					continue
				}

				result = append(result, []interface{}{"r", key, val})

			case "w":
				key := int(transaction[1].(float64))
				keyStr := strconv.Itoa(key)
				val := int(transaction[2].(float64))

				for {
					var oldVal int
					oldVal, err := kv.ReadInt(ctx, keyStr)
					if err != nil {
						kv.CompareAndSwap(ctx, keyStr, 0, 0, true)
						oldVal = 0
					}

					// check if there has been a concurrent write samlam (oldVal no longer same as current val)
					err = kv.CompareAndSwap(ctx, keyStr, oldVal, val, false)
					if err == nil {
						break
					}
				}

				result = append(result, []interface{}{"w", key, val})
			}
		}

		body["type"] = "txn_ok"
		body["txn"] = result

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
