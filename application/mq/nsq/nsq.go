package nsq

import (
	"context"
	"encoding/json"
	"go-nsq/application/mq"
	"go-nsq/store"
	"log"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type INSQClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

type Message struct {
	Timestamp    string
	FileObjectID string
	FileName     string
}

type NSQMessageHandler struct {
	dbstore store.Store
}

// TODO : Apply message format from MQ to update the specified document at mongodb
func processMessage(body []byte) error {
	log.Printf("Receiving message from NSQ with payload : %v ", string(body))
	return nil
}

func (h *NSQMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	var response mq.MQSubscribeMessage
	// do whatever actual message processing is desired
	err := json.Unmarshal(m.Body, &response)
	if err != nil {
		log.Printf("Error when unmarshalling json at NSQMessagehandler with error : %v", err)
		return err
	}
	log.Printf("logging from test mq %v \n", response.BoundingBoxes)

	err = h.dbstore.DocumentStore().UpdateData(context.TODO(), response)
	if err != nil {
		log.Printf("[NSQMessageHandler.HandleMessage] error when update data with error %v \n", err)
		return err
	}

	log.Println("Logging message from NSQMessageHandler")
	log.Println(response.FileObjectID)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

type NSQClient struct {
	config        nsq.Config
	msgCounter    prometheus.Counter
	msgCounterVec prometheus.CounterVec
	mqLatency     prometheus.Histogram
	dbstore       store.Store
}

func NewNSQClient(store store.Store) INSQClient {
	config := nsq.NewConfig()
	// after adding config.DialTimeout, NSQ will not throw i/o timeout anymore
	config.DialTimeout = 3 * time.Second
	reg := prometheus.NewRegistry()
	msgCounter := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "NSQ_message_pumped_count",
		Help: "Number of message pumped by NSQ",
	})
	regMsgCounterVec := prometheus.NewRegistry()
	msgCounterVec := promauto.With(regMsgCounterVec).NewCounterVec(prometheus.CounterOpts{
		Name: "NSQ_msg_pumped_vec_counter",
		Help: "Number of message pumped by NSQ in vector",
	}, []string{"code", "method"})
	histogramReg := prometheus.NewRegistry()
	msgHistogram := promauto.With(histogramReg).NewHistogram(
		prometheus.HistogramOpts{
			Name:    "NSQ_latency_seconds",
			Help:    "Latency of NSQ in seconds",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
	)
	// Register msgCounter metric
	err := prometheus.Register(msgCounter)
	if err != nil {
		log.Printf("Fail to register NSQ message counter with error: %v", err)
		return nil
	}
	err = prometheus.Register(msgCounterVec)
	if err != nil {
		log.Printf("Fail to register NSQ message countervec with error: %v", err)
		return nil
	}
	err = prometheus.Register(msgHistogram)
	if err != nil {
		log.Printf("Fail to register NSQ message latency with error: %v", err)
		return nil
	}
	return &NSQClient{
		config:        *config,
		msgCounter:    msgCounter,
		msgCounterVec: *msgCounterVec,
		mqLatency:     msgHistogram,
		dbstore:       store,
	}
}

func (n NSQClient) Publish(topic string, message []byte) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		n.mqLatency.Observe(v)
	}))
	defer timer.ObserveDuration()
	publisher, err := nsq.NewProducer("127.0.0.1:4150", &n.config)
	if err != nil {
		return err
	}

	err = publisher.Publish(topic, message)
	if err != nil {
		return err
	}
	n.msgCounter.Inc()
	counterVec := n.msgCounterVec.WithLabelValues("200", "POST")
	counterVec.Inc()
	return nil
}

func (n NSQClient) Subscribe(topic string) error {
	nsqSubscriber, err := nsq.NewConsumer(topic, "channel", &n.config)
	if err != nil {
		return err
	}
	nsqSubscriber.AddHandler(&NSQMessageHandler{n.dbstore})

	// either localhost or 127.0.0.1 as address are acceptable, but prefere 127.0.0.1 for consistency
	nsqSubscriber.ConnectToNSQLookupd("127.0.0.1:4161")

	return nil
}
