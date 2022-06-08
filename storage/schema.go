package storage

import (
	"errors"
	"fmt"

	"github.com/google/logger"
	"github.com/jinzhu/gorm"
)

func (s *Storage) DropTablesIfExist() *gorm.DB {
	return s.driver.DropTableIfExists(&Consumption{}, &Resource{})
}

func (s *Storage) CreateTables() {
	var tables = map[string]interface{}{
		"resources":   &Resource{},
		"consumption": &Consumption{},
	}

	for _, v := range tables {
		_, err := s.CreateTable(v)

		if err != nil {
			logger.Errorln(err.Error())
		}
	}
}

func (s *Storage) CreateTable(t interface{}) (string, error) {
	clone := s.driver.CreateTable(t)

	if clone.Error != nil {
		var errMsg = fmt.Sprintf("Table %+v already exists", t)
		return errMsg, errors.New(errMsg)
	} else {
		var successMsg = fmt.Sprintf("Table %+v correctly created", t)
		return successMsg, nil
	}
}

func (s *Storage) AutoMigrate() {
	s.driver.AutoMigrate(&Resource{})
}

func (s *Storage) ConfigureSchema() {
	s.Driver().Model(
		&Consumption{}).AddForeignKey(
		"resource_id", "resources(id)", "CASCADE", "CASCADE")
}

func (s *Storage) CreateSchema() {
	s.DropTablesIfExist()
	s.CreateTables()
	s.AutoMigrate()
	s.ConfigureSchema()

}

type Resource struct {
	gorm.Model
	CPU     int8   `json:"CPU"`
	RAM     int8   `json:"RAM"` //in mb
	Address string `json:"Address"`
}

func (c Resource) String() string {
	return fmt.Sprintf("[ID: %d, CPU: %d, RAM:%d, Address: %s]", c.ID, c.CPU, c.RAM, c.Address)
}

type Consumption struct {
	gorm.Model
	ResourceID uint `json:"ResourceID"`
	CPU        int8 `json:"CPU"`
	RAM        int8 `json:"RAM"` //in mb
}

func (c Consumption) String() string {
	return fmt.Sprintf("[ID: %d, ResourceID: %d, CPU: %d, RAM:%d]", c.ID, c.ResourceID, c.CPU, c.RAM)
}
