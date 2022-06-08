package storage

func (s *Storage) SaveResource(r *Resource) error {
	return s.driver.Save(&r).Error
}

func (s *Storage) RetrieveResource(resourceID uint) (*Resource, error) {
	var resource Resource
	err := s.driver.First(&resource, resourceID).Error
	return &resource, err
}

func (s *Storage) RetrieveResources() ([]*Resource, error) {
	var resources []*Resource

	err := s.driver.Find(&resources).Error

	return resources, err
}
