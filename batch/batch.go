package batch

import (
	"context"
	"log"
	"sync"
	"time"
)

/*
	we have service, that serve some elems in batches (max n elem in p interval)
	this is a client for that service, for optimal batching input elems
*/

type client struct {
	s           Service
	n           uint64
	p           time.Duration
	mu          sync.Mutex // mutex for items and lastprocess
	items       []Item
	lastProcess time.Time
	processMu   sync.Mutex // mutex for only one process
}

func NewClient(s Service) *client {
	n, p := s.GetLimits()
	return &client{
		s:           s,
		n:           n,
		p:           p,
		items:       []Item{},
		lastProcess: time.Now().Add(-2 * p),
	}
}

/**
 * @notice add some Item in process queue
 */
func (c *client) AddItem(item Item) {
	c.mu.Lock()
	c.items = append(c.items, item)
	c.mu.Unlock()
	log.Printf("add item: %v", item.Id)
	go c.Process()
}

/**
 * @notice sends Items from queue to the service for processing
 * if another process already waiting for service - just return
 * else sends <= n items for processing and then delete them
 */
func (c *client) Process() {
	if c.processMu.TryLock() {
		c.mu.Lock()
		if (uint64(len(c.items))) > 0 && time.Since(c.lastProcess) > c.p+5 {
			c.lastProcess = time.Now()
			var num uint64
			if uint64(len(c.items)) > c.n {
				num = c.n
			} else {
				num = uint64(len(c.items))
			}
			var batch Batch = c.items[:num]
			c.mu.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), c.p+5)
			defer cancel()
			log.Printf("Run process with %v\n", num)
			if err := c.s.Process(ctx, batch); err != nil {
				c.processMu.Unlock()
			}
			c.Delete(num)
		} else {
			c.mu.Unlock()
			c.processMu.Unlock()
			log.Println("Wait for service or item")
		}
	} else {
		log.Println("Wait for another process")
	}
}

/**
 * @notice delete processed items from queue and run another process
 */
func (c *client) Delete(num uint64) {
	log.Printf("delete %v\n", num)
	c.mu.Lock()
	c.items = c.items[num:]
	c.mu.Unlock()
	c.processMu.Unlock()
	go c.Process()
}
