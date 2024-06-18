package api

import (
	"errors"
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

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.Email); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register customer successfully. please confirm OTP sent to your email",
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

	if car.Status != model.CarStatusActive {
		responseError(c, errors.New("invalid car status. require active"))
		return
	}

	rentPrice := car.Price * (req.EndDate.Day() - req.StartDate.Day())
	insuranceAmount := rentPrice / 10
	contract := &model.CustomerContract{
		CustomerID:              customer.ID,
		CarID:                   req.CarID,
		RentPrice:               rentPrice,
		StartDate:               req.StartDate,
		EndDate:                 req.EndDate,
		Status:                  model.CustomerContractStatusWaitingContractAgreement,
		InsuranceAmount:         insuranceAmount,
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

type customerSignContractRequest struct {
	CustomerContractID int `json:"customer_contract_id"`
}

func (s *Server) HandleCustomerAgreeContract(c *gin.Context) {
	req := customerSignContractRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseError(c, err)
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

	c.JSON(http.StatusOK, gin.H{"status": "agree contract successfully"})
}

type customerGetContractsRequest struct {
	*Pagination
	ContractStatus string `form:"contract_status"`
}

func (s *Server) HandleCustomerGetContracts(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := customerGetContractsRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
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

func (s *Server) HandleCustomerGetContractDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
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

	if contract.CustomerID != acct.ID {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	c.JSON(http.StatusOK, contract)
}
