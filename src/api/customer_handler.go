package api

import (
	"errors"
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
		responseCustomErr(c, ErrCodeInvalidRegisterCustomerRequest, err)
		return
	}

	hashedPassword, err := s.hashVerifier.Hash(req.Password)
	if err != nil {
		responseCustomErr(c, ErrCodeHashingPassword, err)
		return
	}

	customer := &model.Account{
		RoleID:      model.RoleIDCustomer,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    hashedPassword,
		Status:      model.AccountStatusWaitingConfirmPhoneNumber,
	}

	if err := s.store.AccountStore.Create(customer); err != nil {
		responseGormErr(c, err)
		return
	}

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.PhoneNumber); err != nil {
		responseCustomErr(c, ErrCodeSendOTP, err)
		return
	}

	responseSuccess(c, gin.H{
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
		responseCustomErr(c, ErrCodeInvalidFindCarsRequest, err)
		return
	}

	if !validateStartEndDate(c, req.StartDate, req.EndDate) {
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
				responseCustomErr(c, ErrCodeInvalidFindCarsRequest, err)
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
		responseGormErr(c, err)
		return
	}
	respCars := make([]*carResponse, len(foundCars))
	for i, car := range foundCars {
		respCars[i], err = s.newCarResponse(car)
		if err != nil {
			responseGormErr(c, err)
			return
		}
	}

	responseSuccess(c, respCars)
}

type customerRentCarRequest struct {
	CarID          int                  `json:"car_id" binding:"required"`
	StartDate      time.Time            `json:"start_date" binding:"required"`
	EndDate        time.Time            `json:"end_date" binding:"required"`
	CollateralType model.CollateralType `json:"collateral_type" binding:"required"`
}

func validateStartEndDate(c *gin.Context, startDate, endDate time.Time) bool {
	if time.Now().After(startDate) {
		responseCustomErr(c, ErrCodeInvalidRentCarRequest, errors.New("start_date must be greater than now"))
		return false
	}

	if startDate.After(endDate) {
		responseCustomErr(c, ErrCodeInvalidRentCarRequest, errors.New("start_date must be less than end_date"))
		return false
	}

	if endDate.Sub(startDate) < 24*time.Hour {
		responseCustomErr(c, ErrCodeInvalidRentCarRequest, errors.New("rent period must be at least 1 day"))
		return false
	}

	return true
}

func (s *Server) HandleCustomerRentCar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := customerRentCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidRentCarRequest, err)
		return
	}

	if !validateStartEndDate(c, req.StartDate, req.EndDate) {
		return
	}

	// Check not overlap with other contracts
	isOverlap, err := s.store.CustomerContractStore.IsOverlap(req.CarID, req.StartDate, req.EndDate)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if isOverlap {
		responseCustomErr(c, ErrCodeInvalidRentCarRequest, errors.New("start_date and end_date is overlap with other contracts"))
		return
	}

	customer, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	isEmpty := func(str string) bool { return len(str) == 0 }

	if isEmpty(customer.QRCodeURL) &&
		(isEmpty(customer.BankName) || isEmpty(customer.BankNumber) || isEmpty(customer.BankOwner)) {
		responseCustomErr(c, ErrCodeMissingPaymentInformation, nil)
		return
	}

	drivingLicenseImgs, err := s.store.DrivingLicenseImageStore.Get(customer.ID, model.DrivingLicenseImageStatusActive, MaxNumberDrivingLicenseFiles)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if len(drivingLicenseImgs) != MaxNumberDrivingLicenseFiles {
		responseCustomErr(c, ErrCodeMissingDrivingLicence, nil)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.Status != model.CarStatusActive {
		responseCustomErr(c, ErrCodeInvalidCarStatus, errors.New("invalid car status. require active"))
		return
	}

	rule, err := s.store.ContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	pricing := calculateRentPrice(car, rule, req.StartDate, req.EndDate)
	contract := &model.CustomerContract{
		CustomerID:              customer.ID,
		CarID:                   req.CarID,
		RentPrice:               pricing.TotalRentPriceAmount,
		StartDate:               req.StartDate,
		EndDate:                 req.EndDate,
		Status:                  model.CustomerContractStatusWaitingContractAgreement,
		InsuranceAmount:         pricing.TotalInsuranceAmount,
		CollateralType:          req.CollateralType,
		CollateralCashAmount:    rule.CollateralCashAmount,
		InsurancePercent:        rule.InsurancePercent,
		PrepayPercent:           rule.PrepayPercent,
		RevenueSharingPercent:   rule.RevenueSharingPercent,
		BankName:                customer.BankName,
		BankNumber:              customer.BankNumber,
		BankOwner:               customer.BankOwner,
		IsReturnCollateralAsset: false,
	}
	if err := s.store.CustomerContractStore.Create(contract); err != nil {
		responseGormErr(c, err)
		return
	}

	go func() {
		_ = s.RenderCustomerContractPDF(customer, car, contract)
	}()

	responseSuccess(c, contract)
}

