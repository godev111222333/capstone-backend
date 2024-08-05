package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/service"
	"github.com/godev111222333/capstone-backend/src/token"
)

const (
	layoutDateMonthYear = "01/02/2006"
)

type getCarsRequest struct {
	Pagination
	CarStatus   string `form:"car_status"`
	SearchParam string `form:"search_param"`
}

func (s *Server) HandleAdminGetCars(c *gin.Context) {
	req := getCarsRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeGetCarsRequest, err)
		return
	}

	status := model.CarStatusNoFilter
	if len(req.CarStatus) > 0 {
		status = model.CarStatus(req.CarStatus)
	}

	cars, err := s.store.CarStore.SearchCars(req.Offset, req.Limit, status, req.SearchParam)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	total, err := s.store.CarStore.CountByStatus(status)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"cars":  cars,
		"total": total,
	})
}

func (s *Server) HandleGetCarDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responseCustomErr(c, ErrCodeGetCarDetailRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(id)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	resp, err := s.newCarResponse(car)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, resp)
}

type getGarageConfigResponse struct {
	Max4Seats      int `json:"max_4_seats"`
	Max7Seats      int `json:"max_7_seats"`
	Max15Seats     int `json:"max_15_seats"`
	Total          int `json:"total"`
	Current4Seats  int `json:"current_4_seats"`
	Current7Seats  int `json:"current_7_seats"`
	Current15Seats int `json:"current_15_seats"`
	CurrentTotal   int `json:"current_total"`
}

func (s *Server) HandleGetGarageConfigs(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin {
		responseCustomErr(c, ErrCodeInvalidRole, nil)
		return
	}

	configs, err := s.store.GarageConfigStore.Get()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	countCurrentSeats := func(seatType int) int {
		counter, err := s.store.CarStore.CountBySeats(seatType, model.ParkingLotGarage, []model.CarStatus{model.CarStatusActive, model.CarStatusWaitingDelivery})
		if err != nil {
			responseGormErr(c, err)
			return -1
		}

		return counter
	}

	cur4Seats := countCurrentSeats(4)
	cur7Seats := countCurrentSeats(7)
	cur15Seats := countCurrentSeats(15)

	responseSuccess(c, getGarageConfigResponse{
		Max4Seats:  configs[model.GarageConfigTypeMax4Seats],
		Max7Seats:  configs[model.GarageConfigTypeMax7Seats],
		Max15Seats: configs[model.GarageConfigTypeMax15Seats],
		Total: configs[model.GarageConfigTypeMax4Seats] +
			configs[model.GarageConfigTypeMax7Seats] +
			configs[model.GarageConfigTypeMax15Seats],
		Current4Seats:  cur4Seats,
		Current7Seats:  cur7Seats,
		Current15Seats: cur15Seats,
		CurrentTotal:   cur4Seats + cur7Seats + cur15Seats,
	})
}

type updateGarageConfigRequest struct {
	Max4Seats  int `json:"max_4_seats"`
	Max7Seats  int `json:"max_7_seats"`
	Max15Seats int `json:"max_15_seats"`
}

func (s *Server) HandleUpdateGarageConfigs(c *gin.Context) {
	req := updateGarageConfigRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUpdateGarageConfigRequest, err)
		return
	}

	checkValidOfSeats := func(seatType, maxSeat int) bool {
		counter, err := s.store.CarStore.CountBySeats(seatType, model.ParkingLotGarage, []model.CarStatus{model.CarStatusActive, model.CarStatusWaitingDelivery})
		if err != nil {
			responseGormErr(c, err)
			return false
		}

		if counter > maxSeat {
			responseCustomErr(c, ErrCodeInvalidSeat, errors.New(fmt.Sprintf("invalid type %d seat. Must at least %d", seatType, counter)))
			return false
		}

		return true
	}

	if !checkValidOfSeats(4, req.Max4Seats) || !checkValidOfSeats(7, req.Max7Seats) || !checkValidOfSeats(15, req.Max15Seats) {
		return
	}

	updateParams := map[model.GarageConfigType]int{
		model.GarageConfigTypeMax4Seats:  req.Max4Seats,
		model.GarageConfigTypeMax7Seats:  req.Max7Seats,
		model.GarageConfigTypeMax15Seats: req.Max15Seats,
	}

	if err := s.store.GarageConfigStore.Update(updateParams); err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "update garage configs successfully"})
}

type ApplicationAction string

const (
	ApplicationActionApproveRegister ApplicationAction = "approve_register"
	ApplicationActionApproveDelivery ApplicationAction = "approve_delivery"
	ApplicationActionReject          ApplicationAction = "reject"
)

type adminApproveOrRejectRequest struct {
	CarID  int               `json:"car_id"`
	Action ApplicationAction `json:"action"`
}

