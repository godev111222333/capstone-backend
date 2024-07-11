package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type StatisticRequest struct {
	TotalCustomerContractsBackOffDay int `form:"total_customer_contracts_back_off_day" binding:"required"`
	TotalActivePartnersBackOffDay    int `form:"total_active_partners_back_off_day" binding:"required"`
	TotalActiveCustomersBackOffDay   int `form:"total_active_customers_back_off_day" binding:"required"`
	RevenueBackOffDay                int `form:"revenue_back_off_day" binding:"required"`
	RentedCarsBackOffDay             int `form:"rented_cars_back_off_day" binding:"required"`
}

type StatisticResponse struct {
	TotalCustomerContracts int                `json:"total_customer_contracts,omitempty"`
	TotalActivePartners    int                `json:"total_active_partners,omitempty"`
	TotalActiveCustomers   int                `json:"total_active_customers,omitempty"`
	Revenue                float64            `json:"revenue,omitempty"`
	RentedCars             []*store.RentedCar `json:"rented_cars,omitempty"`
}

func (s *Server) HandleAdminGetStatistic(c *gin.Context) {
	req := StatisticRequest{}
	if err := c.Bind(&req); err != nil {
		responseCustomErr(c, ErrCodeInvalidStatisticRequest, err)
		return
	}

	totalCustomerContracts, err :=
		s.store.CustomerContractStore.CountTotalValidCustomerContracts(
			dayToDuration(req.TotalCustomerContractsBackOffDay),
		)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	totalActivePartners, err := s.store.AccountStore.CountActiveByRole(
		model.RoleIDPartner,
		dayToDuration(req.TotalActivePartnersBackOffDay),
	)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	totalActiveCustomers, err := s.store.AccountStore.CountActiveByRole(
		model.RoleIDCustomer,
		dayToDuration(req.TotalActiveCustomersBackOffDay),
	)
	if err != nil {
		responseGormErr(c, err)
		return
	}

	revenue, err := s.store.CustomerContractStore.SumRevenueForCompletedContracts(dayToDuration(req.RevenueBackOffDay))
	if err != nil {
		responseGormErr(c, err)
		return
	}

	rentedCars, err := s.store.CustomerContractStore.CountRentedCars(dayToDuration(req.RentedCarsBackOffDay))
	if err != nil {
		responseGormErr(c, err)
		return
	}

	responseSuccess(c, StatisticResponse{
		TotalCustomerContracts: totalCustomerContracts,
		TotalActivePartners:    totalActivePartners,
		TotalActiveCustomers:   totalActiveCustomers,
		Revenue:                revenue,
		RentedCars:             rentedCars,
	})
}

func dayToDuration(day int) time.Duration {
	return time.Hour * 24 * time.Duration(day)
}
