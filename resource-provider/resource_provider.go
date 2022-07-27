package resourceprovider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bmatcuk/go-vagrant"
	uuid "github.com/satori/go.uuid"
	"github.com/viniciusbds/arrebol-pb-resource-manager/internal"
	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
)

var (
	numNodes = 0
)

func AddNode(vcpu, memory float64) (string, error) {
	nodeName := fmt.Sprintf("node%d", numNodes+1)

	vagrantfilePath := path.Join(internal.VAGRANT_PATH, nodeName)

	err := os.Mkdir(vagrantfilePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	input, err := ioutil.ReadFile(path.Join(internal.VAGRANTFILE_TEMPLATE_PATH))
	if err != nil {
		return "", err
	}

	output := bytes.Replace(input, []byte("VBOX_NAME"), []byte(nodeName), -1)
	output = bytes.Replace(output, []byte("MEMORY"), []byte(fmt.Sprintf("%v", memory)), -1)
	output = bytes.Replace(output, []byte("CPUS"), []byte(fmt.Sprintf("%v", vcpu)), -1)
	output = bytes.Replace(output, []byte("ARREBOL_PUB_KEY_PATH"), []byte(os.Getenv("ARREBOL_PUB_KEY_PATH")), -1)

	if err = ioutil.WriteFile(path.Join(vagrantfilePath, "Vagrantfile"), output, os.ModePerm); err != nil {
		return "", err
	}

	client, err := vagrant.NewVagrantClient(vagrantfilePath)
	if err != nil {
		return "", err
	}

	upcmd := client.Up()
	upcmd.Verbose = true
	if err := upcmd.Run(); err != nil {
		return "", err
	}
	if upcmd.Error != nil {
		return "", err
	}
	numNodes++

	resource := &storage.Resource{
		Name:    nodeName,
		CPU:     vcpu,
		RAM:     memory,
		Address: "localhost",
	}

	vagrantID, err := getVagrantID(nodeName)
	if err != nil {
		return "", err
	}

	resource.ID, err = uuid.FromString(vagrantID)
	if err != nil {
		return "", err
	}

	if err := storage.DB.SaveResource(resource); err != nil {
		return "", err
	}

	return resource.ID.String(), nil
}

func getVagrantID(nodeName string) (string, error) {
	indexUuidPath := path.Join(internal.VAGRANT_PATH, nodeName, ".vagrant", "machines", "default", "virtualbox", "index_uuid")
	file, err := os.Open(indexUuidPath)
	if err != nil {
		return "", err
	}
	uuidBytes := make([]byte, 32)
	bytesRead, err := file.Read(uuidBytes)
	if err != nil {
		return "", err
	}
	return string(uuidBytes[:bytesRead]), nil
}

func RemoveNode(nodeName string) error {
	vagrantfilePath := path.Join(internal.VAGRANT_PATH, nodeName)

	client, err := vagrant.NewVagrantClient(vagrantfilePath)
	if err != nil {
		return err
	}
	destroycmd := client.Destroy()
	destroycmd.Verbose = true
	if err := destroycmd.Run(); err != nil {
		return err
	}
	if destroycmd.Error != nil {
		return err
	}
	if err := os.RemoveAll(path.Join(internal.VAGRANT_PATH, nodeName)); err != nil {
		return err
	}

	err = storage.DB.DeleteResource(nodeName)
	if err != nil {
		return err
	}

	return nil
}
