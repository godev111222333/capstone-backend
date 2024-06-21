package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

func (s *Server) RegisterPartner(c *gin.Context) {
	req := struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
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
		RoleID:      model.RoleIDPartner,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    hashedPassword,
		Status:      model.AccountStatusWaitingConfirmEmail,
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
	CarModelID   int              `json:"car_model_id" binding:"required"`
	Motion       model.Motion     `json:"motion_code" binding:"required"`
	Fuel         model.Fuel       `json:"fuel_code" binding:"required"`
	ParkingLot   model.ParkingLot `json:"parking_lot" binding:"required"`
	PeriodCode   string           `json:"period_code" binding:"required"`
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

	req := registerCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	period, err := strconv.Atoi(req.PeriodCode)
	if err != nil {
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
		Period:       period,
		Status:       model.CarStatusPendingApplicationPendingCarImages,
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

	if car.Status != model.CarStatusPendingApplicationPendingPrice {
		responseError(c, errors.New("invalid state"))
		return
	}

	if car.Account.Email != authPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"price":  req.NewPrice,
		"status": model.MoveNextCarState(car.Status),
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"car_id":    car.ID,
		"new_price": req.NewPrice,
	})
}

type Pagination struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type carResponse struct {
	ID           int              `json:"id"`
	PartnerID    int              `json:"partner_id"`
	CarModel     model.CarModel   `json:"car_model"`
	LicensePlate string           `json:"license_plate"`
	ParkingLot   model.ParkingLot `json:"parking_lot"`
	Description  string           `json:"description"`
	Fuel         model.Fuel       `json:"fuel"`
	Motion       model.Motion     `json:"motion"`
	Price        int              `json:"price"`
	Status       model.CarStatus  `json:"status"`
	Images       []string         `json:"images"`
	Caveats      []string         `json:"caveats"`
	Rating       float64          `json:"rating"`
	TotalTrip    int              `json:"total_trip"`
	PeriodCode   int              `json:"period_code"`
}

func (s *Server) newCarResponse(car *model.Car) (*carResponse, error) {
	images, err := s.store.CarDocumentStore.GetCarDocuments(car.ID, model.DocumentCategoryCarImages, 5)
	if err != nil {
		return nil, err
	}

	caveats, err := s.store.CarDocumentStore.GetCarDocuments(car.ID, model.DocumentCategoryCaveat, 2)
	if err != nil {
		return nil, err
	}

	return &carResponse{
		ID:           car.ID,
		PartnerID:    car.PartnerID,
		CarModel:     car.CarModel,
		LicensePlate: car.LicensePlate,
		ParkingLot:   car.ParkingLot,
		Description:  car.Description,
		Fuel:         car.Fuel,
		Motion:       car.Motion,
		Price:        car.Price,
		Status:       car.Status,
		Images:       images,
		Caveats:      caveats,
		Rating:       5.0,
		TotalTrip:    1000,
		PeriodCode:   car.Period,
	}, nil
}

type getRegisteredCarsResponse struct {
	Cars []*carResponse `json:"cars"`
}

func (s *Server) HandleGetRegisteredCars(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		Pagination
		CarStatus string `form:"car_status"`
	}{}

	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseError(c, err)
		return
	}

	var status model.CarStatus
	if len(req.CarStatus) == 0 {
		status = model.CarStatusNoFilter
	} else {
		status = model.CarStatus(req.CarStatus)
	}
	cars, err := s.store.CarStore.GetByPartner(acct.ID, req.Offset, req.Limit, status)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	carResp := make([]*carResponse, 0, len(cars))
	for _, car := range cars {
		r, err := s.newCarResponse(car)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}
		carResp = append(carResp, r)
	}
	c.JSON(http.StatusOK, getRegisteredCarsResponse{Cars: carResp})
}

type partnerSignContractRequest struct {
	CarID int `json:"car_id"`
}

func (s *Server) HandlePartnerAgreeContract(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNamePartner {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	req := partnerSignContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	partner, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseError(c, err)
		return
	}

	if car.PartnerID != partner.ID {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	contract, err := s.store.PartnerContractStore.GetByCarID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if contract.Status != model.PartnerContractStatusWaitingForAgreement {
		responseError(c, errors.New("invalid contract status"))
		return
	}

	if err := s.store.PartnerContractStore.Update(contract.ID, map[string]interface{}{
		"status": string(model.PartnerContractStatusAgreed),
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"status": string(model.CarStatusWaitingDelivery),
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "agree contract successfully"})
}

type getContractRequest struct {
	CarID int `form:"car_id"`
}

func (s *Server) HandleGetPartnerContractDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin && authPayload.Role != model.RoleNamePartner {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	req := getContractRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if authPayload.Role == model.RoleNamePartner && car.Account.Email != authPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	contract, err := s.store.PartnerContractStore.GetByCarID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (s *Server) fromUUIDToURL(uuid, extension string) string {
	return s.s3store.Config.BaseURL + uuid + "." + extension
}
