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

func (o *Order) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *OrderCreatedEvent:
			o.items, o.state = e.Items, OrderCreated
			for _, v := range e.Items {
				o.price += Money(float64(v.Price) * float64(v.Quantity))
			}
		case *OrderCanceledEvent:
			o.state = OrderCanceled
		case *OrderPaymetCreatedEvent:
			o.state = OrderPaymetCreated
		case *OrderPaymetCompletedEvent:
			o.state = OrderPaymetCompleted
		case *OrderPaymetFailedEvent:
			o.state = OrderPaymetFailed
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	o.SetVersion(len(events))
}

func (o *Order) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateOrderCommand:
		event = o.processCreateOrderCommand(c)
	case *CancelOrderCommand:
		event = o.processCancelOrderCommand(c)
	case *CreateOrderPaymetCommand:
		event = o.processCreateOrderPaymetCommand(c)
	case *CompleteOrderPaymetCommand:
		event = o.processCompleteOrderPaymetCommand(c)
	case *FailOrderPaymetCommand:
		event = o.processFailOrderPaymetCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (o *Order) processCreateOrderCommand(command *CreateOrderCommand) es.Event {
	if len(command.Items) == 0 {
		panic(fmt.Errorf("order items is nil"))
	}
	return &OrderCreatedEvent{Items: command.Items}
}

func (o *Order) processCancelOrderCommand(command *CancelOrderCommand) es.Event {
	if o.state != OrderCreated {
		panic(fmt.Errorf("Can't process CancelOrderCommand of state:%s", o.state))
	}
	return &OrderCanceledEvent{}
}

func (o *Order) processCreateOrderPaymetCommand(command *CreateOrderPaymetCommand) es.Event {
	if o.state != OrderCreated {
		panic(fmt.Errorf("Can't process CreateOrderPaymetCommand of state:%s", o.state))
	}
	return &OrderPaymetCreatedEvent{
		User:           command.User,
		UserAccount:    command.UserAccount,
		ManagedAccount: command.ManagedAccount,
		Price:          o.price,
	}
}

func (o *Order) processCompleteOrderPaymetCommand(command *CompleteOrderPaymetCommand) es.Event {
	if o.state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process CompleteOrderPaymetCommand of state:%s", o.state))
	}
	return &OrderPaymetCompletedEvent{OrderItems: o.items, User: command.User}
}

func (o *Order) processFailOrderPaymetCommand(command *FailOrderPaymetCommand) es.Event {
	if o.state != OrderPaymetCreated {
		panic(fmt.Errorf("Can't process FailOrderPaymetCommand of state:%s", o.state))
	}
	return &OrderPaymetFailedEvent{OrderItems: o.items, User: command.User}
}

func NewOrder() es.Aggregate {
	return &Order{}
}
