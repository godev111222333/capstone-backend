package api

import (
	"errors"
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
	CarID          int                  `json:"car_id" binding:"required"`
	StartDate      time.Time            `json:"start_date" binding:"required"`
	EndDate        time.Time            `json:"end_date" binding:"required"`
	CollateralType model.CollateralType `json:"collateral_type" binding:"required"`
}

func (s *Server) HandleCustomerRentCar(c *gin.Context) {
	req := customerRentCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	if req.StartDate.After(req.EndDate) {
		responseError(c, errors.New("start_date must be less than end_date"))
		return
	}

	// TODO: check time range between start_date and end_date (at least 1 day?)

	// Check not overlap with other contracts
}
