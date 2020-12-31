package controller

import (
	"context"
	"time"
)

type ConnectMock struct {
	OutError error
}

type DisconnectMock struct {
	InCtx    context.Context
	OutError error
}

// MockClient is a mock implementation for graceful.Client interface.
type MockClient struct {
	StringOutString string

	ConnectCounter int
	ConnectMocks   []ConnectMock

	DisconnectCounter int
	DisconnectMocks   []DisconnectMock
}

func (m *MockClient) String() string {
	return m.StringOutString
}

func (m *MockClient) Connect() error {
	i := m.ConnectCounter
	m.ConnectCounter++
	return m.ConnectMocks[i].OutError
}

func (m *MockClient) Disconnect(ctx context.Context) error {
	i := m.DisconnectCounter
	m.DisconnectCounter++
	m.DisconnectMocks[i].InCtx = ctx
	return m.DisconnectMocks[i].OutError
}

type CheckHealthMock struct {
	InCtx    context.Context
	OutError error
}

// MockChecker is a mock implementation for health.Checker interface.
type MockChecker struct {
	StringOutString string

	CheckHealthCounter int
	CheckHealthMocks   []CheckHealthMock
}

func (m *MockChecker) String() string {
	return m.StringOutString
}

func (m *MockChecker) CheckHealth(ctx context.Context) error {
	i := m.CheckHealthCounter
	m.CheckHealthCounter++
	m.CheckHealthMocks[i].InCtx = ctx
	return m.CheckHealthMocks[i].OutError
}

type GetValueMock struct {
	InCtx          context.Context
	InPhraseID     string
	InLanguageCode string
	OutString      string
	OutError       error
}

// MockTranslateGateway is a mock implementation for gateway.TranslateGateway interface.
type MockTranslateGateway struct {
	MockClient
	MockChecker

	StringOutString string

	GetStringCounter int
	GetValueMocks    []GetValueMock
}

func (m *MockTranslateGateway) String() string {
	return m.StringOutString
}

func (m *MockTranslateGateway) GetValue(ctx context.Context, phraseID, languageCode string) (string, error) {
	i := m.GetStringCounter
	m.GetStringCounter++
	m.GetValueMocks[i].InCtx = ctx
	m.GetValueMocks[i].InPhraseID = phraseID
	m.GetValueMocks[i].InLanguageCode = languageCode
	return m.GetValueMocks[i].OutString, m.GetValueMocks[i].OutError
}

type CreateMock struct {
	InCtx        context.Context
	InGreeting   string
	OutTimestamp time.Time
	OutError     error
}

type GetMock struct {
	InCtx       context.Context
	InTimestamp time.Time
	OutGreeting string
	OutError    error
}

// MockGreetingRepository is a mock implementation for repository.GreetingRepository interface.
type MockGreetingRepository struct {
	MockClient
	MockChecker

	StringOutString string

	CreateCounter int
	CreateMocks   []CreateMock

	GetCounter int
	GetInMocks []GetMock
}

func (m *MockGreetingRepository) String() string {
	return m.StringOutString
}

func (m *MockGreetingRepository) Create(ctx context.Context, greeting string) (time.Time, error) {
	i := m.CreateCounter
	m.CreateCounter++
	m.CreateMocks[i].InCtx = ctx
	m.CreateMocks[i].InGreeting = greeting
	return m.CreateMocks[i].OutTimestamp, m.CreateMocks[i].OutError
}

func (m *MockGreetingRepository) Get(ctx context.Context, timestamp time.Time) (string, error) {
	i := m.GetCounter
	m.GetCounter++
	m.GetInMocks[i].InCtx = ctx
	m.GetInMocks[i].InTimestamp = timestamp
	return m.GetInMocks[i].OutGreeting, m.GetInMocks[i].OutError
}
