package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	MaxUploadFileSize             = 5 * 1 << 20
	MaxNumberFiles                = 5
	MaxNumberDrivingLicenseFiles  = 2
	MaxNumberCollateralAssetFiles = 6
	MaxNumberReceivingCarImages   = 6
)

func (s *Server) HandleUploadAvatar(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		File *multipart.FileHeader `form:"file"`
	}{}

	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUploadDocumentRequest, err)
		return
	}

	header := req.File
	if header.Size > MaxUploadFileSize {
		responseCustomErr(c, ErrCodeInvalidFileSize, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, header.Size))
		return
	}

	body, err := header.Open()
	if err != nil {
		responseCustomErr(c, ErrCodeReadingDocumentRequest, err)
		return
	}
	defer body.Close()

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	url, err := s.uploadDocument(body, header.Filename)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.AccountStore.Update(acct.ID, map[string]interface{}{
		"avatar_url": url,
	}); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"url": url,
	})
}

func (s *Server) HandleUploadCarDocuments(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		CarImageCategory model.CarImageCategory  `form:"document_category"`
		CarID            int                     `form:"car_id"`
		Files            []*multipart.FileHeader `form:"files"`
	}{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUploadDocumentRequest, err)
		return
	}

	car, err := s.store.CarStore.GetByID(req.CarID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if !strings.Contains(string(car.Status), string(model.CarStatusPendingApplication)) {
		responseCustomErr(c, ErrCodeInvalidCarStatus, errors.New("invalid car state"))
		return
	}

	if (req.CarImageCategory == model.CarImageCategoryImages && car.Status != model.CarStatusPendingApplicationPendingCarImages) ||
		req.CarImageCategory == model.CarImageCategoryCaveat && car.Status != model.CarStatusPendingApplicationPendingCarCaveat {
		responseCustomErr(c, ErrCodeInvalidDocumentCategory, errors.New("invalid document category with current car state"))
		return
	}

	if car.Account.PhoneNumber != authPayload.PhoneNumber {
		responseCustomErr(c, ErrCodeInvalidOwnership, nil)
		return
	}

	if len(req.Files) > MaxNumberFiles {
		responseCustomErr(c, ErrCodeInvalidNumberOfFiles, fmt.Errorf("exceed maximum number of files, max %d, has %d", MaxNumberFiles, len(req.Files)))
		return
	}

	images := make([]*model.CarImage, 0)
	for _, f := range req.Files {
		if f.Size > MaxUploadFileSize {
			responseCustomErr(c, ErrCodeInvalidFileSize, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, f.Size))
			return
		}

		body, err := f.Open()
		if err != nil {
			responseCustomErr(c, ErrCodeReadingDocumentRequest, err)
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
			responseInternalServerError(c, err)
			return
		}

		images = append(images, &model.CarImage{
			CarID:    req.CarID,
			URL:      s.s3store.Config.BaseURL + key,
			Category: req.CarImageCategory,
			Status:   model.CarImageStatusActive,
		})
	}

	// TODO: check enough image
	if err := s.store.CarImageStore.Create(images); err != nil {
		responseInternalServerError(c, err)
		return
	}

	if err := s.store.CarStore.Update(car.ID, map[string]interface{}{"status": model.MoveNextCarState(car.Status)}); err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{
		"status": "upload images successfully",
	})
}

func (s *Server) HandleUploadDrivingLicenseImages(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	req := struct {
		Files []*multipart.FileHeader `form:"files"`
	}{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUploadDocumentRequest, err)
		return
	}

	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if len(req.Files) > MaxNumberDrivingLicenseFiles {
		responseCustomErr(c, ErrCodeInvalidNumberOfFiles, fmt.Errorf("exceed maximum number of files, max %d, has %d", MaxNumberDrivingLicenseFiles, len(req.Files)))
		return
	}

	images := make([]*model.DrivingLicenseImage, 0)
	for _, f := range req.Files {
		if f.Size > MaxUploadFileSize {
			responseCustomErr(c, ErrCodeInvalidFileSize, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, f.Size))
			return
		}

		body, err := f.Open()
		if err != nil {
			responseCustomErr(c, ErrCodeReadingDocumentRequest, err)
			return
		}
		defer body.Close()

		url, err := s.uploadDocument(body, f.Filename)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}

		images = append(images, &model.DrivingLicenseImage{
			ID:        0,
			AccountID: acct.ID,
			URL:       url,
			Status:    model.DrivingLicenseImageStatusActive,
		})
	}

	if err := s.store.DrivingLicenseImageStore.Create(images); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "upload driving license images successfully"})
}

func (s *Server) HandleGetDrivingLicenseImages(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	images, err := s.store.DrivingLicenseImageStore.Get(acct.ID, model.DrivingLicenseImageStatusActive, 2)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	urls := make([]string, len(images))
	for i, d := range images {
		urls[i] = d.URL
	}

	responseSuccess(c, urls)
}

func (s *Server) HandleAdminUploadCustomerContractDocument(c *gin.Context) {
	req := struct {
		CustomerContractID            int                                 `form:"customer_contract_id"`
		CustomerContractImageCategory model.CustomerContractImageCategory `form:"document_category"`
		Files                         []*multipart.FileHeader             `form:"files"`
	}{}

	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidUploadDocumentRequest, err)
		return
	}

	contract, err := s.store.CustomerContractStore.FindByID(req.CustomerContractID)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	if contract.Status != model.CustomerContractStatusOrdered &&
		contract.Status != model.CustomerContractStatusAppraisingCarApproved {
		responseCustomErr(c, ErrCodeInvalidCustomerContractStatus, errors.New(
			fmt.Sprintf("invalid customer contract status, found %s", string(contract.Status))),
		)
		return
	}

	maxFile := MaxNumberCollateralAssetFiles
	if req.CustomerContractImageCategory == model.CustomerContractImageCategoryReceivingCarImages {
		maxFile = MaxNumberReceivingCarImages
	}

	if len(req.Files) > maxFile {
		responseCustomErr(c, ErrCodeInvalidNumberOfFiles, fmt.Errorf("exceed maximum number of files, max %d, has %d", maxFile, len(req.Files)))
		return
	}

	images := make([]*model.CustomerContractImage, 0)
	for _, f := range req.Files {
		if f.Size > MaxUploadFileSize {
			responseCustomErr(c, ErrCodeInvalidFileSize, fmt.Errorf("exceed maximum file size, max %d, has %d", MaxUploadFileSize, f.Size))
			return
		}

		body, err := f.Open()
		if err != nil {
			responseCustomErr(c, ErrCodeReadingDocumentRequest, err)
			return
		}
		defer body.Close()

		url, err := s.uploadDocument(body, f.Filename)
		if err != nil {
			responseInternalServerError(c, err)
			return
		}

		images = append(images, &model.CustomerContractImage{
			CustomerContractID: req.CustomerContractID,
			URL:                url,
			Category:           req.CustomerContractImageCategory,
			Status:             model.CustomerContractImageStatusActive,
		})
	}

	if err := s.store.CustomerContractImageStore.Create(images); err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, gin.H{"status": "upload customer contract document successfully"})
}

func (s *Server) uploadDocument(
	reader io.Reader,
	fileName string,
) (string, error) {
	extension := strings.Split(fileName, ".")[1]
	key := strings.Join([]string{uuid.NewString(), extension}, ".")

	_, err := s.s3store.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.s3store.Config.Bucket),
		Body:   reader,
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	url := s.s3store.Config.BaseURL + key
	return url, nil
}
