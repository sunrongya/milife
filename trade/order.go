package trade

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type OrderState string

const (
	OrderCreated         = OrderState("OrderCreated")
	OrderCanceled        = OrderState("OrderCanceled")
	OrderPaymetCreated   = OrderState("OrderPaymetCreated")
	OrderPaymetCompleted = OrderState("OrderPaymetCompleted")
	OrderPaymetFailed    = OrderState("OrderPaymetFailed")
)

type Order struct {
	es.BaseAggregate
	_items []OrderItem
	_price Money
	_state OrderState
}

type OrderItem struct {
	Goods    es.Guid
	Purchase es.Guid
	Name     string
	Price    Money
	Quantity Quantity
}

var _ es.Aggregate = (*Order)(nil)

func NewOrder() es.Aggregate {
	return &Order{}
}

func (this *Order) ProcessCreateOrderCommand(command *CreateOrderCommand) []es.Event {
	if len(command.Items) == 0 {
		panic(fmt.Errorf("order items is nil"))
	}
	return []es.Event{&OrderCreatedEvent{Items: command.Items}}
}

func (this *Order) ProcessCancelOrderCommand(command *CancelOrderCommand) []es.Event {
	if this._state != OrderCreated {
		panic(fmt.Errorf("Can't process CancelOrderCommand of state:%s", this._state))
	}
	return []es.Event{&OrderCanceledEvent{}}
}

func (this *Order) ProcessCreateOrderPaymetCommand(command *CreateOrderPaymetCommand) []es.Event {
	if this._state != OrderCreated {
		panic(fmt.Errorf("Can't process CreateOrderPaymetCommand of state:%s", this._state))
	}
	return []es.Event{
		&OrderPaymetCreatedEvent{
			User:           command.User,
			UserAccount:    command.UserAccount,
			ManagedAccount: command.ManagedAccount,
			Price:          this._price,
		},
	}
}

func (this *Order) ProcessCompleteOrderPaymetCommand(command *CompleteOrderPaymetCommand) []es.Event {
	if this._state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process CompleteOrderPaymetCommand of state:%s", this._state))
	}
	return []es.Event{&OrderPaymetCompletedEvent{OrderItems: this._items, User: command.User}}
}

func (this *Order) ProcessFailOrderPaymetCommand(command *FailOrderPaymetCommand) []es.Event {
	if this._state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process FailOrderPaymetCommand of state:%s", this._state))
	}
	return []es.Event{&OrderPaymetFailedEvent{OrderItems: this._items, User: command.User}}
}

func (this *Order) HandleOrderCreatedEvent(event *OrderCreatedEvent) {
	this._items, this._state = event.Items, OrderCreated
	for _, v := range event.Items {
		this._price += Money(float64(v.Price) * float64(v.Quantity))
	}
}

func (this *Order) HandleOrderCanceledEvent(event *OrderCanceledEvent) {
	this._state = OrderCanceled
}

func (this *Order) HandleOrderPaymetCreatedEvent(event *OrderPaymetCreatedEvent) {
	this._state = OrderPaymetCreated
}

func (this *Order) HandleOrderPaymetCompletedEvent(event *OrderPaymetCompletedEvent) {
	this._state = OrderPaymetCompleted
}

func (this *Order) HandleOrderPaymetFailedEvent(event *OrderPaymetFailedEvent) {
	this._state = OrderPaymetFailed
}
