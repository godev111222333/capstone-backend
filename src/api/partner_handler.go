package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

func (s *Server) RegisterPartner(c *gin.Context) {
	req := struct {
		FirstName                string `json:"first_name"`
		LastName                 string `json:"last_name"`
		PhoneNumber              string `json:"phone_number"`
		Email                    string `json:"email"`
		IdentificationCardNumber string `json:"identification_card_number"`
		Password                 string `json:"password"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	hashedPassword, err := s.hashVerifier.Hash(req.Password)
	if err != nil {
		responseError(c, err)
		return
	}

	partner := &model.Account{
		RoleID:                   model.RoleIDPartner,
		FirstName:                req.FirstName,
		LastName:                 req.LastName,
		PhoneNumber:              req.PhoneNumber,
		Email:                    req.Email,
		IdentificationCardNumber: req.IdentificationCardNumber,
		Password:                 hashedPassword,
		Status:                   model.AccountStatusWaitingConfirmEmail,
	}

	if err := s.store.AccountStore.Create(partner); err != nil {
		responseError(c, err)
		return
	}

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.Email); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register partner successfully. please confirm OTP sent to your email",
	})
}

type registerCarRequest struct {
	LicensePlate string           `json:"license_plate" binding:"required"`
	CarModelID   int              `json:"car_model_id"`
	Motion       model.Motion     `json:"motion_code"`
	Fuel         model.Fuel       `json:"fuel_code"`
	ParkingLot   model.ParkingLot `json:"parking_lot"`
	PeriodCode   string           `json:"period_code"`
	Description  string           `json:"description"`
}

func (s *Server) HandleRegisterCar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNamePartner {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if acct.Status != model.AccountStatusActive {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("account is not active")))
		return
	}

	req := registerCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car := &model.Car{
		PartnerID:    acct.ID,
		CarModelID:   req.CarModelID,
		LicensePlate: req.LicensePlate,
		ParkingLot:   req.ParkingLot,
		Description:  req.Description,
		Fuel:         req.Fuel,
		Motion:       req.Motion,
		Price:        0,
		Status:       model.CarStatusPendingApproval,
	}

	if err := s.store.CarStore.Create(car); err != nil {
		responseInternalServerError(c, err)
		return
	}

	insertedCar, err := s.store.CarStore.GetByID(car.ID)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register car successfully",
		"car":    insertedCar,
	})
}

type updateRentalPriceRequest struct {
	CarID    int `json:"car_id" binding:"required"`
	NewPrice int `json:"new_price" binding:"required"`
}

func (s *Server) HandleUpdateRentalPrice(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := updateRentalPriceRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if car.Account.Email != authPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"price": req.NewPrice,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"car_id":    car.ID,
		"new_price": req.NewPrice,
	})
}
