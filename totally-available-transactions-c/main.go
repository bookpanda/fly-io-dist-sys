package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	log.SetOutput(os.Stderr)

	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	node.Handle("txn", func(msg maelstrom.Message) error {
		log.Printf("HELLO")
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return nil
		}

		rawTxn := body["txn"]
		txnArr, ok := rawTxn.([]interface{})
		if !ok {
			return nil
		}

		now := time.Now().UnixMilli()
		result := []interface{}{}
		writeBuffer := make(map[int]int)

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
			key := int(transaction[1].(float64))
			switch operation {
			case "r":
				rawVersions, err := kv.Read(ctx, strconv.Itoa(key))
				if err != nil {
					result = append(result, []interface{}{"r", key, nil})
					continue
				}

				versions, ok := rawVersions.([]interface{})
				if !ok {
					return nil
				}

				committedVersions := []interface{}{}
				for _, version := range versions {
					v, ok := version.([]interface{})
					if !ok {
						return nil
					}
					timestamp := int64(v[0].(float64))
					if timestamp < now {
						committedVersions = append(committedVersions, version)
					}
				}

				if len(committedVersions) == 0 {
					result = append(result, []interface{}{"r", key, nil})
					continue
				}

				latest := committedVersions[len(committedVersions)-1]
				val := int(latest.([]interface{})[1].(float64))
				log.Printf("key: %d, val: %v", key, val)

				result = append(result, []interface{}{"r", key, val})
			case "w":
				val := int(transaction[2].(float64))
				writeBuffer[key] = val
				result = append(result, []interface{}{"w", key, val})
			}
		}

		ctx := context.Background()
		for key, val := range writeBuffer {
			keyStr := strconv.Itoa(key)
			for {
				rawVersions, err := kv.Read(ctx, keyStr)
				if err != nil {
					kv.CompareAndSwap(ctx, keyStr, [][]int{}, [][]int{}, true)
					rawVersions = []interface{}{}
				}

				versions, ok := rawVersions.([]interface{})
				if !ok {
					return nil
				}

				newVersion := []interface{}{float64(now), float64(val)}
				updatedVersions := append(versions, newVersion)

				log.Printf("key: %d, versions: %v", key, updatedVersions)

				err = kv.CompareAndSwap(ctx, keyStr, versions, updatedVersions, false)
				if err == nil {
					break
				}
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
