package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

func (s *Server) HandleVerifyOTP(c *gin.Context) {
	req := struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	account, err := s.store.AccountStore.GetByEmail(req.Email)
	if err != nil {
		responseError(c, err)
		return
	}

	if account.Status != model.AccountStatusWaitingConfirmEmail {
		responseError(c, errors.New("invalid account status"))
		return
	}

	isValidOTP, err := s.otpService.VerifyOTP(model.OTPTypeRegister, req.Email, req.OTP)
	if err != nil {
		responseError(c, err)
		return
	}

	if !isValidOTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "invalid OTP or OTP was expired",
		})
		return
	}

	if err := s.store.OTPStore.UpdateStatus(req.Email, model.OTPTypeRegister, model.OTPStatusVerified); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "verify account successfully",
	})
}
