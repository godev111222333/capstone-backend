package service

import (
	"time"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/redis/go-redis/v9"
)

const IncomingRentingKey = "incoming_renting"

type BackgroundService struct {
	cfg   *misc.BackgroundJobConfig
	db    *store.DbStore
	cache *redis.Client
}

func NewBackgroundService(cfg *misc.BackgroundJobConfig, db *store.DbStore, cache *redis.Client) *BackgroundService {
	return &BackgroundService{cfg, db, cache}
}

func (s *BackgroundService) ScanIncomingRentingContracts() {
	ticker := time.NewTicker(s.cfg.RentingBackoff)
	for {
		select {
		case <-ticker.C:
			s.processRentingContracts()
		}
	}
}

func (s *BackgroundService) processRentingContracts() {
	return
}
