package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"seat-app-backend/model"
)

type SeatMapRepository interface {
	GetSeatMapData() (*model.SeatMapResponse, error)
}

type JSONSeatMapRepository struct {
	filePath string
	data     *model.SeatMapResponse
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
	jsonFile, err := os.Open(r.filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	r.data = &model.SeatMapResponse{}
	err = json.Unmarshal(byteValue, r.data)
	if err != nil {
		return err
	}
	log.Printf("Seat map data loaded successfully from %s.", r.filePath)
	return nil
}

func (r *JSONSeatMapRepository) GetSeatMapData() (*model.SeatMapResponse, error) {

	if r.data == nil {
		return nil, os.ErrNotExist
	}
	return r.data, nil
}
