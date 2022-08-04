package storage

import uuid "github.com/satori/go.uuid"

func (s *Storage) SaveConsumption(c *Consumption) error {
	return s.driver.Save(&c).Error
}

func (s *Storage) RetrieveConsumptionByResource(resourceId uuid.UUID) ([]Consumption, error) {
	var consumptions []Consumption
	err := s.driver.Where("resource_id = ?", resourceId).Find(&consumptions).Error
	return consumptions, err
}

func (s *Storage) DeleteConsumption(workerID string) error {
	return s.driver.Where("worker_id = ?", workerID).Delete(&Consumption{}).Error
}
