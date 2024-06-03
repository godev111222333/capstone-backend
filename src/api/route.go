package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RoutePing              = "ping"
	RouteTestAuthorization = "test_authorization"
	RouteRegisterPartner   = "register_partner"
	RouteUploadAvatar      = "upload_avatar"
	RouteRawLogin          = "login"
	RouteRenewAccessToken  = "renew_access_token"
	RouteUpdateProfile     = "update_profile"
	RouteGetAllCarModels   = "all_car_models"
)

type RouteInfo = struct {
	Path        string
	Method      string
	Handler     func(c *gin.Context)
	RequireAuth bool
}

func (s *Server) AllRoutes() map[string]RouteInfo {
	return map[string]RouteInfo{
		RoutePing: {
			Path:   "/ping",
			Method: http.MethodGet,
			Handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, "pong")
			},
			RequireAuth: false,
		},
		RouteTestAuthorization: {
			Path:   "/test_author",
			Method: http.MethodPost,
			Handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "authed",
				})
			},
			RequireAuth: true,
		},
		RouteRegisterPartner: {
			Path:        "/partner/register",
			Method:      http.MethodPost,
			Handler:     s.RegisterPartner,
			RequireAuth: false,
		},
		RouteUploadAvatar: {
			Path:        "/user/avatar/upload",
			Method:      http.MethodPost,
			Handler:     s.HandleUploadAvatar,
			RequireAuth: false,
		},
		RouteRawLogin: {
			Path:        "/login",
			Method:      http.MethodPost,
			Handler:     s.HandleRawLogin,
			RequireAuth: false,
		},
		RouteRenewAccessToken: {
			Path:        "/renew",
			Method:      http.MethodPost,
			Handler:     s.HandleRenewAccessToken,
			RequireAuth: false,
		},
		RouteUpdateProfile: {
			Path:        "/profile",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdateProfile,
			RequireAuth: true,
		},
		RouteGetAllCarModels: {
			Path:        "/models",
			Method:      http.MethodGet,
			Handler:     s.HandleGetAllCarModels,
			RequireAuth: false,
		},
	}
}
