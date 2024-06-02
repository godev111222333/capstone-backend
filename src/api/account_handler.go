package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

type verifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (s *Server) HandleVerifyOTP(c *gin.Context) {
	req := verifyOTPRequest{}
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

type rawLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type rawLoginResponse struct {
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	User                  *accountResponse `json:"user"`
}

type accountResponse struct {
	ID                       int    `json:"id"`
	Role                     string `json:"role"`
	FirstName                string `json:"first_name"`
	LastName                 string `json:"last_name"`
	PhoneNumber              string `json:"phone_number"`
	Email                    string `json:"email"`
	IdentificationCardNumber string `json:"identification_card_number"`
	AvatarUrl                string `json:"avatar_url"`
}

func newAccountResponse(acct *model.Account) *accountResponse {
	return &accountResponse{
		ID:                       acct.ID,
		Role:                     acct.Role.RoleName,
		FirstName:                acct.FirstName,
		LastName:                 acct.LastName,
		PhoneNumber:              acct.PhoneNumber,
		Email:                    acct.Email,
		IdentificationCardNumber: acct.IdentificationCardNumber,
		AvatarUrl:                acct.AvatarURL,
	}
}

func (s *Server) HandleRawLogin(c *gin.Context) {
	req := rawLoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(req.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if acct == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email not found",
		})
	}

	if s.hashVerifier.Compare(acct.Password, req.Password) != nil {
		responseError(c, errors.New("password is not matched"))
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(req.Email, acct.Role.RoleName, s.cfg.AccessTokenDuration)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(req.Email, acct.Role.RoleName, s.cfg.AccessTokenDuration)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.SessionStore.Create(&model.Session{
		ID:           accessTokenPayload.ID,
		Email:        req.Email,
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
		User:                  newAccountResponse(acct),
	})
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
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

	if session.Email != refreshPayload.Email {
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

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.Email, refreshPayload.Role, s.cfg.AccessTokenDuration)
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
