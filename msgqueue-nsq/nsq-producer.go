package main

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"log"
	//"time"
)

func main() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("192.168.6.109:4150", config)

	for i := 0; i < 1000000; i++ {
		err := w.Publish("test", []byte(fmt.Sprintf("Hello NSQ! %d", i)))
		if err != nil {
			log.Panic("Could not connect")
			break
		}
		//time.Sleep(time.Second) /*interval of 1s*/
	}

	w.Stop()
}
