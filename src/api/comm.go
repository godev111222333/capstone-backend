package api

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ErrorCode int

const (
	ErrCodeSuccess                                            ErrorCode = 10000
	ErrCodeRecordNotFound                                     ErrorCode = 10001
	ErrCodeDuplicatedKey                                      ErrorCode = 10002
	ErrCodePrimaryKeyRequired                                 ErrorCode = 10003
	ErrCodeForeignKeyViolated                                 ErrorCode = 10004
	ErrCodeCheckConstraintViolated                            ErrorCode = 10005
	ErrCodeInvalidOwnership                                   ErrorCode = 10006
	ErrCodeInternalServerError                                ErrorCode = 10009
	ErrCodeInvalidVerifyOTPRequest                            ErrorCode = 10010
	ErrCodeInvalidAccountStatus                               ErrorCode = 10011
	ErrCodeInvalidOTP                                         ErrorCode = 10012
	ErrCodeInvalidLoginRequest                                ErrorCode = 10013
	ErrCodeWrongPhoneNumberOrPassword                         ErrorCode = 10014
	ErrCodeAccountNotActive                                   ErrorCode = 10015
	ErrCodeInvalidUpdateProfileRequest                        ErrorCode = 10016
	ErrCodeMissingAuthorizeHeader                             ErrorCode = 10017
	ErrCodeInvalidAuthorizeHeaderFormat                       ErrorCode = 10018
	ErrCodeInvalidRole                                        ErrorCode = 10019
	ErrCodeInvalidAuthorizeType                               ErrorCode = 10020
	ErrCodeVerifyAccessToken                                  ErrorCode = 10021
	ErrCodeUpdatePaymentInfoRequest                           ErrorCode = 10022
	ErrCodeInvalidUploadDocumentRequest                       ErrorCode = 10023
	ErrCodeInvalidFileSize                                    ErrorCode = 10024
	ErrCodeReadingDocumentRequest                             ErrorCode = 10025
	ErrCodeGetCarsRequest                                     ErrorCode = 10026
	ErrCodeGetCarDetailRequest                                ErrorCode = 10027
	ErrCodeInvalidUpdateGarageConfigRequest                   ErrorCode = 10028
	ErrCodeInvalidSeat                                        ErrorCode = 10029
	ErrCodeInvalidAdminApproveOrRejectCarRequest              ErrorCode = 10030
	ErrCodeNotEnoughSlotAtGarage                              ErrorCode = 10031
	ErrCodeInvalidCarStatus                                   ErrorCode = 10032
	ErrCodeInvalidPartnerContractStatus                       ErrorCode = 10033
	ErrCodeInvalidGetCustomerContractRequest                  ErrorCode = 10034
	ErrCodeInvalidAdminApproveOrRejectCustomerContractRequest ErrorCode = 10035
	ErrCodeInvalidCustomerContractStatus                      ErrorCode = 10036
	ErrCodeInvalidGetAccountsRequest                          ErrorCode = 10037
	ErrCodeInvalidSetAccountStatusRequest                     ErrorCode = 10038
	ErrCodeInvalidGetAccountDetailRequest                     ErrorCode = 10039
	ErrCodeInvalidGetCustomerPaymentRequest                   ErrorCode = 10040
	ErrCodeInvalidCreateCustomerPaymentRequest                ErrorCode = 10041
	ErrCodeInvalidGenerateCustomerPaymentQRCode               ErrorCode = 10042
	ErrCodeGenerateQRCode                                     ErrorCode = 10043
	ErrCodeInvalidGetParkingLotRequest                        ErrorCode = 10044
	ErrCodeInvalidRegisterCustomerRequest                     ErrorCode = 10045
	ErrCodeHashingPassword                                    ErrorCode = 10046
	ErrCodeSendOTP                                            ErrorCode = 10047
	ErrCodeInvalidFindCarsRequest                             ErrorCode = 10048
	ErrCodeInvalidRentCarRequest                              ErrorCode = 10049
	ErrCodeInvalidCustomerAgreeContractRequest                ErrorCode = 10050
	ErrCodeInvalidGetCustomerContractDetailRequest            ErrorCode = 10051
	ErrCodeInvalidCalculateRentingPriceRequest                ErrorCode = 10052
	ErrCodeGetLastPaymentTypeRequest                          ErrorCode = 10053
	ErrCodeInvalidDocumentCategory                            ErrorCode = 10054
	ErrCodeInvalidNumberOfFiles                               ErrorCode = 10055
	ErrCodeInvalidRegisterPartnerRequest                      ErrorCode = 10056
	ErrCodeInvalidRegisterCarRequest                          ErrorCode = 10057
	ErrCodeInvalidUpdateRentalPriceRequest                    ErrorCode = 10058
	ErrCodeInvalidGetRegisteredCarsRequest                    ErrorCode = 10059
	ErrCodeInvalidPartnerAgreeContractRequest                 ErrorCode = 10060
	ErrCodeInvalidGetPartnerContractDetailRequest             ErrorCode = 10061
	ErrCodeDatabaseError                                      ErrorCode = 10062
	ErrCodeInvalidCompleteCustomerContractRequest             ErrorCode = 10063
	ErrCodeExistPendingPayments                               ErrorCode = 10064
	ErrCodeInvalidUpdateCustomerContractImageStatusRequest    ErrorCode = 10065
	ErrCodeInvalidCustomerGetActivitiesRequest                ErrorCode = 10066
	ErrCodeInvalidGiveFeedbackRequest                         ErrorCode = 10067
	ErrCodeMissingPaymentInformation                          ErrorCode = 10068
	ErrCodeInvalidPartnerGetActivityDetailRequest             ErrorCode = 10069
	ErrCodeInvalidAdminGetFeedbackRequest                     ErrorCode = 10070
	ErrCodeInvalidAdminUpdateFeedbackStatusRequest            ErrorCode = 10071
	ErrCodeInvalidAdminCancelCustomerPaymentRequest           ErrorCode = 10072
	ErrCodeUnableUpgradeWebsocket                             ErrorCode = 10073
	ErrCodeInvalidAdminGetConversationsRequest                ErrorCode = 10074
	ErrCodeInvalidAdminGetMessagesRequest                     ErrorCode = 10075
	ErrCodeInvalidGenerateMultipleCustomerPaymentsRequest     ErrorCode = 10076
	ErrCodeMissingDrivingLicence                              ErrorCode = 10077
	ErrCodeInvalidGetFeedBackByCarRequest                     ErrorCode = 10078
	ErrCodeInvalidUpdateIsReturnCollateralAsset               ErrorCode = 10079
	ErrCodeInvalidRegisterExpoPushTokenRequest                ErrorCode = 10080
	ErrCodeInvalidStatisticRequest                            ErrorCode = 10081
	ErrCodeInvalidGetNotificationHistoryRequest               ErrorCode = 10082
	ErrCodeCacheError                                         ErrorCode = 10083
	ErrCodeInvalidAdminMakeMonthlyPaymentRequest              ErrorCode = 10084
	ErrCodeInvalidAdminGetMonthlyPartnerPaymentRequest        ErrorCode = 10085
	ErrCodeInvalidGenerateMultiplePartnerPaymentsRequest      ErrorCode = 10086
	ErrCodeInvalidAdminChangeCarRequest                       ErrorCode = 10087
	ErrCodeOverlapOtherContract                               ErrorCode = 10088
	ErrCodeInvalidUpdateWarningCounterRequest                 ErrorCode = 10089
	ErrCodeInvalidCreateCustomerContractRuleRequest           ErrorCode = 10090
	ErrCodeInvalidCreatePartnerContractRuleRequest            ErrorCode = 10091
	ErrCodeInvalidPartnerGetPendingApprovalRequest            ErrorCode = 10092
	ErrCodeInvalidPartnerApproveCustomerContractRequest       ErrorCode = 10093
	ErrCodeInvalidAppraisingCarRequest                        ErrorCode = 10094
	ErrCodeInvalidReturnCarCustomerContractRequest            ErrorCode = 10095
	ErrCodeInvalidTechAppraisingReturnCarRequest              ErrorCode = 10096
	ErrCodeInvalidGetCarModelsRequest                         ErrorCode = 10097
	ErrCodeInvalidCreateCarModelsRequest                      ErrorCode = 10098
	ErrCodeInvalidUpdateCarModelsRequest                      ErrorCode = 10099
	ErrCodeInvalidSetCustomerContractResolveStatusRequest     ErrorCode = 10100
	ErrCodeInvalidInactiveCarRequest                          ErrorCode = 100101
	ErrCodeInvalidPartnerGetRevenueRequest                    ErrorCode = 100102
	ErrCodeOverWithOtherContractRequest                       ErrorCode = 100103
)

