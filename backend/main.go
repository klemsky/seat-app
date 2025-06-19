package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"seat-app-backend/repository"
	"seat-app-backend/service"
	"time"
)

var seatMapServiceInstance service.SeatMapService

func init() {

	seatMapRepo := repository.NewJSONSeatMapRepository("SeatMapResponse.json")

	seatMapServiceInstance = service.NewSeatMapService(seatMapRepo)
	log.Println("Backend services initialized.")
}

func getSeatMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data, err := seatMapServiceInstance.GetFullSeatMap()
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

func selectSeat(w http.ResponseWriter, r *http.Request) {
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

	updatedData, err := seatMapServiceInstance.SelectSeat(req)
	if err != nil {
		log.Printf("Error selecting seat %s: %v", req.SeatCode, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updatedData)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/seatmap", getSeatMap)
	router.HandleFunc("/select-seat", selectSeat)

	router.HandleFunc("/api/seatmap", getSeatMap)
	router.HandleFunc("/api/select-seat", selectSeat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Go backend server listening on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}
