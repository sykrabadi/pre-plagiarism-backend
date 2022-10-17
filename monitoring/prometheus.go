package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MonitoringMetric struct{
	NSQMessageEmmitedCounter prometheus.Counter
	HttpRequestGauge prometheus.Gauge
}

func InitMonitoring() *MonitoringMetric{
	reg := prometheus.NewRegistry()
	nsqMsgCounter := promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "NSQ_message_pumped_count",
			Help: "Number of message pumped by NSQ",
		})
	httpRequestGauge := promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "HTTP_requests_gauge",
			Help: "Number of concurrent requests",
		})
	return &MonitoringMetric{
		NSQMessageEmmitedCounter: nsqMsgCounter,
		HttpRequestGauge: httpRequestGauge,
	}
}
