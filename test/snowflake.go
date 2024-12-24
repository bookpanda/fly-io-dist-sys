package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

func timestamp() {
	// timestamp := uint64(time.Now().UnixNano())
	// timestamp := uint64(time.Now().UnixMicro())
	timestamp := uint64(time.Now().UnixMilli())

	// The first 41 bits are the timestamp
	fmt.Printf("timestamp: %d\n", timestamp)
	fmt.Printf("timestamp: %b\n", timestamp)
	fmt.Printf("timestamp: %064b\n", timestamp)
	// no. of preceeding 0s
	timeString := fmt.Sprintf("%064b", timestamp)
	for i := 0; i < 64; i++ {
		if timeString[i] == '1' {
			fmt.Printf("no. of preceeding 0s: %d\n", i)
			fmt.Printf("timestamp size: %d\n", 64-i)
			break
		}
	}
}

func print64Bits(num int64) {
	fmt.Printf("%064b\n", num)
	sign := num >> 63
	fmt.Printf("sign: %d\n", sign)
	timestamp := num >> 22
	fmt.Printf("timestamp: %064b\n", timestamp)
	node := num >> 20 & 3
	fmt.Printf("node: %064b\n", node)
	seq := num & 0xFFFFF
	fmt.Printf("seq: %064b\n", seq)
}

func binToDec(bin string) int64 {
	dec, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		log.Fatalf("Error converting binary to decimal: %v", err)
		panic(err)
	}
	return dec
}

func snowflake() {
	var id int64
	timestamp := int64(time.Now().UnixMilli())
	id |= timestamp << 22
	fmt.Printf("id: %064b\n", id)

	re := regexp.MustCompile("[0-9]+")
	match := re.FindString("n1")
	node, err := strconv.ParseInt(match, 10, 64)
	if err != nil {
		log.Fatalf("Error converting string to int64: %v", err)
		panic(err)
	}
	nodeBit0 := node & 1
	nodeBit1 := (node >> 1) & 1
	id |= nodeBit0 << 20
	id |= nodeBit1 << 21
	fmt.Printf("id: %064b\n", id)

	var seq int64 = 11
	id |= seq
	fmt.Printf("id: %064b\n", id)

	// print64Bits(id)

	// fmt.Printf("id: %d\n", id)

	// print64Bits(7277244150523625000)
	println(binToDec("0110010011111110000010010001001100111111100000000000000000000000"))
	println(binToDec("0110010011111110000010010001001100111111100000000000000000000001"))
	//7277264026151682049
	//7277264026151682000
	//7277266422153084928
	//7277266422153085000
}
