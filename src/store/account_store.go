package store

import (
	"errors"
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type AccountStore struct {
	db *gorm.DB
}

func NewAccountStore(db *gorm.DB) *AccountStore {
	return &AccountStore{db: db}
}

func (s *AccountStore) Create(acct *model.Account) error {
	if err := s.db.Create(acct).Error; err != nil {
		fmt.Printf("AccountStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *AccountStore) Update(accountID int, values map[string]interface{}) error {
	if err := s.db.Model(&model.Account{ID: accountID}).Updates(values).Error; err != nil {
		fmt.Printf("AccountStore: %v\n", err)
		return err
	}

	return nil
}

func (s *AccountStore) GetByEmail(email string) (*model.Account, error) {
	res := &model.Account{}
	if err := s.db.Where("email = ?", email).Preload("Role").First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		fmt.Printf("AccountStore: GetByEmail %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *AccountStore) GetByID(id int) (*model.Account, error) {
	res := &model.Account{}
	if err := s.db.Where("id = ?", id).Preload("Role").First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		fmt.Printf("AccountStore: GetByID %v\n", err)
		return nil, err
	}

	return res, nil
}