var customErrMapping = map[ErrorCode]CommResponse{
	ErrCodeInvalidOTP:                 {ErrCodeInvalidOTP, "invalid OTP or OTP was expired", nil},
	ErrCodeWrongPhoneNumberOrPassword: {ErrCodeWrongPhoneNumberOrPassword, "wrong phone number or password", nil},
	ErrCodeAccountNotActive:           {ErrCodeInvalidAccountStatus, "active is not active", nil},
	ErrCodeInvalidRole:                {ErrCodeInvalidRole, "invalid role", nil},
	ErrCodeInvalidOwnership:           {ErrCodeInvalidOwnership, "invalid ownership", nil},
	ErrCodeExistPendingPayments:       {ErrCodeExistPendingPayments, "exist pending payments for this contract", nil},
	ErrCodeMissingPaymentInformation:  {ErrCodeMissingPaymentInformation, "missing bank information", nil},
	ErrCodeMissingDrivingLicence:      {ErrCodeMissingDrivingLicence, "missing driving license images", nil},
	ErrCodeSuccess:                    {ErrCodeSuccess, "success", nil},
}

var gormErrMapping = map[error]CommResponse{
	gorm.ErrRecordNotFound:          {ErrCodeRecordNotFound, gorm.ErrRecordNotFound.Error(), nil},
	gorm.ErrDuplicatedKey:           {ErrCodeDuplicatedKey, gorm.ErrDuplicatedKey.Error(), nil},
	gorm.ErrPrimaryKeyRequired:      {ErrCodePrimaryKeyRequired, gorm.ErrPrimaryKeyRequired.Error(), nil},
	gorm.ErrForeignKeyViolated:      {ErrCodeForeignKeyViolated, gorm.ErrForeignKeyViolated.Error(), nil},
	gorm.ErrCheckConstraintViolated: {ErrCodeCheckConstraintViolated, gorm.ErrCheckConstraintViolated.Error(), nil},
}

