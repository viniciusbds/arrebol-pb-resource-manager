package storage

import (
	"fmt"

	"github.com/google/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Storage struct {
	driver *gorm.DB
}

const dbDialect string = "postgres"

var DB *Storage

func New(host string, port string, user string, dbname string, password string) *Storage {
	dbConfig := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	driver, err := gorm.Open(dbDialect, dbConfig)

	if err != nil {
		logger.Errorln(err.Error())
	}

	err = driver.DB().Ping()

	if err != nil {
		logger.Errorln(err.Error())
	}

	DB = &Storage{
		driver,
	}

	return DB
}

func (s *Storage) Setup() {
	s.CreateSchema()
}

func (s *Storage) Driver() *gorm.DB {
	return s.driver
}
