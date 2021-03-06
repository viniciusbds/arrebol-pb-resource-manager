package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/viniciusbds/arrebol-pb-resource-manager/internal"
)

func CreateWorker(workerID string, queueID string, vcpu float64, ram float64, node string) error {

	vagrantfilePath := path.Join(internal.VAGRANT_PATH, node)

	cmd := exec.Command("cp", "../launcher/scripts/startup_worker.sh", vagrantfilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("bash", path.Join(vagrantfilePath, "startup_worker.sh"),
		vagrantfilePath,
		fmt.Sprintf("%f", vcpu),
		fmt.Sprintf("%f", ram),
		workerID,
		queueID,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