type CommResponse struct {
	ErrorCode ErrorCode   `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func responseGormErr(c *gin.Context, err error) {
	respCode := http.StatusBadRequest
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respCode = http.StatusNotFound
	}

	code := ErrCodeDatabaseError
	e, ok := gormErrMapping[err]
	if ok {
		code = e.ErrorCode
	}

	c.AbortWithStatusJSON(respCode, CommResponse{
		ErrorCode: code,
		Message:   err.Error(),
	})
	return
}

// responseCustomErr if err is nil, use the error message from customErrMapping
func responseCustomErr(c *gin.Context, errCode ErrorCode, err error) {
	var errMsg string
	if err == nil {
		msg, ok := customErrMapping[errCode]
		if ok {
			errMsg = msg.Message
		}
	} else {
		errMsg = err.Error()
	}
	c.AbortWithStatusJSON(errCodeToResponseCode(errCode), CommResponse{
		ErrorCode: errCode,
		Message:   errMsg,
	})
	return
}

func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, CommResponse{
		ErrorCode: ErrCodeSuccess,
		Message:   "success",
		Data:      data,
	})
}

func responseInternalServerError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, CommResponse{
		ErrorCode: ErrCodeInternalServerError,
		Message:   err.Error(),
	})
}

func errCodeToResponseCode(errCode ErrorCode) int {
	authorizeCode := []ErrorCode{
		ErrCodeInvalidOwnership,
		ErrCodeInvalidAuthorizeType,
		ErrCodeVerifyAccessToken,
	}
	for _, c := range authorizeCode {
		if errCode == c {
			return http.StatusUnauthorized
		}
	}

	return http.StatusBadRequest
}
