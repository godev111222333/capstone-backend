package api

import (
	"errors"
	"fmt"
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
	CarStatus string `form:"car_status"`
}

func (s *Server) HandleAdminGetCars(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	req := getCarsRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	status := model.CarStatusNoFilter
	if len(req.CarStatus) > 0 {
		status = model.CarStatus(req.CarStatus)
	}

	cars, err := s.store.CarStore.GetAll(req.Offset, req.Limit, status)
	if err != nil {
		responseError(c, err)
		return
	}

	total, err := s.store.CarStore.CountByStatus(status)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cars":  cars,
		"total": total,
	})
}

func (s *Server) HandleGetCarDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(id)
	if err != nil {
		responseError(c, err)
		return
	}

	resp, err := s.newCarResponse(car)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
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
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	configs, err := s.store.GarageConfigStore.Get()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	countCurrentSeats := func(seatType int) int {
		counter, err := s.store.CarStore.CountBySeats(seatType)
		if err != nil {
			responseInternalServerError(c, err)
			return -1
		}

		return counter
	}

	cur4Seats := countCurrentSeats(4)
	cur7Seats := countCurrentSeats(7)
	cur15Seats := countCurrentSeats(15)

	c.JSON(http.StatusOK, getGarageConfigResponse{
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
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	req := updateGarageConfigRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	checkValidOfSeats := func(seatType, maxSeat int) bool {
		counter, err := s.store.CarStore.CountBySeats(seatType)
		if err != nil {
			responseInternalServerError(c, err)
			return false
		}

		if counter > maxSeat {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New(fmt.Sprintf("invalid type %d seat. Must at least %d", seatType, counter))))
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

	c.JSON(http.StatusOK, gin.H{"status": "update garage configs successfully"})
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
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	req := adminApproveOrRejectRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	newStatus := string(model.CarStatusRejected)
	if req.Action == ApplicationActionApproveRegister {
		if car.Status != model.CarStatusPendingApproval {
			c.JSON(http.StatusBadRequest, errorResponse(
				fmt.Errorf(
					"invalid car status, require %s,"+
						" found %s",
					string(model.CarStatusPendingApproval),
					string(car.Status),
				),
			))
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
			partner, err := s.store.AccountStore.GetByEmail(car.Account.Email)
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
			c.JSON(http.StatusBadRequest, errorResponse(
				fmt.Errorf(
					"invalid car status, require %s, found %s",
					string(model.CarStatusWaitingDelivery),
					string(car.Status),
				),
			))
			return
		}

		contract, err := s.store.PartnerContractStore.GetByCarID(car.ID)
		if err != nil {
			responseError(c, err)
			return
		}

		if contract.Status != model.PartnerContractStatusAgreed {
			responseError(c, errors.New("partner must agree the contract first"))
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

	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("%s car successfully", req.Action)})
}

func (s *Server) RenderPartnerContractPDF(partner *model.Account, car *model.Car) error {
	now := time.Now()
	year, month, date := now.Date()

	contract, err := s.store.PartnerContractStore.GetByCarID(car.ID)
	if err != nil {
		fmt.Println(err)
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
		"price":                 strconv.Itoa(car.Price * car.Period * 30),
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
}

func (s *Server) HandleAdminGetCustomerContracts(c *gin.Context) {
	req := adminGetContractRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	status := model.CustomerContractStatusNoFilter
	if reqStatus := req.CustomerContractStatus; len(reqStatus) > 0 {
		status = model.CustomerContractStatus(reqStatus)
	}

	contracts, err := s.store.CustomerContractStore.GetByStatus(status, req.Offset, req.Limit)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	total, err := s.store.CustomerContractStore.CountByStatus(status)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"contracts": contracts, "total": total})
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
		responseError(c, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseError(c, err)
		return
	}

	newStatus := string(model.CustomerContractStatusCancel)
	if req.Action == CustomerContractActionApprove {
		if contract.Status != model.CustomerContractStatusOrdered {
			responseError(c, errors.New(
				fmt.Sprintf("invalid customer contract status. expect %s, found %s",
					string(model.CustomerContractStatusOrdered), string(contract.Status))))
			return
		}

		newStatus = string(model.CustomerContractStatusRenting)
	}

	if err := s.store.CustomerContractStore.Update(contract.ID, map[string]interface{}{"status": newStatus}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "approve/reject customer contract successfully"})
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
		responseError(c, err)
		return
	}

	status := model.AccountStatusNoFilter
	if len(req.Status) > 0 {
		status = model.AccountStatus(req.Status)
	}

	accounts, err := s.store.AccountStore.Get(status, req.Role, req.SearchParam, req.Offset, req.Limit)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	respAccts := make([]*accountResponse, len(accounts))
	for i, acct := range accounts {
		respAccts[i] = s.newAccountResponse(acct)
	}

	c.JSON(http.StatusOK, respAccts)
}

type adminSetAccountStatusRequest struct {
	AccountID int                 `json:"account_id"`
	Status    model.AccountStatus `json:"status"`
}

func (s *Server) HandleAdminSetAccountStatus(c *gin.Context) {
	req := adminSetAccountStatusRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	if err := s.store.AccountStore.Update(req.AccountID, map[string]interface{}{"status": string(req.Status)}); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "update account status successfully"})
}

func (s *Server) HandleAdminGetAccountDetail(c *gin.Context) {
	id := c.Param("account_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByID(idInt)
	if err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, s.newAccountResponse(acct))
}
