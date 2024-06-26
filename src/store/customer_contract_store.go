package store

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CustomerContractStore struct {
	db *gorm.DB
}

func NewCustomerContractStore(db *gorm.DB) *CustomerContractStore {
	return &CustomerContractStore{db: db}
}

func (s *CustomerContractStore) Create(c *model.CustomerContract) error {
	if err := s.db.Create(c).Error; err != nil {
		fmt.Printf("CustomerContractStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerContractStore) IsOverlap(carID int, desiredStartDate time.Time, desiredEndDate time.Time) (bool, error) {
	// 1. case start_date >= desiredStartDate
	counter := struct {
		Count int `json:"count"`
	}{}
	raw := `select count(*) as count from customer_contracts where car_id = ? and start_date >= ? and ? >= start_date and (status = 'ordered' or status = 'renting' or status = 'completed')`
	if err := s.db.Raw(raw, carID, desiredStartDate, desiredEndDate).Scan(&counter).Error; err != nil {
		fmt.Printf("CustomerContractStore: CheckOverlap case 1 %v\n", err)
		return false, err
	}

	if counter.Count > 0 {
		return true, nil
	}

	// 2. case desiredStartDate >= start_date
	raw = `select count(*) as count from customer_contracts where car_id = ? and ? >= start_date and end_date >= ? and (status = 'ordered' or status = 'renting' or status = 'completed')`
	if err := s.db.Raw(raw, carID, desiredStartDate, desiredStartDate).Scan(&counter).Error; err != nil {
		fmt.Printf("CustomerContractStore: CheckOverlap case 2 %v\n", err)
		return false, err
	}

	return counter.Count > 0, nil
}

func (s *CustomerContractStore) FindByID(id int) (*model.CustomerContract, error) {
	res := &model.CustomerContract{}
	if err := s.db.Where("id = ?", id).Preload("Customer").Preload("Car").Find(res).Error; err != nil {
		fmt.Printf("CustomerContractStore: FindByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(&model.CustomerContract{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CustomerContractStore: Update %v\n", err)
		return err
	}
	return nil
}

func (s *CustomerContractStore) GetByCustomerID(cusID int, status model.CustomerContractStatus, offset, limit int) ([]*model.CustomerContract, error) {
	var res []*model.CustomerContract
	if limit == 0 {
		limit = 1000
	}

	if status == model.CustomerContractStatusNoFilter {
		if err := s.db.Where("customer_id = ?", cusID).Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: GetByCustomerID %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("customer_id = ? and status like ?", cusID, "%"+string(status)+"%").Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: GetByCustomerID %v\n", err)
			return nil, err
		}
	}

	return res, nil
}

func (s *CustomerContractStore) GetByStatus(
	status model.CustomerContractStatus,
	offset,
	limit int,
	searchParam string,
) ([]*model.CustomerContract, int, error) {
	var err error
	if limit == 0 {
		limit = 1000
	}

	if len(searchParam) == 0 {
		res := []*model.CustomerContract{}
		count := int64(-1)
		if status == model.CustomerContractStatusNoFilter {
			err = s.db.Preload("Customer").Preload("Car").Preload("Car.CarModel").Offset(offset).Limit(limit).Find(&res).Error
			if err == nil {
				if err := s.db.Model(&model.CustomerContract{}).Count(&count).Error; err != nil {
					fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
					return nil, -1, err
				}
			}
		} else {
			err = s.db.Where("status = ?", string(status)).Preload("Customer").Preload("Car").Preload("Car.CarModel").Offset(offset).Limit(limit).Find(&res).Error
			if err == nil {
				if err := s.db.Model(&model.CustomerContract{}).Where("status = ?", string(status)).Count(&count).Error; err != nil {
					fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
					return nil, -1, err
				}
			}
		}

		if err != nil {
			fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
			return nil, -1, err
		}

		return res, int(count), nil
	}

	joinModel := []struct {
		CustomerContract *model.CustomerContract `gorm:"embedded"`
		Account          *model.Account          `gorm:"embedded"`
		Car              *model.Car              `gorm:"embedded"`
	}{}
	counter := struct {
		Count int
	}{}
	if status == model.CustomerContractStatusNoFilter {
		rawSql := `
select *
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where concat(accounts.last_name, ' ', accounts.first_name) like ? or cars.license_plate = ? order by customer_contracts.id desc offset ? limit ?
`
		err = s.db.Raw(rawSql, likeQuery(searchParam), searchParam, offset, limit).Scan(&joinModel).Error
		if err == nil {
			countSql := `
select count(*) as count
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where concat(accounts.last_name, ' ', accounts.first_name) like ? or cars.license_plate = ? group by customer_contracts.id;
  `
			err = s.db.Raw(countSql, likeQuery(searchParam), searchParam).Scan(&counter).Error
		}
	} else {
		rawSql := `
select *
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where customer_contracts.status = ?
  and (concat(accounts.last_name, ' ', accounts.first_name) like ? or cars.license_plate = ?) order by customer_contracts.id desc offset ? limit ?
`
		err = s.db.Raw(rawSql, string(status), likeQuery(searchParam), searchParam, offset, limit).Scan(&joinModel).Error
		if err == nil {
			countSql := `
select count(*) as count
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where customer_contracts.status = ?
  and (concat(accounts.last_name, ' ', accounts.first_name) like ? or cars.license_plate = ?) group by customer_contracts.id
`
			err = s.db.Raw(countSql, string(status), likeQuery(searchParam), searchParam).Scan(&counter).Error
		}
	}

	if err != nil {
		fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
		return nil, -1, err
	}

	res := make([]*model.CustomerContract, len(joinModel))
	for index, join := range joinModel {
		res[index] = join.CustomerContract
		res[index].Car = *join.Car
		res[index].Customer = *join.Account
	}
	return res, counter.Count, nil
}

func (s *CustomerContractStore) CountByStatus(status model.CustomerContractStatus) (int, error) {
	var count int64
	var err error
	if status == model.CustomerContractStatusNoFilter {
		err = s.db.Model(model.CustomerContract{}).Count(&count).Error
	} else {
		err = s.db.Model(model.CustomerContract{}).Where("status = ?", string(status)).Count(&count).Error
	}

	if err != nil {
		fmt.Printf("CustomerContractStore: CountByStatus %v\n", err)
		return -1, err
	}

	return int(count), nil
}
