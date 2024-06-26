package api

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

const layoutDateMonthYear = "01/02/2006"

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

		// Create new partner contract record
		contract := &model.PartnerContract{
			CarID:     req.CarID,
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, car.Period, 0),
			Status:    model.PartnerContractStatusWaitingForAgreement,
		}
		if err := s.store.PartnerContractStore.Create(contract); err != nil {
			responseInternalServerError(c, err)
			return
		}

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

		contract, err := s.store.PartnerContractStore.GetByCarID(car.ID)
		if err != nil {
			responseGormErr(c, err)
			return
		}

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

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"status": newStatus,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": fmt.Sprintf("%s car successfully", req.Action)})
}

func (s *Server) RenderPartnerContractPDF(partner *model.Account, car *model.Car) error {
	now := time.Now()
	year, month, date := now.Date()

	contract, err := s.store.PartnerContractStore.GetByCarID(car.ID)
	if err != nil {
		return err
	}
	startYear, startMonth, startDate := contract.StartDate.Date()
	endYear, endMonth, endDate := contract.EndDate.Date()

	docUUID, err := s.pdfService.Render(RenderTypePartner, map[string]string{
		"now_date":              strconv.Itoa(date),
		"now_month":             strconv.Itoa(int(month)),
		"now_year":              strconv.Itoa(year),
		"partner_fullname":      partner.LastName + " " + partner.FirstName,
		"partner_date_of_birth": partner.DateOfBirth.Format(layoutDateMonthYear),
		"partner_id_card":       partner.IdentificationCardNumber,
		"partner_address":       "",
		"brand_model":           car.CarModel.Brand + " " + car.CarModel.Model,
		"license_plate":         car.LicensePlate,
		"number_of_seats":       strconv.Itoa(car.CarModel.NumberOfSeats),
		"car_year":              strconv.Itoa(car.CarModel.Year),
		"price":                 strconv.Itoa(car.Price),
		"period":                strconv.Itoa(car.Period),
		"period_start_date":     strconv.Itoa(startDate),
		"period_start_month":    strconv.Itoa(int(startMonth)),
		"period_start_year":     strconv.Itoa(startYear),
		"period_end_date":       strconv.Itoa(endDate),
		"period_end_month":      strconv.Itoa(int(endMonth)),
		"period_end_year":       strconv.Itoa(endYear),
	})
	if err != nil {
		fmt.Printf("error when rendering partner contract %v\n", err)
		return err
	}

	if err := s.store.PartnerContractStore.Update(
		contract.ID,
		map[string]interface{}{"url": s.fromUUIDToURL(docUUID, model.ExtensionPDF)},
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
	nowYear, nowMonth, nowDate := time.Now().Date()
	startDate, endDate := contract.StartDate, contract.EndDate
	startHour, startDay, startMonth, startYear := startDate.Hour(), startDate.Day(), int(startDate.Month()), startDate.Year()
	endHour, endDay, endMonth, endYear := endDate.Hour(), endDate.Day(), int(endDate.Month()), endDate.Year()
	docUUID, err := s.pdfService.Render(RenderTypeCustomer, map[string]string{
		"now_date":               strconv.Itoa(nowDate),
		"now_month":              strconv.Itoa(int(nowMonth)),
		"now_year":               strconv.Itoa(nowYear),
		"customer_fullname":      customer.LastName + " " + customer.FirstName,
		"customer_date_of_birth": customer.DateOfBirth.Format(layoutDateMonthYear),
		"customer_id_card":       customer.IdentificationCardNumber,
		"customer_address":       "",
		"brand_model":            car.CarModel.Brand + " " + car.CarModel.Model,
		"license_plate":          car.LicensePlate,
		"number_of_seats":        strconv.Itoa(car.CarModel.NumberOfSeats),
		"car_year":               strconv.Itoa(car.CarModel.Year),
		"price":                  strconv.Itoa(contract.RentPrice),
		"start_hour":             strconv.Itoa(startHour),
		"start_date":             strconv.Itoa(startDay),
		"start_month":            strconv.Itoa(startMonth),
		"start_year":             strconv.Itoa(startYear),
		"end_hour":               strconv.Itoa(endHour),
		"end_date":               strconv.Itoa(endDay),
		"end_month":              strconv.Itoa(endMonth),
		"end_year":               strconv.Itoa(endYear),
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

		newStatus = string(model.CustomerContractStatusRenting)
	}

	if err := s.store.CustomerContractStore.Update(contract.ID, map[string]interface{}{"status": newStatus}); err != nil {
		responseGormErr(c, err)
		return
	}

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

type adminSetAccountStatusRequest struct {
	AccountID int                 `json:"account_id"`
	Status    model.AccountStatus `json:"status"`
}

func (s *Server) HandleAdminSetAccountStatus(c *gin.Context) {
	req := adminSetAccountStatusRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidSetAccountStatusRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByID(req.AccountID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if acct.RoleID == model.RoleIDPartner {
		carStatus := model.CarStatusActive
		if req.Status == model.AccountStatusInactive {
			carStatus = model.CarStatusInactive
		}

		if err := s.store.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.store.CarStore.UpdateByPartnerID(tx, acct.ID, map[string]interface{}{"status": string(carStatus)}); err != nil {
				return err
			}

			return s.store.AccountStore.UpdateTx(tx, req.AccountID, map[string]interface{}{"status": string(req.Status)})
		}); err != nil {
			responseInternalServerError(c, err)
			return
		}
	} else if acct.RoleID == model.RoleIDCustomer {
		if err := s.store.AccountStore.Update(req.AccountID, map[string]interface{}{"status": string(req.Status)}); err != nil {
			responseInternalServerError(c, err)
			return
		}
	}

	responseSuccess(c, acct)
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

func newCustomerPaymentResponse(p *model.CustomerPayment) *customerPaymentResponse {
	r := &customerPaymentResponse{CustomerPayment: p}
	payer := model.RoleNameCustomer
	if p.PaymentType == model.PaymentTypeReturnCollateralCash {
		payer = model.RoleNameAdmin
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

type adminCreateCustomerPaymentRequest struct {
	CustomerContractID int               `json:"customer_contract_id" binding:"required"`
	PaymentType        model.PaymentType `json:"payment_type" binding:"required"`
	Amount             int               `json:"amount" binding:"required"`
	Note               string            `json:"note"`
}

func (s *Server) HandleAdminCreateCustomerPayment(c *gin.Context) {
	req := adminCreateCustomerPaymentRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidCreateCustomerPaymentRequest, err)
		return
	}

	payment := &model.CustomerPayment{
		CustomerContractID: req.CustomerContractID,
		PaymentType:        req.PaymentType,
		Amount:             req.Amount,
		Note:               req.Note,
		Status:             model.PaymentStatusPending,
	}

	if err := s.store.CustomerPaymentStore.Create(payment); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, payment)
}

type generateCustomerPaymentQRCode struct {
	CustomerPaymentID int    `json:"customer_payment_id" binding:"required"`
	ReturnURL         string `json:"return_url" binding:"required"`
}

func (s *Server) HandleAdminGenerateCustomerPaymentQRCode(c *gin.Context) {
	req := generateCustomerPaymentQRCode{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidGenerateCustomerPaymentQRCode, err)
		return
	}

	customerPayment, err := s.store.CustomerPaymentStore.GetByID(req.CustomerPaymentID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	contract := customerPayment.CustomerContract
	originURL, err := s.generateCustomerContractPaymentQRCode(
		contract.ID,
		customerPayment.Amount,
		customerPayment.PaymentType,
		req.ReturnURL,
	)
	if err != nil {
		responseCustomErr(c, ErrCodeGenerateQRCode, err)
		return
	}

	responseSuccess(c, gin.H{"payment_url": originURL})
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