func (s *Server) HandleAdminApproveOrRejectCar(c *gin.Context) {
	req := adminApproveOrRejectRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminApproveOrRejectCarRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.ParkingLot == model.ParkingLotGarage && (req.Action == ApplicationActionApproveRegister || req.Action == ApplicationActionApproveDelivery) {
		validSeat, err := s.checkIfInsertableNewSeat(car.CarModel.NumberOfSeats)
		if err != nil {
			responseGormErr(c, err)
			return
		}

		if !validSeat {
			responseCustomErr(c, ErrCodeNotEnoughSlotAtGarage, errors.New("not enough slot at garage"))
			return
		}
	}

	newStatus := string(model.CarStatusRejected)
	if req.Action == ApplicationActionApproveRegister {
		if car.Status != model.CarStatusPendingApproval {
			responseCustomErr(c, ErrCodeInvalidCarStatus,
				fmt.Errorf(
					"invalid car status, require %s,"+
						" found %s",
					string(model.CarStatusPendingApproval),
					string(car.Status),
				),
			)
			return
		}

		newStatus = string(model.CarStatusApproved)
		go func() {
			partner, err := s.store.AccountStore.GetByPhoneNumber(car.Account.PhoneNumber)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := s.RenderPartnerContractPDF(partner, car); err != nil {
				fmt.Println(err)
			}
		}()
	}

	if req.Action == ApplicationActionApproveDelivery {
		if car.Status != model.CarStatusWaitingDelivery {
			responseCustomErr(c, ErrCodeInvalidCarStatus,
				fmt.Errorf(
					"invalid car status, require %s, found %s",
					string(model.CarStatusWaitingDelivery),
					string(car.Status),
				),
			)
			return
		}

		contract := car.ToPartnerContract()
		if contract == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "contract not found"})
			return
		}

		if contract.Status != model.PartnerContractStatusAgreed {
			responseCustomErr(c, ErrCodeInvalidPartnerContractStatus, errors.New("partner must agree the contract first"))
			return
		}

		newStatus = string(model.CarStatusActive)
	}

	updateValues := map[string]interface{}{
		"status": newStatus,
	}
	if req.Action == ApplicationActionApproveRegister {
		updateValues["partner_contract_status"] = string(model.PartnerContractStatusWaitingForAgreement)
	}

	if err := s.store.CarStore.Update(car.ID, updateValues); err != nil {
		responseInternalServerError(c, err)
		return
	}

	go func() {
		carID, phone, expoToken := car.ID, car.Account.PhoneNumber, s.getExpoToken(car.Account.PhoneNumber)
		var msg *service.PushMessage
		switch req.Action {
		case ApplicationActionReject:
			msg = s.notificationPushService.NewRejectCarMsg(car.ID, expoToken, phone)
			if car.Status == model.CarStatusActive {
				msg = s.notificationPushService.NewRejectPartnerContractMsg(carID, expoToken, phone)
			}
			break
		case ApplicationActionApproveDelivery:
			msg = s.notificationPushService.NewApproveCarDeliveryMsg(car.ID, expoToken, phone)
			break
		case ApplicationActionApproveRegister:
			msg = s.notificationPushService.NewApproveCarRegisterMsg(carID, expoToken, phone)
			break
		}

		if msg != nil {
			_ = s.notificationPushService.Push(car.Account.ID, msg)
		}
	}()

	responseSuccess(c, gin.H{"status": fmt.Sprintf("%s car successfully", req.Action)})
}

func seatNumberToGarageConfigType(seatNumber int) model.GarageConfigType {
	seatCode := model.GarageConfigTypeMax4Seats
	if seatNumber == 7 {
		seatCode = model.GarageConfigTypeMax7Seats
	} else if seatNumber == 15 {
		seatCode = model.GarageConfigTypeMax15Seats
	}

	return seatCode
}

func (s *Server) checkIfInsertableNewSeat(seatNumber int) (bool, error) {
	garageCfg, err := s.store.GarageConfigStore.Get()
	if err != nil {
		return false, err
	}

	cur, err := s.store.CarStore.CountBySeats(seatNumber, model.ParkingLotGarage, []model.CarStatus{model.CarStatusActive, model.CarStatusWaitingDelivery})
	if err != nil {
		return false, err
	}

	return cur < garageCfg[seatNumberToGarageConfigType(seatNumber)], nil
}

func convertUTCToGmt7(t time.Time) time.Time {
	return t.Add(7 * time.Hour)
}

