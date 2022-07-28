package autoscaler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/viniciusbds/arrebol-pb-resource-manager/launcher"
	resourceProvider "github.com/viniciusbds/arrebol-pb-resource-manager/resource-provider"
	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
	"github.com/viniciusbds/arrebol-pb-resource-manager/utils"
)

var (
	RUNNING  = true
	channels map[string](chan string)
)

const (
	PUBLIC_KEY      = "Public-Key"
	SERVER_ENDPOINT = "SERVER_ENDPOINT"
	RM_PAYLOAD      = "RM_PAYLOAD"
	DEFAULT_RAM     = 1024
	DEFAULT_CPU     = 1
)

func Start() error {
	fmt.Println("Starting autoscaler...")

	channels = make(map[string](chan string)) // workerID ----> chan (string)

	interval, err := strconv.Atoi(os.Getenv("BALANCE_CHECK_INTERVAL"))
	if err != nil {
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
		return AddWorker(queueID, DEFAULT_CPU, DEFAULT_RAM)
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

func AddWorker(queueId string, vcpu float64, ram float64) (err error) {
	nodeID, err := firstAvailableNode(vcpu, ram)
	if err != nil {
		return err
	}

	if nodeID == "" {
		nodeID, err = resourceProvider.AddNode(DEFAULT_CPU, DEFAULT_RAM)
		if err != nil {
			return err
		}
	}

	workerId, err := RequestWorkerId()
	if err != nil {
		return err
	}

	channels[workerId] = make(chan string)

	err = launcher.CreateWorker(workerId, queueId, vcpu, ram, nodeID, channels[workerId])
	if err != nil {
		return err
	}

	nodeUUID, err := uuid.FromString(nodeID)
	if err != nil {
		return err
	}

	err = storage.DB.SaveConsumption(&storage.Consumption{
		ResourceID: nodeUUID,
		WorkerID:   workerId,
		CPU:        vcpu,
		RAM:        ram,
	})
	if err != nil {
		return err
	}

	return nil
}

func RequestWorkerId() (string, error) {
	publicKey, err := utils.GetBase64PubKey()
	if err != nil {
		return "", err
	}

	headers := http.Header{}
	headers.Set(PUBLIC_KEY, publicKey)

	keyId := "resource-manager"
	payload := os.Getenv(RM_PAYLOAD)
	endpoint := os.Getenv(SERVER_ENDPOINT) + "/workers/id"

	httpResponse, err := utils.RequestWorkerId(keyId, payload, headers, endpoint)
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

func RemoveWorker(workerID string) error {

	// triggers removal of worker process
	<-channels[workerID]

	fmt.Printf("Deleting consumption for worker [%s]\n", workerID)
	if err := storage.DB.DeleteConsumption(workerID); err != nil {
		return err
	}

	// verify if can RemoveNode();
	return nil
}

func firstAvailableNode(vcpu float64, ram float64) (string, error) {
	resources, err := storage.DB.RetrieveResources()
	if err != nil {
		return "", err
	}

	for _, resource := range resources {
		hasAvailable, err := hasAvailableResources(resource, vcpu, ram)
		if err != nil {
			return "", err
		}
		if hasAvailable {
			return resource.ID.String(), nil
		}
	}

	return "", nil
}

func hasAvailableResources(resource *storage.Resource, cpuNeeded, ramNeeded float64) (bool, error) {
	consumptions, err := storage.DB.RetrieveConsumptionByResource(resource.ID)
	if err != nil {
		return false, err
	}

	var (
		total_cpu_used float64
		total_ram_used float64
	)

	for _, c := range consumptions {
		total_cpu_used += c.CPU
		total_ram_used += c.RAM
	}

	hasCPU := resource.CPU-total_cpu_used >= cpuNeeded
	hasRAM := resource.RAM-total_ram_used >= ramNeeded

	return hasCPU && hasRAM, nil

}
