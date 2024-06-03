package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
)

type OptionResponse struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

var (
	RegisterCarPeriods = []OptionResponse{
		{
			Code: "1",
			Text: "1 tháng",
		},
		{
			Code: "3",
			Text: "3 tháng",
		},
		{
			Code: "6",
			Text: "6 tháng",
		},
		{
			Code: "12",
			Text: "12 tháng",
		},
	}
	RegisterCarFuels = []OptionResponse{
		{
			Code: string(model.FuelGas),
			Text: "Xăng",
		},
		{
			Code: string(model.FuelOil),
			Text: "Dầu",
		},
		{
			Code: string(model.FuelElectricity),
			Text: "Điện",
		},
	}
	RegisterCarMotions = []OptionResponse{
		{
			Code: string(model.MotionAutomaticTransmission),
			Text: "Số tự động",
		},
		{
			Code: string(model.MotionManualTransmission),
			Text: "Số sàn",
		},
	}
)

type registerCarMetadataResponse struct {
	Models  []*model.CarModel `json:"models"`
	Periods []OptionResponse  `json:"periods"`
	Fuels   []OptionResponse  `json:"fuels"`
	Motions []OptionResponse  `json:"motions"`
}

func (s *Server) HandleGetRegisterCarMetadata(c *gin.Context) {
	models, err := s.store.CarModelStore.GetAll()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, registerCarMetadataResponse{
		Models:  models,
		Periods: RegisterCarPeriods,
		Fuels:   RegisterCarFuels,
		Motions: RegisterCarMotions,
	})
}