func (s *Server) RenderPartnerContractPDF(partner *model.Account, car *model.Car) error {
	now := convertUTCToGmt7(time.Now())
	year, month, date := now.Date()

	car, err := s.store.CarStore.GetByID(car.ID)
	if err != nil {
		return err
	}

	contract := car.ToPartnerContract()
	startYear, startMonth, startDate := contract.StartDate.Date()
	endYear, endMonth, endDate := contract.EndDate.Date()

	docUUID, err := s.pdfService.Render(service.RenderTypePartner, map[string]string{
		"now_date":                strconv.Itoa(date),
		"now_month":               strconv.Itoa(int(month)),
		"now_year":                strconv.Itoa(year),
		"partner_fullname":        partner.LastName + " " + partner.FirstName,
		"partner_date_of_birth":   partner.DateOfBirth.Format(layoutDateMonthYear),
		"partner_id_card":         partner.IdentificationCardNumber,
		"brand_model":             car.CarModel.Brand + " " + car.CarModel.Model,
		"license_plate":           car.LicensePlate,
		"number_of_seats":         strconv.Itoa(car.CarModel.NumberOfSeats),
		"car_year":                strconv.Itoa(car.CarModel.Year),
		"period":                  strconv.Itoa(car.Period),
		"period_start_date":       strconv.Itoa(startDate),
		"period_start_month":      strconv.Itoa(int(startMonth)),
		"period_start_year":       strconv.Itoa(startYear),
		"period_end_date":         strconv.Itoa(endDate),
		"period_end_month":        strconv.Itoa(int(endMonth)),
		"period_end_year":         strconv.Itoa(endYear),
		"partner_revenue_percent": strconv.Itoa(int(100 - contract.RevenueSharingPercent)),
		"partner_bank_number":     contract.BankNumber,
		"partner_bank_owner":      contract.BankOwner,
		"partner_bank_name":       contract.BankName,
	})
	if err != nil {
		fmt.Printf("error when rendering partner contract %v\n", err)
		return err
	}

	if err := s.store.CarStore.Update(
		car.ID,
		map[string]interface{}{"partner_contract_url": s.fromUUIDToURL(docUUID, model.ExtensionPDF)},
	); err != nil {
		fmt.Printf("error when update partner contract URL %v\n", err)
		return err
	}

	return nil
}

func (s *Server) RenderCustomerContractPDF(
	customer *model.Account, car *model.Car,
	contract *model.CustomerContract,
) error {
	nowYear, nowMonth, nowDate := convertUTCToGmt7(time.Now()).Date()
	startDate, endDate := convertUTCToGmt7(contract.StartDate), convertUTCToGmt7(contract.EndDate)
	startHour, startDay, startMonth, startYear := startDate.Hour(), startDate.Day(), int(startDate.Month()), startDate.Year()
	endHour, endDay, endMonth, endYear := endDate.Hour(), endDate.Day(), int(endDate.Month()), endDate.Year()

	collateralAmount := contract.CustomerContractRule.CollateralCashAmount
	if collateralAmount == 0 {
		rule, err := s.store.CustomerContractRuleStore.GetLast()
		if err != nil {
			return err
		}

		collateralAmount = rule.CollateralCashAmount
	}

	docUUID, err := s.pdfService.Render(service.RenderTypeCustomer, map[string]string{
		"now_date":               strconv.Itoa(nowDate),
		"now_month":              strconv.Itoa(int(nowMonth)),
		"now_year":               strconv.Itoa(nowYear),
		"customer_fullname":      customer.LastName + " " + customer.FirstName,
		"customer_date_of_birth": customer.DateOfBirth.Format(layoutDateMonthYear),
		"customer_id_card":       customer.IdentificationCardNumber,
		"brand_model":            car.CarModel.Brand + " " + car.CarModel.Model,
		"license_plate":          car.LicensePlate,
		"number_of_seats":        strconv.Itoa(car.CarModel.NumberOfSeats),
		"car_year":               strconv.Itoa(car.CarModel.Year),
		"price":                  strconv.Itoa(contract.RentPrice),
		"prepay_percent":         fmt.Sprintf("%.2f", contract.CustomerContractRule.PrepayPercent),
		"insurance_percent":      fmt.Sprintf("%.2f", contract.CustomerContractRule.InsurancePercent),
		"start_hour":             strconv.Itoa(startHour),
		"start_date":             strconv.Itoa(startDay),
		"start_month":            strconv.Itoa(startMonth),
		"start_year":             strconv.Itoa(startYear),
		"end_hour":               strconv.Itoa(endHour),
		"end_date":               strconv.Itoa(endDay),
		"end_month":              strconv.Itoa(endMonth),
		"end_year":               strconv.Itoa(endYear),
		"bank_number":            contract.BankNumber,
		"bank_name":              contract.BankName,
		"bank_owner":             contract.BankOwner,
		"collateral_amount_1":    strconv.Itoa(collateralAmount),
		"collateral_amount_2":    strconv.Itoa(collateralAmount),
	})
	if err != nil {
		fmt.Printf("error when rendering customer contract %v\n", err)
		return err
	}

	if err := s.store.CustomerContractStore.Update(
		contract.ID,
		map[string]interface{}{"url": s.fromUUIDToURL(docUUID, model.ExtensionPDF)},
	); err != nil {
		fmt.Printf("error when update customer contract URL %v\n", err)
		return err
	}

	return nil
}

type adminGetContractRequest struct {
	Pagination
	CustomerContractStatus string `form:"customer_contract_status"`
	SearchParam            string `form:"search_param"`
}

