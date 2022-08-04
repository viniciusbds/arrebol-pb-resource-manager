package launcher

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
)

func CreateWorker(workerID string, vcpu float64, ram float64, resourceID string, c chan string) error {

	vagrantID := strings.Replace(resourceID, "-", "", -1)

	cmd := exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "mkdir %s"`, vagrantID, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "git clone https://github.com/viniciusbds/arrebol-pb-worker %s"`, vagrantID, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "mv %s/.env.example2 %s/.env"`, vagrantID, workerID, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "cp server.pub %s/certs"`, vagrantID, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "cd %s && ./create_worker_conf.sh  %f %f %s"`, vagrantID, workerID, vcpu, ram, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "cd %s && /usr/local/go/bin/go build"`, vagrantID, workerID))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`vagrant ssh %s -c "cd %s && ./arrebol-pb-worker"`, vagrantID, workerID))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	go func() {
		defer quit(workerID)
		if err := cmd.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	go func() {
		c <- fmt.Sprintf("finishing worker [%s]", workerID)
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			fmt.Println(err.Error())
		}
	}()
	return nil
}

func quit(workerID string) {
	fmt.Printf("Deleting consumption for worker [%s]\n", workerID)
	storage.DB.DeleteConsumption(workerID)
}
