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
	items []OrderItem
	price Money
	state OrderState
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

func (o *Order) ProcessCreateOrderCommand(command *CreateOrderCommand) []es.Event {
	if len(command.Items) == 0 {
		panic(fmt.Errorf("order items is nil"))
	}
	return []es.Event{&OrderCreatedEvent{Items: command.Items}}
}

func (o *Order) ProcessCancelOrderCommand(command *CancelOrderCommand) []es.Event {
	if o.state != OrderCreated {
		panic(fmt.Errorf("Can't process CancelOrderCommand of state:%s", o.state))
	}
	return []es.Event{&OrderCanceledEvent{}}
}

func (o *Order) ProcessCreateOrderPaymetCommand(command *CreateOrderPaymetCommand) []es.Event {
	if o.state != OrderCreated {
		panic(fmt.Errorf("Can't process CreateOrderPaymetCommand of state:%s", o.state))
	}
	return []es.Event{
		&OrderPaymetCreatedEvent{
			User:           command.User,
			UserAccount:    command.UserAccount,
			ManagedAccount: command.ManagedAccount,
			Price:          o.price,
		},
	}
}

func (o *Order) ProcessCompleteOrderPaymetCommand(command *CompleteOrderPaymetCommand) []es.Event {
	if o.state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process CompleteOrderPaymetCommand of state:%s", o.state))
	}
	return []es.Event{&OrderPaymetCompletedEvent{OrderItems: o.items, User: command.User}}
}

func (o *Order) ProcessFailOrderPaymetCommand(command *FailOrderPaymetCommand) []es.Event {
	if o.state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process FailOrderPaymetCommand of state:%s", o.state))
	}
	return []es.Event{&OrderPaymetFailedEvent{OrderItems: o.items, User: command.User}}
}

func (o *Order) HandleOrderCreatedEvent(event *OrderCreatedEvent) {
	o.items, o.state = event.Items, OrderCreated
	for _, v := range event.Items {
		o.price += Money(float64(v.Price) * float64(v.Quantity))
	}
}

func (o *Order) HandleOrderCanceledEvent(event *OrderCanceledEvent) {
	o.state = OrderCanceled
}

func (o *Order) HandleOrderPaymetCreatedEvent(event *OrderPaymetCreatedEvent) {
	o.state = OrderPaymetCreated
}

func (o *Order) HandleOrderPaymetCompletedEvent(event *OrderPaymetCompletedEvent) {
	o.state = OrderPaymetCompleted
}

func (o *Order) HandleOrderPaymetFailedEvent(event *OrderPaymetFailedEvent) {
	o.state = OrderPaymetFailed
}
