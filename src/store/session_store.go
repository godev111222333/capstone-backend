package store

import (
	"errors"
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

func NewSessionStore(db *gorm.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) Create(session *model.Session) error {
	if err := s.db.Create(session).Error; err != nil {
		fmt.Printf("SessionStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *SessionStore) GetSession(id uuid.UUID) (*model.Session, error) {
	res := &model.Session{}
	if err := s.db.Where("id = ?", id).First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		fmt.Printf("SessionStore: GetSession %v\n", err)
		return nil, err
	}
	return res, nil
}
