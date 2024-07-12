package store

import (
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{db: db}
}

func (s *NotificationStore) Create(n *model.Notification) error {
	if err := s.db.Create(n).Error; err != nil {
		fmt.Printf("NotificationStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *NotificationStore) GetByAcctID(acctID, offset, limit int) ([]*model.Notification, error) {
	if limit == 0 {
		limit = 1000
	}
	var res []*model.Notification
	if err := s.db.Model(model.Notification{}).Where("account_id = ? and status = ?", acctID, string(model.NotificationStatusActive)).
		Order("id desc").
		Offset(offset).
		Limit(limit).Scan(&res).Error; err != nil {
		fmt.Printf("NotificationStore: GetByAcctID %v\n", err)
		return nil, err
	}

	return res, nil
}