type customerContractResponse struct {
	*model.CustomerContract
	ReceivingCarImages    []*model.CustomerContractImage `json:"receiving_car_images"`
	CollateralAssetImages []*model.CustomerContractImage `json:"collateral_asset_images"`
}

func (s *Server) newCustomerContractResponse(contract *model.CustomerContract) *customerContractResponse {
	resp := &customerContractResponse{
		CustomerContract: contract,
	}
	collaterals, err := s.store.CustomerContractImageStore.Get(
		contract.ID,
		model.CustomerContractImageCategoryCollateralAssets,
		MaxNumberCollateralAssetFiles,
		model.CustomerContractImageStatusActive,
	)
	if err == nil {
		resp.CollateralAssetImages = collaterals
	}

	receivingImages, err := s.store.CustomerContractImageStore.Get(
		contract.ID,
		model.CustomerContractImageCategoryReceivingCarImages,
		MaxNumberReceivingCarImages,
		model.CustomerContractImageStatusActive,
	)
	if err == nil {
		resp.ReceivingCarImages = receivingImages
	}

	return resp
}

func (s *Server) HandleAdminGetCustomerContracts(c *gin.Context) {
	req := adminGetContractRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetCustomerContractRequest, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if reqStatus := req.CustomerContractStatus; len(reqStatus) > 0 {
		status = model.CustomerContractStatus(reqStatus)
	}

	contracts, counter, err := s.store.CustomerContractStore.GetByStatus(status, req.Offset, req.Limit, req.SearchParam)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	contractResp := make([]*customerContractResponse, len(contracts))
	for i, contract := range contracts {
		contractResp[i] = s.newCustomerContractResponse(contract)
	}

	responseSuccess(c, gin.H{"contracts": contractResp, "total": counter})
}

type CustomerContractAction string

const (
	CustomerContractActionApprove CustomerContractAction = "approve"
	CustomerContractActionReject  CustomerContractAction = "reject"
)

type adminApproveOrRejectCustomerContractRequest struct {
	CustomerContractID int                    `json:"customer_contract_id" binding:"required"`
	Action             CustomerContractAction `json:"action" binding:"required"`
}

func (s *Server) HandleAdminApproveOrRejectCustomerContract(c *gin.Context) {
	req := adminApproveOrRejectCustomerContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminApproveOrRejectCustomerContractRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	newStatus := string(model.CustomerContractStatusCancel)
	if req.Action == CustomerContractActionApprove {
		if contract.Status != model.CustomerContractStatusOrdered {
			responseCustomErr(c, ErrCodeInvalidCustomerContractStatus, errors.New(
				fmt.Sprintf("invalid customer contract status. expect %s, found %s",
					string(model.CustomerContractStatusOrdered), string(contract.Status))))
			return
		}

		car, err := s.store.CarStore.GetByID(contract.CarID)
		if err != nil {
			responseGormErr(c, err)
			return
		}

		if car.Status != model.CarStatusActive {
			responseCustomErr(c, ErrCodeInvalidCarStatus, err)
			return
		}

		newStatus = string(model.CustomerContractStatusRenting)
	}

	if err := s.store.CustomerContractStore.Update(contract.ID, map[string]interface{}{"status": newStatus}); err != nil {
		responseGormErr(c, err)
		return
	}

	go func() {
		phone, expoToken := contract.Customer.PhoneNumber, s.getExpoToken(contract.Customer.PhoneNumber)
		msg := s.notificationPushService.NewApproveRentingCarRequestMsg(contract.ID, expoToken, phone)
		if req.Action == CustomerContractActionReject {
			msg = s.notificationPushService.NewRejectRentingCarRequestMsg(expoToken, phone)
		}

		_ = s.notificationPushService.Push(contract.CustomerID, msg)
	}()

	responseSuccess(c, contract)
}

type adminGetAccountsRequest struct {
	Pagination
	Role        string `form:"role"`
	Status      string `form:"status"`
	SearchParam string `form:"search_param"`
}

func (s *Server) HandleAdminGetAccounts(c *gin.Context) {
	req := adminGetAccountsRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetAccountsRequest, err)
		return
	}

	status := model.AccountStatusNoFilter
	if len(req.Status) > 0 {
		status = model.AccountStatus(req.Status)
	}

	accounts, err := s.store.AccountStore.Get(status, req.Role, req.SearchParam, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	respAccts := make([]*accountResponse, len(accounts))
	for i, acct := range accounts {
		respAccts[i] = s.newAccountResponse(acct)
	}

	responseSuccess(c, respAccts)
}

