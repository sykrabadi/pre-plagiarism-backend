package transport

import (
	"go-nsq/application/prefalsification"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// func InitRoutes(Mongo *db.Mongo) *mux.Router {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/uploadFile", uploadFile).Methods(http.MethodPost)
// 	return router
// }

// func uploadFile(w http.ResponseWriter, r *http.Request) {
// 	db := clien
// 	response := []byte
// 	w.Write(response)
// }

type server struct {
	prefalsificationService prefalsification.IPrefalsificationService
}

func NewHTTPServer(
	prefalsificationService prefalsification.IPrefalsificationService,
) *mux.Router {
	router := mux.NewRouter()
	server := server{
		prefalsificationService: prefalsificationService,
	}
	router.HandleFunc("/sendDocument", server.SendDocument).Methods(http.MethodPost)

	return router
}

func (s *server) SendDocument(w http.ResponseWriter, r *http.Request) {
	err := s.prefalsificationService.SendData()
	if err != nil {
		log.Println("Error sending data")
	}
	log.Println("Upload Document Success")
}
