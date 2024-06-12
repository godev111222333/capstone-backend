package api

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/model"
)

type WaitingTime string

const (
	WaitingTimeNoTime   WaitingTime = "no"
	WaitingTimeTwoHours WaitingTime = "two_hours"
)

type customerFindCarsRequest struct {
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	Brand       string      `json:"brand"`
	Fuel        string      `json:"fuel"`
	Motion      string      `json:"motion"`
	Seats       int         `json:"seats"`
	WaitingTime WaitingTime `json:"waiting_time"`
}

func (s *Server) HandleCustomerFindCars(c *gin.Context) {

}

type customerRentCarRequest struct {
	CarID          int
	StartDate      time.Time
	EndDate        time.Time
	CollateralType model.CollateralType
}

func (s *Server) HandleCustomerRentCar(c *gin.Context) {
}
