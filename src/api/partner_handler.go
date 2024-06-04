package api

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

func (s *Server) RegisterPartner(c *gin.Context) {
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

	partner := &model.Account{
		RoleID:      model.RoleIDPartner,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    hashedPassword,
		Status:      model.AccountStatusWaitingConfirmEmail,
	}

	if err := s.store.AccountStore.Create(partner); err != nil {
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

type registerCarRequest struct {
	LicensePlate string           `json:"license_plate" binding:"required"`
	CarModelID   int              `json:"car_model_id" binding:"required"`
	Motion       model.Motion     `json:"motion_code" binding:"required"`
	Fuel         model.Fuel       `json:"fuel_code" binding:"required"`
	ParkingLot   model.ParkingLot `json:"parking_lot" binding:"required"`
	PeriodCode   string           `json:"period_code" binding:"required"`
	Description  string           `json:"description"`
}

func (s *Server) HandleRegisterCar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != model.RoleNamePartner {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if acct.Status != model.AccountStatusActive {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("account is not active")))
		return
	}

	req := registerCarRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car := &model.Car{
		PartnerID:    acct.ID,
		CarModelID:   req.CarModelID,
		LicensePlate: req.LicensePlate,
		ParkingLot:   req.ParkingLot,
		Description:  req.Description,
		Fuel:         req.Fuel,
		Motion:       req.Motion,
		Price:        0,
		Status:       model.CarStatusPendingApproval,
	}

	if err := s.store.CarStore.Create(car); err != nil {
		responseInternalServerError(c, err)
		return
	}

	insertedCar, err := s.store.CarStore.GetByID(car.ID)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "register car successfully",
		"car":    insertedCar,
	})
}

type updateRentalPriceRequest struct {
	CarID    int `json:"car_id" binding:"required"`
	NewPrice int `json:"new_price" binding:"required"`
}

func (s *Server) HandleUpdateRentalPrice(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := updateRentalPriceRequest{}
	if err := c.BindJSON(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if car.Account.Email != authPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{
		"price": req.NewPrice,
	}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"car_id":    car.ID,
		"new_price": req.NewPrice,
	})
}

func (s *Server) HandleUploadCarDocuments(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		DocumentCategory model.DocumentCategory  `form:"document_category"`
		CarID            int                     `form:"car_id"`
		Files            []*multipart.FileHeader `form:"files"`
	}{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if car.Account.Email != authPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid ownership")))
		return
	}

	if len(req.Files) > MaxNumberFiles {
		responseError(c, fmt.Errorf("exceed maximum number of files, max %d, has %d", MaxNumberFiles, len(req.Files)))
		return
	}

	for _, f := range req.Files {
		if f.Size > MaxUploadFileSize {
			responseError(c, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, f.Size))
			return
		}

		body, err := f.Open()
		if err != nil {
			responseError(c, err)
			return
		}
		defer body.Close()

		extension := strings.Split(f.Filename, ".")[1]
		key := strings.Join([]string{uuid.NewString(), extension}, ".")
		_, err = s.s3store.Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(s.s3store.Config.Bucket),
			Body:   body,
			Key:    aws.String(key),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err != nil {
			responseError(c, err)
			return
		}

		url := s.s3store.Config.BaseURL + key

		document := &model.Document{
			AccountID: car.Account.ID,
			Url:       url,
			Extension: extension,
			Category:  req.DocumentCategory,
			Status:    model.DocumentStatusActive,
		}

		if err := s.store.CarDocumentStore.Create(car.ID, document); err != nil {
			responseInternalServerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "upload images successfully",
	})
}
