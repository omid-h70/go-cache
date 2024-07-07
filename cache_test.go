package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestValuesWithTTL(t *testing.T) {

	myHashMap := NewHashMap(10)
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("john%d", i)
		val := fmt.Sprintf("doe%d", i)
		myHashMap.Set(key, val, 2*time.Second)
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("john%d", i)
		val1 := myHashMap.Get(key)
		require.Equal(t, val1, "", "Keys with TTL  > 0 must be expired in desired Time")
	}

	myHashMap.Release()
}

func TestConcurrentSet(t *testing.T) {
	var wg sync.WaitGroup
	myHashMap := NewHashMap(5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			key := fmt.Sprintf("john%d", i)
			val := fmt.Sprintf("doe%d", i)

			myHashMap.Set(key, val, 0*time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()

	//Get Concurrent Values
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("john%d", i)
		val1 := fmt.Sprintf("doe%d", i)
		val2 := myHashMap.Get(key)
		require.Equal(t, val1, val2)
	}
}

func TestDeleteFunc(t *testing.T) {
	myHashMap := NewHashMap(1)
	myHashMap.Set("key1", "val1", 0*time.Second)
	myHashMap.Delete("key1")
	val := myHashMap.Get("key1")
	require.Equal(t, val, "", "Deleted Key must be gone !")
}

func TestSetAndGetValue(t *testing.T) {
	myHashMap := NewHashMap(10)

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("john%d", i)
		val := fmt.Sprintf("doe%d", i)
		myHashMap.Set(key, val, 0*time.Second)
		val1 := myHashMap.Get(key)
		require.Equal(t, val, val1, "Set Value And Get Value Must be the same")
	}
	myHashMap.Release()
}

func TestLRUCache(t *testing.T) {

	t.Log("Values Over below 5 must be over written because of LRU cache")
	myHashMap := NewHashMap(5)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("john%d", i)
		val := fmt.Sprintf("doe%d", i)
		myHashMap.Set(key, val, 0*time.Second)
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("john%d", i)
		val := myHashMap.Get(key)
		if i < 5 {
			require.Equal(t, val, "", "key must be evicted because of lru cache")
		} else {
			require.NotEqual(t, val, "", "new values must be replaced with old ones")
		}
	}

	myHashMap.Release()
}

func TestLRUCacheIfOlderValueIsUsed(t *testing.T) {

	t.Log("newer Value must be overwritten, because we're using old value !")
	myHashMap := NewHashMap(2)

	myHashMap.Set("key1", "val1", 0*time.Second)
	myHashMap.Set("key2", "val2", 0*time.Second)
	myHashMap.Get("key1")
	myHashMap.Set("key3", "val3", 0*time.Second)

	val := myHashMap.Get("key2")

	require.Equal(t, val, "", "Newer Value must be overwritten ! ")
	myHashMap.Release()
}
