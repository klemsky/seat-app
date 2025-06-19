package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"seat-app-backend/model"
	"sync"
)

type SeatMapRepository interface {
	GetSeatMapData() (*model.SeatMapResponse, error)
	UpdateSeat(seatCode string) (*model.SeatMapResponse, error)
	AddSelectedSeat(seat model.Seat) (*model.SeatMapResponse, error)
}

type JSONSeatMapRepository struct {
	filePath string
	data     *model.SeatMapResponse
	mu       sync.Mutex
}

func NewJSONSeatMapRepository(filePath string) *JSONSeatMapRepository {
	repo := &JSONSeatMapRepository{
		filePath: filePath,
	}

	err := repo.loadData()
	if err != nil {
		log.Fatalf("Failed to load seat map data from %s: %v", filePath, err)
	}
	return repo
}

func (r *JSONSeatMapRepository) loadData() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	byteValue, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file %s: %w", r.filePath, err)
	}

	r.data = &model.SeatMapResponse{}
	err = json.Unmarshal(byteValue, r.data)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON from %s: %w", r.filePath, err)
	}
	log.Printf("Seat map data loaded successfully from %s.", r.filePath)
	return nil
}

func (r *JSONSeatMapRepository) saveDataInternal() error {
	if r.data == nil {
		return fmt.Errorf("no data to save")
	}

	updatedJSON, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling updated data: %w", err)
	}

	err = os.WriteFile(r.filePath, updatedJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated data to file %s: %w", r.filePath, err)
	}
	log.Printf("Seat map data saved successfully to %s.", r.filePath)
	return nil
}

func (r *JSONSeatMapRepository) GetSeatMapData() (*model.SeatMapResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.data == nil {
		return nil, os.ErrNotExist
	}

	dataCopy := *r.data
	return &dataCopy, nil
}

func (r *JSONSeatMapRepository) UpdateSeat(seatCode string) (*model.SeatMapResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, fmt.Errorf("seat map data not loaded")
	}

	found := false
	for _, part := range r.data.SeatsItineraryParts {
		for _, segmentMap := range part.SegmentSeatMaps {
			for _, passengerMap := range segmentMap.PassengerSeatMaps {
				for _, cabin := range passengerMap.SeatMap.Cabins {
					for _, row := range cabin.SeatRows {
						for i := range row.Seats {
							seat := &row.Seats[i]
							if seat.Code == seatCode && seat.Available && seat.StorefrontSlotCode == "SEAT" {
								seat.Available = false

								seat.Prices = nil
								seat.Taxes = nil
								seat.Total = nil
								found = true
								log.Printf("Seat %s marked as taken in memory.", seatCode)
								break
							}
						}
						if found {
							break
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("seat %s not found or not available for selection", seatCode)
	}

	err := r.saveDataInternal()
	if err != nil {
		return nil, fmt.Errorf("failed to save updated seat map: %w", err)
	}

	dataCopy := *r.data
	return &dataCopy, nil
}

func (r *JSONSeatMapRepository) AddSelectedSeat(seat model.Seat) (*model.SeatMapResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, fmt.Errorf("seat map data not loaded")
	}

	var seatInterface interface{}
	seatJSON, err := json.Marshal(seat)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal seat for selectedSeats: %w", err)
	}
	err = json.Unmarshal(seatJSON, &seatInterface)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal seat to interface{} for selectedSeats: %w", err)
	}

	r.data.SelectedSeats = append(r.data.SelectedSeats, seatInterface)
	log.Printf("Seat %s added to selectedSeats in memory.", seat.Code)

	err = r.saveDataInternal()
	if err != nil {
		return nil, fmt.Errorf("failed to save updated seat map with selected seats: %w", err)
	}

	dataCopy := *r.data
	return &dataCopy, nil
}