func (s *Server) HandleAdminGetAccountDetail(c *gin.Context) {
	id := c.Param("account_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		responseCustomErr(c, ErrCodeInvalidGetAccountDetailRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByID(idInt)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, s.newAccountResponse(acct))
}

type adminAdminGetCustomerPaymentsRequest struct {
	Pagination
	CustomerContractID int    `form:"customer_contract_id" binding:"required"`
	PaymentStatus      string `form:"payment_status"`
}

type customerPaymentResponse struct {
	*model.CustomerPayment
	Payer string `json:"payer"`
}

var AdminAsPayer = []model.PaymentType{model.PaymentTypeReturnCollateralCash, model.PaymentTypeReturnPrepay}

func newCustomerPaymentResponse(p *model.CustomerPayment) *customerPaymentResponse {
	r := &customerPaymentResponse{CustomerPayment: p}
	payer := model.RoleNameCustomer
	for _, paymentType := range AdminAsPayer {
		if p.PaymentType == paymentType {
			payer = model.RoleNameAdmin
			break
		}
	}

	r.Payer = payer
	return r
}

func (s *Server) HandleAdminGetCustomerPayments(c *gin.Context) {
	req := adminAdminGetCustomerPaymentsRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetCustomerPaymentRequest, err)
		return
	}

	status := model.PaymentStatusNoFilter
	if len(req.PaymentStatus) > 0 {
		status = model.PaymentStatus(req.PaymentStatus)
	}
	payments, err := s.store.CustomerPaymentStore.GetByCustomerContractID(req.CustomerContractID, status, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	resp := make([]*customerPaymentResponse, len(payments))
	for index, p := range payments {
		resp[index] = newCustomerPaymentResponse(p)
	}

	responseSuccess(c, resp)
}

type generateCustomerPaymentQRCode struct {
	CustomerContractID int               `json:"customer_contract_id" binding:"required"`
	ReturnURL          string            `json:"return_url" binding:"required"`
	PaymentType        model.PaymentType `json:"payment_type" binding:"required"`
	Amount             int               `json:"amount" binding:"required"`
	Note               string            `json:"note"`
}

func (s *Server) HandleAdminGenerateCustomerPaymentQRCode(c *gin.Context) {
	req := generateCustomerPaymentQRCode{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGenerateCustomerPaymentQRCode, err)
		return
	}

	originURL, err := s.generateCustomerContractPaymentQRCode(
		req.CustomerContractID,
		req.Amount,
		req.PaymentType,
		req.ReturnURL,
		req.Note,
	)
	if err != nil {
		responseCustomErr(c, ErrCodeGenerateQRCode, err)
		return
	}

	go func() {
		contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
		if err != nil {
			return
		}

		phone, expoToken := contract.Customer.PhoneNumber, s.getExpoToken(contract.Customer.PhoneNumber)
		msg := s.notificationPushService.NewCustomerAdditionalPaymentMsg(contract.ID, expoToken, phone)
		_ = s.notificationPushService.Push(contract.CustomerID, msg)
	}()

	responseSuccess(c, gin.H{"payment_url": originURL.PaymentURL})
}

type generateMultipleCustomerPaymentQRCode struct {
	CustomerPaymentIDs []int  `json:"customer_payment_ids" binding:"required"`
	ReturnURL          string `json:"return_url"`
}

func (s *Server) HandleAdminGenerateMultipleCustomerPayments(c *gin.Context) {
	req := generateMultipleCustomerPaymentQRCode{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGenerateMultipleCustomerPaymentsRequest, err)
		return
	}

	pendingPayments, err := s.store.CustomerPaymentStore.GetPendingBatch(req.CustomerPaymentIDs)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	amt := 0
	ids := make([]int, len(pendingPayments))
	for i, p := range pendingPayments {
		amt += p.Amount
		ids[i] = p.ID
	}

	url, err := s.paymentService.GeneratePaymentURL(ids, amt, time.Now().Format("02150405"), req.ReturnURL)
	if err != nil {
		responseCustomErr(c, ErrCodeGenerateQRCode, err)
		return
	}

	responseSuccess(c, gin.H{"payment_url": url})
}

type generateMultiplePartnerPaymentQRCode struct {
	PartnerPaymentIDs []int  `json:"partner_payment_ids" binding:"required"`
	ReturnURL         string `json:"return_url" binding:"required"`
}

func (s *Server) HandleAdminGenerateMultiplePartnerPayments(c *gin.Context) {
	req := generateMultiplePartnerPaymentQRCode{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGenerateMultiplePartnerPaymentsRequest, err)
		return
	}

	pendingPayments, err := s.store.PartnerPaymentHistoryStore.GetPendingBatch(req.PartnerPaymentIDs)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	amt := 0
	ids := make([]int, len(pendingPayments))
	for i, p := range pendingPayments {
		amt += p.Amount
		ids[i] = p.ID
	}

	txnRef := fmt.Sprintf("%s__%s", PrefixPartnerPayment, time.Now().Format("02150405"))
	url, err := s.paymentService.GeneratePaymentURL(ids, amt, txnRef, req.ReturnURL)
	if err != nil {
		responseCustomErr(c, ErrCodeGenerateQRCode, err)
		return
	}

	responseSuccess(c, gin.H{"payment_url": url})
}

type completeCustomerContractRequest struct {
	CustomerContractID int `json:"customer_contract_id" binding:"required"`
}

