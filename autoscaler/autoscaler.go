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

func RemoveWorker(queueID string) error {
	// verify if can RemoveNode();
	return nil
}
