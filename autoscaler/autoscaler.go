package autoscaler

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
)

var (
	RUNNING = true
)

func Start() error {

	interval, err := strconv.Atoi(os.Getenv("BALANCE_CHECK_INTERVAL"))
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	for RUNNING {
		err = Balance()
		if err != nil {
			logger.Errorln(err.Error())
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return nil
}

func Stop() error {
	logger.Infoln("Stopping autoscaler service ...")
	RUNNING = false
	return nil
}

func Balance() error {
	queueID := "1"
	// get SERVER INFO : current resource availability matches the current workload?

	re, err := CheckUnbalance()
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}

	if re > 0 {
		fmt.Println("Balancing ...")
		return AddWorker(queueID)
	}

	if re < 0 {
		fmt.Println("Balancing ...")
		return RemoveWorker(queueID)
	}

	fmt.Println("No need balance ...")
	return nil
}

func CheckUnbalance() (int, error) {
	fmt.Println("Checking ...")
	return 0, nil
}

func AddWorker(queueID string) error {
	// if need AddNode();
	return nil
}

func GetWorkerId() (string, error) {
	publicKey, err := utils.GetBase64PubKey()
	if err != nil {
		return "", err
	}

	headers := http.Header{}
	headers.Set(PUBLIC_KEY, publicKey)

	httpResponse, err := utils.RequestWorkerId(
		"resource-manager",
		os.Getenv(RM_PAYLOAD),
		headers,
		os.Getenv(SERVER_ENDPOINT)+"/workers/id",
	)
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

func RemoveWorker(queueID string) error {
	// verify if can RemoveNode();
	return nil
}