func (s *Server) HandleAdminCompleteCustomerContract(c *gin.Context) {
	req := completeCustomerContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCompleteCustomerContractRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.Status != model.CustomerContractStatusRenting {
		responseCustomErr(
			c,
			ErrCodeInvalidCustomerContractStatus,
			errors.New(fmt.Sprintf("invalid customer contract status, require %s, found %s", string(model.CustomerContractStatusRenting), string(contract.Status))),
		)
		return
	}

	// check if any pending payment
	payments, err := s.store.CustomerPaymentStore.GetByCustomerContractID(req.CustomerContractID, model.PaymentStatusPending, 0, 100000)
	if err != nil {
		responseGormErr(c, err)
		return
	}
	if len(payments) > 0 {
		responseCustomErr(c, ErrCodeExistPendingPayments, nil)
		return
	}

	if err := s.store.CustomerContractStore.Update(
		req.CustomerContractID,
		map[string]interface{}{"status": model.CustomerContractStatusCompleted},
	); err != nil {
		responseGormErr(c, err)
		return
	}

	contract, err = s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	_ = s.notificationPushService.Push(contract.CustomerID, s.notificationPushService.NewCompletedCustomerContract(
		contract.ID,
		s.getExpoToken(contract.Customer.PhoneNumber),
		contract.Customer.PhoneNumber,
	))

	responseSuccess(c, contract)
}

type adminUpdateCustomerContractImageStatusRequest struct {
	CustomerContractImageID int                               `json:"customer_contract_image_id" binding:"required"`
	NewStatus               model.CustomerContractImageStatus `json:"new_status" binding:"required"`
}

func (s *Server) HandleAdminUpdateCustomerContractImageStatus(c *gin.Context) {
	req := adminUpdateCustomerContractImageStatusRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUpdateCustomerContractImageStatusRequest, err)
		return
	}

	if err := s.store.CustomerContractImageStore.Update(req.CustomerContractImageID, req.NewStatus); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "update image status successfully"})
}

type adminGetFeedbackRequest struct {
	Pagination
}

func (s *Server) HandleAdminGetFeedbacks(c *gin.Context) {
	req := adminGetFeedbackRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminGetFeedbackRequest, err)
		return
	}

	feedbacks, total, err := s.store.CustomerContractStore.GetFeedbacks(req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"total": total, "feedbacks": feedbacks})
}

type adminUpdateFeedbackStatus struct {
	CustomerContractID int                  `json:"customer_contract_id"`
	NewStatus          model.FeedBackStatus `json:"new_status"`
}

func (s *Server) HandleAdminUpdateFeedbackStatus(c *gin.Context) {
	req := adminUpdateFeedbackStatus{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminUpdateFeedbackStatusRequest, err)
		return
	}

	if err := s.store.CustomerContractStore.Update(req.CustomerContractID, map[string]interface{}{"feedback_status": string(req.NewStatus)}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "update feedback status successfully"})
}

type adminCancelCustomerPayment struct {
	CustomerPaymentID int `json:"customer_payment_id,omitempty" binding:"required"`
}

func (s *Server) HandleAdminCancelCustomerPayment(c *gin.Context) {
	req := adminCancelCustomerPayment{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminCancelCustomerPaymentRequest, err)
		return
	}

	if err := s.store.CustomerPaymentStore.Update(
		req.CustomerPaymentID,
		map[string]interface{}{"status": string(model.PaymentStatusCanceled)},
	); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "cancel payment successfully"})
}

func (s *Server) HandleAdminGetConversations(c *gin.Context) {
	req := Pagination{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminGetConversationsRequest, err)
		return
	}

	conversations, err := s.store.ConversationStore.Get(req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, conversations)
}

type adminGetMessagesRequest struct {
	Pagination
	ConversationID int `form:"conversation_id" binding:"required"`
}

func (s *Server) HandleAdminGetMessages(c *gin.Context) {
	req := adminGetMessagesRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminGetMessagesRequest, err)
		return
	}

	msgs, err := s.store.MessageStore.GetByConversationID(req.ConversationID, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, msgs)
}

type adminUpdateIsReturnCollateralAssetRequest struct {
	CustomerContractID int  `json:"customer_contract_id"`
	NewStatus          bool `json:"new_status"`
}

func (s *Server) HandleAdminUpdateReturnCollateralAsset(c *gin.Context) {
	req := adminUpdateIsReturnCollateralAssetRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUpdateIsReturnCollateralAsset, err)
		return
	}

	if err := s.store.CustomerContractStore.Update(
		req.CustomerContractID,
		map[string]interface{}{"is_return_collateral_asset": req.NewStatus},
	); err != nil {
		responseGormErr(c, err)
		return
	}

	if req.NewStatus {
		contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
		if err != nil {
			responseGormErr(c, err)
			return
		}

		acct, err := s.store.AccountStore.GetByID(contract.CustomerID)
		if err != nil {
			responseGormErr(c, err)
			return
		}

		_ = s.notificationPushService.Push(acct.ID, s.notificationPushService.NewReturnCollateralAssetMsg(
			contract.ID,
			s.getExpoToken(acct.PhoneNumber),
			acct.PhoneNumber,
		))
	}

	responseSuccess(c, gin.H{"status": "update is_return_collateral_asset successfully"})
}

