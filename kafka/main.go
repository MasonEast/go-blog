package main

import (
	"kafka/consumer"
	"kafka/producer"
	"time"
)
func main() {

	go producer.Put()
	go consumer.Get(1)

	for {
		time.Sleep(time.Hour * 60)
	}
}
