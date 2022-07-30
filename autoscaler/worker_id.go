package autoscaler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/viniciusbds/arrebol-pb-resource-manager/utils"
)

func RequestWorkerId() (string, error) {
	headers := http.Header{}

	keyId := RESOURCE_MANAGER_KEY_NAME
	message := os.Getenv(RM_AUTH_MESSAGE)
	endpoint := os.Getenv(WORKER_API_ENDPOINT) + "/workers/id"

	httpResponse, err := utils.Post(keyId, message, headers, endpoint)
	if err != nil {
		log.Fatal("Communication error with the server: " + err.Error())
	}

	workerId, err := HandleGetWorkerIdResponse(httpResponse)
	if err != nil {
		return "", err
	}

	return workerId, nil
}

func HandleGetWorkerIdResponse(response *utils.HttpResponse) (string, error) {
	if response.StatusCode != 200 {
		log.Fatal("The work could not be subscribed. Status Code: " + strconv.Itoa(response.StatusCode))
	}

	var parsedBody map[string]string
	err := json.Unmarshal(response.Body, &parsedBody)

	if err != nil {
		return "", err
	}

	workerId, ok := parsedBody["worker-id"]

	if !ok {
		log.Fatal("The worker-id is not in the response body")
	}

	return workerId, nil
}
