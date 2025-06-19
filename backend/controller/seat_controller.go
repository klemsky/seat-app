package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"seat-app-backend/service"
)

type SeatMapController struct {
	SeatMapService service.SeatMapService
}

func NewSeatMapController(svc service.SeatMapService) *SeatMapController {
	return &SeatMapController{
		SeatMapService: svc,
	}
}

func (c *SeatMapController) GetSeatMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data, err := c.SeatMapService.GetFullSeatMap()
	if err != nil {
		log.Printf("Error fetching seat map data: %v", err)
		http.Error(w, "Error retrieving seat map data", http.StatusInternalServerError)
		return
	}

	if data == nil {
		http.Error(w, "Seat map data not available", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func (c *SeatMapController) SelectSeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.SelectSeatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received seat selection request for seat: %s by passenger: %s", req.SeatCode, req.PassengerName)

	updatedData, err := c.SeatMapService.SelectSeat(req)
	if err != nil {
		log.Printf("Error selecting seat %s: %v", req.SeatCode, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updatedData)
}
