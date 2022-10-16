package transport

import (
	"encoding/json"
	"go-nsq/application/entrypoint"
	"go-nsq/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	entryPointService entrypoint.IEntryPointService
}

func NewHTTPServer(
	router *mux.Router,
	entryPointService entrypoint.IEntryPointService,
) http.Handler {
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
