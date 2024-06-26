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
	ParkingLotMetadata = []OptionResponse{
		{
			Code: string(model.ParkingLotHome),
			Text: "Tại nhà",
		},
		{
			Code: string(model.ParkingLotGarage),
			Text: "Bãi đỗ MinhHungCar",
		},
	}
)

type registerCarMetadataResponse struct {
	Models     []*model.CarModel `json:"models"`
	Periods    []OptionResponse  `json:"periods"`
	Fuels      []OptionResponse  `json:"fuels"`
	Motions    []OptionResponse  `json:"motions"`
	ParkingLot []OptionResponse  `json:"parking_lot"`
}

func (s *Server) HandleGetRegisterCarMetadata(c *gin.Context) {
	models, err := s.store.CarModelStore.GetAll()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, registerCarMetadataResponse{
		Models:     models,
		Periods:    RegisterCarPeriods,
		Fuels:      RegisterCarFuels,
		Motions:    RegisterCarMotions,
		ParkingLot: ParkingLotMetadata,
	})
}

type getParkingLotMetadataRequest struct {
	SeatType int `form:"seat_type"`
}

func (s *Server) HandleGetParkingLotMetadata(c *gin.Context) {
	req := getParkingLotMetadataRequest{}
	if err := c.Bind(&req); err != nil {
		responseError(c, err)
		return
	}

	totalCarInGarage, err := s.store.CarStore.CountBySeats(req.SeatType, model.ParkingLotGarage, []model.CarStatus{model.CarStatusActive, model.CarStatusWaitingDelivery})
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	garageConfig, err := s.store.GarageConfigStore.Get()
	if err != nil {
		responseInternalServerError(c, err)
		return
	}

	typeCode := model.GarageConfigTypeMax4Seats
	if req.SeatType == 7 {
		typeCode = model.GarageConfigTypeMax7Seats
	} else if req.SeatType == 15 {
		typeCode = model.GarageConfigTypeMax15Seats
	}

	if totalCarInGarage >= garageConfig[typeCode] {
		c.JSON(http.StatusOK, []OptionResponse{
			{
				Code: string(model.ParkingLotHome),
				Text: "Tại nhà",
			},
		})
		return
	}

	c.JSON(http.StatusOK, ParkingLotMetadata)
}