func (s *Server) HandleGetNotificationHistory(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := Pagination{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGetNotificationHistoryRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	notis, err := s.store.NotificationStore.GetByAcctID(acct.ID, req.Offset, req.Limit)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, notis)
}

type AdminMakeMonthlyPartnerPayments struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	ReturnURL string    `json:"return_url" binding:"required"`
}

func (s *Server) HandleAdminMakeMonthlyPartnerPayments(c *gin.Context) {
	req := AdminMakeMonthlyPartnerPayments{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminMakeMonthlyPaymentRequest, err)
		return
	}

	completedContracts, err := s.store.CustomerContractStore.
		GetByStatusEndTimeInRange(req.StartDate, req.EndDate, model.CustomerContractStatusCompleted)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	// partnerID -> list of customer contract IDs
	partnerPayments := make(map[int][]int, 0)

	// partnerID -> needed pay amount
	amounts := make(map[int]int, 0)

	for _, contract := range completedContracts {
		partnerID := contract.Car.PartnerID
		_, existed := partnerPayments[partnerID]
		if !existed {
			partnerPayments[partnerID] = []int{}
		}
		partnerPayments[partnerID] = append(partnerPayments[partnerID], contract.ID)

		_, existed = amounts[partnerID]
		if !existed {
			amounts[partnerID] = 0
		}
		amounts[partnerID] += contract.RentPrice * int(100-contract.Car.PartnerContractRule.RevenueSharingPercent) / 100
	}

	for partnerID, cusContractIds := range partnerPayments {
		history := &model.PartnerPaymentHistory{
			PartnerID: partnerID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Amount:    amounts[partnerID],
			Status:    model.PartnerPaymentHistoryStatusPending,
		}
		if err := s.store.PartnerPaymentHistoryStore.Create(history, cusContractIds); err != nil {
			responseGormErr(c, err)
			return
		}

		if err := s.generatePartnerPaymentQRCode(history.ID, history.Amount, req.ReturnURL); err != nil {
			responseCustomErr(c, ErrCodeGenerateQRCode, err)
			return
		}
	}

	responseSuccess(c, gin.H{"status": "generate monthly partner payment successfully"})
}

const PrefixPartnerPayment = "partner_payment"

func (s *Server) generatePartnerPaymentQRCode(partnerPaymentID, amount int, returnURL string) error {
	txnRef := fmt.Sprintf("%s__%s", PrefixPartnerPayment, time.Now().Format("02150405"))
	url, err := s.paymentService.GeneratePaymentURL([]int{partnerPaymentID}, amount, txnRef, returnURL)
	if err != nil {
		return err
	}

	return s.store.PartnerPaymentHistoryStore.Update(partnerPaymentID, map[string]interface{}{
		"payment_url": url,
	})
}

type AdminGetMonthlyPartnerPayments struct {
	Pagination
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
	Status    string    `form:"status"`
}

func (s *Server) HandleAdminGetMonthlyPartnerPayments(c *gin.Context) {
	req := AdminGetMonthlyPartnerPayments{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminGetMonthlyPartnerPaymentRequest, err)
		return
	}

	status := model.PartnerPaymentHistoryStatusNoFilter
	if len(req.Status) > 0 {
		status = model.PartnerPaymentHistoryStatus(req.Status)
	}

	payments, err := s.store.PartnerPaymentHistoryStore.GetInTimeRange(
		req.StartDate, req.EndDate, status, req.Offset, req.Limit,
	)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, payments)
}

type AdminChangeCarRequest struct {
	CustomerContractID int `json:"customer_contract_id" binding:"required"`
	NewCarID           int `json:"new_car_id" binding:"required"`
}

func (s *Server) HandleAdminChangeCar(c *gin.Context) {
	req := AdminChangeCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAdminChangeCarRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.CarID == req.NewCarID {
		responseCustomErr(c, ErrCodeInvalidAdminChangeCarRequest, err)
		return
	}

	if contract.Status != model.CustomerContractStatusOrdered {
		responseCustomErr(c, ErrCodeInvalidCustomerContractStatus, err)
		return
	}

	isOverlap, err := s.store.CustomerContractStore.IsOverlap(req.NewCarID, contract.StartDate, contract.EndDate)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if isOverlap {
		responseCustomErr(c, ErrCodeOverlapOtherContract, err)
		return
	}

	if err := s.store.CustomerContractStore.Update(req.CustomerContractID, map[string]interface{}{"car_id": req.NewCarID}); err != nil {
		responseGormErr(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.NewCarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	newPartnerMsg := s.notificationPushService.NewReplaceByCar(
		car.ID,
		req.CustomerContractID,
		s.getExpoToken(car.Account.PhoneNumber),
		car.Account.PhoneNumber,
	)
	_ = s.notificationPushService.Push(car.Account.ID, newPartnerMsg)

	oldPartner := contract.Car.Account
	oldPartnerMsg := s.notificationPushService.NewReplaceByOtherCar(
		contract.CarID,
		s.getExpoToken(oldPartner.PhoneNumber),
		oldPartner.PhoneNumber,
	)
	_ = s.notificationPushService.Push(oldPartner.ID, oldPartnerMsg)

	go func() {
		_ = s.RenderCustomerContractPDF(contract.Customer, car, contract)
	}()

	responseSuccess(c, gin.H{"status": "change car successfully"})
}

func (s *Server) HandleAdminGetCustomerContractRule(c *gin.Context) {
	rule, err := s.store.CustomerContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, rule)
}

func (s *Server) HandleAdminGetPartnerContractRule(c *gin.Context) {
	rule, err := s.store.PartnerContractRuleStore.GetLast()
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, rule)
}

