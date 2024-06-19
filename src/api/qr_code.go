package api

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(url string) ([]byte, error) {
	image, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Printf("error when generating qr code %v\n", err)
		return nil, err
	}

	return image, nil
}
