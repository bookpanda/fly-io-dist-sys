package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

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

		now := time.Now().UnixMilli()
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
				rawVersions, err := kv.Read(ctx, strconv.Itoa(key))
				if err != nil {
					result = append(result, []interface{}{"r", key, nil})
					continue
				}

				versions, ok := rawVersions.([]interface{})
				if !ok {
					return nil
				}

				high := len(versions) - 1
				low := 0
				idx := -1

				for low <= high {
					mid := low + (high-low)/2

					version, ok := versions[mid].([]interface{})
					if !ok {
						return nil
					}
					timestamp := int64(version[0].(float64))

					if timestamp >= now {
						high = mid - 1
					} else if timestamp < now {
						low = mid + 1
						idx = mid
					}
				}

				version, ok := versions[idx].([]interface{})
				if !ok {
					return nil
				}
				val := int(version[1].(float64))

				result = append(result, []interface{}{"r", key, val})

			case "w":
				key := int(transaction[1].(float64))
				keyStr := strconv.Itoa(key)
				val := int(transaction[2].(float64))

				for {
					rawVersions, err := kv.Read(ctx, keyStr)
					if err != nil {
						kv.CompareAndSwap(ctx, keyStr, [][]int{}, [][]int{}, true)
						rawVersions = [][]int{}
					}

					versions, ok := rawVersions.([]interface{})
					if !ok {
						return nil
					}

					versions = append(versions, []interface{}{int(now), int(val)})

					// check if there has been a concurrent write samlam (oldVal no longer same as current val)
					err = kv.CompareAndSwap(ctx, keyStr, rawVersions, versions, false)
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
