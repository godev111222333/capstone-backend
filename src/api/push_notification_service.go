package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var _ INotificationPushService = (*NotificationPushService)(nil)

type INotificationPushService interface {
	Push(m *PushMessage) error
	NewApproveCarRegisterMsg(carID int, expoToken, toPhone string) *PushMessage
	NewApproveCarDeliveryMsg(carID int, expoToken, toPhone string) *PushMessage
	NewRejectCarMsg(carID int, expoToken, toPhone string) *PushMessage
	NewInactiveCarMsg(carID int, expoToken, toPhone string) *PushMessage
	NewWarningCountMsg(carID, curCount, maxCount int, expoToken, toPhone string) *PushMessage
	NewChatMsg(expoToken, toPhone string) *PushMessage
	NewReceivingPaymentMsg(amount int, expoToken, toPhone string) *PushMessage
	NewRejectPartnerContractMsg(carID int, expoToken, toPhone string) *PushMessage
	NewRejectRentingCarRequestMsg(expoToken, toPhone string) *PushMessage
	NewApproveRentingCarRequestMsg(contractID int, expoToken, toPhone string) *PushMessage
	NewCustomerAdditionalPaymentMsg(contractID int, expoToken, toPhone string) *PushMessage
	NewReturnCollateralAssetMsg(contractID int, expoToken, toPhone string) *PushMessage
	NewCompletedCustomerContract(contractID int, expoToken, toPhone string) *PushMessage
}

type PushMessage struct {
	To    []string    `json:"to"`
	Title string      `json:"title,omitempty"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data,omitempty"`
}

type NotificationPushService struct {
	FrontendURL string
	ExpoURL     string
}

func NewNotificationPushService(feURL string) *NotificationPushService {
	return &NotificationPushService{
		FrontendURL: feURL,
		ExpoURL:     "https://exp.host/--/api/v2/push/send",
	}
}

func (s *NotificationPushService) Push(m *PushMessage) error {
	bz, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("NotificationPushService: Push %v\n", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.ExpoURL, bytes.NewReader(bz))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if code := resp.StatusCode; code < 200 || code > 299 {
		err := fmt.Errorf("HTTP error. Status code %d", code)
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *NotificationPushService) NewApproveCarRegisterMsg(carID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Xe của bạn đã được duyệt!",
		Body:  "Vui lòng xem và xác nhận thông tin hợp đồng!",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewApproveCarDeliveryMsg(carID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Xe giao tới garage của bạn đã được duyệt!",
		Body:  "MinhHungCar đã chấp nhận xe giao tới garage của bạn",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewRejectCarMsg(carID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Xe của bạn không được duyệt!",
		Body:  "MinhHungCar từ chối yêu cầu đăng kí xe của bạn",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewInactiveCarMsg(carID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Hợp đồng bị hủy",
		Body:  "Hợp đồng cho thuê xe của bạn đã bị hủy",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewWarningCountMsg(
	carID,
	curCount,
	maxCount int,
	expoToken,
	toPhone string,
) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Xe của bạn bị cảnh báo",
		Body:  fmt.Sprintf("Xe của bạn bị cảnh báo do đi trễ. Nếu vượt quá số lần tối đi, hợp đồng cho thuê xe sẽ bị huỷ. Số lần đi trễ hiện tại: %d, tối đa: %d", curCount, maxCount),
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewChatMsg(expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Bạn có tin nhắn mới",
		Body:  "Bạn nhận được tin nhắn mới từ MinhHungCar",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/chat", s.FrontendURL),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewReceivingPaymentMsg(amount int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Nhận tiền từ MinhHungCar",
		Body:  fmt.Sprintf("MinhHungCar thanh toán tiền tháng này là %d VNĐ", amount),
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/index", s.FrontendURL),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewRejectPartnerContractMsg(carID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Hợp đồng bị hủy",
		Body:  "Hợp đồng cho thuê xe của bạn đã bị hủy",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detail/%d", s.FrontendURL, carID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewRejectRentingCarRequestMsg(expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Yêu cầu thuê xe bị từ chối",
		Body:  "MinhHungCar đã từ chối yêu cầu thuê xe của bạn",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/trip", s.FrontendURL),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewApproveRentingCarRequestMsg(contractID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Yêu cầu thuê xe đã được chấp nhận",
		Body:  "MinhHungCar đã chấp nhận yêu cầu thuê xe của bạn",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detailTrip?contractID=%d", s.FrontendURL, contractID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewCustomerAdditionalPaymentMsg(contractID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Phát sinh thêm thanh toán",
		Body:  "MinhHungCar vừa thêm 1 khoản thanh toán cho chuyến xe của bạn",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detailTrip?contractID=%d", s.FrontendURL, contractID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewReturnCollateralAssetMsg(contractID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Hoàn trả thế chấp",
		Body:  "MinhHungCar đã hoàn trả thế chấp cho chuyến xe của bạn!",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detailTrip?contractID=%d", s.FrontendURL, contractID),
			"phone_number": toPhone,
		},
	}
}

func (s *NotificationPushService) NewCompletedCustomerContract(contractID int, expoToken, toPhone string) *PushMessage {
	return &PushMessage{
		To:    []string{expoToken},
		Title: "Hoàn thành chuyến xe",
		Body:  "Chuyến xe của bạn đã hoàn thành cùng với các chi phí đã hoàn tất thanh toán!",
		Data: map[string]interface{}{
			"screen":       fmt.Sprintf("%s/detailTrip?contractID=%d", s.FrontendURL, contractID),
			"phone_number": toPhone,
		},
	}
}
