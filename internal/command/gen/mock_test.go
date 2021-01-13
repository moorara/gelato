package gen

type (
	GenerateMock struct {
		InPath    string
		InMock    bool
		InFactory bool
		OutError  error
	}

	MockGenerateService struct {
		GenerateIndex int
		GenerateMocks []GenerateMock
	}
)

func (m *MockGenerateService) Generate(path string, mock, factory bool) error {
	i := m.GenerateIndex
	m.GenerateIndex++
	m.GenerateMocks[i].InPath = path
	m.GenerateMocks[i].InMock = mock
	m.GenerateMocks[i].InFactory = factory
	return m.GenerateMocks[i].OutError
}
