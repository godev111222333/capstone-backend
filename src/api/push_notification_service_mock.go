// Code generated by MockGen. DO NOT EDIT.
// Source: src/api/push_notification_service.go
//
// Generated by this command:
//
//	mockgen -source src/api/push_notification_service.go -destination=src/api/push_notification_service_mock.go -package=api
//

// Package api is a generated GoMock package.
package api

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockINotificationPushService is a mock of INotificationPushService interface.
type MockINotificationPushService struct {
	ctrl     *gomock.Controller
	recorder *MockINotificationPushServiceMockRecorder
}

// MockINotificationPushServiceMockRecorder is the mock recorder for MockINotificationPushService.
type MockINotificationPushServiceMockRecorder struct {
	mock *MockINotificationPushService
}

// NewMockINotificationPushService creates a new mock instance.
func NewMockINotificationPushService(ctrl *gomock.Controller) *MockINotificationPushService {
	mock := &MockINotificationPushService{ctrl: ctrl}
	mock.recorder = &MockINotificationPushServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINotificationPushService) EXPECT() *MockINotificationPushServiceMockRecorder {
	return m.recorder
}

// NewApproveCarDeliveryMsg mocks base method.
func (m *MockINotificationPushService) NewApproveCarDeliveryMsg(carID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewApproveCarDeliveryMsg", carID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewApproveCarDeliveryMsg indicates an expected call of NewApproveCarDeliveryMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewApproveCarDeliveryMsg(carID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewApproveCarDeliveryMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewApproveCarDeliveryMsg), carID, expoToken, toPhone)
}

// NewApproveCarRegisterMsg mocks base method.
func (m *MockINotificationPushService) NewApproveCarRegisterMsg(carID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewApproveCarRegisterMsg", carID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewApproveCarRegisterMsg indicates an expected call of NewApproveCarRegisterMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewApproveCarRegisterMsg(carID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewApproveCarRegisterMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewApproveCarRegisterMsg), carID, expoToken, toPhone)
}

// NewApproveRentingCarRequestMsg mocks base method.
func (m *MockINotificationPushService) NewApproveRentingCarRequestMsg(contractID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewApproveRentingCarRequestMsg", contractID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewApproveRentingCarRequestMsg indicates an expected call of NewApproveRentingCarRequestMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewApproveRentingCarRequestMsg(contractID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewApproveRentingCarRequestMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewApproveRentingCarRequestMsg), contractID, expoToken, toPhone)
}

// NewChatMsg mocks base method.
func (m *MockINotificationPushService) NewChatMsg(expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewChatMsg", expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewChatMsg indicates an expected call of NewChatMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewChatMsg(expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewChatMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewChatMsg), expoToken, toPhone)
}

// NewCompletedCustomerContract mocks base method.
func (m *MockINotificationPushService) NewCompletedCustomerContract(contractID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCompletedCustomerContract", contractID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewCompletedCustomerContract indicates an expected call of NewCompletedCustomerContract.
func (mr *MockINotificationPushServiceMockRecorder) NewCompletedCustomerContract(contractID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCompletedCustomerContract", reflect.TypeOf((*MockINotificationPushService)(nil).NewCompletedCustomerContract), contractID, expoToken, toPhone)
}

// NewCustomerAdditionalPaymentMsg mocks base method.
func (m *MockINotificationPushService) NewCustomerAdditionalPaymentMsg(contractID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCustomerAdditionalPaymentMsg", contractID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewCustomerAdditionalPaymentMsg indicates an expected call of NewCustomerAdditionalPaymentMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewCustomerAdditionalPaymentMsg(contractID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCustomerAdditionalPaymentMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewCustomerAdditionalPaymentMsg), contractID, expoToken, toPhone)
}

// NewInactiveCarMsg mocks base method.
func (m *MockINotificationPushService) NewInactiveCarMsg(carID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewInactiveCarMsg", carID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewInactiveCarMsg indicates an expected call of NewInactiveCarMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewInactiveCarMsg(carID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewInactiveCarMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewInactiveCarMsg), carID, expoToken, toPhone)
}

// NewReceivingPaymentMsg mocks base method.
func (m *MockINotificationPushService) NewReceivingPaymentMsg(amount int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReceivingPaymentMsg", amount, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewReceivingPaymentMsg indicates an expected call of NewReceivingPaymentMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewReceivingPaymentMsg(amount, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReceivingPaymentMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewReceivingPaymentMsg), amount, expoToken, toPhone)
}

// NewRejectCarMsg mocks base method.
func (m *MockINotificationPushService) NewRejectCarMsg(carID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRejectCarMsg", carID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewRejectCarMsg indicates an expected call of NewRejectCarMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewRejectCarMsg(carID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRejectCarMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewRejectCarMsg), carID, expoToken, toPhone)
}

// NewRejectPartnerContractMsg mocks base method.
func (m *MockINotificationPushService) NewRejectPartnerContractMsg(carID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRejectPartnerContractMsg", carID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewRejectPartnerContractMsg indicates an expected call of NewRejectPartnerContractMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewRejectPartnerContractMsg(carID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRejectPartnerContractMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewRejectPartnerContractMsg), carID, expoToken, toPhone)
}

// NewRejectRentingCarRequestMsg mocks base method.
func (m *MockINotificationPushService) NewRejectRentingCarRequestMsg(expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRejectRentingCarRequestMsg", expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewRejectRentingCarRequestMsg indicates an expected call of NewRejectRentingCarRequestMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewRejectRentingCarRequestMsg(expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRejectRentingCarRequestMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewRejectRentingCarRequestMsg), expoToken, toPhone)
}

// NewReturnCollateralAssetMsg mocks base method.
func (m *MockINotificationPushService) NewReturnCollateralAssetMsg(contractID int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReturnCollateralAssetMsg", contractID, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewReturnCollateralAssetMsg indicates an expected call of NewReturnCollateralAssetMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewReturnCollateralAssetMsg(contractID, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReturnCollateralAssetMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewReturnCollateralAssetMsg), contractID, expoToken, toPhone)
}

// NewWarningCountMsg mocks base method.
func (m *MockINotificationPushService) NewWarningCountMsg(carID, curCount, maxCount int, expoToken, toPhone string) *PushMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWarningCountMsg", carID, curCount, maxCount, expoToken, toPhone)
	ret0, _ := ret[0].(*PushMessage)
	return ret0
}

// NewWarningCountMsg indicates an expected call of NewWarningCountMsg.
func (mr *MockINotificationPushServiceMockRecorder) NewWarningCountMsg(carID, curCount, maxCount, expoToken, toPhone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWarningCountMsg", reflect.TypeOf((*MockINotificationPushService)(nil).NewWarningCountMsg), carID, curCount, maxCount, expoToken, toPhone)
}

// Push mocks base method.
func (m_2 *MockINotificationPushService) Push(accID int, m *PushMessage) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Push", accID, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push.
func (mr *MockINotificationPushServiceMockRecorder) Push(accID, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockINotificationPushService)(nil).Push), accID, m)
}
