package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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
	Base
	Name    string  `json:"Name"`
	CPU     float64 `json:"CPU"`
	RAM     float64 `json:"RAM"` //in mb
	Address string  `json:"Address"`
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (c Resource) String() string {
	return fmt.Sprintf("[ID: %d, Name: %s, CPU: %f, RAM:%ff, Address: %s]", c.ID, c.Name, c.CPU, c.RAM, c.Address)
}

type Consumption struct {
	gorm.Model
	WorkerID   string    `json:"WorkerID"`
	QueueID    string    `json:"QueueID"`
	ResourceID uuid.UUID `gorm:"type:uuid;foreign_key;"`
	CPU        float64   `json:"CPU"`
	RAM        float64   `json:"RAM"` //in mb
}

func (c Consumption) String() string {
	return fmt.Sprintf("[ID: %d,  WorkerId: %s,ResourceID: %d, CPU: %f, RAM:%f]", c.ID, c.WorkerID, c.ResourceID, c.CPU, c.RAM)
}
