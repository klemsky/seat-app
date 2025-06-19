package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"seat-app-backend/controller"
	"seat-app-backend/repository"
	"seat-app-backend/service"
	"time"
)

var seatMapControllerInstance *controller.SeatMapController

func init() {

	seatMapRepo := repository.NewJSONSeatMapRepository("SeatMapResponse.json")

	seatMapService := service.NewSeatMapService(seatMapRepo)

	seatMapControllerInstance = controller.NewSeatMapController(seatMapService)
	log.Println("Backend services and controllers initialized.")
}

func main() {
	router := http.NewServeMux()

	controller.RegisterRoutes(router, seatMapControllerInstance)

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
