package service

import "quic_upload/api/internal/storage"

type Service struct {
	storage storage.IStorage
}

func NewService(storage storage.IStorage) IService {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Create() {

}

func (s *Service) Delete() {

}
