package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentHandler_GeneratePaymentURL(t *testing.T) {
	t.Parallel()

	paymentService := NewVnPayService(TestConfig.VNPay)
	url, err := paymentService.GeneratePaymentURL(1, 100_000, "1")
	require.NoError(t, err)
	fmt.Println(url)
}
