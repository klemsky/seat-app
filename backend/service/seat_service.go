package service

import (
	"seat-app-backend/model"
	"seat-app-backend/repository"
)

type SeatMapService interface {
	GetFullSeatMap() (*model.SeatMapResponse, error)
}

type seatMapService struct {
	repo repository.SeatMapRepository
}

func NewSeatMapService(repo repository.SeatMapRepository) SeatMapService {
	return &seatMapService{repo: repo}
}

func (s *seatMapService) GetFullSeatMap() (*model.SeatMapResponse, error) {
	return s.repo.GetSeatMapData()
}
