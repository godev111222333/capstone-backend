package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/godev111222333/capstone-backend/src/model"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/google/go-querystring/query"
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
	GeneratePaymentURL(paymentID, amount int, txnRef, returnURL string) (string, error)
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
	paymentID, amount int,
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
		OrderInfo:   encodeOrderInfo(paymentID, amount),
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
		responseError(c, err)
		return
	}

	if req.ResponseCode != "00" {
		c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
		return
	}

	paymentID, _ := decodeOrderInfo(req.OrderInfo)
	if err := s.store.CustomerPaymentStore.Update(
		paymentID,
		map[string]interface{}{"status": string(model.PaymentStatusPaid)},
	); err != nil {
		c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
		return
	}

	// Update contract status to Ordered
	payment, err := s.store.CustomerPaymentStore.GetByID(paymentID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
		return
	}

	if err := s.store.CustomerContractStore.Update(
		payment.CustomerContractID,
		map[string]interface{}{"status": string(model.CustomerContractStatusOrdered)},
	); err != nil {
		c.JSON(http.StatusOK, gin.H{"RspCode": "97", "Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
}

func (s *Server) HandleVnPayReturnURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func encodeOrderInfo(paymentID, amount int) string {
	return fmt.Sprintf("%d.%d", paymentID, amount)
}

// decodeOrderInfo return paymentID, amount
func decodeOrderInfo(s string) (int, int) {
	arr := strings.Split(s, ".")
	paymentID, err := strconv.Atoi(arr[0])
	if err != nil {
		return -1, -1
	}
	amount, err := strconv.Atoi(arr[1])
	if err != nil {
		return -1, -1
	}

	return paymentID, amount
}
