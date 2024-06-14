package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RoutePing                               = "ping"
	RouteTestAuthorization                  = "test_authorization"
	RouteRegisterPartner                    = "register_partner"
	RouteRegisterCustomer                   = "register_customer"
	RouteVerifyOTP                          = "verify_otp"
	RouteUploadAvatar                       = "upload_avatar"
	RouteRawLogin                           = "login"
	RouteRenewAccessToken                   = "renew_access_token"
	RouteUpdateProfile                      = "update_profile"
	RouteGetRegisterCarMetadata             = "register_car_metadata"
	RouteRegisterCar                        = "register_car"
	RouteUpdateRentalPrice                  = "update_rental_price"
	RouteUploadCarDocuments                 = "upload_car_images"
	RouteGetRegisteredCars                  = "get_registered_cars"
	RouteGetBankMetadata                    = "get_bank_metadata"
	RouteUpdatePaymentInformation           = "update_payment_information"
	RouteGetPaymentInformation              = "get_payment_information"
	RouteGetProfile                         = "get_profile"
	RouteUploadQRCode                       = "upload_qr_code"
	RouteGetGarageConfigs                   = "get_garage_configs"
	RouteUpdateGarageConfigs                = "update_garage_configs"
	RouteAdminGetCars                       = "admin_get_cars"
	RouteAdminGetCarDetails                 = "admin_get_car_details"
	RouteAdminApproveCar                    = "admin_approve_car"
	RoutePartnerSignContract                = "partner_sign_contract"
	RouteGetPartnerContractDetails          = "get_partner_contract_detail"
	RouteCustomerFindCars                   = "customer_find_cars"
	RouteCustomerRentCar                    = "customer_rent_car"
	RouteCustomerUploadDrivingLicenseImages = "customer_upload_driving_license_images"
	RouteCustomerGetDrivingLicenseImages    = "customer_get_driving_license_images"
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
		RouteGetProfile: {
			Path:        "/profile",
			Method:      http.MethodGet,
			Handler:     s.HandleGetProfile,
			RequireAuth: true,
		},
		RouteGetRegisterCarMetadata: {
			Path:        "/register_car_metadata",
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
		RouteGetBankMetadata: {
			Path:        "/banks",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPaymentInformationMetadata,
			RequireAuth: false,
		},
		RouteGetPaymentInformation: {
			Path:        "/payment_info",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPaymentInformation,
			RequireAuth: true,
		},
		RouteUpdatePaymentInformation: {
			Path:        "/payment_info",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdatePaymentInformation,
			RequireAuth: true,
		},
		RouteUploadQRCode: {
			Path:        "/payment_info/qr",
			Method:      http.MethodPost,
			Handler:     s.HandleUpdateQRCodeImage,
			RequireAuth: true,
		},
		RouteGetGarageConfigs: {
			Path:        "/garage_config",
			Method:      http.MethodGet,
			Handler:     s.HandleGetGarageConfigs,
			RequireAuth: true,
		},
		RouteUpdateGarageConfigs: {
			Path:        "/garage_config",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdateGarageConfigs,
			RequireAuth: true,
		},
		RouteAdminGetCars: {
			Path:        "/admin/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCars,
			RequireAuth: true,
		},
		RouteAdminGetCarDetails: {
			Path:        "/admin/car/:id",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCarDetails,
			RequireAuth: true,
		},
		RouteAdminApproveCar: {
			Path:        "/admin/car_application",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminApproveOrRejectCar,
			RequireAuth: true,
		},
		RoutePartnerSignContract: {
			Path:        "/partner/contract/agree",
			Method:      http.MethodPut,
			Handler:     s.HandlePartnerAgreeContract,
			RequireAuth: true,
		},
		RouteGetPartnerContractDetails: {
			Path:        "/partner/contract",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPartnerContractDetails,
			RequireAuth: true,
		},
		RouteCustomerFindCars: {
			Path:        "/customer/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerFindCars,
			RequireAuth: false,
		},
		RouteCustomerRentCar: {
			Path:        "/customer/rent",
			Method:      http.MethodPost,
			Handler:     s.HandleCustomerRentCar,
			RequireAuth: true,
		},
		RouteRegisterCustomer: {
			Path:        "/customer/register",
			Method:      http.MethodPost,
			Handler:     s.HandleRegisterCustomer,
			RequireAuth: false,
		},
		RouteCustomerUploadDrivingLicenseImages: {
			Path:        "/customer/driving_license",
			Method:      http.MethodPost,
			Handler:     s.HandleUploadDrivingLicenseImages,
			RequireAuth: true,
		},
		RouteCustomerGetDrivingLicenseImages: {
			Path:        "/customer/driving_license",
			Method:      http.MethodGet,
			Handler:     s.HandleGetDrivingLicenseImages,
			RequireAuth: true,
		},
	}
}
