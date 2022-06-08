package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	ErrProcReq   = errors.New("error while trying to process response")
	ErrEncodeRes = errors.New("error while trying encode response")
)

type Version struct {
	Tag  string `json:"Tag"`
	Name string `json:"Name"`
}

func Write(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println(ErrEncodeRes)
	}
}

func (a *HttpApi) GetVersion(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v1/version getVersion
	//
	// Retrieve the system version
	// ---
	// consumes:
	// - application/json
	// produces:
	// - text/plain
	// responses:
	//   '200':
	//     description: The system version
	//     type: string
	Write(w, http.StatusOK, Version{Tag: os.Getenv("VERSION_TAG"), Name: os.Getenv("VERSION_NAME")})
}
