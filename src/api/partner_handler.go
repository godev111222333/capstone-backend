package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

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
		responseCustomErr(c, ErrCodeInvalidRegisterPartnerRequest, err)
		return
	}

	hashedPassword, err := s.hashVerifier.Hash(req.Password)
	if err != nil {
		responseCustomErr(c, ErrCodeHashingPassword, err)
		return
	}

	partner := &model.Account{
		RoleID:      model.RoleIDPartner,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    hashedPassword,
		Status:      model.AccountStatusWaitingConfirmPhoneNumber,
	}

	if err := s.store.AccountStore.Create(partner); err != nil {
		responseGormErr(c, err)
		return
	}

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.PhoneNumber); err != nil {
		responseCustomErr(c, ErrCodeSendOTP, err)
		return
	}

	responseSuccess(c, gin.H{
		"status": "register partner successfully. please confirm OTP sent to your phone",
	})
}

type registerCarRequest struct {
	LicensePlate string           `json:"license_plate" binding:"required,license_plate"`
	CarModelID   int              `json:"car_model_id" binding:"required"`
	Motion       model.Motion     `json:"motion_code" binding:"required"`
	Fuel         model.Fuel       `json:"fuel_code" binding:"required"`
	ParkingLot   model.ParkingLot `json:"parking_lot" binding:"required"`
	PeriodCode   string           `json:"period_code" binding:"required"`
	Description  string           `json:"description"`
}

func (s *Server) HandleRegisterCar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	req := registerCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidRegisterCarRequest, err)
		return
	}

	period, err := strconv.Atoi(req.PeriodCode)
	if err != nil {
		responseCustomErr(c, ErrCodeInvalidRegisterCarRequest, err)
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
		responseGormErr(c, err)
		return
	}

	insertedCar, err := s.store.CarStore.GetByID(car.ID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"car": insertedCar,
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
		responseCustomErr(c, ErrCodeInvalidUpdateRentalPriceRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.Status != model.CarStatusPendingApplicationPendingPrice {
		responseCustomErr(c, ErrCodeInvalidCarStatus, errors.New("invalid state"))
		return
	}

	if car.Account.PhoneNumber != authPayload.PhoneNumber {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"price":  req.NewPrice,
		"status": model.MoveNextCarState(car.Status),
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{
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
	images, err := s.store.CarImageStore.GetByCategory(car.ID, model.CarImageCategoryImages, model.CarImageStatusActive, 5)
	if err != nil {
		return nil, err
	}

	caveats, err := s.store.CarImageStore.GetByCategory(car.ID, model.CarImageCategoryCaveat, model.CarImageStatusActive, 2)
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
		responseCustomErr(c, ErrCodeInvalidGetRegisteredCarsRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
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
		responseGormErr(c, err)
		return
	}

	carResp := make([]*carResponse, 0, len(cars))
	for _, car := range cars {
		r, err := s.newCarResponse(car)
		if err != nil {
			responseGormErr(c, err)
			return
		}
		carResp = append(carResp, r)
	}

	responseSuccess(c, getRegisteredCarsResponse{Cars: carResp})
}

type partnerAgreeContractRequest struct {
	CarID int `json:"car_id"`
}

func (s *Server) HandlePartnerAgreeContract(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := partnerAgreeContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidPartnerAgreeContractRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	partner, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.PartnerID != partner.ID {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	contract, err := s.store.PartnerContractStore.GetByCarID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.Status != model.PartnerContractStatusWaitingForAgreement {
		responseCustomErr(c, ErrCodeInvalidPartnerContractStatus, errors.New("invalid contract status"))
		return
	}

	if contract.Car.ParkingLot == model.ParkingLotGarage {
		isValid, err := s.checkIfInsertableNewSeat(car.CarModel.NumberOfSeats)
		if err != nil {
			responseGormErr(c, err)
			return
		}

		if !isValid {
			responseCustomErr(c, ErrCodeNotEnoughSlotAtGarage, errors.New("not enough slot at garage"))
			return
		}
	}

	if err := s.store.PartnerContractStore.Update(contract.ID, map[string]interface{}{
		"status": string(model.PartnerContractStatusAgreed),
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	status := string(model.CarStatusWaitingDelivery)
	if car.ParkingLot == model.ParkingLotHome {
		status = string(model.CarStatusActive)
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"status": status,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "agree contract successfully"})
}

type getContractRequest struct {
	CarID int `form:"car_id"`
}

func (s *Server) HandleGetPartnerContractDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := getContractRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetPartnerContractDetailRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if authPayload.Role == model.RoleNamePartner && car.Account.PhoneNumber != authPayload.PhoneNumber {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	contract, err := s.store.PartnerContractStore.GetByCarID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, contract)
}

func (s *Server) fromUUIDToURL(uuid, extension string) string {
	return s.s3store.Config.BaseURL + uuid + "." + extension
}