type customerAgreeContractRequest struct {
	CustomerContractID int    `json:"customer_contract_id" binding:"required"`
	ReturnURL          string `json:"return_url" binding:"required"`
}

func (s *Server) HandleCustomerAgreeContract(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	req := customerAgreeContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCustomerAgreeContractRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.CustomerID != acct.ID {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	if contract.Status != model.CustomerContractStatusWaitingContractAgreement {
		responseCustomErr(c, ErrCodeInvalidCustomerContractStatus, errors.New("invalid customer contract status"))
		return
	}

	if err := s.store.CustomerContractStore.Update(
		req.CustomerContractID,
		map[string]interface{}{"status": string(model.CustomerContractStatusWaitingContractPayment)},
	); err != nil {
		responseGormErr(c, err)
		return
	}

	rule, err := s.store.ContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	pricing := calculateRentPrice(&contract.Car, rule, contract.StartDate, contract.EndDate)
	prepayPayment, err := s.generateCustomerContractPaymentQRCode(
		contract.ID,
		pricing.PrepaidAmount,
		model.PaymentTypePrePay,
		req.ReturnURL,
		"",
	)
	if err != nil {
		responseCustomErr(c, ErrCodeGenerateQRCode, err)
		return
	}

	// if this is collateral cash type, combine prepay + collateral pay into one payment_url
	if contract.CollateralType == model.CollateralTypeCash && contract.CollateralCashAmount > 0 {
		collateralPayment, err := s.generateCustomerContractPaymentQRCode(
			contract.ID, contract.CollateralCashAmount, model.PaymentTypeCollateralCash, req.ReturnURL, "")
		if err != nil {
			responseCustomErr(c, ErrCodeGenerateQRCode, err)
			return
		}

		combined, err := s.paymentService.GeneratePaymentURL(
			[]int{prepayPayment.ID, collateralPayment.ID},
			prepayPayment.Amount+collateralPayment.Amount,
			time.Now().Format("02150405"),
			req.ReturnURL,
		)
		if err != nil {
			responseCustomErr(c, ErrCodeGenerateQRCode, err)
			return
		}

		responseSuccess(c, gin.H{"payment_url": combined})
		return
	}

	responseSuccess(c, gin.H{"payment_url": prepayPayment.PaymentURL})
}

func (s *Server) generateCustomerContractPaymentQRCode(
	contractID int,
	amount int,
	paymentType model.PaymentType,
	returnURL string,
	note string,
) (*model.CustomerPayment, error) {
	payment := &model.CustomerPayment{
		CustomerContractID: contractID,
		PaymentType:        paymentType,
		Amount:             amount,
		Status:             model.PaymentStatusPending,
		Note:               note,
	}
	if err := s.store.CustomerPaymentStore.Create(payment); err != nil {
		return nil, err
	}

	url, err := s.paymentService.GeneratePaymentURL([]int{payment.ID}, amount, time.Now().Format("02150405"), returnURL)
	if err != nil {
		return nil, err
	}

	if err := s.store.CustomerPaymentStore.Update(payment.ID, map[string]interface{}{"payment_url": url}); err != nil {
		return nil, err
	}

	updatedPayment, err := s.store.CustomerPaymentStore.GetByID(payment.ID)
	if err != nil {
		return nil, err
	}

	return updatedPayment, nil
}

type customerGetContractsRequest struct {
	Pagination
	ContractStatus string `form:"contract_status"`
}

func (s *Server) HandleCustomerGetContracts(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := customerGetContractsRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetCustomerContractRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if len(req.ContractStatus) > 0 {
		status = model.CustomerContractStatus(req.ContractStatus)
	}

	contracts, err := s.store.CustomerContractStore.GetByCustomerID(acct.ID, status, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, contracts)
}

func (s *Server) HandleCustomerAdminGetCustomerContractDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	id := c.Param("customer_contract_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		responseCustomErr(c, ErrCodeInvalidGetCustomerContractDetailRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(idInt)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if authPayload.Role == model.RoleNameCustomer && contract.CustomerID != acct.ID {
		responseCustomErr(c, ErrCodeInvalidOwnership, err)
		return
	}

	if acct.Role.RoleName == model.RoleNameAdmin {
		responseSuccess(c, s.newCustomerContractResponse(contract))
		return
	}

	responseSuccess(c, contract)
}

