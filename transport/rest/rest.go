package rest

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type RestClient struct{
	Client *resty.Client
	URL string
	MsgCounter prometheus.Counter
	ClientLatency prometheus.Histogram
}

// type IRestClient interface{
// 	SendToRest([]byte) error
// }

// func NewRestClient(client *resty.Client, url string) (IRestClient, error){
// 	latencyReg := prometheus.NewRegistry()
// 	latencyMetric := promauto.With(latencyReg).NewHistogram(
// 		prometheus.HistogramOpts{
// 			Name: "REST_server_latency",
// 			Help: "Latency of REST server in seconds",
// 			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
// 		},
// 	)
// 	prometheus.Register(latencyMetric)

// 	counterReg := prometheus.NewRegistry()
// 	counterMetric := promauto.With(counterReg).NewCounter(
// 		prometheus.CounterOpts{
// 			Name: "REST_message_pumped_count",
// 			Help: "Number of message pumped to REST server",
// 		},
// 	)
// 	prometheus.Register(counterMetric)
// 	return &RestClient{
// 		URL: url,
// 		Client: client,
// 		MsgCounter: counterMetric,
// 		ClientLatency: latencyMetric,
// 	}, nil
// }

type RestRequest struct {
	Client *resty.Client
	Body interface{}
	Header map[string]string
	Method string
	URL  string
}

func ExecuteRestyRequest(req *RestRequest) error{
	client := req.Client
	if client == nil{
		log.Fatal("Resty Client Not Set")
	}

	restyReq := client.R().
	SetHeaders(req.Header).
	SetBody(req.Body)

	_, err := restyReq.Execute(req.Method, req.URL)
	if err != nil {
		log.Fatalf("error sending to %v, with error %v \n", req.URL, err)
		return err
	}
	return nil
}
