package api

import "github.com/godev111222333/capstone-backend/src/service"

var _ service.IPDFService = (*MockPDFService)(nil)

type MockPDFService struct{}

func (m *MockPDFService) Render(tz service.RenderType, payload map[string]string) (string, error) {
	return "rendered_url", nil
}
