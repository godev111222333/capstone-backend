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

const FakePDF = "https://rentalcar-capstone.s3.ap-southeast-2.amazonaws.com/vib_contract.pdf"

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

	c.JSON(http.StatusOK, cars)
}

func (s *Server) HandleAdminGetCarDetails(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNameAdmin {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

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
	Max4Seats  int `json:"max_4_seats"`
	Max7Seats  int `json:"max_7_seats"`
	Max15Seats int `json:"max_15_seats"`
	Total      int `json:"total"`
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

	c.JSON(http.StatusOK, getGarageConfigResponse{
		Max4Seats:  configs[model.GarageConfigTypeMax4Seats],
		Max7Seats:  configs[model.GarageConfigTypeMax7Seats],
		Max15Seats: configs[model.GarageConfigTypeMax15Seats],
		Total: configs[model.GarageConfigTypeMax4Seats] +
			configs[model.GarageConfigTypeMax7Seats] +
			configs[model.GarageConfigTypeMax15Seats],
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
			Url:       FakePDF,
		}
		if err := s.store.PartnerContractStore.Create(contract); err != nil {
			responseInternalServerError(c, err)
			return
		}
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
