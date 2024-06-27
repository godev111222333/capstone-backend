package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

type verifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

func (s *Server) HandleVerifyOTP(c *gin.Context) {
	req := verifyOTPRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	account, err := s.store.AccountStore.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}

	if account.Status != model.AccountStatusWaitingConfirmEmail {
		responseError(c, errors.New("invalid account status"))
		return
	}

	isValidOTP, err := s.otpService.VerifyOTP(model.OTPTypeRegister, req.PhoneNumber, req.OTP)
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

	if err := s.store.OTPStore.UpdateStatus(req.PhoneNumber, model.OTPTypeRegister, model.OTPStatusVerified); err != nil {
		responseError(c, err)
		return
	}

	if err := s.store.AccountStore.Update(account.ID, map[string]interface{}{
		"status": model.AccountStatusActive,
	}); err != nil {
		responseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "verify account successfully",
	})
}

type rawLoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type rawLoginResponse struct {
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	User                  *accountResponse `json:"user"`
}

type accountResponse struct {
	ID                       int       `json:"id"`
	Role                     string    `json:"role"`
	FirstName                string    `json:"first_name"`
	LastName                 string    `json:"last_name"`
	PhoneNumber              string    `json:"phone_number"`
	Email                    string    `json:"email"`
	IdentificationCardNumber string    `json:"identification_card_number"`
	AvatarUrl                string    `json:"avatar_url"`
	DrivingLicense           string    `json:"driving_license"`
	DateOfBirth              time.Time `json:"date_of_birth"`
	Status                   string    `json:"status"`
	DrivingLicenseImages     []string  `json:"driving_license_images"`
}

func (s *Server) newAccountResponse(acct *model.Account) *accountResponse {
	resp := &accountResponse{
		ID:                       acct.ID,
		Role:                     acct.Role.RoleName,
		FirstName:                acct.FirstName,
		LastName:                 acct.LastName,
		PhoneNumber:              acct.PhoneNumber,
		Email:                    acct.Email,
		IdentificationCardNumber: acct.IdentificationCardNumber,
		AvatarUrl:                acct.AvatarURL,
		DrivingLicense:           acct.DrivingLicense,
		DateOfBirth:              acct.DateOfBirth,
		Status:                   string(acct.Status),
	}
	if acct.RoleID == model.RoleIDCustomer {
		drivingLicenseImages, err := s.store.DocumentStore.GetByCategory(acct.ID, model.DocumentCategoryDrivingLicense, 2)
		if err != nil {
			fmt.Println(err)
		} else {
			if len(drivingLicenseImages) >= 2 {
				resp.DrivingLicenseImages = []string{drivingLicenseImages[0].Url, drivingLicenseImages[1].Url}
			}
		}
	}
	return resp
}

func (s *Server) HandleRawLogin(c *gin.Context) {
	req := rawLoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if acct == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email not found",
		})
		return
	}

	if s.hashVerifier.Compare(acct.Password, req.Password) != nil {
		responseError(c, errors.New("invalid email or password"))
		return
	}

	if acct.Status != model.AccountStatusActive {
		responseError(c, errors.New("active is not active"))
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(req.PhoneNumber, acct.Role.RoleName, s.cfg.AccessTokenDuration)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(req.PhoneNumber, acct.Role.RoleName, s.cfg.RefreshTokenDuration)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.SessionStore.Create(&model.Session{
		ID:           refreshTokenPayload.ID,
		PhoneNumber:  req.PhoneNumber,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIP:     c.ClientIP(),
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rawLoginResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User:                  s.newAccountResponse(acct),
	})
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) HandleRenewAccessToken(c *gin.Context) {
	req := &renewAccessTokenRequest{}
	if err := c.BindJSON(req); err != nil {
		responseError(c, err)
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		responseError(c, err)
		return
	}

	session, err := s.store.SessionStore.GetSession(refreshPayload.ID)
	if err != nil {
		responseError(c, err)
		return
	}

	if session.PhoneNumber != refreshPayload.PhoneNumber {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("incorrect session email")))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("mismatch session token")))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("expired session")))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.PhoneNumber, refreshPayload.Role, s.cfg.AccessTokenDuration)
	if err != nil {
		responseError(c, err)
		return
	}
	resp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	c.JSON(http.StatusOK, resp)
}

type updateProfileRequest struct {
	FirstName                string    `json:"first_name"`
	LastName                 string    `json:"last_name"`
	Email                    string    `json:"email"`
	DateOfBirth              time.Time `json:"date_of_birth"`
	IdentificationCardNumber string    `json:"identification_card_number" binding:"id_card"`
	DrivingLicense           string    `json:"driving_license" binding:"driving_license"`
	Password                 string    `json:"password"`
}

func (s *Server) HandleGetProfile(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, s.newAccountResponse(acct))
}

func (s *Server) HandleUpdateProfile(c *gin.Context) {
	req := updateProfileRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseError(c, err)
		return
	}
	if acct == nil || acct.PhoneNumber != authPayload.PhoneNumber {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("mismatch token or account not found")))
		return
	}

	updateParams := map[string]interface{}{
		"first_name":                 req.FirstName,
		"last_name":                  req.LastName,
		"email":                      req.Email,
		"identification_card_number": req.IdentificationCardNumber,
		"driving_license":            req.DrivingLicense,
		"date_of_birth":              req.DateOfBirth,
	}
	if len(req.Password) > 0 {
		h, err := s.hashVerifier.Hash(req.Password)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}
		updateParams["password"] = h
	}

	if err := s.store.AccountStore.Update(acct.ID, updateParams); err != nil {
		responseInternalServerError(c, err)
		return
	}

	updatedAcct, err := s.store.AccountStore.GetByID(acct.ID)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, accountResponse{
		ID:                       updatedAcct.ID,
		Role:                     updatedAcct.Role.RoleName,
		FirstName:                updatedAcct.FirstName,
		LastName:                 updatedAcct.LastName,
		PhoneNumber:              updatedAcct.PhoneNumber,
		Email:                    updatedAcct.Email,
		IdentificationCardNumber: updatedAcct.IdentificationCardNumber,
		AvatarUrl:                updatedAcct.AvatarURL,
		DateOfBirth:              updatedAcct.DateOfBirth,
	})
}

type updatePaymentInfoRequest struct {
	BankNumber string `json:"bank_number" binding:"required"`
	BankOwner  string `json:"bank_owner" binding:"required"`
	BankName   string `json:"bank_name" binding:"required"`
}

func (s *Server) HandleUpdatePaymentInformation(c *gin.Context) {
	req := updatePaymentInfoRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.PaymentInformationStore.Update(acct.ID, map[string]interface{}{
		"bank_number": req.BankNumber,
		"bank_owner":  req.BankOwner,
		"bank_name":   req.BankName,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated payment information successfully",
	})
}

func (s *Server) HandleGetPaymentInformation(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	p, err := s.store.PaymentInformationStore.GetByAcctID(acct.ID)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (s *Server) HandleUpdateQRCodeImage(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	req := struct {
		File *multipart.FileHeader `form:"file"`
	}{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	if req.File.Size > MaxUploadFileSize {
		responseError(c, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, req.File.Size))
		return
	}

	file, err := req.File.Open()
	if err != nil {
		responseError(c, err)
		return
	}
	defer file.Close()

	doc, err := s.uploadDocument(file, acct.ID, req.File.Filename, model.DocumentCategoryPersonalQRCodeImage)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.DocumentStore.Create(doc); err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.PaymentInformationStore.Update(acct.ID, map[string]interface{}{
		"qr_code_url": doc.Url,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code_url": doc.Url,
	})
}
