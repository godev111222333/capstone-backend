package store

import (
	"fmt"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
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
	if err := s.db.Where("id = ?", id).Preload("Customer").Preload("Car").Preload("Car.Account").Preload("Car.CarModel").Preload("CustomerContractRule").First(res).Error; err != nil {
		fmt.Printf("CustomerContractStore: FindByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractStore) FindByCarID(
	carID int,
	status model.CustomerContractStatus,
	offset,
	limit int,
) ([]*model.CustomerContract, error) {
	if limit == 0 {
		limit = 1000
	}

	var res []*model.CustomerContract
	if status == model.CustomerContractStatusNoFilter {
		if err := s.db.Where("car_id = ?", carID).
			Preload("Customer").
			Preload("Car").
			Preload("Car.CarModel").
			Preload("Car.PartnerContractRule").
			Preload("CustomerContractRule").Order("id desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: FindByCarID %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("car_id = ? and status = ?", carID, string(status)).
			Preload("Customer").
			Preload("Car").
			Preload("Car.CarModel").
			Preload("Car.PartnerContractRule").
			Preload("CustomerContractRule").Order("id desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: FindByCarID %v\n", err)
			return nil, err
		}
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
		if err := s.db.Where("customer_id = ?", cusID).Preload("Customer").Preload("Car").Preload("Car.CarModel").Preload("CustomerContractRule").Order("start_date desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: GetByCustomerID %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("customer_id = ? and status like ?", cusID, "%"+string(status)+"%").Preload("Customer").Preload("Car").Preload("Car.CarModel").Preload("CustomerContractRule").Order("start_date desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
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
			err = s.db.Preload("Customer").Preload("Car").Preload("Car.CarModel").Preload("CustomerContractRule").Order("end_date desc").Offset(offset).Limit(limit).Find(&res).Error
			if err == nil {
				if err := s.db.Model(&model.CustomerContract{}).Count(&count).Error; err != nil {
					fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
					return nil, -1, err
				}
			}
		} else {
			err = s.db.Where("status = ?", string(status)).Preload("Customer").Preload("Car").Preload("Car.CarModel").Preload("CustomerContractRule").Order("end_date desc").Offset(offset).Limit(limit).Find(&res).Error
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
where accounts.phone_number like ? or cars.license_plate = ? order by customer_contracts.id desc offset ? limit ?
`
		err = s.db.Raw(rawSql, likeQuery(searchParam), searchParam, offset, limit).Scan(&joinModel).Error
		if err == nil {
			countSql := `
select count(*) as count
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where accounts.phone_number like ? or cars.license_plate = ? group by customer_contracts.id;
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
  and (accounts.phone_number like ? or cars.license_plate = ?) order by customer_contracts.id desc offset ? limit ?
`
		err = s.db.Raw(rawSql, string(status), likeQuery(searchParam), searchParam, offset, limit).Scan(&joinModel).Error
		if err == nil {
			countSql := `
select count(*) as count
from customer_contracts
         join accounts on customer_contracts.customer_id = accounts.id
         join cars on customer_contracts.car_id = cars.id
where customer_contracts.status = ?
  and (accounts.phone_number like ? or cars.license_plate = ?) group by customer_contracts.id
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
		res[index].Customer = join.Account
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

func (s *CustomerContractStore) CountTotalValidCustomerContracts(backoff time.Duration) (int, error) {
	var count int64
	err := s.db.Model(model.CustomerContract{}).
		Where("status != ? and start_date >= ?", string(model.CustomerContractStatusCancel), time.Now().Add(-backoff)).
		Count(&count).Error
	if err != nil {
		fmt.Printf("CustomerContractStore: CountByStatus %v\n", err)
		return -1, err
	}

	return int(count), nil
}

func (s *CustomerContractStore) SumRevenueForCompletedContracts(startTime, endTime time.Time) (float64, error) {
	sum := struct {
		Sum float64 `json:"sum"`
	}{}
	sql := `select SUM(customer_contracts.rent_price * pcr.revenue_sharing_percent / 100)
from customer_contracts join cars on customer_contracts.car_id = cars.id join partner_contract_rules pcr on cars.partner_contract_rule_id = pcr.id
where customer_contracts.status = 'completed' and customer_contracts.end_date >= ?
  and customer_contracts.end_date < ?`

	if err := s.db.Raw(sql, startTime, endTime).Scan(&sum).Error; err != nil {
		fmt.Printf("CustomerContractStore: SumRevenueForCompletedContracts %v\n", err)
		return -1, err
	}

	return sum.Sum, nil
}

type RentedCar struct {
	CarBrandModel string `json:"car_brand_model"`
	Count         int    `json:"count"`
}

func (s *CustomerContractStore) CountRentedCars(backoff time.Duration) ([]*RentedCar, error) {
	res := make([]*RentedCar, 0)
	sql := `
select concat(car_models.brand, ' ', car_models.model) as car_brand_model, count(customer_contracts.id)
from customer_contracts
         join cars on cars.id = customer_contracts.car_id
         join car_models on cars.car_model_id = car_models.id
where customer_contracts.status = 'completed' and customer_contracts.end_date >= ?
group by concat(car_models.brand, ' ', car_models.model)
order by count(customer_contracts.id) desc;

`
	if err := s.db.Raw(sql, time.Now().Add(-backoff)).Scan(&res).Error; err != nil {
		fmt.Printf("CustomerContractStore: CountRentedCars %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractStore) GetFeedbacks(offset, limit int) ([]*model.CustomerContract, int, error) {
	var res []*model.CustomerContract
	if err := s.db.Model(model.CustomerContract{}).Where("length(feedback_content) > 0 and feedback_rating > 0").Preload("Customer").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		fmt.Printf("CustomerContracStore: GetFeedbacks %v\n", err)
		return nil, 1, err
	}

	var counter int64
	if err := s.db.Model(model.CustomerContract{}).Where("length(feedback_content) > 0 and feedback_rating > 0").Count(&counter).Error; err != nil {
		fmt.Printf("CustomerContracStore: GetFeedbacks %v\n", err)
		return nil, 1, err
	}

	return res, int(counter), nil
}

func (s *CustomerContractStore) GetFeedbacksByCarID(carID, offset, limit int, status model.FeedBackStatus) ([]*model.CustomerContract, int, error) {
	if limit == 0 {
		limit = 1000
	}
	var res []*model.CustomerContract
	if err := s.db.Model(model.CustomerContract{}).
		Where("car_id = ? and length(feedback_content) > 0 and feedback_rating > 0 and feedback_status = ?", carID, string(status)).
		Order("id desc").
		Offset(offset).
		Limit(limit).
		Preload("Customer").
		Find(&res).Error; err != nil {
		fmt.Printf("CustomerContracStore: GetFeedbacks %v\n", err)
		return nil, -1, err
	}

	var counter int64
	if err := s.db.Model(model.CustomerContract{}).
		Where("car_id = ? and length(feedback_content) > 0 and feedback_rating > 0 and feedback_status = ?", carID, string(status)).
		Count(&counter).Error; err != nil {
		fmt.Printf("CustomerContracStore: GetFeedbacks %v\n", err)
		return nil, -1, err
	}

	return res, int(counter), nil
}

func (s *CustomerContractStore) GetByStatusEndTimeInRange(
	fromDate,
	toDate time.Time, status model.CustomerContractStatus,
) ([]*model.CustomerContract, error) {
	var res []*model.CustomerContract
	if err := s.db.Model(model.CustomerContract{}).Where("status = ? and end_date >= ? and end_date < ?", string(status), fromDate, toDate).
		Preload("Car").
		Preload("Car.PartnerContractRule").
		Find(&res).Error; err != nil {
		fmt.Printf("CustomerContractStore: GetByStatusEndTimeInRange %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractStore) GetAverageRating(carID int) (float64, error) {
	var avg struct {
		Avg float64 `json:"avg"`
	}
	rawSql := `
select avg(feedback_rating) from customer_contracts where status = ? and car_id = ? and feedback_rating > 0
`
	if err := s.db.Raw(rawSql, model.CustomerContractStatusCompleted, carID).Scan(&avg).Error; err != nil {
		fmt.Printf("CustomerContractStore: GetAverageRating %v\n", err)
		return -1, err
	}

	if avg.Avg == 0.0 {
		return 5.0, nil
	}

	return avg.Avg, nil
}

func (s *CustomerContractStore) GetTotalCompletedContracts(carID int) (int, error) {
	var count int64
	err := s.db.Model(model.CustomerContract{}).
		Where("car_id = ? and status = ?", carID, string(model.CustomerContractStatusCompleted)).
		Count(&count).Error
	if err != nil {
		fmt.Printf("CustomerContractStore: GetTotalCompletedContracts %v\n", err)
		return -1, err
	}

	return int(count), nil
}

func (s *CustomerContractStore) GetIncomingRentingCustomerContracts(
	backoff time.Duration,
) ([]*model.CustomerContract, error) {
	res := make([]*model.CustomerContract, 0)
	now := time.Now()
	if err := s.db.Where("start_date >= ? and status = ?", now.Add(-backoff), model.CustomerContractStatusOrdered).Find(&res).Error; err != nil {
		fmt.Printf("CustomerContractStore: GetIncomingRentingCustomerContracts %v\n", err)
		return nil, err
	}

	return res, nil
}
