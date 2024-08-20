package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
)

type PayRequest struct {
	Version     string `url:"vnp_Version"`
	Command     string `url:"vnp_Command"`
	TmnCode     string `url:"vnp_TmnCode"`
	Amount      int    `url:"vnp_Amount"`
	CreatedDate string `url:"vnp_CreateDate"`
	ExpireDate  string `url:"vnp_ExpireDate"`
	CurrCode    string `url:"vnp_CurrCode"`
	IpAddress   string `url:"vnp_IpAddr"`
	Locale      string `url:"vnp_Locale"`
	OrderInfo   string `url:"vnp_OrderInfo"`
	OrderType   string `url:"vnp_OrderType"`
	ReturnURL   string `url:"vnp_ReturnUrl"`
	TxnRef      string `url:"vnp_TxnRef"`
	BankCode    string `url:"vnp_BankCode"`
}

var _ IPaymentService = (*VnPayService)(nil)

type IPaymentService interface {
	GeneratePaymentURL(paymentIDs []int, amount int, txnRef, returnURL string) (string, error)
}

type VnPayService struct {
	cfg    *misc.VNPayConfig
	signer hash.Hash
}

func NewVnPayService(cfg *misc.VNPayConfig) *VnPayService {
	signer := hmac.New(sha512.New, []byte(cfg.HashSecret))
	return &VnPayService{cfg: cfg, signer: signer}
}

func (s *VnPayService) GeneratePaymentURL(
	paymentIDs []int, amount int,
	txnRef, returnURL string,
) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.cfg.PayURL, nil)
	if err != nil {
		fmt.Printf("error when generating payment url %v\n", err)
		return "", err
	}

	now := time.Now()
	layoutyyyyMMddHHmmss := "20060102150405"
	reqBody := PayRequest{
		Version:     s.cfg.Version,
		Command:     s.cfg.Command,
		TmnCode:     s.cfg.TMNCode,
		Amount:      amount * 100,
		CreatedDate: now.Format(layoutyyyyMMddHHmmss),
		ExpireDate:  now.AddDate(0, 0, 7).Format(layoutyyyyMMddHHmmss),
		CurrCode:    "VND",
		IpAddress:   "::1",
		Locale:      s.cfg.Locale,
		OrderInfo:   encodeOrderInfo(paymentIDs),
		OrderType:   "other",
		ReturnURL:   returnURL,
		TxnRef:      txnRef,
		BankCode:    s.cfg.BankCode,
	}

	values, err := query.Values(reqBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	signData := values.Encode()

	s.signer.Reset()
	s.signer.Write([]byte(signData))
	secureHash := hex.EncodeToString(s.signer.Sum(nil))

	req.URL.RawQuery = signData + "&vnp_SecureHash=" + secureHash
	return req.URL.String(), nil
}

type VnPayIPNRequest struct {
	TmnCode           string `form:"vnp_TmnCode"`
	Amount            int    `form:"vnp_Amount"`
	BankCode          string `form:"vnp_BankCode"`
	BankTranNo        string `form:"vnp_BankTranNo"`
	CardType          string `form:"vnp_CardType"`
	PayDate           string `form:"vnp_PayDate"`
	OrderInfo         string `form:"vnp_OrderInfo"`
	TransactionNo     string `form:"vnp_TransactionNo"`
	ResponseCode      string `form:"vnp_ResponseCode"`
	TransactionStatus string `form:"vnp_TransactionStatus"`
	TxnRef            string `form:"vnp_TxnRef"`
	SecureHashType    string `form:"vnp_SecureHashType"`
	SecureHash        string `form:"vnp_SecureHash"`
}

