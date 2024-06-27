package store

import (
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type CustomerPaymentStore struct {
	db *gorm.DB
}

func NewCustomerPaymentStore(db *gorm.DB) *CustomerPaymentStore {
	return &CustomerPaymentStore{db: db}
}

func (s *CustomerPaymentStore) Create(m *model.CustomerPayment) error {
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CustomerPaymentStore) CreatePaymentDocument(
	customerPaymentID,
	docID int,
	paymentURL string,
) error {
	m := &model.CustomerPaymentDocument{
		CustomerPaymentID: customerPaymentID,
		DocumentID:        docID,
		PaymentURL:        paymentURL,
	}
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("CustomerContractStore: CreatePaymentDocument %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerPaymentStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.CustomerPayment{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: Update %v\n", err)
		return err
	}
	return nil
}

func (s *CustomerPaymentStore) GetByID(id int) (*model.CustomerPayment, error) {
	res := &model.CustomerPayment{}
	if err := s.db.Where("id = ?", id).Preload("CustomerContract").Find(res).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: GetByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerPaymentStore) GetByCustomerContractID(
	cusContractID int, status model.PaymentStatus, offset, limit int,
) ([]*model.CustomerPayment, error) {
	if limit == 0 {
		limit = 1000
	}
	res := []*model.CustomerPayment{}

	if status == model.PaymentStatusNoFilter {
		if err := s.db.Where("customer_contract_id = ?", cusContractID).Preload("CustomerContract").Order("id desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerPaymentStore: GetByCustomerContractID %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("customer_contract_id = ? and status = ?", cusContractID, string(status)).Preload("CustomerContract").Order("id desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerPaymentStore: GetByCustomerContractID %v\n", err)
			return nil, err
		}
	}

	return res, nil
}

type CustomerPaymentDetail struct {
	CustomerPayment         *model.CustomerPayment         `gorm:"embedded" json:"customer_payment" `
	CustomerPaymentDocument *model.CustomerPaymentDocument `gorm:"embedded" json:"customer_payment_document"`
	Document                *model.Document                `gorm:"embedded" json:"document"`
}

func (s *CustomerPaymentStore) GetLastByPaymentType(
	cusContractID int,
	paymentType model.PaymentType,
) (*CustomerPaymentDetail, error) {
	rawSql := `
select *
from customer_payments
         join customer_payment_documents on customer_payments.id = customer_payment_documents.customer_payment_id
         join documents on customer_payment_documents.document_id = documents.id
where customer_payments.customer_contract_id = ? and customer_payments.payment_type = ? order by customer_payments.id desc
`
	res := &CustomerPaymentDetail{}
	if err := s.db.Raw(rawSql, cusContractID, string(paymentType)).First(res).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: GetLastByPaymentType %v\n", err)
		return nil, err
	}

	return res, nil
}
