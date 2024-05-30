package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

func (s *Server) RegisterPartner(c *gin.Context) {
	req := struct {
		FirstName                string `json:"first_name"`
		LastName                 string `json:"last_name"`
		PhoneNumber              string `json:"phone_number"`
		Email                    string `json:"email"`
		IdentificationCardNumber string `json:"identification_card_number"`
		Password                 string `json:"password"`
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

	partner := &model.Partner{
		Account: model.Account{
			RoleID:                   model.RoleIDPartner,
			FirstName:                req.FirstName,
			LastName:                 req.LastName,
			PhoneNumber:              req.PhoneNumber,
			Email:                    req.Email,
			IdentificationCardNumber: req.IdentificationCardNumber,
			Password:                 hashedPassword,
			Status:                   model.AccountStatusWaitingConfirmEmail,
		},
	}

	if err := s.store.PartnerStore.Create(partner); err != nil {
		responseError(c, err)
		return
	}

	if err := s.otpService.SendOTP(model.OTPTypeRegister, req.Email); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register partner successfully. please confirm OTP sent to your email",
	})
}
