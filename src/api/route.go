package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RoutePing                   = "ping"
	RouteTestAuthorization      = "test_authorization"
	RouteRegisterPartner        = "register_partner"
	RouteVerifyOTP              = "verify_otp"
	RouteUploadAvatar           = "upload_avatar"
	RouteRawLogin               = "login"
	RouteRenewAccessToken       = "renew_access_token"
	RouteUpdateProfile          = "update_profile"
	RouteGetRegisterCarMetadata = "register_car_metadata"
	RouteRegisterCar            = "register_car"
	RouteUpdateRentalPrice      = "update_rental_price"
	RouteUploadCarDocuments     = "upload_car_images"
	RouteGetRegisteredCars      = "get_registered_cars"
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
			RequireAuth: true,
		},
		RouteVerifyOTP: {
			Path:        "/user/otp",
			Method:      http.MethodPost,
			Handler:     s.HandleVerifyOTP,
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
		RouteGetRegisterCarMetadata: {
			Path:        "/models",
			Method:      http.MethodGet,
			Handler:     s.HandleGetRegisterCarMetadata,
			RequireAuth: false,
		},
		RouteRegisterCar: {
			Path:        "/car",
			Method:      http.MethodPost,
			Handler:     s.HandleRegisterCar,
			RequireAuth: true,
		},
		RouteUpdateRentalPrice: {
			Path:        "/car/price",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdateRentalPrice,
			RequireAuth: true,
		},
		RouteUploadCarDocuments: {
			Path:        "/car/document",
			Method:      http.MethodPost,
			Handler:     s.HandleUploadCarDocuments,
			RequireAuth: true,
		},
		RouteGetRegisteredCars: {
			Path:        "/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleGetRegisteredCars,
			RequireAuth: true,
		},
	}
}
