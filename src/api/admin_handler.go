package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

type getCarsRequest struct {
	*Pagination
	CarStatus string `form:"car_status"`
}

func (s *Server) HandleAdminGetCars(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	auth, err := s.checkValidRole(authPayload, model.RoleIDAdmin)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if !auth {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid permission")))
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

func (s *Server) checkValidRole(authPayload *token.Payload, role model.RoleID) (bool, error) {
	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		return false, err
	}

	return acct.Role.ID == int(role), nil
}
