package api

import (
	"errors"
	"strconv"
	"time"

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

	isEmpty := func(s string) bool {
		return len(s) == 0
	}

	if (isEmpty(acct.BankName) || isEmpty(acct.BankOwner) || isEmpty(acct.BankNumber)) && isEmpty(acct.QRCodeURL) {
		responseCustomErr(c, ErrCodeMissingPaymentInformation, err)
		return
	}

	req := registerCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidRegisterCarRequest, err)
		return
	}

	rule, err := s.store.ContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	period, err := strconv.Atoi(req.PeriodCode)
	if err != nil {
		responseCustomErr(c, ErrCodeInvalidRegisterCarRequest, err)
		return
	}

	now := time.Now()
	car := &model.Car{
		PartnerID:             acct.ID,
		CarModelID:            req.CarModelID,
		LicensePlate:          req.LicensePlate,
		ParkingLot:            req.ParkingLot,
		Description:           req.Description,
		Fuel:                  req.Fuel,
		Motion:                req.Motion,
		Price:                 0,
		RevenueSharingPercent: rule.RevenueSharingPercent,
		BankName:              acct.BankName,
		BankNumber:            acct.BankNumber,
		BankOwner:             acct.BankOwner,
		StartDate:             now,
		EndDate:               now.AddDate(0, period, 0),
		PartnerContractStatus: model.PartnerContractStatusWaitingForApproval,
		Status:                model.CarStatusPendingApplicationPendingCarImages,
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

	go func() {
		adminIds, err := s.store.AccountStore.GetAllAdminIDs()
		if err == nil {
			for _, id := range adminIds {
				s.adminNotificationQueue <- s.NewCarRegisterNotificationMsg(id, car.ID)
			}
		}
	}()

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

	avgRating, err := s.store.CustomerContractStore.GetAverageRating(car.ID)
	if err != nil {
		return nil, err
	}

	totalTrip, err := s.store.CustomerContractStore.GetTotalCompletedContracts(car.ID)
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
		Rating:       avgRating,
		TotalTrip:    totalTrip,
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

	if car.PartnerContractStatus != model.PartnerContractStatusWaitingForAgreement {
		responseCustomErr(c, ErrCodeInvalidPartnerContractStatus, errors.New("invalid contract status"))
		return
	}

	if car.ParkingLot == model.ParkingLotGarage {
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

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"partner_contract_status": string(model.PartnerContractStatusAgreed),
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

	go func() {
		adminIds, err := s.store.AccountStore.GetAllAdminIDs()
		if err == nil {
			for _, id := range adminIds {
				msg := s.NewCarDeliveryNotificationMsg(id, car.ID, car.LicensePlate)
				if status == string(model.CarStatusActive) {
					msg = s.NewCarActiveNotificationMsg(id, car.ID, car.LicensePlate)
				}

				s.adminNotificationQueue <- msg
			}
		}
	}()

	responseSuccess(c, gin.H{"status": "agree contract successfully"})
}

type getContractRequest struct {
	CarID int `form:"car_id"`
}

func (s *Server) HandleGetPartnerContractDetail(c *gin.Context) {
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

	responseSuccess(c, car.ToPartnerContract())
}

type partnerGetActivityDetailRequest struct {
	Pagination
	CarID                  int    `form:"car_id" binding:"required"`
	CustomerContractStatus string `form:"customer_contract_status"`
}

func (s *Server) HandlePartnerGetActivityDetail(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	req := partnerGetActivityDetailRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidPartnerGetActivityDetailRequest, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if len(req.CustomerContractStatus) > 0 {
		status = model.CustomerContractStatus(req.CustomerContractStatus)
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.PartnerID != acct.ID {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	contracts, err := s.store.CustomerContractStore.FindByCarID(req.CarID, status, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	type contractWNetReceive struct {
		*model.CustomerContract
		NetReceive int `json:"net_receive"`
	}
	resp := make([]*contractWNetReceive, len(contracts))
	for i, contract := range contracts {
		wNetReceive := &contractWNetReceive{
			contract,
			contract.RentPrice * int(100.0-contract.RevenueSharingPercent) / 100}
		resp[i] = wNetReceive
	}

	responseSuccess(c, resp)
}

func (s *Server) HandlePartnerGetRevenue(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	payments, err := s.store.PartnerPaymentHistoryStore.GetRevenue(acct.ID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	revenue := 0
	for _, p := range payments {
		revenue += p.Amount
	}

	responseSuccess(c, gin.H{"total_revenue": revenue, "payments:": payments})
}

func (s *Server) fromUUIDToURL(uuid, extension string) string {
	return s.s3store.Config.BaseURL + uuid + "." + extension
}
