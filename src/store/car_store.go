package store

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
)

const BufferAtHomeTime = 2 * time.Hour
const BufferAtGarage = time.Hour

var NotAvailableForRentStatuses = []string{
	string(model.CustomerContractStatusOrdered),
	string(model.CustomerContractStatusAppraisingCarApproved),
	string(model.CustomerContractStatusAppraisingCarRejected),
	string(model.CustomerContractStatusRenting),
	string(model.CustomerContractStatusReturnedCar),
	string(model.CustomerContractStatusAppraisedReturnCar),
	string(model.CustomerContractStatusPendingResolve),
	string(model.CustomerContractStatusCompleted),
}

type CarStore struct {
	db *gorm.DB
}

func NewCarStore(db *gorm.DB) *CarStore {
	return &CarStore{db: db}
}

func (s *CarStore) Create(car *model.Car) error {
	if err := s.db.Create(car).Error; err != nil {
		fmt.Printf("CarStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CarStore) CreateBatch(cars []*model.Car) error {
	if err := s.db.Create(cars).Error; err != nil {
		fmt.Printf("CarStore: CreateBatch %v\n", err)
		return err
	}
	return nil
}

func (s *CarStore) SearchCars(offset, limit int, status model.CarStatus, searchParam string) ([]*model.Car, error) {
	if limit == 0 {
		limit = 1000
	}

	res := make([]*model.Car, 0)

	if len(searchParam) == 0 {
		if status == model.CarStatusNoFilter {
			if err := s.db.
				Offset(offset).Limit(limit).
				Order("ID desc").Preload("Account").Preload("PartnerContractRule").Preload("CarModel").Find(&res).Error; err != nil {
				fmt.Printf("CarStore: SearchCars %v\n", err)
				return nil, err
			}
		} else {
			if err := s.db.Where("status like ?", string(status)+"%").
				Offset(offset).Limit(limit).
				Order("ID desc").Preload("Account").Preload("PartnerContractRule").Preload("CarModel").Find(&res).Error; err != nil {
				fmt.Printf("CarStore: SearchCars %v\n", err)
				return nil, err
			}
		}

		return res, nil
	}

	joinModel := []struct {
		Car                 *model.Car                 `gorm:"embedded"`
		Account             *model.Account             `gorm:"embedded"`
		CarModel            *model.CarModel            `gorm:"embedded"`
		PartnerContractRule *model.PartnerContractRule `gorm:"embedded"`
	}{}
	var err error

	if status != model.CarStatusNoFilter {
		rawSql := `select *
from cars
         join accounts on cars.partner_id = accounts.id join car_models on cars.car_model_id = car_models.id join partner_contract_rules on partner_contract_rules.id = cars.partner_contract_rule_id
where cars.status like ?
  and (car_models.brand = ? or car_models.model = ? or cars.license_plate = ? or concat(accounts.last_name, ' ', accounts.first_name) like ?) order by cars.id desc offset ? limit ?`
		err = s.db.Raw(rawSql, string(status)+"%", searchParam, searchParam, searchParam, likeQuery(searchParam), offset, limit).Scan(&joinModel).Error
	} else {
		rawSql := `select *
from cars
         join accounts on cars.partner_id = accounts.id join car_models on cars.car_model_id = car_models.id join partner_contract_rules on partner_contract_rules.id = cars.partner_contract_rule_id
where car_models.brand = ? or car_models.model = ? or cars.license_plate = ? or concat(accounts.last_name, ' ', accounts.first_name) like ? order by cars.id desc offset ? limit ?`
		err = s.db.Raw(rawSql, searchParam, searchParam, searchParam, likeQuery(searchParam), offset, limit).Scan(&joinModel).Error
	}

	if err != nil {
		fmt.Printf("CarStore: SearchCars %v\n", err)
		return nil, err
	}

	searchRes := make([]*model.Car, len(joinModel))
	for i, record := range joinModel {
		searchRes[i] = record.Car
		searchRes[i].CarModel = *record.CarModel
		searchRes[i].PartnerContractRule = *record.PartnerContractRule
		searchRes[i].Account = record.Account
	}

	return searchRes, nil
}

func (s *CarStore) GetByID(id int) (*model.Car, error) {
	res := &model.Car{}
	if err := s.db.Where("id = ?", id).Preload("Account").Preload("CarModel").Preload("PartnerContractRule").First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CarStore) GetByPartner(partnerID, offset, limit int, status model.CarStatus) ([]*model.Car, error) {
	var res []*model.Car
	if limit == 0 {
		limit = 1000
	}

	var row *gorm.DB
	if status == model.CarStatusNoFilter {
		row = s.db.Where("partner_id = ?", partnerID).Preload("Account").Preload("CarModel").Preload("PartnerContractRule").Order("id desc").Offset(offset).Limit(limit).Find(&res)
	} else {
		row = s.db.Where("partner_id = ? and status like ?", partnerID, "%"+string(status)+"%").Preload("Account").Preload("CarModel").Preload("PartnerContractRule").Order("id desc").Offset(offset).Limit(limit).Find(&res)
	}
	if err := row.Error; err != nil {
		fmt.Printf("CarStore: GetByPartner %v\n", err)
		return nil, err
	}

	if row.RowsAffected == 0 {
		return []*model.Car{}, nil
	}

	return res, nil
}

func (s *CarStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(&model.Car{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CarStore: Update %v\n", err)
		return err
	}
	return nil
}

func (s *CarStore) UpdateByPartnerID(tx *gorm.DB, partnerID int, values map[string]interface{}) error {
	if err := tx.Model(&model.Car{}).Where("partner_id = ?", partnerID).Updates(values).Error; err != nil {
		fmt.Printf("CarStore: UpdateByPartnerID %v\n", err)
		return err
	}

	return nil
}

func (s *CarStore) CountByStatus(status model.CarStatus) (int, error) {
	var count int64
	var err error
	if status == model.CarStatusNoFilter {
		err = s.db.Model(model.Car{}).Count(&count).Error
	} else {
		err = s.db.Model(model.Car{}).Where("status = ?", string(status)).Count(&count).Error
	}

	if err != nil {
		fmt.Printf("CarStore: CountByStatus %v\n", err)
		return -1, err
	}

	return int(count), nil
}

func (s *CarStore) CountBySeats(seatType int, parkingLot model.ParkingLot, statuses []model.CarStatus) (int, error) {
	res := struct {
		Count int `json:"count"`
	}{}
	raw := `select count(*)
				from cars inner join car_models cm on cars.car_model_id = cm.id
				where cm.number_of_seats = ? and cars.status in ? and cars.parking_lot = ?`

	if err := s.db.Raw(raw, seatType, statuses, string(parkingLot)).Scan(&res).Error; err != nil {
		fmt.Printf("CarStore: CountBySeats %v\n", err)
		return -1, err
	}

	return res.Count, nil
}

func (s *CarStore) CountByParkingLot(parkingLot model.ParkingLot, status model.CarStatus) (int, error) {
	var count int64
	if err := s.db.Model(model.Car{}).Where("parking_lot = ? and status = ?", string(parkingLot), string(status)).Count(&count).Error; err != nil {
		fmt.Printf("CarStore: CountByParkingLot %v\n", err)
		return -1, err
	}

	return int(count), nil
}

func (s *CarStore) FindCars(
	startDate, endDate time.Time,
	optionParams map[string]interface{},
) ([]*model.Car, error) {
	optionsSQL := make([]string, 0, len(optionParams))
	if values, ok := optionParams["brands"].([]string); ok {
		optionsSQL = append(optionsSQL, fmt.Sprintf("brand in %s", quoteString(values)))
	}
	if values, ok := optionParams["fuels"].([]string); ok {
		optionsSQL = append(optionsSQL, fmt.Sprintf("fuel in %s", quoteString(values)))
	}
	if values, ok := optionParams["motions"].([]string); ok {
		optionsSQL = append(optionsSQL, fmt.Sprintf("motion in %s", quoteString(values)))
	}
	if values, ok := optionParams["number_of_seats"].([]int); ok {
		optionsSQL = append(optionsSQL, fmt.Sprintf("number_of_seats in %s", quoteInt(values)))
	}

	opt := strings.Join(optionsSQL, " and ")
	if len(opt) > 0 {
		opt = opt + ` and`
	}

	searchByParkingLot := func(parkingLot model.ParkingLot, startDate, endDate time.Time) ([]*model.Car, error) {
		cars := make([]*model.CarJoinCarModel, 0)
		rawSql := `select cars.*, cars.id as car_id, cm.*
				from cars inner join car_models cm on cars.car_model_id = cm.id 
				where ` + opt + ` cars.parking_lot = ? and cars.status = ? and cars.id not in (select car_id from customer_contracts where (customer_contracts.start_date >= ? and ? >= customer_contracts.start_date and (customer_contracts.status in ?)) or cars.end_date < ?)`
		if err := s.db.Raw(rawSql, string(parkingLot), string(model.CarStatusActive), startDate, endDate, NotAvailableForRentStatuses, endDate).Preload("CarModel").Find(&cars).Error; err != nil {
			fmt.Printf("CarStore: FindCars %v\n", err)
			return nil, err
		}

		cars2 := make([]*model.CarJoinCarModel, 0)
		rawSql = `select cars.*, cars.id as car_id, cm.*
				from cars inner join car_models cm on cars.car_model_id = cm.id
				where ` + opt + ` cars.parking_lot = ? and cars.status = ? and cars.id not in (select car_id from customer_contracts where (? >= customer_contracts.start_date and customer_contracts.end_date >= ? and (customer_contracts.status in ?)) or cars.end_date < ?)`
		if err := s.db.Raw(rawSql, string(parkingLot), string(model.CarStatusActive), startDate, startDate, NotAvailableForRentStatuses, endDate).Preload("CarModel").Find(&cars2).Error; err != nil {
			fmt.Printf("CarStore: FindCars %v\n", err)
			return nil, err
		}

		res := takeDuplicatedCars(cars, cars2)
		resCars := make([]*model.Car, len(res))
		for i, r := range res {
			resCars[i] = r.ToCar()
		}

		return resCars, nil
	}

	validCarsAtHome, err := searchByParkingLot(model.ParkingLotHome, startDate.Add(-BufferAtHomeTime), endDate.Add(BufferAtHomeTime))
	if err != nil {
		return nil, err
	}

	validCarsAtGarage, err := searchByParkingLot(model.ParkingLotGarage, startDate.Add(-BufferAtGarage), endDate.Add(BufferAtGarage))
	if err != nil {
		return nil, err
	}

	searchParkingLots, ok := optionParams["parking_lots"].([]string)
	if !ok {
		return append(validCarsAtGarage, validCarsAtHome...), nil
	}

	resCars := make([]*model.Car, 0)
	for _, p := range searchParkingLots {
		if p == string(model.ParkingLotGarage) {
			resCars = append(resCars, validCarsAtGarage...)
		}

		if p == string(model.ParkingLotHome) {
			resCars = append(resCars, validCarsAtHome...)
		}
	}

	return resCars, nil
}

func takeDuplicatedCars(cars1, cars2 []*model.CarJoinCarModel) []*model.CarJoinCarModel {
	m1, m2, checkedID := make(map[int]struct{}, 0), make(map[int]struct{}, 0), make(map[int]struct{}, 0)
	for _, c := range cars1 {
		m1[c.CarID] = struct{}{}
	}
	for _, c := range cars2 {
		m2[c.CarID] = struct{}{}
	}

	res := []*model.CarJoinCarModel{}
	cars := append(cars1, cars2...)
	for _, c := range cars {
		if _, ok := m1[c.CarID]; !ok {
			continue
		}
		if _, ok := m2[c.CarID]; !ok {
			continue
		}

		if _, ok := checkedID[c.CarID]; ok {
			continue
		}
		res = append(res, c)
		checkedID[c.CarID] = struct{}{}
	}

	return res
}

func quoteString(src []string) string {
	if len(src) == 0 {
		return ""
	}

	for i, s := range src {
		src[i] = "'" + s + "'"
	}

	return "(" + strings.Join(src, ",") + ")"
}

func quoteInt(src []int) string {
	if len(src) == 0 {
		return ""
	}

	res := []string{}
	for _, s := range src {
		res = append(res, strconv.Itoa(s))
	}

	return "(" + strings.Join(res, ",") + ")"
}

func likeQuery(param string) string {
	return "%" + param + "%"
}
