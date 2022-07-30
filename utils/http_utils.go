package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HTTPBody struct {
	Signature []byte
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	SIGNATURE_KEY_PATTERN = "Signature"
)

var (
	Client       HTTPClient                                        = &http.Client{}
	GetSignature func(payload interface{}, workerId string) []byte = getSignature
)

type HttpResponse struct {
	Body       []byte
	Headers    http.Header
	StatusCode int
}

func getSignature(message interface{}, keyId string) []byte {
	parsedMessage, err := json.Marshal(message)

	if err != nil {
		log.Fatal("Error on marshalling the payload")
	}

	signature, _ := SignMessage(GetPrivateKey(keyId), parsedMessage)

	return signature
}

func Post(keyId string, message string, headers http.Header, endpoint string) (*HttpResponse, error) {
	signature := GetSignature(message, keyId)
	requestBody, err := json.Marshal(HTTPBody{Signature: signature})

	if err != nil {
		log.Fatal("Unable to marshal body")
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &HttpResponse{Body: respBody, Headers: resp.Header, StatusCode: resp.StatusCode}, nil
}
