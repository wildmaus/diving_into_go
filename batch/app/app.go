package main

import (
	"batch"
	"fmt"
	"time"
)

func main() {
	service := batch.NewService(2, 10*time.Second)
	client := batch.NewClient(service)

	// add one item, must process with them only
	client.AddItem(batch.Item{Id: 1})
	time.Sleep(time.Second)

	// add 3 items, their just added in queue
	client.AddItem(batch.Item{Id: 2})
	client.AddItem(batch.Item{Id: 3})
	client.AddItem(batch.Item{Id: 4})
	time.Sleep(10*time.Second - 10)

	// add items arount time, where previos process ends
	client.AddItem(batch.Item{Id: 5})
	time.Sleep(10 * time.Nanosecond)
	client.AddItem(batch.Item{Id: 6})
	time.Sleep(10 * time.Nanosecond)
	client.AddItem(batch.Item{Id: 7})
	// wait for all items in queue processed
	time.Sleep(30 * time.Second)

	// add two items in empty queue
	// they must be processed right away
	client.AddItem(batch.Item{Id: 8})
	time.Sleep(15 * time.Second)
	client.AddItem(batch.Item{Id: 9})
	time.Sleep(11 * time.Second)
	fmt.Println("All done!")
}
