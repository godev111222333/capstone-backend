package misc

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

type HashVerifier struct {
}

func NewHashVerifier() *HashVerifier {
	return &HashVerifier{}
}

func (h *HashVerifier) Hash(src string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

func (h *HashVerifier) Compare(hashedPassword string, password string) error {
	bzHash, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(bzHash, []byte(password))
}
