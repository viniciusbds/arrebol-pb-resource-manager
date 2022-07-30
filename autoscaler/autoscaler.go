package autoscaler

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/viniciusbds/arrebol-pb-resource-manager/launcher"
	resourceProvider "github.com/viniciusbds/arrebol-pb-resource-manager/resource-provider"
	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
)

var (
	RUNNING  = true
	channels map[string](chan string)
	balancer Balancer
)

const (
	PUBLIC_KEY                = "Public-Key"
	RESOURCE_MANAGER_KEY_NAME = "resource-manager"
	SERVER_ENDPOINT           = "SERVER_ENDPOINT"
	RM_AUTH_MESSAGE           = "RM_AUTH_MESSAGE"
	DEFAULT_RAM               = 1024
	DEFAULT_CPU               = 1
)

func Start() error {
	fmt.Println("Starting autoscaler...")

	balancer = NewDefaultBalancer()

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
	fmt.Println("Balancing ...")

	queuesState, err := RequestStateSummary()
	if err != nil {
		return err
	}

	for _, qs := range queuesState {
		re, err := balancer.Check(qs)
		if err != nil {
			return err
		}

		if re > 0 {

			for i := 0; i < re; i++ {
				if err := AddWorker(qs.QueueID, DEFAULT_CPU, DEFAULT_RAM); err != nil {
					return err
				}
			}

		} else if re < 0 {

			numWorkersToRemove := -1 * re
			workers, err := storage.DB.RetrieveConsumptionByQueue(qs.QueueID)
			if err != nil {
				return err
			}

			if len(workers) < numWorkersToRemove {
				return errors.New("error: there are less workers than the number of desired to remove")
			}

			for i := 0; i < numWorkersToRemove; i++ {
				worker := workers[i]
				if err := RemoveWorker(worker.WorkerID); err != nil {
					return err
				}
			}

		} else {
			fmt.Println("No need balance ...")
		}
	}

	return nil
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
		QueueID:    queueId,
		CPU:        vcpu,
		RAM:        ram,
	})
	if err != nil {
		return err
	}

	return nil
}

func RequestWorkerId() (string, error) {
	headers := http.Header{}

	keyId := RESOURCE_MANAGER_KEY_NAME
	message := os.Getenv(RM_AUTH_MESSAGE)
	endpoint := os.Getenv(SERVER_ENDPOINT) + "/workers/id"

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
