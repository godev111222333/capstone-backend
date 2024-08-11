package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

const (
	RoutePing                                        = "ping"
	RouteTestAuthorization                           = "test_authorization"
	RouteRegisterPartner                             = "register_partner"
	RouteRegisterCustomer                            = "register_customer"
	RouteVerifyOTP                                   = "verify_otp"
	RouteUploadAvatar                                = "upload_avatar"
	RouteRawLogin                                    = "login"
	RouteUpdateProfile                               = "update_profile"
	RouteGetRegisterCarMetadata                      = "register_car_metadata"
	RouteGetParkingLotMetadata                       = "get_parking_lot_metadata"
	RouteRegisterCar                                 = "register_car"
	RouteUpdateRentalPrice                           = "update_rental_price"
	RouteUploadCarDocuments                          = "upload_car_images"
	RouteGetRegisteredCars                           = "get_registered_cars"
	RouteGetBankMetadata                             = "get_bank_metadata"
	RouteUpdatePaymentInformation                    = "update_payment_information"
	RouteGetPaymentInformation                       = "get_payment_information"
	RouteGetProfile                                  = "get_profile"
	RouteUploadQRCode                                = "upload_qr_code"
	RouteGetGarageConfigs                            = "get_garage_configs"
	RouteUpdateGarageConfigs                         = "update_garage_configs"
	RouteAdminGetCars                                = "admin_get_cars"
	RouteGetCarDetail                                = "admin_get_car_details"
	RouteAdminApproveCar                             = "admin_approve_car"
	RouteAdminGetCustomerContracts                   = "admin_get_customer_contracts"
	RouteAdminUploadCustomerContractDocument         = "admin_upload_customer_contract_document"
	RouteAdminApproveRejectCustomerContract          = "admin_approve_reject_customer_contract"
	RouteAdminGetAccounts                            = "admin_get_accounts"
	RouteAdminGetAccountDetail                       = "admin_get_account_detail"
	RouteAdminGetPartnerContractDetail               = "admin_get_partner_contract_detail"
	RouteAdminGetCustomerPayments                    = "admin_get_customer_payments"
	RouteAdminGenerateCustomerPaymentQRCode          = "admin_generate_customer_payment_qr_code"
	RouteAdminGenerateMultipleCustomerPaymentsQRCode = "admin_generate_multiple_customer_payments_qr_code"
	RouteAdminGenerateMultiplePartnerPaymentsQRCode  = "admin_generate_multiple_partner_payments_qr_code"
	RouteAdminCompleteCustomerContract               = "admin_complete_customer_contract"
	RouteAdminUpdateCustomerContractImageStatus      = "admin_update_customer_contract_image_status"
	RouteAdminGetFeedbacks                           = "admin_get_feedbacks"
	RouteAdminUpdateFeedbackStatus                   = "admin_update_feedback_status"
	RouteAdminCancelCustomerPayment                  = "admin_cancel_customer_payment"
	RouteAdminGetConversations                       = "admin_get_conversations"
	RouteAdminGetConversationMessage                 = "admin_get_conversation_message"
	RouteAdminUpdateIsReturnCollateralAsset          = "admin_update_collateral_asset"
	RouteAdminSubscribeNotification                  = "admin_subscribe_notification"
	RouteAdminSubscribeNewConversation               = "admin_subscribe_new_conversation"
	RouteTechSubscribeNotification                   = "tech_subscribe_notification"
	RouteAdminGetStatistic                           = "admin_get_statistic"
	RouteAdminMakeMonthlyPartnerPayments             = "admin_make_monthly_partner_payments"
	RouteAdminGetMonthlyPartnerPayments              = "admin_get_monthly_partner_payments"
	RouteAdminFindCars                               = "admin_find_cars"
	RouteAdminChangeCar                              = "admin_change_car"
	RouteAdminUpdateWarningCount                     = "admin_update_warning_count"
	RouteAdminGetCustomerContractRule                = "admin_get_customer_contract_rule"
	RouteAdminGetPartnerContractRule                 = "admin_get_partner_contract_rule"
	RouteAdminCreateCustomerContractRule             = "admin_create_customer_contract_rule"
	RouteAdminCreatePartnerContractRule              = "admin_create_partner_contract_rule"
	RoutePartnerAgreeContract                        = "partner_agree_contract"
	RouteGetPartnerContractDetail                    = "get_partner_contract_detail"
	RoutePartnerGetActivityDetail                    = "partner_get_activity_detail"
	RouteGetPartnerRevenue                           = "partner_get_revenue"
	RoutePartnerGetPendingCustomerContracts          = "partner_get_pending_application"
	RoutePartnerApproveCustomerContract              = "partner_approve_customer_contract"
	RouteCustomerFindCars                            = "customer_find_cars"
	RouteCustomerRentCar                             = "customer_rent_car"
	RouteCustomerUploadDrivingLicenseImages          = "customer_upload_driving_license_images"
	RouteCustomerGetDrivingLicenseImages             = "customer_get_driving_license_images"
	RouteCustomerGetContracts                        = "customer_get_contracts"
	RouteCustomerAdminGetContractDetail              = "customer_get_contract_detail"
	RouteCustomerAgreeContract                       = "customer_agree_contract"
	RouteCustomerGetLastPaymentDetail                = "customer_get_payment_document_detail"
	RouteCustomerCalculateRentingPrice               = "customer_calculate_renting_price"
	RouteCustomerGetActivities                       = "customer_get_activities"
	RouteCustomerGiveFeedback                        = "customer_give_feedback"
	RouteCustomerGetSuggestedCars                    = "customer_get_suggested_cars"
	RouteCustomerPartnerGetFeedbacksByCar            = "customer_partner_get_feedbacks_by_car"
	RouteTechAppraisingCarOfCustomerContract         = "tech_appraising_car_of_cus_contract"
	RouteChat                                        = "chat"
	RouteVNPayIPNURL                                 = "vn_pay_ipn_url"
	RouteVNPayReturnURL                              = "vn_pay_return_url"
	RouteRegisterExpoPushToken                       = "register_expo_push_token"
	RouteGetNotificationHistory                      = "get_notification_history"
)

