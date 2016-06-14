package trade

import(
    "testing"
    es "github.com/sunrongya/eventsourcing"
    "github.com/sunrongya/eventsourcing/utiltest"
)

func TestOrderServiceDoCreateOrder(t *testing.T) {
    orderItems := []OrderItem{
        {
            Goods: es.NewGuid(),
            Name: "abc",
            Price: 30,
            Quantity: 4,
        },
        {
            Goods: es.NewGuid(),
            Name: "ccd",
            Price: 25,
            Quantity: 3,
        },
    }
    
    utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command{
        gs := OrderService{ Service: service}
        guid := gs.CreateOrder(orderItems)
        return &CreateOrderCommand{ WithGuid:es.WithGuid{guid}, Items:orderItems }
    })
}

func TestOrderServiceDoCancelOrder(t *testing.T) {
    utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command{
        guid := es.NewGuid()
        gs := OrderService{ Service: service}
        gs.CancelOrder(guid)
        return &CancelOrderCommand{ WithGuid:es.WithGuid{guid} }
    })
}

func TestOrderServiceDoPaymetOrder(t *testing.T) {
    utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command{
        guid, user := es.NewGuid(), es.NewGuid()
        userAccount, managedAccount := BankAccount("95533999494"), BankAccount("955882333")
        gs := OrderService{ Service: service}
        gs.PaymetOrder(guid, user, userAccount, managedAccount)
        return &CreateOrderPaymetCommand{
            WithGuid: es.WithGuid{guid}, 
            User: user, 
            UserAccount: userAccount, 
            ManagedAccount: managedAccount,
        }
    })
}

