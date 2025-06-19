package service

import (
	"errors"
	"seat-app-backend/model"
	"seat-app-backend/repository"
)

type SelectSeatRequest struct {
	SeatCode      string `json:"seatCode"`
	PassengerName string `json:"passengerName"`
}

type SeatMapService interface {
	GetFullSeatMap() (*model.SeatMapResponse, error)
	SelectSeat(req SelectSeatRequest) (*model.SeatMapResponse, error)
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

func (s *seatMapService) SelectSeat(req SelectSeatRequest) (*model.SeatMapResponse, error) {
	currentData, err := s.repo.GetSeatMapData()
	if err != nil {
		return nil, err
	}

	var selectedSeatData model.Seat
	foundSeat := false

	for _, part := range currentData.SeatsItineraryParts {
		for _, segmentMap := range part.SegmentSeatMaps {
			for _, passengerMap := range segmentMap.PassengerSeatMaps {
				for _, cabin := range passengerMap.SeatMap.Cabins {
					for _, row := range cabin.SeatRows {
						for _, seat := range row.Seats {
							if seat.Code == req.SeatCode {
								if seat.StorefrontSlotCode != "SEAT" {
									return nil, errors.New("cannot select non-seat slot")
								}
								if !seat.Available {
									return nil, errors.New("seat is already taken or unavailable")
								}
								if seat.FreeOfCharge {
									return nil, errors.New("free seats cannot be explicitly selected via this endpoint")
								}

								selectedSeatData = seat
								foundSeat = true
								break
							}
						}
						if foundSeat {
							break
						}
					}
					if foundSeat {
						break
					}
				}
				if foundSeat {
					break
				}
			}
			if foundSeat {
				break
			}
		}
		if foundSeat {
			break
		}
	}

	if !foundSeat {
		return nil, errors.New("seat not found with the given code")
	}

	_, err = s.repo.UpdateSeat(req.SeatCode)
	if err != nil {
		return nil, err
	}

	finalData, err := s.repo.AddSelectedSeat(selectedSeatData)
	if err != nil {
		return nil, err
	}

	return finalData, nil
}
