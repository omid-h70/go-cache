package main

import (
	"time"
)

func main() {
	// Create a new hashmap with size 10
	myHashMap := NewHashMap(5)

	// Insert key-value pairs

	myHashMap.Set("joe", "dalton", 0*time.Second)
	myHashMap.Set("jack", "dalton", 2*time.Second)
	myHashMap.Set("william", "dalton", 0*time.Second)
	myHashMap.Set("april", "dalton", 0*time.Second)
	myHashMap.Set("lucky", "luke", 0*time.Second)

	myHashMap.Delete("foo")
	myHashMap.Release()
}
