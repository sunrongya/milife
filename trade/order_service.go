package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type OrderService struct {
	es.Service
}

func NewOrderService(store es.EventStore) *OrderService {
	service := &OrderService{
		Service: es.NewService(store, NewOrder),
	}
	return service
}

func (this *OrderService) CreateOrder(orderItems []OrderItem) es.Guid {
	guid := es.NewGuid()
	c := &CreateOrderCommand{
		WithGuid: es.WithGuid{guid},
		Items:    orderItems,
	}
	this.PublishCommand(c)
	return guid
}

func (this *OrderService) CancelOrder(guid es.Guid) {
	c := &CancelOrderCommand{
		WithGuid: es.WithGuid{guid},
	}
	this.PublishCommand(c)
}

func (this *OrderService) PaymetOrder(guid, user es.Guid, userAccount, managedAccount BankAccount) {
	c := &CreateOrderPaymetCommand{
		WithGuid:       es.WithGuid{guid},
		User:           user,
		UserAccount:    userAccount,
		ManagedAccount: managedAccount,
	}
	this.PublishCommand(c)
}
