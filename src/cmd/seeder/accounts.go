package seeder

import (
	"os"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type Account struct {
	ID                       int                 `csv:"id"`
	RoleID                   model.RoleID        `csv:"role_id"`
	FirstName                string              `csv:"first_name"`
	LastName                 string              `csv:"last_name"`
	PhoneNumber              string              `csv:"phone_number"`
	Email                    string              `csv:"email"`
	IdentificationCardNumber string              `csv:"identification_card_number"`
	Password                 string              `csv:"password"`
	AvatarURL                string              `csv:"avatar_url"`
	DrivingLicense           string              `csv:"driving_license"`
	Status                   model.AccountStatus `csv:"status"`
	DateOfBirth              DateTime            `csv:"date_of_birth"`
	BankNumber               string              `csv:"bank_number"`
	BankOwner                string              `csv:"bank_owner"`
	BankName                 string              `csv:"bank_name"`
	QRCodeURL                string              `csv:"qr_code_url"`
	CreatedAt                DateTime            `csv:"created_at"`
	UpdatedAt                DateTime            `csv:"updated_at"`
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

func SeedAccounts(dbStore *store.DbStore) error {
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
