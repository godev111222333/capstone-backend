package api

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func (s *Server) updateAdminReturnURL(c *gin.Context) {
	req := struct {
		NewReturnURL string `json:"new_return_url" binding:"required"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	s.feCfg.AdminReturnURL = req.NewReturnURL
	file, err := os.OpenFile(s.feCfg.Path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(s.feCfg); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	responseSuccess(c, gin.H{"status": "update admin return url successfully"})
}

func (s *Server) seedImage(c *gin.Context) {
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

	url, err := s.uploadDocument(body, header.Filename)
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	responseSuccess(c, gin.H{"url": url})
}
