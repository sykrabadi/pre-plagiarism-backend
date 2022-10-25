package transport

import (
	"encoding/json"
	"go-nsq/application/entrypoint"
	"go-nsq/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	entryPointService entrypoint.IEntryPointService
}

var reg = prometheus.NewRegistry()
var sendDocumentLatency = promauto.With(reg).NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_sendDocument_latency",
		Help: "Latency of sendDocument endpoint",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
	},
	[]string{"status"},
)

func NewHTTPServer(
	router *mux.Router,
	entryPointService entrypoint.IEntryPointService,
) http.Handler {
	prometheus.Register(sendDocumentLatency)
	server := server{
		entryPointService: entryPointService,
	}
	router.HandleFunc("/sendDocument", server.SendDocument).Methods(http.MethodPost)
	router.Handle("/metrics", promhttp.Handler())

	return router
}

func httpWriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *server) SendDocument(w http.ResponseWriter, r *http.Request) {
	var status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		sendDocumentLatency.WithLabelValues(status).Observe(f)
	}))
	defer func(){
		timer.ObserveDuration()
	}()
	if err := r.ParseMultipartForm(4096); err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	err = s.entryPointService.SendData(fileHeader)

	if err != nil {
		log.Println("Error sending data")
		httpWriteResponse(w, &model.ServerResponse{
			Message: "Error Sending Data",
		})
	}
	log.Println("Upload Document Success")
	httpWriteResponse(w, &model.ServerResponse{
		Message: "Success",
	})
}
