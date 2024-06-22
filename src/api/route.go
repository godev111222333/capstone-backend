package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
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
	RouteGetCarDetail                       = "admin_get_car_details"
	RouteAdminApproveCar                    = "admin_approve_car"
	RouteAdminGetContracts                  = "admin_get_contracts"
	RouteAdminApproveRejectCustomerContract = "admin_approve_reject_customer_contract"
	RouteAdminGetAccounts                   = "admin_get_accounts"
	RoutePartnerAgreeContract               = "partner_agree_contract"
	RouteGetPartnerContractDetails          = "get_partner_contract_detail"
	RouteCustomerFindCars                   = "customer_find_cars"
	RouteCustomerRentCar                    = "customer_rent_car"
	RouteCustomerUploadDrivingLicenseImages = "customer_upload_driving_license_images"
	RouteCustomerGetDrivingLicenseImages    = "customer_get_driving_license_images"
	RouteCustomerGetContracts               = "customer_get_contracts"
	RouteCustomerAdminGetContractDetails    = "customer_get_contract_details"
	RouteCustomerAgreeContract              = "customer_agree_contract"
	RouteVNPayIPNURL                        = "vn_pay_ipn_url"
	RouteVNPayReturnURL                     = "vn_pay_return_url"
)

var (
	AuthRolePartner       = []string{model.RoleNamePartner}
	AuthRoleAdmin         = []string{model.RoleNameAdmin}
	AuthRoleCustomer      = []string{model.RoleNameCustomer}
	AuthRoleCustomerAdmin = []string{model.RoleNameCustomer, model.RoleNameAdmin}
)

type RouteInfo = struct {
	Path        string
	Method      string
	Handler     func(c *gin.Context)
	RequireAuth bool
	AuthRoles   []string
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
			AuthRoles:   AuthRolePartner,
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
			Path:        "/partner/car",
			Method:      http.MethodPost,
			Handler:     s.HandleRegisterCar,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteUpdateRentalPrice: {
			Path:        "/partner/car/price",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdateRentalPrice,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteUploadCarDocuments: {
			Path:        "/partner/car/document",
			Method:      http.MethodPost,
			Handler:     s.HandleUploadCarDocuments,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteGetRegisteredCars: {
			Path:        "/partner/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleGetRegisteredCars,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
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
			Path:        "/admin/garage_config",
			Method:      http.MethodGet,
			Handler:     s.HandleGetGarageConfigs,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteUpdateGarageConfigs: {
			Path:        "/admin/garage_config",
			Method:      http.MethodPut,
			Handler:     s.HandleUpdateGarageConfigs,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetCars: {
			Path:        "/admin/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCars,
			RequireAuth: true,
		},
		RouteGetCarDetail: {
			Path:        "/car/:id",
			Method:      http.MethodGet,
			Handler:     s.HandleGetCarDetail,
			RequireAuth: false,
		},
		RouteAdminApproveCar: {
			Path:        "/admin/car_application",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminApproveOrRejectCar,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetContracts: {
			Path:        "/admin/contracts",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetContracts,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminApproveRejectCustomerContract: {
			Path:        "/admin/contract",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminApproveOrRejectCustomerContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RoutePartnerAgreeContract: {
			Path:        "/partner/contract/agree",
			Method:      http.MethodPut,
			Handler:     s.HandlePartnerAgreeContract,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteGetPartnerContractDetails: {
			Path:        "/partner/contract",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPartnerContractDetails,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteCustomerFindCars: {
			Path:        "/customer/cars",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerFindCars,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerRentCar: {
			Path:        "/customer/rent",
			Method:      http.MethodPost,
			Handler:     s.HandleCustomerRentCar,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
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
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerGetDrivingLicenseImages: {
			Path:        "/customer/driving_license",
			Method:      http.MethodGet,
			Handler:     s.HandleGetDrivingLicenseImages,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerGetContracts: {
			Path:        "/customer/contracts",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerGetContracts,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerAdminGetContractDetails: {
			Path:        "/contract/:customer_contract_id",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerAdminGetContractDetails,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerAdmin,
		},
		RouteCustomerAgreeContract: {
			Path:        "/customer/contract/agree",
			Method:      http.MethodPut,
			Handler:     s.HandleCustomerAgreeContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteVNPayIPNURL: {
			Path:        "/vnpay/ipn",
			Method:      http.MethodGet,
			Handler:     s.HandleVnPayIPN,
			RequireAuth: false,
		},
		RouteVNPayReturnURL: {
			Path:        "/vnpay/return_url",
			Method:      http.MethodGet,
			Handler:     s.HandleVnPayReturnURL,
			RequireAuth: false,
		},
		RouteAdminGetAccounts: {
			Path:        "/admin/accounts",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetAccounts,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
	}
}
