package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/godev111222333/capstone-backend/src/misc"
)

var _ IPDFService = (*PDFService)(nil)

type RenderType string

const (
	RenderTypeCustomer RenderType = "customer"
	RenderTypePartner  RenderType = "partner"
)

type IPDFService interface {
	Render(tz RenderType, payload map[string]string) (string, error)
}

type PDFService struct {
	cfg    *misc.PDFServiceConfig
	client *http.Client
}

func NewPDFService(cfg *misc.PDFServiceConfig) *PDFService {
	client := http.DefaultClient
	client.Timeout = cfg.Timeout
	return &PDFService{cfg, client}
}

func (s *PDFService) Render(tz RenderType, payload map[string]string) (string, error) {
	url := s.cfg.Url + "/render_customer_contract"
	if tz == RenderTypePartner {
		url = s.cfg.Url + "/render_partner_contract"
	}
	bz, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewReader(bz))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("render contract error :%v\n", err)
		return "", nil
	}

	bz, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()

	respPayload := struct {
		UUID string `json:"uuid"`
	}{}
	if err := json.Unmarshal(bz, &respPayload); err != nil {
		fmt.Println(err)
		return "", err
	}

	return respPayload.UUID, nil
}
