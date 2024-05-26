package usecase

import (
	"testing"
	"time"

	"github.com/samuelferpim/go-clean-arch-challenge/internal/entity"
	"github.com/samuelferpim/go-clean-arch-challenge/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of OrderRepositoryInterface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

// MockEvent is a mock implementation of EventInterface
type MockEvent struct {
	mock.Mock
	payload interface{}
}

func (m *MockEvent) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEvent) GetDateTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockEvent) GetPayload() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockEvent) SetPayload(payload interface{}) {
	m.Called(payload)
	m.payload = payload
}

// MockEventDispatcher is a mock implementation of EventDispatcherInterface
type MockEventDispatcher struct {
	mock.Mock
}

func (m *MockEventDispatcher) Register(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventDispatcher) Dispatch(event events.EventInterface) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventDispatcher) Remove(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventDispatcher) Has(eventName string, handler events.EventHandlerInterface) bool {
	args := m.Called(eventName, handler)
	return args.Bool(0)
}

func (m *MockEventDispatcher) Clear() {
	m.Called()
}

func TestCreateOrderUseCase_Execute(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	orderCreated := new(MockEvent)
	eventDispatcher := new(MockEventDispatcher)

	// Initialize the use case with mocks
	useCase := NewCreateOrderUseCase(orderRepo, orderCreated, eventDispatcher)

	input := OrderInputDTO{
		ID:    "123",
		Price: 100.0,
		Tax:   10.0,
	}

	expectedOrder := entity.Order{
		ID:    "123",
		Price: 100.0,
		Tax:   10.0,
	}
	expectedOrder.CalculateFinalPrice()

	orderRepo.On("Save", &expectedOrder).Return(nil)
	orderCreated.On("SetPayload", mock.Anything).Return()
	orderCreated.On("GetDateTime").Return(time.Now())
	orderCreated.On("GetName").Return("OrderCreated")
	orderCreated.On("GetPayload").Return(expectedOrder)
	eventDispatcher.On("Dispatch", mock.Anything).Return(nil)

	output, err := useCase.Execute(input)

	assert.Nil(t, err)
	assert.Equal(t, "123", output.ID)
	assert.Equal(t, 100.0, output.Price)
	assert.Equal(t, 10.0, output.Tax)
	assert.Equal(t, 110.0, output.FinalPrice)

	orderRepo.AssertCalled(t, "Save", &expectedOrder)
	orderCreated.AssertCalled(t, "SetPayload", output)
	eventDispatcher.AssertCalled(t, "Dispatch", orderCreated)
}
