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
		Amount:      amount * 100,
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
	fmt.Println(secureHash)
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

	c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "success"})
}

func (s *Server) HandleVnPayReturnURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

//https://sandbox.vnpayment.vn/paymentv2/vpcpay.html?vnp_Amount=1806000&vnp_Command=pay&vnp_CreateDate=20210801153333&vnp_CurrCode=VND&vnp_IpAddr=127.0.0.1&vnp_Locale=vn&vnp_OrderInfo=Thanh+toan+don+hang+%3A5&vnp_OrderType=other&vnp_ReturnUrl=https%3A%2F%2Fdomainmerchant.vn%2FReturnUrl&vnp_TmnCode=DEMOV210&vnp_TxnRef=5&vnp_Version=2.1.0&vnp_SecureHash=3e0d61a0c0534b2e36680b3f7277743e8784cc4e1d68fa7d276e79c23be7d6318d338b477910a27992f5057bb1582bd44bd82ae8009ffaf6d141219218625c42
//https://sandbox.vnpayment.vn/paymentv2/vpcpay.html?vnp_Amount=10000000&vnp_Command=pay&vnp_CreateDate=20241806211631&vnp_CurrCode=VND&vnp_ExpireDate=20241906211631&vnp_IpAddr=123.123.123.123&vnp_Locale=vn&vnp_OrderInfo=Thanh+toan+cho+payment+%231+tai+MinhHungCar.+Tong+so+tien+100000&vnp_OrderType=other&vnp_ReturnUrl=https%3A%2F%2Fminhhungcar.xyz%2Fvnpay%2Freturn_url&vnp_SecureHash=&vnp_TmnCode=UPUEB83F&vnp_TxnRef=1&vnp_Version=2.1.0&vnp_SecureHash=3eaab587793416839df5aed76561451c561d80187d95028c29b8261ad27ac0fa8eb1e39903149254f67db59968464e90853994a56ba0f7389e1a55ec1af56d4c
