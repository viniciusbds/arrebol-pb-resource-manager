package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
)

type HttpApi struct {
	storage *storage.Storage
	server  *http.Server
}

func New(storage *storage.Storage) *HttpApi {
	return &HttpApi{
		storage: storage,
	}
}

func (a *HttpApi) Start(port string) error {
	a.server = &http.Server{
		Addr:    ":" + port,
		Handler: a.bootRouter(),
	}
	log.Println("Service available! Running on port " + port)
	return a.server.ListenAndServe()
}

func (a *HttpApi) Shutdown() error {
	return a.server.Shutdown(context.Background())
}

func (a *HttpApi) bootRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/version", a.GetVersion).Methods(http.MethodGet)

	return router
}
