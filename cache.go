package main

import (
	"log"
	"sync"
	"time"
)

// InMemCache To implement Repository Pattern For Testing purposes
type InMemCache interface {
	Set(key string, value string, ttl time.Duration) bool
	Get(key string) string
	Delete(key string)
	Release()
}

type Node struct {
	Key   string
	Value string
	Next  *Node
	TTL   <-chan time.Time
	Rank  int
}

type HashMap struct {
	buckets  []*Node
	cnt      int
	size     int
	shutdown chan bool
	wg       sync.WaitGroup
	op       sync.Mutex
	ranker   int
}

//var ranker int = 0

func NewHashMap(size int) InMemCache {
	hm := &HashMap{
		buckets:  make([]*Node, size),
		size:     size,
		shutdown: make(chan bool),
		op:       sync.Mutex{},
	}

	hm.wg.Add(1)
	go func() {
		defer func() {
			close(hm.shutdown)
			hm.wg.Done()
			log.Println("Cleaning Up")
		}()

		var i int = 0
		for {

			if len(hm.buckets) > 0 {

				if hm.buckets[i] != nil && hm.buckets[i].TTL != nil {
					select {
					case <-hm.buckets[i].TTL:
						log.Printf("%s is Expired", hm.buckets[i].Key)
						hm.Delete(hm.buckets[i].Key)
					default:
					}
				}
				//To Implement Gracefull Shutdown
				select {
				case <-hm.shutdown:
					return
				default:
				}

				if i == size-1 {
					i = 0
				} else {
					i++
				}
			}
		}
	}()

	return hm
}

func hashFunction(key string, size int) uint {
	return uint(len(key) % size)
}

func (hm *HashMap) Set(key string, value string, ttl time.Duration) bool {

	hm.op.Lock()
	defer hm.op.Unlock()

	var nodePtr *Node
	// Calc int index as hash
	index := hashFunction(key, hm.size)

	// Here we create our node with key and value
	node := &Node{Key: key, Value: value}

	if hm.cnt < hm.size {

		if hm.buckets[index] == nil {
			hm.buckets[index], nodePtr = node, node
		} else {

			//For key Collision Cases
			current := hm.buckets[index]
			for current.Next != nil {
				current = current.Next
			}
			current.Next, nodePtr = node, node
		}

		if nodePtr != nil {
			//Keep Cnt of Elements Because of LRU cache
			hm.cnt++
		}
	} else {

		var found int = -1
		var rank int = -1
		var firstElement bool = false

		for i, bucket := range hm.buckets {
			if bucket != nil {
				if !firstElement {
					rank = bucket.Rank
					firstElement = true
					found = i
				} else if bucket.Rank < rank && firstElement {
					rank = bucket.Rank
					found = i
				}
			}
		}
		//log.Printf("First Rank %d", rank)

		if found >= 0 {
			nodePtr = hm.buckets[found]
			rank = nodePtr.Rank
			//var lruNode *Node
			for lruNode := nodePtr; ; {
				if lruNode.Next == nil {
					break
				}

				lruNode = lruNode.Next
				if lruNode.Rank < rank {
					rank = lruNode.Rank
					nodePtr = lruNode
				}
			}
			//log.Printf("New %s Evicted LRU USED ===== >>>> Key %s Val %s", node.Key, nodePtr.Key, nodePtr.Value)
			if nodePtr != nil {

				if index == uint(found) {
					node.Next = nodePtr.Next
					*nodePtr = *node
				} else {
					//node.Rank = ranker
					hm.buckets[index] = node
					//Evict Old Node
					if nodePtr.Next != nil {
						*nodePtr = *nodePtr.Next
					}
					//Update the Last Reference
					nodePtr = node
				}
			}
		}
	}

	if nodePtr != nil {
		if ttl > 0 {
			nodePtr.TTL = time.After(ttl)
		}
		hm.ranker++
		nodePtr.Rank = hm.ranker
	}

	//log.Printf("Inserted %s => %s", key, value)
	// we always return true because we support capacity and LRU together
	return true
}

func (hm *HashMap) Get(key string) string {
	var current *Node
	defer func() {
		if current != nil {
			hm.ranker++
			current.Rank = hm.ranker
		}
	}()
	// Calculate the index with hashFunction to know where we should look.
	index := hashFunction(key, hm.size)

	// We get the first node in this index and assign it to a variable.
	current = hm.buckets[index]

	for current != nil {
		/* Check Also Next siblings for Collision Cases */
		//log.Printf("%s == %s", current.Key, key)
		if current.Key == key {
			return current.Value
		}
		current = current.Next
	}
	/* Key Not Found*/
	return ""
}

func (hm *HashMap) Delete(key string) {
	hm.op.Lock()
	defer func() {
		hm.op.Unlock()
		log.Printf("Try to Clear key %s", key)
	}()
	// Calculate the index of the key using the hash function
	index := hashFunction(key, hm.size)
	// Get the first node at the calculated index
	current := hm.buckets[index]
	// Initialize a pointer to keep track of the previous node
	var prev *Node

	for current != nil {
		// If the current node's key matches the key to be deleted
		if current.Key == key {
			// If the previous pointer is nil, it means the node to delete is the first node
			if prev == nil {

				hm.buckets[index] = current.Next
			} else {
				/* Skip the current node  */
				prev.Next = current.Next
			}
			// Exit the method after deletion
			if hm.cnt > 0 {
				hm.cnt--
			}
			return
		}
		/* get next node in the linked list */
		prev = current
		current = current.Next
	}
}

func (hm *HashMap) Release() {
	hm.shutdown <- true
	hm.wg.Wait()
}
