package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleGetPaymentInformationMetadata(c *gin.Context) {
	responseSuccess(c, gin.H{
		"banks": s.bankMetadata,
	})
}
