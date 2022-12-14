package preplagiarism

import (
	"go-nsq/transport/rest"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *PrePlagiraism) SendToRest(payload []byte) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64){
		c.ClientLatency.Observe(v)
	}))
	defer timer.ObserveDuration()
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	request := rest.RestRequest{
		Client: c.client,
		Header: header,
		Method: resty.MethodPost,
		Body: payload,
		URL: c.url + "/sendData",
	}
	
	err := rest.ExecuteRestyRequest(&request)
	if err != nil {
		log.Printf("fail to send to REST server with error %v \n", err)
		return err
	}

	c.ClientMsgCounter.Inc()
	return nil
}
