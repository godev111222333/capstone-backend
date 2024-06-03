package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) HandleGetAllCarModels(c *gin.Context) {
	models, err := s.store.CarModelStore.GetAll()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, models)
}
