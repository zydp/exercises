package main

import (
	"log"
	"sync"
	"time"
	"github.com/bitly/go-nsq"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("test", "channel_one", config)	/*指定topic & channel*/
	rfunc := func(message *nsq.Message) error {
		log.Printf("Body: %s\n", message.Body)
		log.Printf("ID: %s\n", message.ID)
		log.Printf("NSQDAddress: %s\n", message.NSQDAddress)
		log.Printf("Attempts: %d\n", message.Attempts)
		//wg.Done()  //don't exit
		return nil
	}
	q.AddHandler(nsq.HandlerFunc(rfunc))
	//err := q.ConnectToNSQD("192.168.6.109:4150")   /*直接连接nsqd*/
	err := q.ConnectToNSQLookupd("192.168.6.109:4161")	/*连接nsqlookupd来获取nsqd地址*/
	if err != nil {
		log.Panic("Could not connect")
	}
	time.Sleep(5 * time.Second)
	wg.Wait()
}
