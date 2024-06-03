package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	RegisterCarPeriod = []int{1, 3, 6, 12}
)

func (s *Server) HandleGetRegisterCarMetadata(c *gin.Context) {
	models, err := s.store.CarModelStore.GetAll()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models":  models,
		"periods": RegisterCarPeriod,
	})
}
