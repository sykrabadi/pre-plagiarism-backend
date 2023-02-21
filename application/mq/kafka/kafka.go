package kafka

import (
	"context"
	"encoding/json"
	"go-nsq/application/mq"
	"go-nsq/store"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type IKafkaClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

type KafkaClient struct {
	client sarama.SyncProducer
	consumer sarama.Consumer
	msgCounter prometheus.Counter
	mqLatency prometheus.Histogram
	dbstore store.Store
}

func NewKafkaClient(store store.Store) (IKafkaClient, error){
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5

	consumerConf := sarama.NewConfig()
	consumerConf.Consumer.Return.Errors = true
	url := os.Getenv("KAFKA_BROKER_ADDR")
	conn, err := sarama.NewSyncProducer([]string{url}, conf)
	if err != nil {
		log.Printf("Fail to initialize Kafka producer with error: %v", err)
		return nil, err
	}
	consumer, err := sarama.NewConsumer([]string{url}, consumerConf)
	if err != nil {
		log.Printf("Fail to initialize Kafka consumer with error: %v", err)
		return nil, err
	}
	reg := prometheus.NewRegistry()
	msgCounter := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name:      "Kafka_message_pumped_count",
		Help:      "Number of message pumped by Kafka",
	})
	histogramReg := prometheus.NewRegistry()
	msgHistogram := promauto.With(histogramReg).NewHistogram(
		prometheus.HistogramOpts{
			Name: "Kafka_latency_seconds",
			Help: "Latency of Kafka in seconds",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
	)
	// Register msgCounter metric
	prometheus.Register(msgCounter)
	prometheus.Register(msgHistogram)
	return &KafkaClient{
		client: conn,
		consumer: consumer,
		msgCounter: msgCounter,
		mqLatency: msgHistogram,
		dbstore: store,
	}, nil
}

func (k *KafkaClient) Publish(topic string, message []byte) error{
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64){
		k.mqLatency.Observe(v)
	}))
	defer timer.ObserveDuration()
	val := sarama.ByteEncoder(message)
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: val,
	}
	_, _, err := k.client.SendMessage(&msg)
	if err != nil {
		log.Printf("Error publish message from Kafka with error: %v", err)
		return err
	}
	k.msgCounter.Inc()
	return nil
}

func (k *KafkaClient) Subscribe(topic string) error{
	subscriber, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Fail to consume partition from Kafka with error:%v", err)
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var response mq.MQSubscribeMessage
	
	doneCh := make(chan struct{})
	go func() error{
		for{
			select{
			case err := <-subscriber.Errors():
				log.Printf("Fail to Subscribe from Kafka with error:%v", err)
				return err
			case msg := <-subscriber.Messages():
				log.Printf("Consuming message from topic:%v | message: %v", string(msg.Topic), string(msg.Value))
				err = json.Unmarshal(msg.Value, &response)
				if err != nil {
					log.Printf("Error when unmarshalling json at [KafkaClient.Subscribe] with error : %v", err)
					return err
				}
				err = k.dbstore.DocumentStore().UpdateData(context.TODO(), response)
				if err != nil {
					log.Printf("[NSQMessageHandler.HandleMessage] error when update data with error %v \n", err)
					return err
				}
			case <-sigchan:
				err = k.consumer.Close()
				if err != nil {
					log.Fatalf("Error shutting down Kafka consumer gracefully with error:%v", err)
				}
				log.Printf("Shutting down Kafka conumser \n")
				doneCh <- struct{}{}
			}
		}
	}()

	return nil
}
