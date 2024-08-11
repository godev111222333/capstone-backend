package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

type AppraisingCarAction string

const (
	AppraisingCarActionApprove AppraisingCarAction = "approve"
	AppraisingCarActionReject  AppraisingCarAction = "reject"
)

type techAppraisingCar struct {
	CustomerContractID int                 `json:"customer_contract_id"`
	Action             AppraisingCarAction `json:"action"`
}

func (s *Server) HandleTechnicianAppraisingCarOfCusContract(c *gin.Context) {
	req := techAppraisingCar{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidAppraisingCarRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.Status != model.CustomerContractStatusAppraisingCar {
		responseCustomErr(c,
			ErrCodeInvalidCustomerContractStatus,
			fmt.Errorf("customer contract status require %s, found %s", model.CustomerContractStatusAppraisingCar, contract.Status))
		return
	}

	nextStatus := model.CustomerContractStatusRenting
	if req.Action == AppraisingCarActionReject {
		nextStatus = model.CustomerContractStatusAppraisingCarFailed
	}

	go func() {
		if nextStatus == model.CustomerContractStatusRenting {
			phone, expoToken := contract.Customer.PhoneNumber, s.getExpoToken(contract.Customer.PhoneNumber)
			msg := s.notificationPushService.NewApproveRentingCarRequestMsg(contract.ID, expoToken, phone)
			_ = s.notificationPushService.Push(contract.CustomerID, msg)
		}
	}()

	if err := s.store.CustomerContractStore.Update(req.CustomerContractID, map[string]interface{}{
		"status": string(nextStatus),
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "appraising car successfully"})
}
