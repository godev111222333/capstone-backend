package main

import (
	"os"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type Account struct {
	ID                       int                 `json:"id" csv:"id"`
	RoleID                   model.RoleID        `json:"role_id" csv:"role_id"`
	FirstName                string              `json:"first_name" csv:"first_name"`
	LastName                 string              `json:"last_name" csv:"last_name"`
	PhoneNumber              string              `json:"phone_number" csv:"phone_number"`
	Email                    string              `json:"email" csv:"email"`
	IdentificationCardNumber string              `json:"identification_card_number" csv:"identification_card_number"`
	Password                 string              `json:"password" csv:"password"`
	AvatarURL                string              `json:"avatar_url" csv:"avatar_url"`
	DrivingLicense           string              `json:"driving_license" csv:"driving_license"`
	Status                   model.AccountStatus `json:"status" csv:"status"`
	DateOfBirth              DateTime            `json:"date_of_birth" csv:"date_of_birth"`
	BankNumber               string              `json:"bank_number" csv:"bank_number"`
	BankOwner                string              `json:"bank_owner" csv:"bank_owner"`
	BankName                 string              `json:"bank_name" csv:"bank_name"`
	QRCodeURL                string              `json:"qr_code_url" csv:"qr_code_url"`
	CreatedAt                DateTime            `json:"created_at" csv:"created_at"`
	UpdatedAt                DateTime            `json:"updated_at" csv:"updated_at"`
}

func (a *Account) ToDbAccount() *model.Account {
	return &model.Account{
		ID:                       a.ID,
		RoleID:                   a.RoleID,
		FirstName:                a.FirstName,
		LastName:                 a.LastName,
		PhoneNumber:              a.PhoneNumber,
		Email:                    a.Email,
		IdentificationCardNumber: a.IdentificationCardNumber,
		Password:                 a.Password,
		AvatarURL:                a.AvatarURL,
		DrivingLicense:           a.DrivingLicense,
		Status:                   a.Status,
		DateOfBirth:              a.DateOfBirth.Time,
		BankNumber:               a.BankNumber,
		BankOwner:                a.BankOwner,
		BankName:                 a.BankName,
		QRCodeURL:                a.QRCodeURL,
		CreatedAt:                a.CreatedAt.Time,
		UpdatedAt:                a.UpdatedAt.Time,
	}
}

func seedAccounts(dbStore *store.DbStore) error {
	accounts := make([]*Account, 0)
	accountFile, err := os.OpenFile(toFilePath(AccountsFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer accountFile.Close()

	if err := gocsv.UnmarshalFile(accountFile, &accounts); err != nil {
		return err
	}

	accts := make([]*model.Account, len(accounts))
	for i, a := range accounts {
		accts[i] = a.ToDbAccount()
	}

	return dbStore.AccountStore.CreateBatch(accts)
}
