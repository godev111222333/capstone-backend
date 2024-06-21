package store

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
)

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

func (s *CarStore) GetAll(offset, limit int, status model.CarStatus) ([]*model.Car, error) {
	if limit == 0 {
		limit = 1000
	}

	res := make([]*model.Car, 0)
	if status == model.CarStatusNoFilter {
		if err := s.db.
			Offset(offset).Limit(limit).
			Order("ID desc").Preload("Account").Preload("CarModel").Find(&res).Error; err != nil {
			fmt.Printf("CarStore: GetAll %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("status like ?", "%"+string(status)+"%").
			Offset(offset).Limit(limit).
			Order("ID desc").Preload("Account").Preload("CarModel").Find(&res).Error; err != nil {
			fmt.Printf("CarStore: GetAll %v\n", err)
			return nil, err
		}
	}

	return res, nil
}

func (s *CarStore) GetByID(id int) (*model.Car, error) {
	res := &model.Car{}
	if err := s.db.Where("id = ?", id).Preload("Account").Preload("CarModel").Find(res).Error; err != nil {
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
		row = s.db.Where("partner_id = ?", partnerID).Preload("Account").Preload("CarModel").Order("id desc").Offset(offset).Limit(limit).Find(&res)
	} else {
		row = s.db.Where("partner_id = ? and status like ?", partnerID, "%"+string(status)+"%").Preload("Account").Preload("CarModel").Order("id desc").Offset(offset).Limit(limit).Find(&res)
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
	if values, ok := optionParams["parking_lots"].([]string); ok {
		optionsSQL = append(optionsSQL, fmt.Sprintf("parking_lot in %s", quoteString(values)))
	}

	opt := strings.Join(optionsSQL, " and ")
	if len(opt) > 0 {
		opt = opt + ` and`
	}

	cars := make([]*model.CarJoinCarModel, 0)
	rawSql := `select *
				from cars inner join car_models cm on cars.car_model_id = cm.id
				where ` + opt + ` status = ? and cars.id not in (select car_id from customer_contracts where start_date >= ? and ? >= start_date and (status = 'ordered' or status = 'renting' or status = 'completed'))`
	if err := s.db.Debug().Raw(rawSql, string(model.CarStatusActive), startDate, endDate).Preload("CarModel").Find(&cars).Error; err != nil {
		fmt.Printf("CarStore: FindCars %v\n", err)
		return nil, err
	}

	cars2 := make([]*model.CarJoinCarModel, 0)
	rawSql = `select *
				from cars inner join car_models cm on cars.car_model_id = cm.id
				where ` + opt + ` status = ? and cars.id not in (select car_id from customer_contracts where ? >= start_date and end_date >= ? and (status = 'ordered' or status = 'renting' or status = 'completed'))`
	if err := s.db.Debug().Raw(rawSql, string(model.CarStatusActive), startDate, startDate).Preload("CarModel").Find(&cars2).Error; err != nil {
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

func takeDuplicatedCars(cars1, cars2 []*model.CarJoinCarModel) []*model.CarJoinCarModel {
	m1, m2, checkedID := make(map[int]struct{}, 0), make(map[int]struct{}, 0), make(map[int]struct{}, 0)
	for _, c := range cars1 {
		m1[c.ID] = struct{}{}
	}
	for _, c := range cars2 {
		m2[c.ID] = struct{}{}
	}

	res := []*model.CarJoinCarModel{}
	cars := append(cars1, cars2...)
	for _, c := range cars {
		if _, ok := m1[c.ID]; !ok {
			continue
		}
		if _, ok := m2[c.ID]; !ok {
			continue
		}

		if _, ok := checkedID[c.ID]; ok {
			continue
		}
		res = append(res, c)
		checkedID[c.ID] = struct{}{}
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
