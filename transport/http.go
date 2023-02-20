package transport

import (
	"context"
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
var sendDocumentCounterReg = prometheus.NewRegistry()
var sendDocumentLatency = promauto.With(reg).NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_sendDocument_latency",
		Help: "Latency of sendDocument endpoint",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
	},
	[]string{"status"},
)
var sendDocumentCounter = promauto.With(sendDocumentCounterReg).NewCounter(prometheus.CounterOpts{
	Name: "sendDocument_client_message_pumped_count",
	Help: "Number of message pumped from client",
})

func NewHTTPServer(
	router *mux.Router,
	entryPointService entrypoint.IEntryPointService,
) http.Handler {
	prometheus.Register(sendDocumentLatency)
	prometheus.Register(sendDocumentCounter)
	server := server{
		entryPointService: entryPointService,
	}
	router.HandleFunc("/send-document", server.SendDocument).Methods(http.MethodPost)
	router.HandleFunc("/show-document/{documentName}", server.ShowDocument).Methods(http.MethodGet)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	file, fileHeader, err := r.FormFile("file")
	r.Body = http.MaxBytesReader(w, r.Body, 5 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileContentType := fileHeader.Header.Get("Content-Type")
	if fileContentType != "application/pdf"{
		http.Error(w, "content-type must be application/pdf", http.StatusInternalServerError)
		return
	}
	filename, err := s.entryPointService.SendData(fileHeader)
	if err != nil {
		log.Println("Error sending data")
		httpWriteResponse(w, &model.ServerResponse{
			Message: "Error Sending Data",
		})
	}
	sendDocumentCounter.Inc()
	log.Println("Upload Document Success")
	httpWriteResponse(w, &model.ServerResponse{
		Message: "Success",
		File_Name: filename,
	})
}

func (s *server) ShowDocument(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	documentName := params["documentName"]
	res, err := s.entryPointService.GetDocument(context.TODO(), documentName)
	if err != nil{
		log.Printf("[server.ShowDocument] unable to GetDocument with error %v \n", err)
		return
	}
	// TODO : fix data marshalling
	httpWriteResponse(w, &model.ServerResponse{
		Message: "Success",
		Data: res,
	})
}
