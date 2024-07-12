package api

import (
	"fmt"
	"mime/multipart"
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
		responseCustomErr(c, ErrCodeInvalidVerifyOTPRequest, err)
		return
	}

	account, err := s.store.AccountStore.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if account.Status != model.AccountStatusWaitingConfirmPhoneNumber {
		responseCustomErr(c, ErrCodeInvalidAccountStatus, nil)
		return
	}

	isValidOTP, err := s.otpService.VerifyOTP(model.OTPTypeRegister, req.PhoneNumber, req.OTP)
	if err != nil {
		responseCustomErr(c, ErrCodeCacheError, err)
		return
	}

	if !isValidOTP {
		responseCustomErr(c, ErrCodeInvalidOTP, nil)
		return
	}

	if err := s.store.AccountStore.Update(account.ID, map[string]interface{}{
		"status": model.AccountStatusActive,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, account)
}

type rawLoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type rawLoginResponse struct {
	AccessToken          string           `json:"access_token"`
	AccessTokenExpiresAt time.Time        `json:"access_token_expires_at"`
	User                 *accountResponse `json:"user"`
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
		drivingLicenseImages, err := s.store.DrivingLicenseImageStore.Get(acct.ID, model.DrivingLicenseImageStatusActive, 2)
		if err != nil {
			fmt.Println(err)
		} else {
			if len(drivingLicenseImages) >= 2 {
				resp.DrivingLicenseImages = []string{drivingLicenseImages[0].URL, drivingLicenseImages[1].URL}
			}
		}
	}
	return resp
}

func (s *Server) HandleRawLogin(c *gin.Context) {
	req := rawLoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidLoginRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(req.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if s.hashVerifier.Compare(acct.Password, req.Password) != nil {
		responseCustomErr(c, ErrCodeWrongPhoneNumberOrPassword, nil)
		return
	}

	if acct.Status != model.AccountStatusActive {
		responseCustomErr(c, ErrCodeAccountNotActive, nil)
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(req.PhoneNumber, acct.Role.RoleName, s.cfg.AccessTokenDuration)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, rawLoginResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
		User:                 s.newAccountResponse(acct),
	})
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

	responseSuccess(c, s.newAccountResponse(acct))
}

func (s *Server) HandleUpdateProfile(c *gin.Context) {
	req := updateProfileRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUpdateProfileRequest, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}
	if acct.PhoneNumber != authPayload.PhoneNumber {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
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

	responseSuccess(c, accountResponse{
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
		responseCustomErr(c, ErrCodeUpdatePaymentInfoRequest, nil)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if err := s.store.AccountStore.Update(acct.ID, map[string]interface{}{
		"bank_number": req.BankNumber,
		"bank_owner":  req.BankOwner,
		"bank_name":   req.BankName,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, acct)
}

func (s *Server) HandleGetPaymentInformation(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"bank_number": acct.BankNumber,
		"bank_owner":  acct.BankOwner,
		"bank_name":   acct.BankName,
		"qr_code_url": acct.QRCodeURL,
	})
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
		responseCustomErr(c, ErrCodeInvalidUploadDocumentRequest, nil)
		return
	}

	if req.File.Size > MaxUploadFileSize {
		responseCustomErr(c, ErrCodeInvalidFileSize, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, req.File.Size))
		return
	}

	file, err := req.File.Open()
	if err != nil {
		responseCustomErr(c, ErrCodeReadingDocumentRequest, err)
		return
	}
	defer file.Close()

	url, err := s.uploadDocument(file, req.File.Filename)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.AccountStore.Update(acct.ID, map[string]interface{}{
		"qr_code_url": url,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"qr_code_url": url,
	})
}

type RegisterExpoPushTokenRequest struct {
	ExpoPushToken string `json:"expo_push_token"`
}

func (s *Server) HandleRegisterExpoPushToken(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := RegisterExpoPushTokenRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidRegisterExpoPushTokenRequest, err)
		return
	}

	s.expoPushTokens.Store(authPayload.PhoneNumber, req.ExpoPushToken)
	responseSuccess(c, gin.H{"status": "register expo push token successfully"})
}
