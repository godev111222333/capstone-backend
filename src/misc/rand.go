package misc

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomOTP(len int) string {
	rand.Seed(time.Now().UnixNano())

	s := ""
	for i := 0; i < len; i++ {
		s += fmt.Sprintf("%d", rand.Intn(9)+1)
	}

	return s
}
