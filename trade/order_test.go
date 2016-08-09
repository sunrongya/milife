package trade

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestOrderRestore(t *testing.T) {
	orderItems := []OrderItem{
		{
			Goods:    es.NewGuid(),
			Name:     "abc",
			Price:    30,
			Quantity: 4,
		},
		{
			Goods:    es.NewGuid(),
			Name:     "ccd",
			Price:    25,
			Quantity: 3,
		},
	}
	order := &Order{}
	order.HandleOrderCreatedEvent(&OrderCreatedEvent{Items: orderItems})
	order.HandleOrderCanceledEvent(&OrderCanceledEvent{})

	assert.Equal(t, orderItems, order._items, "Items error")
	assert.Equal(t, Money(195), order._price, "price error")
	assert.Equal(t, OrderCanceled, order._state, "state error")
}

func TestCreateOrderCommand(t *testing.T) {
	orderItems := []OrderItem{
		{
			Goods:    es.NewGuid(),
			Name:     "abc",
			Price:    30,
			Quantity: 4,
		},
		{
			Goods:    es.NewGuid(),
			Name:     "ccd",
			Price:    25,
			Quantity: 3,
		},
	}
	command := &CreateOrderCommand{Items: orderItems}
	events := []es.Event{&OrderCreatedEvent{Items: orderItems}}
	order := &Order{}

	assert.Equal(t, events, order.ProcessCreateOrderCommand(command))
}

func TestCancelOrderCommand(t *testing.T) {
	command := &CancelOrderCommand{}
	events := []es.Event{&OrderCanceledEvent{}}
	order := &Order{_state: OrderCreated}

	assert.Equal(t, events, order.ProcessCancelOrderCommand(command))
}

func TestCreateOrderPaymetCommand(t *testing.T) {
	orderItems := []OrderItem{
		{
			Goods:    es.NewGuid(),
			Name:     "abc",
			Price:    30,
			Quantity: 4,
		},
		{
			Goods:    es.NewGuid(),
			Name:     "ccd",
			Price:    25,
			Quantity: 3,
		},
	}
	user := es.NewGuid()
	command := &CreateOrderPaymetCommand{User: user, UserAccount: "95588333", ManagedAccount: "93388388"}
	events := []es.Event{&OrderPaymetCreatedEvent{User: user, UserAccount: "95588333", ManagedAccount: "93388388", Price: 195}}
	order := &Order{_state: OrderCreated, _items: orderItems, _price: 195}

	assert.Equal(t, events, order.ProcessCreateOrderPaymetCommand(command))
}

func TestCompleteOrderPaymetCommand(t *testing.T) {
	orderItems := []OrderItem{
		{
			Goods:    es.NewGuid(),
			Name:     "abc",
			Price:    30,
			Quantity: 4,
		},
		{
			Goods:    es.NewGuid(),
			Name:     "ccd",
			Price:    25,
			Quantity: 3,
		},
	}
	user := es.NewGuid()
	command := &CompleteOrderPaymetCommand{User: user, UserAccount: "95588333", ManagedAccount: "93388388", Price: 195}
	events := []es.Event{&OrderPaymetCompletedEvent{OrderItems: orderItems, User: user}}
	order := &Order{_state: OrderPaymetCreated, _items: orderItems, _price: 195}

	assert.Equal(t, events, order.ProcessCompleteOrderPaymetCommand(command))
}

func TestFailOrderPaymetCommand(t *testing.T) {
	orderItems := []OrderItem{
		{
			Goods:    es.NewGuid(),
			Name:     "abc",
			Price:    30,
			Quantity: 4,
		},
		{
			Goods:    es.NewGuid(),
			Name:     "ccd",
			Price:    25,
			Quantity: 3,
		},
	}
	user := es.NewGuid()
	command := &FailOrderPaymetCommand{User: user, UserAccount: "95588333", ManagedAccount: "93388388", Price: 195}
	events := []es.Event{&OrderPaymetFailedEvent{OrderItems: orderItems, User: user}}
	order := &Order{_state: OrderPaymetCreated, _items: orderItems, _price: 195}

	assert.Equal(t, events, order.ProcessFailOrderPaymetCommand(command))
}

func TestCreateOrderCommand_Panic(t *testing.T) {
	assert.Panics(t, func() { new(Order).ProcessCreateOrderCommand(&CreateOrderCommand{Items: []OrderItem{}}) })
}

func TestCancelOrderCommand_Panic(t *testing.T) {
	assert.Panics(t, func() { new(Order).ProcessCancelOrderCommand(&CancelOrderCommand{}) })
}

func TestCreateOrderPaymetCommand_Panic(t *testing.T) {
	assert.Panics(t, func() { new(Order).ProcessCreateOrderPaymetCommand(&CreateOrderPaymetCommand{}) })
}

func TestCompleteOrderPaymetCommand_Panic(t *testing.T) {
	assert.Panics(t, func() { new(Order).ProcessCompleteOrderPaymetCommand(&CompleteOrderPaymetCommand{}) })
}

func TestFailOrderPaymetCommand_Panic(t *testing.T) {
	assert.Panics(t, func() { new(Order).ProcessFailOrderPaymetCommand(&FailOrderPaymetCommand{}) })
}
