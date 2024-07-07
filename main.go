package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Create a new hashmap with size 10
	myHashMap := NewHashMap(10)

	// Insert key-value pairs
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("john%d", i)
		val := fmt.Sprintf("doe%d", i)
		//t := time.Unix(0, time.Now().UnixNano())
		//elapsed := time.Since(t)
		//log.Printf("wtf %s", elapsed.String())

		myHashMap.Set(key, val, 0*time.Second)
	}

	myHashMap.Set("jack", "bar", 2*time.Second)

	// Get and print values
	key := "john14"
	value := myHashMap.Get(key)
	log.Printf("Value for key %s: %s", key, value)

	/* If we try to get the value for key "foo" we will get an empty string. (You can return a
	   proper error or a flag in your get method) */
	time.Sleep(10 * time.Second)

	// Delete a key
	log.Printf("About To Clear Foo")
	myHashMap.Delete("foo")
	myHashMap.Release()
}
