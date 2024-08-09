package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

const PendingPartnerApprovalKey = "Pending_Partner_Approval"

type BackgroundService struct {
	cfg                  *misc.BackgroundJobConfig
	db                   *store.DbStore
	cache                *redis.Client
	newPartnerApprovalCh chan int
}

func NewBackgroundService(
	cfg *misc.BackgroundJobConfig,
	db *store.DbStore,
	cache *redis.Client,
	newPartnerApprovalCh chan int,
) *BackgroundService {
	return &BackgroundService{cfg, db, cache, newPartnerApprovalCh}
}

func (s *BackgroundService) RunPartnerApprovalChecker() error {
	if err := s.bootstrapPartnerApprovalChecker(context.Background()); err != nil {
		return err
	}

	ticker := time.NewTicker(s.cfg.CheckWaitingPartnerApprovalInterval)
	for {
		select {
		case <-ticker.C:
			_ = s.processWaitingPartnerApproval()
		case contractID := <-s.newPartnerApprovalCh:
			_ = s.appendNewPendingPartnerApproval(contractID)
		}
	}
}

func (s *BackgroundService) appendNewPendingPartnerApproval(contractID int) error {
	old, err := s.loadMap(PendingPartnerApprovalKey)
	if err != nil {
		return err
	}

	old[contractID] = struct{}{}
	return s.saveMap(PendingPartnerApprovalKey, old)
}

func (s *BackgroundService) bootstrapPartnerApprovalChecker(ctx context.Context) error {
	// reset old records
	if err := s.cache.Del(ctx, PendingPartnerApprovalKey).Err(); err != nil {
		fmt.Printf("BackgroundService: del old keys %v\n", err)
		return err
	}

	pendingContracts, _, err := s.db.CustomerContractStore.GetByStatus(
		model.CustomerContractStatusWaitingPartnerApproval,
		0,
		0,
		"",
	)
	if err != nil {
		return err
	}

	ids := make(map[int]struct{})
	for _, c := range pendingContracts {
		ids[c.ID] = struct{}{}
	}

	return s.saveMap(PendingPartnerApprovalKey, ids)
}

func (s *BackgroundService) processWaitingPartnerApproval() error {
	ids, err := s.loadMap(PendingPartnerApprovalKey)
	if err != nil {
		return err
	}

	idArr := make([]int, 0)
	for key, _ := range ids {
		idArr = append(idArr, key)
	}

	contracts, err := s.db.CustomerContractStore.FindBatchByID(idArr)
	if err != nil {
		return err
	}

	delIds := make([]int, 0)
	for _, c := range contracts {
		if time.Now().After(c.CreatedAt.Add(s.cfg.MaxPartnerWaitingApprovalTime)) {
			delIds = append(delIds, c.ID)
		}
	}

	if err := s.db.CustomerContractStore.UpdateBatch(
		delIds, map[string]interface{}{"status": model.CustomerContractStatusCancel}); err != nil {
		return err
	}

	return s.delElements(PendingPartnerApprovalKey, delIds)
}

func (s *BackgroundService) loadMap(key string) (map[int]struct{}, error) {
	value, err := s.cache.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	old := make(map[int]struct{})
	if err := json.Unmarshal([]byte(value), &old); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return old, nil
}

func (s *BackgroundService) saveMap(key string, m map[int]struct{}) error {
	bz, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := s.cache.Set(context.Background(), key, string(bz), time.Duration(0)).Err(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *BackgroundService) delElements(key string, deletedElements []int) error {
	old, err := s.loadMap(key)
	if err != nil {
		return err
	}

	for _, e := range deletedElements {
		delete(old, e)
	}
	return s.saveMap(key, old)
}
