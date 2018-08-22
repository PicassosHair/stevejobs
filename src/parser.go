package main

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func main() {
	fmt.Printf("starting...")

	file, _ := os.Open("/Users/hao/Desktop/data/1964.json")
	iter := jsoniter.Parse(jsoniter.ConfigDefault, file, 1024)
	count := 0

	var value interface{} = iter.Read()
	for value != nil {
		fmt.Printf("%v", count+1)
		value = iter.Read()
	}
	file.Close()
	// iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[ {"a" : [{"b": "c"}], "d": 102 }, "b"]`)
	// iter.ReadArray()
	// iter.Skip()
	// iter.ReadArray()
	// fmt.Printf(iter.ReadString())
}
