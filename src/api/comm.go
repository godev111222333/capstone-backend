package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func responseError(ctx *gin.Context, err error) {
	errTxt := "something went wrong!"
	if err != nil {
		errTxt = err.Error()
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": errTxt,
	})
}