type AdminCreateCustomerContractRuleRequest struct {
	InsurancePercent     float64 `json:"insurance_percent" binding:"required"`
	PrepayPercent        float64 `json:"prepay_percent" binding:"required"`
	CollateralCashAmount int     `json:"collateral_cash_amount" binding:"required"`
}

func (s *Server) HandleAdminCreateCustomerContractRule(c *gin.Context) {
	req := AdminCreateCustomerContractRuleRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCreateCustomerContractRuleRequest, err)
		return
	}

	if err := s.store.CustomerContractRuleStore.Create(&model.CustomerContractRule{
		InsurancePercent:     req.InsurancePercent,
		PrepayPercent:        req.PrepayPercent,
		CollateralCashAmount: req.CollateralCashAmount,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "created customer contract rule successfully"})
}

type AdminCreatePartnerContractRuleRequest struct {
	RevenueSharingPercent float64 `json:"revenue_sharing_percent" binding:"required"`
	MaxWarningCount       int     `json:"max_warning_count" binding:"required"`
}

func (s *Server) HandleAdminCreatePartnerContractRule(c *gin.Context) {
	req := AdminCreatePartnerContractRuleRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCreatePartnerContractRuleRequest, err)
		return
	}

	if err := s.store.PartnerContractRuleStore.Create(&model.PartnerContractRule{
		RevenueSharingPercent: req.RevenueSharingPercent,
		MaxWarningCount:       req.MaxWarningCount,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "created partner contract rule successfully"})
}

type AdminUpdateWarningCounter struct {
	CarID           int `json:"car_id"`
	NewWarningCount int `json:"new_warning_count"`
}

func (s *Server) HandleAdminUpdateWarningCount(c *gin.Context) {
	req := AdminUpdateWarningCounter{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUpdateWarningCounterRequest, err)
		return
	}

	if err := s.store.CarStore.Update(req.CarID, map[string]interface{}{
		"warning_count": req.NewWarningCount,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByID(car.PartnerID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if car.WarningCount > car.PartnerContractRule.MaxWarningCount {
		msg := s.notificationPushService.NewInactiveCarMsg(car.ID, s.getExpoToken(car.Account.PhoneNumber), car.Account.PhoneNumber)
		_ = s.notificationPushService.Push(car.Account.ID, msg)

		if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
			"status": model.CarStatusInactive,
		}); err != nil {
			responseGormErr(c, err)
			return
		}

		responseSuccess(c, gin.H{"status": "update warning count successfully. Car status changed to inactive"})
		return
	}

	msg := s.notificationPushService.NewWarningCountMsg(
		car.ID,
		car.WarningCount,
		car.PartnerContractRule.MaxWarningCount,
		acct.PhoneNumber,
		s.getExpoToken(acct.PhoneNumber),
	)

	_ = s.notificationPushService.Push(car.PartnerID, msg)
	responseSuccess(c, gin.H{"status": "update warning count successfully"})
}

type seedPartnerPDFRequest struct {
	CarID int `json:"car_id"`
}

func (s *Server) HandleSeedPartnerPDF(c *gin.Context) {
	req := seedPartnerPDFRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if err := s.RenderPartnerContractPDF(car.Account, car); err != nil {
		responseCustomErr(c, -1, err)
		return
	}
	responseSuccess(c, gin.H{"status": "seed partner contract OK"})
}

type seedCustomerPDFRequest struct {
	CustomerContractID int `json:"customer_contract_id"`
}

func (s *Server) HandleSeedCustomerPDF(c *gin.Context) {
	req := seedCustomerPDFRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if err := s.RenderCustomerContractPDF(contract.Customer, &contract.Car, contract); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	responseSuccess(c, gin.H{"status": "seed customer contract OK"})
}

func (s *Server) getExpoToken(phone string) string {
	expoToken, err := s.redisClient.Get(
		context.Background(),
		fmt.Sprintf("%s__%s", ExpoPushTokenCacheKey, phone),
	).Result()
	if err != nil {
		return ""
	}

	return expoToken
}
