package controller

import (
	"net/http"
)

func RegisterRoutes(router *http.ServeMux, ctrl *SeatMapController) {

	router.HandleFunc("/seatmap", ctrl.GetSeatMap)
	router.HandleFunc("/select-seat", ctrl.SelectSeat)

	router.HandleFunc("/api/seatmap", ctrl.GetSeatMap)
	router.HandleFunc("/api/select-seat", ctrl.SelectSeat)
}
