package autoscaler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/viniciusbds/arrebol-pb-resource-manager/utils"
)

type QueueState struct {
	QueueID           string `json:"QueueID"`
	NumWorkers        int    `json:"NumWorkers"`
	NumReadToRunTasks int    `json:"NumReadToRunTasks"`
}

func RequestStateSummary() ([]QueueState, error) {
	headers := http.Header{}

	keyId := RESOURCE_MANAGER_KEY_NAME
	message := os.Getenv(RM_AUTH_MESSAGE)
	endpoint := os.Getenv(MAIN_API_ENDPOINT) + "/queues/statesummary"

	httpResponse, err := utils.Post(keyId, message, headers, endpoint)
	if err != nil {
		log.Fatal("Communication error with the server: " + err.Error())
	}

	queuesStateSummary, err := HandleGetQueuesStateSummaryResponse(httpResponse)
	if err != nil {
		return []QueueState{}, err
	}

	return queuesStateSummary, nil
}

func HandleGetQueuesStateSummaryResponse(response *utils.HttpResponse) ([]QueueState, error) {
	if response.StatusCode != 200 {
		log.Fatal("The work could not be subscribed. Status Code: " + strconv.Itoa(response.StatusCode))
	}

	var parsedBody map[string][]QueueState
	err := json.Unmarshal(response.Body, &parsedBody)

	if err != nil {
		return []QueueState{}, err
	}

	workerId, ok := parsedBody["state-summary"]

	if !ok {
		log.Fatal("The worker-id is not in the response body")
	}

	return workerId, nil
}