type calculateRentingPricingRequest struct {
	CarID     int       `form:"car_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
}

func (s *Server) HandleCustomerCalculateRentPricing(c *gin.Context) {
	req := calculateRentingPricingRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCalculateRentingPriceRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	rule, err := s.store.ContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, calculateRentPrice(car, rule, req.StartDate, req.EndDate))
}

type getLastPaymentDetailRequest struct {
	CustomerContractID int               `form:"customer_contract_id" binding:"required"`
	PaymentType        model.PaymentType `form:"payment_type" binding:"required"`
}

func (s *Server) HandleCustomerGetLastPaymentDetail(c *gin.Context) {
	req := getLastPaymentDetailRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeGetLastPaymentTypeRequest, err)
		return
	}

	paymentDetail, err := s.store.CustomerPaymentStore.GetLastByPaymentType(req.CustomerContractID, req.PaymentType)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, paymentDetail)
}

type customerGetActivitiesRequest struct {
	Pagination `json:"pagination"`
	Status     string `form:"status"`
}

func (s *Server) HandleCustomerGetActivities(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := customerGetActivitiesRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCustomerGetActivitiesRequest, err)
		return
	}

	cus, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if len(req.Status) > 0 {
		status = model.CustomerContractStatus(req.Status)
	}

	contracts, err := s.store.CustomerContractStore.GetByCustomerID(cus.ID, status, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, contracts)
}

type customerGiveFeedbackRequest struct {
	CustomerContractID int    `json:"customer_contract_id" binding:"required"`
	Content            string `json:"content" binding:"required,max=1000"`
	Rating             int    `json:"rating" binding:"required,max=5"`
}

func (s *Server) HandleCustomerGiveFeedback(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	req := customerGiveFeedbackRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGiveFeedbackRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.CustomerID != acct.ID {
		responseCustomErr(c, ErrCodeInvalidOwnership, err)
		return
	}

	if contract.Status != model.CustomerContractStatusCompleted {
		responseCustomErr(c, ErrCodeInvalidCustomerContractStatus, errors.New("invalid contract status. require completed"))
		return
	}

	updateParams := map[string]interface{}{
		"feedback_content": req.Content,
		"feedback_rating":  req.Rating,
		"feedback_status":  model.FeedbackStatusActive,
	}

	if err := s.store.CustomerContractStore.Update(req.CustomerContractID, updateParams); err != nil {
		responseGormErr(c, err)
		return
	}

	updateContract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, updateContract)
}

type getFeedbackByCarRequest struct {
	Pagination
	CarID int `form:"car_id" binding:"required"`
}

func (s *Server) HandleGetFeedbackByCar(c *gin.Context) {
	req := getFeedbackByCarRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeMissingDrivingLicence, err)
		return
	}

	feedbacks, counter, err := s.store.CustomerContractStore.GetFeedbacksByCarID(
		req.CarID,
		req.Offset,
		req.Limit,
		model.FeedbackStatusActive,
	)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"total": counter, "feedbacks": feedbacks})
}

type RentPricing struct {
	RentPriceQuotation      int `json:"rent_price_quotation"`
	InsurancePriceQuotation int `json:"insurance_price_quotation"`

	TotalRentPriceAmount int `json:"total_rent_price_amount"`
	TotalInsuranceAmount int `json:"total_insurance_amount"`
	TotalAmount          int `json:"total_amount"`
	PrepaidAmount        int `json:"prepaid_amount"`
}

func calculateRentPrice(car *model.Car, rule *model.ContractRule, startDate, endDate time.Time) *RentPricing {
	totalRentPriceAmount := car.Price * int(((endDate.Sub(startDate)).Hours())/24.0)
	totalInsuranceAmount := float64(totalRentPriceAmount) * rule.InsurancePercent / 100.0
	return &RentPricing{
		RentPriceQuotation:      car.Price,
		InsurancePriceQuotation: int(float64(car.Price) * rule.InsurancePercent / 100.0),
		TotalRentPriceAmount:    totalRentPriceAmount,
		TotalInsuranceAmount:    int(totalInsuranceAmount),
		TotalAmount:             totalRentPriceAmount + int(totalInsuranceAmount),
		PrepaidAmount:           int((float64(totalRentPriceAmount) + totalInsuranceAmount) * rule.PrepayPercent / 100.0),
	}
}
