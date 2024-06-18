package api

var _ IPDFService = (*MockPDFService)(nil)

type MockPDFService struct{}

func (m *MockPDFService) Render(tz RenderType, payload map[string]string) (string, error) {
	return "rendered_url", nil
}
