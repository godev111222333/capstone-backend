package api

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

func (s *Server) HandleRegisterCustomer(c *gin.Context) {
	req := struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required,phone_number"`
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

	customer := &model.Account{
		RoleID:      model.RoleIDCustomer,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    hashedPassword,
		Status:      model.AccountStatusWaitingConfirmEmail,
	}

	if err := s.store.AccountStore.Create(customer); err != nil {
		responseError(c, err)
		return
	}

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.PhoneNumber); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register customer successfully. please confirm OTP sent to your phone",
	})
}

type customerFindCarsRequest struct {
	StartDate     time.Time `form:"start_date" binding:"required"`
	EndDate       time.Time `form:"end_date" binding:"required"`
	Brands        string    `form:"brands"`
	Fuels         string    `form:"fuels"`
	Motions       string    `form:"motions"`
	NumberOfSeats string    `form:"number_of_seats"`
	ParkingLots   string    `form:"parking_lots"`
}

func (s *Server) HandleCustomerFindCars(c *gin.Context) {
	req := customerFindCarsRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	findQueries := make(map[string]interface{}, 0)
	separator := ","
	if len(req.Brands) > 0 {
		findQueries["brands"] = strings.Split(req.Brands, separator)
	}
	if len(req.Fuels) > 0 {
		findQueries["fuels"] = strings.Split(req.Fuels, separator)
	}
	if len(req.Motions) > 0 {
		findQueries["motions"] = strings.Split(req.Motions, separator)
	}
	if len(req.NumberOfSeats) > 0 {
		arr := strings.Split(req.NumberOfSeats, separator)
		arrInt := make([]int, len(arr))
		for i, s := range arr {
			var err error
			arrInt[i], err = strconv.Atoi(s)
			if err != nil {
				responseError(c, err)
				return
			}
		}
		findQueries["number_of_seats"] = arrInt
	}
	if len(req.ParkingLots) > 0 {
		findQueries["parking_lots"] = strings.Split(req.ParkingLots, separator)
	}

	foundCars, err := s.store.CarStore.FindCars(req.StartDate, req.EndDate, findQueries)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}
	respCars := make([]*carResponse, len(foundCars))
	for i, car := range foundCars {
		respCars[i], err = s.newCarResponse(car)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, respCars)
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

	if time.Now().After(req.StartDate) {
		responseError(c, errors.New("start_date must be greater than now"))
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

	customer, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if car.Status != model.CarStatusActive {
		responseError(c, errors.New("invalid car status. require active"))
		return
	}

	pricing := calculateRentPrice(car, req.StartDate, req.EndDate)
	contract := &model.CustomerContract{
		CustomerID:              customer.ID,
		CarID:                   req.CarID,
		RentPrice:               pricing.TotalRentPriceAmount,
		StartDate:               req.StartDate,
		EndDate:                 req.EndDate,
		Status:                  model.CustomerContractStatusWaitingContractAgreement,
		InsuranceAmount:         pricing.TotalInsuranceAmount,
		CollateralType:          req.CollateralType,
		IsReturnCollateralAsset: false,
	}
	if err := s.store.CustomerContractStore.Create(contract); err != nil {
		responseError(c, err)
		return
	}

	go func() {
		_ = s.RenderCustomerContractPDF(customer, car, contract)
	}()

	c.JSON(http.StatusOK, gin.H{"status": "create customer contract successfully", "contract": contract})
}

type customerAgreeContractRequest struct {
	CustomerContractID int `json:"customer_contract_id"`
}

func (s *Server) HandleCustomerAgreeContract(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}

	req := customerAgreeContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseError(c, err)
		return
	}

	if contract.CustomerID != acct.ID {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if contract.Status != model.CustomerContractStatusWaitingContractAgreement {
		responseError(c, errors.New("invalid customer contract status"))
		return
	}

	if err := s.store.CustomerContractStore.Update(
		req.CustomerContractID,
		map[string]interface{}{"status": string(model.CustomerContractStatusWaitingContractPayment)},
	); err != nil {
		responseError(c, err)
		return
	}

	url, rawURL, err := s.generatePrepayQRCode(acct.ID, contract)
	if err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "agree contract successfully", "qr_code_image": url, "payment_url": rawURL})
}

