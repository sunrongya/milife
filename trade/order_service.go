package trade

import(
    es "github.com/sunrongya/eventsourcing"
)

type OrderService struct {
    es.Service
}

func NewOrderService(store es.EventStore) *OrderService{
    service := &OrderService{
        Service: es.NewService(store, NewOrder),
    }
    return service
}

func (o *OrderService) CreateOrder(orderItems []OrderItem) es.Guid {
    guid := es.NewGuid()
    c := &CreateOrderCommand{
        WithGuid: es.WithGuid{guid},
        Items: orderItems,
    }
    o.PublishCommand(c)
    return guid
}

func (o *OrderService) CancelOrder(guid es.Guid) {
    c := &CancelOrderCommand{
        WithGuid: es.WithGuid{guid},
    }
    o.PublishCommand(c)
}

func (o *OrderService) PaymetOrder(guid, user es.Guid, userAccount, managedAccount BankAccount) {
    c := &CreateOrderPaymetCommand{
        WithGuid: es.WithGuid{guid},
        User: user,
        UserAccount: userAccount,
        ManagedAccount: managedAccount,
    }
    o.PublishCommand(c)
}
