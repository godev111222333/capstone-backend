package store

import (
	"fmt"
	"strings"
	"time"

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

func (s *AccountStore) CreateBatch(accts []*model.Account) error {
	if err := s.db.Create(accts).Error; err != nil {
		fmt.Printf("AccountStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *AccountStore) Update(accountID int, values map[string]interface{}) error {
	if err := s.db.Model(&model.Account{ID: accountID}).Updates(values).Error; err != nil {
		fmt.Printf("AccountStore: Update %v\n", err)
		return err
	}

	return nil
}

func (s *AccountStore) UpdateTx(tx *gorm.DB, accountID int, values map[string]interface{}) error {
	if err := tx.Model(&model.Account{ID: accountID}).Updates(values).Error; err != nil {
		fmt.Printf("AccountStore: UpdateTx %v\n", err)
		return err
	}

	return nil
}

func (s *AccountStore) GetByPhoneNumber(phoneNumber string) (*model.Account, error) {
	res := &model.Account{}
	if err := s.db.Where("phone_number = ?", phoneNumber).Preload("Role").First(res).Error; err != nil {
		fmt.Printf("AccountStore: GetByPhoneNumber %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *AccountStore) GetByID(id int) (*model.Account, error) {
	res := &model.Account{}
	if err := s.db.Where("id = ?", id).Preload("Role").First(res).Error; err != nil {
		fmt.Printf("AccountStore: GetByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *AccountStore) Get(status model.AccountStatus, role string, searchParam string, offset, limit int) ([]*model.Account, error) {
	if limit == 0 {
		limit = 1000
	}

	rawSql := `select accounts.*, r.* from accounts join roles r on accounts.role_id = r.id where role_name != 'admin'`

	statusQuery, roleQuery, searchQuery := "", "", ""
	if status != model.AccountStatusNoFilter {
		statusQuery = fmt.Sprintf(`status = '%s'`, string(status))
	}

	if len(role) > 0 {
		roleQuery = fmt.Sprintf(`role_name = '%s'`, role)
	}

	if len(searchParam) > 0 {
		searchQuery = fmt.Sprintf(` (first_name like '%s' or last_name like '%s' or CONCAT(last_name, ' ', first_name) like '%s' or phone_number = '%s' or email = '%s')`, "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%", searchParam, searchParam)
	}

	if len(statusQuery)+len(roleQuery)+len(searchQuery) > 0 {
		rawSql = rawSql + ` and `
	}

	combinedQuery := []string{}
	for _, str := range []string{statusQuery, roleQuery, searchQuery} {
		if len(str) > 0 {
			combinedQuery = append(combinedQuery, str)
		}
	}

	combined := strings.Join(combinedQuery, " and ")
	combined += fmt.Sprintf(` ORDER BY accounts.id OFFSET %d LIMIT %d`, offset, limit)

	var joinModel []struct {
		Account *model.Account `gorm:"embedded"`
		Role    *model.Role    `gorm:"embedded"`
	}
	if err := s.db.Raw(rawSql + combined).Scan(&joinModel).Error; err != nil {
		fmt.Printf("AccountStore: Get %v\n", err)
		return nil, err
	}

	var res []*model.Account
	for _, acct := range joinModel {
		acct.Account.Role = acct.Role
		res = append(res, acct.Account)
	}

	return res, nil
}

func (s *AccountStore) CountActiveByRole(roleID model.RoleID, backoff time.Duration) (int, error) {
	var count int64
	err := s.db.Model(model.Account{}).
		Where("role_id = ? and status = ? and created_at >= ?", roleID, string(model.AccountStatusActive), time.Now().Add(-backoff)).
		Count(&count).Error
	if err != nil {
		fmt.Printf("AccountStore: CountActiveByRole %v\n", err)
		return -1, err
	}

	return int(count), nil
}

func (s *AccountStore) GetAllAdminIDs() ([]int, error) {
	return s.GetAllIdsByRole(model.RoleIDAdmin)
}

func (s *AccountStore) GetAllIdsByRole(roleID model.RoleID) ([]int, error) {
	var res []struct {
		ID int `json:"id"`
	}
	if err := s.db.Model(model.Account{}).
		Select("id").
		Where("role_id = ?", roleID).
		Scan(&res).Error; err != nil {
		fmt.Printf("AccountStore: GetAllAdminIDs %v\n", err)
		return nil, err
	}

	ids := make([]int, len(res))
	for i, r := range res {
		ids[i] = r.ID
	}

	return ids, nil
}
