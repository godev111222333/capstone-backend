package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateQRCode(t *testing.T) {
	t.Skip()

	url := "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html?vnp_Amount=1000000&vnp_BankCode=NCB&vnp_Command=pay&vnp_CreateDate=20240619131144&vnp_CurrCode=VND&vnp_IpAddr=%3A%3A1&vnp_Locale=vn&vnp_OrderInfo=Thanh+toan+cho+payment+%231.+So+Tien%3A+10000&vnp_OrderType=other&vnp_ReturnUrl=http%3A%2F%2F0.0.0.0%3A9876%2Fvnpay%2Freturn_url&vnp_TmnCode=UPUEB83F&vnp_TxnRef=367561&vnp_Version=2.1.0&vnp_SecureHash=866645d68496822daf0e504ee908afc14fa9ddf99b090e65e86db8c64b202607d36cb296e8f9ab81a9fdb2a51837b30d8bfd0a0b2ff3134a4b5a81c4b6fc51ab"
	image, err := GenerateQRCode(url)
	require.NoError(t, err)

	f, err := os.OpenFile("../../etc/qr_code.png", os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	defer f.Close()

	_, err = f.Write(image)
	require.NoError(t, err)
}
