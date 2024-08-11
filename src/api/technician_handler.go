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

	if contract.Status != model.CustomerContractStatusOrdered {
		responseCustomErr(c,
			ErrCodeInvalidCustomerContractStatus,
			fmt.Errorf("customer contract status require %s, found %s", model.CustomerContractStatusOrdered, contract.Status))
		return
	}

	nextStatus := model.CustomerContractStatusAppraisingCarApproved
	if req.Action == AppraisingCarActionReject {
		nextStatus = model.CustomerContractStatusAppraisingCarRejected
	}

	if err := s.store.CustomerContractStore.Update(req.CustomerContractID, map[string]interface{}{
		"status": string(nextStatus),
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "appraising car successfully"})
}