func (s *Server) generatePrepayQRCode(acctID int, contract *model.CustomerContract) (string, string, error) {
	prepayAmt := (contract.RentPrice + contract.InsuranceAmount) * 30 / 100
	payment := &model.CustomerPayment{
		CustomerContractID: contract.ID,
		PaymentType:        model.PaymentTypePrePay,
		Amount:             prepayAmt,
		Status:             model.PaymentStatusPending,
	}
	if err := s.store.CustomerPaymentStore.Create(payment); err != nil {
		return "", "", err
	}
	url, err := s.paymentService.GeneratePaymentURL(payment.ID, prepayAmt, time.Now().Format("02150405"))
	if err != nil {
		return "", "", err
	}

	qrCodeImage, err := GenerateQRCode(url)
	if err != nil {
		return "", "", err
	}

	doc, err := s.uploadDocument(bytes.NewReader(qrCodeImage), acctID, uuid.NewString()+".png", model.DocumentCategoryPrepayQRCodeImage)
	if err != nil {
		return "", "", err
	}

	if err := s.store.DocumentStore.Create(doc); err != nil {
		return "", "", err
	}

	return doc.Url, url, s.store.CustomerPaymentStore.CreatePaymentDocument(payment.ID, doc.ID)
}

type customerGetContractsRequest struct {
	Pagination
	ContractStatus string `form:"contract_status"`
}

func (s *Server) HandleCustomerGetContracts(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := customerGetContractsRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if len(req.ContractStatus) > 0 {
		status = model.CustomerContractStatus(req.ContractStatus)
	}

	contracts, err := s.store.CustomerContractStore.GetByCustomerID(acct.ID, status, req.Offset, req.Limit)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts)
}

func (s *Server) HandleCustomerAdminGetCustomerContractDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}

	id := c.Param("customer_contract_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		responseError(c, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(idInt)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if authPayload.Role == model.RoleNameCustomer && contract.CustomerID != acct.ID {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if acct.Role.RoleName == model.RoleNameAdmin {
		c.JSON(http.StatusOK, s.newCustomerContractResponse(contract))
		return
	}

	c.JSON(http.StatusOK, contract)
}

type calculateRentingPricingRequest struct {
	CarID     int       `form:"car_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
}

func (s *Server) HandleCustomerCalculateRentPricing(c *gin.Context) {
	req := calculateRentingPricingRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, calculateRentPrice(car, req.StartDate, req.EndDate))
}

type RentPricing struct {
	RentPriceQuotation      int `json:"rent_price_quotation"`
	InsurancePriceQuotation int `json:"insurance_price_quotation"`

	TotalRentPriceAmount int `json:"total_rent_price_amount"`
	TotalInsuranceAmount int `json:"total_insurance_amount"`
	TotalAmount          int `json:"total_amount"`
	PrepaidAmount        int `json:"prepaid_amount"`
}

func calculateRentPrice(car *model.Car, startDate, endDate time.Time) *RentPricing {
	totalRentPriceAmount := car.Price * int(((endDate.Sub(startDate)).Hours())/24.0)
	totalInsuranceAmount := totalRentPriceAmount / 10
	return &RentPricing{
		RentPriceQuotation:      car.Price,
		InsurancePriceQuotation: car.Price / 10,
		TotalRentPriceAmount:    totalRentPriceAmount,
		TotalInsuranceAmount:    totalInsuranceAmount,
		TotalAmount:             totalRentPriceAmount + totalInsuranceAmount,
		PrepaidAmount:           (totalRentPriceAmount + totalInsuranceAmount) * 30 / 100,
	}
}
