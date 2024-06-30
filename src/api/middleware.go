package api

import (
	"errors"
	"fmt"
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

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			responseCustomErr(ctx, ErrCodeMissingAuthorizeHeader, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			responseCustomErr(ctx, ErrCodeInvalidAuthorizeHeaderFormat, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			responseCustomErr(ctx, ErrCodeInvalidAuthorizeType, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			responseCustomErr(ctx, ErrCodeVerifyAccessToken, err)
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
			responseGormErr(ctx, err)
			return
		}

		if acct.Status != model.AccountStatusActive {
			responseCustomErr(ctx, ErrCodeAccountNotActive, err)
			return
		}

		ctx.Next()
	}
}

func (s *Server) authRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if authPayload.Role != role {
			responseCustomErr(ctx, ErrCodeInvalidRole, nil)
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
