package consumer

import (
	nsqmq "go-nsq/application/mq/nsq"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

type Handler struct {
}

func InitNSQSubscriber(
	client nsqmq.INSQClient,
) {
	go func() {
		err := client.Subscribe("TESTAGAIN")
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func main() {
	nsqSubscriber, err := nsq.NewConsumer("TESTAGAIN", "channel", nsq.NewConfig())
	if err != nil {
		log.Fatalln(err)
	}
	nsqSubscriber.AddHandler(&nsqmq.NSQMessageHandler{})
	err = nsqSubscriber.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the nsqSubscriber.
	nsqSubscriber.Stop()
}
