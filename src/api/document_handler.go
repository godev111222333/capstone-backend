package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

const (
	MaxUploadFileSize            = 5 * 1 << 20
	MaxNumberFiles               = 5
	MaxNumberDrivingLicenseFiles = 2
)

func (s *Server) HandleUploadAvatar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		File *multipart.FileHeader `form:"file"`
	}{}

	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	header := req.File
	if header.Size > MaxUploadFileSize {
		responseError(c, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, header.Size))
		return
	}

	body, err := header.Open()
	if err != nil {
		responseError(c, err)
		return
	}
	defer body.Close()

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	doc, err := s.uploadDocument(body, acct.ID, header.Filename, model.DocumentCategoryAvatarImage)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.AccountStore.Update(acct.ID, map[string]interface{}{
		"avatar_url": doc.Url,
	}); err != nil {
		responseError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "upload avatar successfully",
		"url":    doc.Url,
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

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseError(c, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseError(c, err)
		return
	}

	if !strings.Contains(string(car.Status), string(model.CarStatusPendingApplication)) {
		responseError(c, errors.New("invalid car state"))
		return
	}

	if (req.DocumentCategory == model.DocumentCategoryCarImages && car.Status != model.CarStatusPendingApplicationPendingCarImages) ||
		req.DocumentCategory == model.DocumentCategoryCaveat && car.Status != model.CarStatusPendingApplicationPendingCarCaveat {
		responseError(c, errors.New("invalid document category with current car state"))
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

		document, err := s.uploadDocument(body, acct.ID, f.Filename, req.DocumentCategory)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}

		if err := s.store.CarDocumentStore.Create(car.ID, document); err != nil {
			responseInternalServerError(c, err)
			return
		}
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{"status": model.MoveNextCarState(car.Status)}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "upload images successfully",
	})
}

func (s *Server) HandleUploadDrivingLicenseImages(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		Files []*multipart.FileHeader `form:"files"`
	}{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if len(req.Files) > MaxNumberDrivingLicenseFiles {
		responseError(c, fmt.Errorf("exceed maximum number of files, max %d, has %d", MaxNumberDrivingLicenseFiles, len(req.Files)))
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

		if _, err := s.uploadDocument(body, acct.ID, f.Filename, model.DocumentCategoryDrivingLicense); err != nil {
			responseInternalServerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "upload driving license images successfully"})
}

func (s *Server) HandleGetDrivingLicenseImages(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByEmail(authPayload.Email)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	docs, err := s.store.DocumentStore.GetByCategory(acct.ID, model.DocumentCategoryDrivingLicense, 2)
	if err != nil {
		responseError(c, err)
		return
	}

	urls := make([]string, len(docs))
	for i, d := range docs {
		urls[i] = d.Url
	}

	c.JSON(http.StatusOK, urls)
}

func (s *Server) uploadDocument(
	reader io.Reader,
	accountID int,
	fileName string,
	category model.DocumentCategory,
) (*model.Document, error) {
	extension := strings.Split(fileName, ".")[1]
	key := strings.Join([]string{uuid.NewString(), extension}, ".")

	_, err := s.s3store.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.s3store.Config.Bucket),
		Body:   reader,
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, err
	}

	url := s.s3store.Config.BaseURL + key
	doc := &model.Document{
		AccountID: accountID,
		Url:       url,
		Extension: extension,
		Category:  category,
		Status:    model.DocumentStatusActive,
	}

	return doc, s.store.DB.Create(doc).Error
}
