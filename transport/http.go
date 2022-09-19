package transport

import (
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

func NewHTTPServer() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/updateDocument", Test).Methods(http.MethodPost)

	return router
}

func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling Test Handler")

	// TODO : Call SendData() from store here
}
