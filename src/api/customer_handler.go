package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

type customerFindCarsRequest struct {
	StartDate     time.Time        `json:"start_date" binding:"required"`
	EndDate       time.Time        `json:"end_date" binding:"required"`
	Brand         string           `json:"brand"`
	Fuel          string           `json:"fuel"`
	Motion        string           `json:"motion"`
	NumberOfSeats int              `json:"number_of_seats"`
	ParkingLot    model.ParkingLot `json:"parking_lot"`
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
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameCustomer {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

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
	isOverlap, err := s.store.CustomerContractStore.IsOverlap(req.CarID, req.StartDate, req.EndDate)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if isOverlap {
		responseError(c, errors.New("start_date and end_date is overlap with other contracts"))
		return
	}

	customer, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	insuranceAmount := car.Price / 10
	contract := &model.CustomerContract{
		CustomerID:              customer.ID,
		CarID:                   req.CarID,
		StartDate:               req.StartDate,
		EndDate:                 req.EndDate,
		Status:                  model.CustomerContractStatusWaitingContractSigning,
		InsuranceAmount:         insuranceAmount,
		CollateralType:          req.CollateralType,
		IsReturnCollateralAsset: false,
	}
	if err := s.store.CustomerContractStore.Create(contract); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "create customer contract successfully"})
}
