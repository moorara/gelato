package gen

type (
	GenerateMock struct {
		InPath   string
		OutError error
	}

	MockGenerateService struct {
		GenerateIndex int
		GenerateMocks []GenerateMock
	}
)

func (m *MockGenerateService) Generate(path string) error {
	i := m.GenerateIndex
	m.GenerateIndex++
	m.GenerateMocks[i].InPath = path
	return m.GenerateMocks[i].OutError
}
