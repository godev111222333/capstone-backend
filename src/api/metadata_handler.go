package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleGetPaymentInformationMetadata(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"banks": s.bankMetadata,
	})
}
