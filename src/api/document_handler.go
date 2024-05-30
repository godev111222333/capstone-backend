package api

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MaxUploadFileSize = 5 * 1 << 20

func (s *Server) HandleUploadAvatar(c *gin.Context) {
	req := struct {
		AccountID int                   `form:"account_id"`
		File      *multipart.FileHeader `form:"file"`
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

	extension := strings.Split(header.Filename, ".")[1]
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
	if err := s.store.AccountStore.Update(req.AccountID, map[string]interface{}{
		"avatar_url": url,
	}); err != nil {
		responseError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "upload avatar successfully",
		"url":    url,
	})
}