var (
	AuthRolePartner         = []string{model.RoleNamePartner}
	AuthRoleAdmin           = []string{model.RoleNameAdmin}
	AuthRoleCustomer        = []string{model.RoleNameCustomer}
	AuthRoleTechnician      = []string{model.RoleNameTechnician}
	AuthRoleCustomerAdmin   = []string{model.RoleNameCustomer, model.RoleNameAdmin}
	AuthRoleCustomerPartner = []string{model.RoleNameCustomer, model.RoleNamePartner}
	AuthRoleAdminTechnician = []string{model.RoleNameAdmin, model.RoleNameTechnician}
	AuthRoleAll             = []string{model.RoleNameCustomer, model.RoleNamePartner, model.RoleNameAdmin, model.RoleNameTechnician}
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
		RouteGetParkingLotMetadata: {
			Path:        "/register_car_metadata/parking_lot",
			Method:      http.MethodGet,
			Handler:     s.HandleGetParkingLotMetadata,
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
		RouteAdminGetCustomerContracts: {
			Path:        "/contracts",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCustomerContracts,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdminTechnician,
		},
		RouteAdminApproveRejectCustomerContract: {
			Path:        "/admin/contract",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminApproveOrRejectCustomerContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetAccounts: {
			Path:        "/admin/accounts",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetAccounts,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetAccountDetail: {
			Path:        "/account/:account_id",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetAccountDetail,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdminTechnician,
		},
		RouteAdminUploadCustomerContractDocument: {
			Path:        "/contract/document",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminUploadCustomerContractDocument,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdminTechnician,
		},
		RouteAdminGetPartnerContractDetail: {
			Path:        "/admin/partner_contract",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPartnerContractDetail,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetCustomerPayments: {
			Path:        "/customer_payments",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCustomerPayments,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerAdmin,
		},
		RouteAdminGenerateCustomerPaymentQRCode: {
			Path:        "/customer_payment/generate_qr",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminGenerateCustomerPaymentQRCode,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGenerateMultipleCustomerPaymentsQRCode: {
			Path:        "/customer_payment/multiple/generate_qr",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminGenerateMultipleCustomerPayments,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerAdmin,
		},
		RouteAdminGenerateMultiplePartnerPaymentsQRCode: {
			Path:        "/admin/monthly_partner_payment/multiple/generate_qr",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminGenerateMultiplePartnerPayments,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetMonthlyPartnerPayments: {
			Path:        "/admin/monthly_partner_payments",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetMonthlyPartnerPayments,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminCompleteCustomerContract: {
			Path:        "/admin/contract/complete",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminCompleteCustomerContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminUpdateCustomerContractImageStatus: {
			Path:        "/admin/contract/image",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminUpdateCustomerContractImageStatus,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetFeedbacks: {
			Path:        "/admin/feedbacks",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetFeedbacks,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminUpdateFeedbackStatus: {
			Path:        "/admin/feedback",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminUpdateFeedbackStatus,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminCancelCustomerPayment: {
			Path:        "/admin/customer_payment/cancel",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminCancelCustomerPayment,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminUpdateIsReturnCollateralAsset: {
			Path:        "/admin/update_is_return_collateral_asset",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminUpdateReturnCollateralAsset,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetStatistic: {
			Path:        "/admin/statistic",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetStatistic,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetConversations: {
			Path:        "/admin/conversations",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetConversations,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetConversationMessage: {
			Path:        "/conversation/messages",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetMessages,
			RequireAuth: true,
			AuthRoles:   AuthRoleAll,
		},
		RouteAdminSubscribeNotification: {
			Path:        "/admin/subscribe_notification",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminSubscribeNotification,
			RequireAuth: false,
		},
		RouteAdminSubscribeNewConversation: {
			Path:        "/admin/subscribe_conversation",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminSubscribeNewConversation,
			RequireAuth: false,
		},
		RouteTechSubscribeNotification: {
			Path:        "/tech/subscribe_notification",
			Method:      http.MethodGet,
			Handler:     s.HandleTechnicianSubscribeNotification,
			RequireAuth: false,
		},
		RoutePartnerAgreeContract: {
			Path:        "/partner/contract/agree",
			Method:      http.MethodPut,
			Handler:     s.HandlePartnerAgreeContract,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RoutePartnerGetActivityDetail: {
			Path:        "/partner/activity",
			Method:      http.MethodGet,
			Handler:     s.HandlePartnerGetActivityDetail,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RoutePartnerGetPendingCustomerContracts: {
			Path:        "/partner/customer_contracts",
			Method:      http.MethodGet,
			Handler:     s.HandlePartnerGetPendingCustomerContracts,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RoutePartnerApproveCustomerContract: {
			Path:        "/partner/customer_contract/approve_reject",
			Method:      http.MethodPut,
			Handler:     s.HandlePartnerApproveCustomerContract,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteGetPartnerContractDetail: {
			Path:        "/partner/contract",
			Method:      http.MethodGet,
			Handler:     s.HandleGetPartnerContractDetail,
			RequireAuth: true,
			AuthRoles:   AuthRolePartner,
		},
		RouteGetPartnerRevenue: {
			Path:        "/partner/revenue",
			Method:      http.MethodGet,
			Handler:     s.HandlePartnerGetRevenue,
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
		RouteCustomerAdminGetContractDetail: {
			Path:        "/contract/:customer_contract_id",
			Method:      http.MethodGet,
			Handler:     s.HandleGetCustomerContractDetails,
			RequireAuth: true,
			AuthRoles:   AuthRoleAll,
		},
		RouteCustomerCalculateRentingPrice: {
			Path:        "/customer/calculate_rent_pricing",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerCalculateRentPricing,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerAgreeContract: {
			Path:        "/customer/contract/agree",
			Method:      http.MethodPut,
			Handler:     s.HandleCustomerAgreeContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerGetLastPaymentDetail: {
			Path:        "/customer/last_payment_detail",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerGetLastPaymentDetail,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerGetActivities: {
			Path:        "/customer/activities",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerGetActivities,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerGiveFeedback: {
			Path:        "/customer/feedback",
			Method:      http.MethodPut,
			Handler:     s.HandleCustomerGiveFeedback,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteCustomerPartnerGetFeedbacksByCar: {
			Path:        "/feedbacks/car",
			Method:      http.MethodGet,
			Handler:     s.HandleGetFeedbackByCar,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerPartner,
		},
		RouteCustomerGetSuggestedCars: {
			Path:        "/customer/suggested_cars",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerGetSuggestedCars,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomer,
		},
		RouteChat: {
			Path:        "/chat",
			Method:      http.MethodGet,
			Handler:     s.HandleChat,
			RequireAuth: false,
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
		RouteRegisterExpoPushToken: {
			Path:        "/expo_push_token",
			Method:      http.MethodPost,
			Handler:     s.HandleRegisterExpoPushToken,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerPartner,
		},
		RouteGetNotificationHistory: {
			Path:        "/notifications",
			Method:      http.MethodGet,
			Handler:     s.HandleGetNotificationHistory,
			RequireAuth: true,
			AuthRoles:   AuthRoleAll,
		},
		RouteAdminMakeMonthlyPartnerPayments: {
			Path:        "/admin/monthly_partner_payments",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminMakeMonthlyPartnerPayments,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminFindCars: {
			Path:        "/admin/find_change_cars",
			Method:      http.MethodGet,
			Handler:     s.HandleCustomerFindCars,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminChangeCar: {
			Path:        "/admin/customer_contract/change_car",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminChangeCar,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminUpdateWarningCount: {
			Path:        "/admin/warning_count",
			Method:      http.MethodPut,
			Handler:     s.HandleAdminUpdateWarningCount,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminGetCustomerContractRule: {
			Path:        "/customer_contract_rule",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetCustomerContractRule,
			RequireAuth: true,
			AuthRoles:   AuthRoleCustomerAdmin,
		},
		RouteAdminGetPartnerContractRule: {
			Path:        "/admin/partner_contract_rule",
			Method:      http.MethodGet,
			Handler:     s.HandleAdminGetPartnerContractRule,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminCreateCustomerContractRule: {
			Path:        "/admin/customer_contract_rule",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminCreateCustomerContractRule,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteAdminCreatePartnerContractRule: {
			Path:        "/admin/partner_contract_rule",
			Method:      http.MethodPost,
			Handler:     s.HandleAdminCreatePartnerContractRule,
			RequireAuth: true,
			AuthRoles:   AuthRoleAdmin,
		},
		RouteTechAppraisingCarOfCustomerContract: {
			Path:        "/customer_contract/appraising_car",
			Method:      http.MethodPut,
			Handler:     s.HandleTechnicianAppraisingCarOfCusContract,
			RequireAuth: true,
			AuthRoles:   AuthRoleTechnician,
		},

		// Temporary API
		"set_admin_return_url": {
			Path:        "/set_admin_return_url",
			Method:      http.MethodPost,
			Handler:     s.updateAdminReturnURL,
			RequireAuth: false,
		},
	}
}
