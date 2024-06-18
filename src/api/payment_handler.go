package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
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
	BankCode    string `url:"vnp_BankCode"`
	CreatedDate string `url:"vnp_CreateDate"`
	CurrCode    string `url:"vnp_CurrCode"`
	IpAddress   string `url:"vnp_IpAddr"`
	Locale      string `url:"vnp_Locale"`
	OrderInfo   string `url:"vnp_OrderInfo"`
	OrderType   string `url:"vnp_OrderType"`
	ReturnURL   string `url:"vnp_ReturnUrl"`
	ExpireDate  string `url:"vnp_ExpireDate"`
	TxnRef      string `url:"vnp_TxnRef"`
	SecureHash  string `url:"vnp_SecureHash"`
}

var _ IPaymentService = (*VnPayService)(nil)

type IPaymentService interface {
	GeneratePaymentURL(paymentID int, amount int, data string) (string, error)
}

type VnPayService struct {
	cfg    *misc.VNPayConfig
	signer hash.Hash
}

func NewVnPayService(cfg *misc.VNPayConfig) *VnPayService {
	signer := hmac.New(sha512.New, []byte(cfg.HashSecret))
	return &VnPayService{cfg: cfg, signer: signer}
}

func (s *VnPayService) GeneratePaymentURL(paymentID, amount int, data string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.cfg.PayURL, nil)
	if err != nil {
		fmt.Printf("error when generating payment url %v\n", err)
		return "", err
	}

	now := time.Now()
	layoutyyyyMMddHHmmss := "20060201150405"
	reqBody := PayRequest{
		Version:     s.cfg.Version,
		Command:     s.cfg.Command,
		TmnCode:     s.cfg.TMNCode,
		Amount:      amount,
		BankCode:    s.cfg.BankCode,
		CreatedDate: now.Format(layoutyyyyMMddHHmmss),
		CurrCode:    "VND",
		IpAddress:   "123.123.123.123",
		Locale:      s.cfg.Locale,
		OrderInfo:   fmt.Sprintf("Thanh toan cho payment #%d tai MinhHungCar. Tong so tien %d", paymentID, amount),
		OrderType:   "other",
		ReturnURL:   s.cfg.ReturnURL,
		ExpireDate:  now.AddDate(0, 0, 1).Format(layoutyyyyMMddHHmmss),
		TxnRef:      data,
	}

	values, err := query.Values(reqBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	s.signer.Reset()
	s.signer.Write([]byte(values.Encode()))
	secureHash := hex.EncodeToString(s.signer.Sum(nil))

	req.URL.RawQuery = values.Encode() + "&vnp_SecureHash=" + secureHash
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

	c.JSON(http.StatusOK, gin.H{"status": "update payment successfully"})
}

func (s *Server) HandleVnPayReturnURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