func (s *Server) HandleVnPayIPN(c *gin.Context) {
	req := VnPayIPNRequest{}
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.ResponseCode != "00" {
		c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
		return
	}

	if strings.HasPrefix(req.TxnRef, PrefixPartnerPayment) {
		s.handlePartnerPayments(c, req)
		return
	}

	paymentIDs := decodeOrderInfo(req.OrderInfo)

	var (
		commCustomerContractID int
		licensePlate           string
	)

	for _, paymentID := range paymentIDs {
		// Update contract status to Ordered
		payment, err := s.store.CustomerPaymentStore.GetByID(paymentID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
			return
		}

		commCustomerContractID = payment.CustomerContractID
		licensePlate = payment.CustomerContract.Car.LicensePlate

		if payment.PaymentType == model.PaymentTypePrePay {
			// check if this car is still available
			good, err := s.checkIfContractStillAvailable(commCustomerContractID)
			if err != nil || !good {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error or not available car for contract"})
				return
			}

			if err := s.store.CustomerContractStore.Update(
				payment.CustomerContractID,
				map[string]interface{}{"status": string(model.CustomerContractStatusOrdered)},
			); err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}

			partner, err := s.store.AccountStore.GetByID(payment.CustomerContract.Car.PartnerID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}

			msg := s.notificationPushService.NewRentingContract(
				payment.CustomerContract.CarID,
				payment.CustomerContractID,
				s.getExpoToken(partner.PhoneNumber),
				partner.PhoneNumber,
			)
			_ = s.notificationPushService.Push(partner.ID, msg)

			adminIds, err := s.store.AccountStore.GetAllAdminIDs()
			if err == nil {
				for _, id := range adminIds {
					s.adminNotificationQueue <- s.NewCustomerContractNotificationMsg(id, commCustomerContractID, licensePlate)
				}
			}

			techIds, err := s.store.AccountStore.GetAllIdsByRole(model.RoleIDTechnician)
			if err == nil {
				for _, id := range techIds {
					s.technicianNotificationQueue <- s.NewAppraisingCarOfCusContract(id, commCustomerContractID)
				}
			}
		} else if payment.PaymentType == model.PaymentTypeCollateralCash {
			contract, err := s.store.CustomerContractStore.FindByID(payment.CustomerContractID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}

			_, err = s.GenerateCustomerContractPaymentQRCode(
				contract.ID,
				payment.Amount,
				model.PaymentTypeReturnCollateralCash,
				s.feCfg.AdminReturnURL+strconv.Itoa(contract.ID), "",
			)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}
		} else if payment.PaymentType == model.PaymentTypeReturnCollateralCash {
			contract, err := s.store.CustomerContractStore.FindByID(payment.CustomerContractID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}

			acct, err := s.store.AccountStore.GetByID(contract.CustomerID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}

			_ = s.notificationPushService.Push(acct.ID, s.notificationPushService.NewReturnCollateralAssetMsg(
				payment.CustomerContractID,
				s.getExpoToken(acct.PhoneNumber),
				acct.PhoneNumber,
			))
			if err := s.store.CustomerContractStore.Update(payment.CustomerContractID, map[string]interface{}{
				"is_return_collateral_asset": true,
			}); err != nil {
				c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
				return
			}
		}
	}

	if err := s.store.CustomerPaymentStore.UpdateMulti(
		paymentIDs,
		map[string]interface{}{"status": string(model.PaymentStatusPaid)},
	); err != nil {
		c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
		return
	}

	go func() {
		adminIds, err := s.store.AccountStore.GetAllAdminIDs()
		if err == nil {
			for _, id := range adminIds {
				s.adminNotificationQueue <- s.NewCustomerContractPaymentNotificationMsg(id, commCustomerContractID, licensePlate)
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
}

func (s *Server) handlePartnerPayments(c *gin.Context, req VnPayIPNRequest) {
	paymentIDs := decodeOrderInfo(req.OrderInfo)
	if err := s.store.PartnerPaymentHistoryStore.UpdateMulti(
		paymentIDs,
		map[string]interface{}{"status": string(model.PartnerPaymentHistoryStatusPaid)},
	); err != nil {
		c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
		return
	}

	for _, paymentID := range paymentIDs {
		payment, err := s.store.PartnerPaymentHistoryStore.GetByID(paymentID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
			return
		}

		acct, err := s.store.AccountStore.GetByID(payment.PartnerID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
			return
		}

		msg := s.notificationPushService.NewReceivingPaymentMsg(
			payment.Amount,
			s.getExpoToken(acct.PhoneNumber),
			acct.PhoneNumber,
		)

		_ = s.notificationPushService.Push(acct.ID, msg)
	}

	c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
}

func (s *Server) HandleVnPayReturnURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func encodeOrderInfo(paymentIDs []int) string {
	strs := make([]string, len(paymentIDs))
	for i, pID := range paymentIDs {
		strs[i] = strconv.Itoa(pID)
	}

	return strings.Join(strs, ".")
}

// decodeOrderInfo return paymentID, amount
func decodeOrderInfo(s string) []int {
	pIDs := strings.Split(s, ".")
	res := make([]int, len(pIDs))
	for i, pID := range pIDs {
		res[i], _ = strconv.Atoi(pID)
	}

	return res
}

func (s *Server) checkIfContractStillAvailable(contractID int) (bool, error) {
	contract, err := s.store.CustomerContractStore.FindByID(contractID)
	if err != nil {
		return false, err
	}

	car := contract.Car
	foundCars, err := s.store.CarStore.FindCars(contract.StartDate, contract.EndDate, map[string]interface{}{
		"brands":          []string{car.CarModel.Brand},
		"fuels":           []string{string(car.Fuel)},
		"motions":         []string{string(car.Motion)},
		"number_of_seats": []int{car.CarModel.NumberOfSeats},
	})
	if err != nil {
		return false, err
	}

	for _, fCar := range foundCars {
		if fCar.ID == car.ID {
			return true, nil
		}
	}

	return false, nil
}
