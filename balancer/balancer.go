package balancer

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
)

func Start() error {

	interval, err := strconv.Atoi(os.Getenv("BALANCE_CHECK_INTERVAL"))
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	for {
		err = Balance()
		if err != nil {
			logger.Errorln(err.Error())
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func Balance() error {
	// TODO
	fmt.Println("Balancing ...")
	return nil
}

func Check() error {
	// TODO
	fmt.Println("Checking ...")
	return nil
}
