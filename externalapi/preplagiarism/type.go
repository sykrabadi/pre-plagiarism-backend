package preplagiarism

import (
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrePlagiraism struct {
	client        *resty.Client
	url           string
	ClientMsgCounter    prometheus.Counter
	ClientLatency prometheus.Histogram
}

type IPrePlagiarism interface{
	SendToRest([]byte) error
}

func NewPrePlagiarismClient(url string, client *resty.Client) (IPrePlagiarism){
	counterReg := prometheus.NewRegistry()
	counterMetric := promauto.With(counterReg).NewCounter(
		prometheus.CounterOpts{
			Name: "REST_client_message_pumped_count",
			Help: "Number of messages sent by REST client",
		},
	)
	prometheus.Register(counterMetric)

	latencyReg := prometheus.NewRegistry()
	latencyMetric := promauto.With(latencyReg).NewHistogram(
		prometheus.HistogramOpts{
			Name: "REST_client_latency",
			Help: "Latency of REST client in seconds",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
	)
	prometheus.Register(latencyMetric)

	return &PrePlagiraism{
		client: client,
		url: url,
		ClientLatency: latencyMetric,
		ClientMsgCounter: counterMetric,
	}
}
