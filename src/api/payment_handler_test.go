package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentHandler_GeneratePaymentURL(t *testing.T) {
	t.Parallel()

	paymentService := NewVnPayService(TestConfig.VNPay)
	url, err := paymentService.GeneratePaymentURL(1, 10_000, strconv.Itoa(rand.Int()%1_000_000), "return_url")
	require.NoError(t, err)
	fmt.Println(url)
}
