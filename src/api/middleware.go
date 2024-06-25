package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func (s *Server) activeAccountMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if acct.Status != model.AccountStatusActive {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("account is not active")))
			return
		}

		ctx.Next()
	}
}

func (s *Server) authRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if authPayload.Role != role {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid role")))
			return
		}

		ctx.Next()
	}
}

func (s *Server) arrayContainsString(arr []string, ele string) bool {
	for _, str := range arr {
		if strings.EqualFold(str, ele) {
			return true
		}
	}
	return false
}
