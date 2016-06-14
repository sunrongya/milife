package trade

import(
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
    es "github.com/sunrongya/eventsourcing"
)

func TestOrderRestore(t *testing.T) {
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
    guid := es.NewGuid()
    order := &Order{}
    order.ApplyEvents([]es.Event{
        &OrderCreatedEvent{ WithGuid:es.WithGuid{guid}, Items:orderItems },
        &OrderCanceledEvent{ WithGuid:es.WithGuid{guid} },
    })
    assert.Equal(t, 2, order.Version(), "version error")
    assert.Equal(t, orderItems, order.items, "Items error")
    assert.Equal(t, Money(195), order.price, "price error")
    assert.Equal(t, OrderCanceled, order.state, "state error")
}

func TestOrderRestoreForErrorEvent(t *testing.T){
    assert.Panics(t, func(){ 
        NewOrder().ApplyEvents([]es.Event{ &struct{es.WithGuid}{} }) 
    }, "restore error event must panic error")
}

func TestCheckOrderApplyEvents(t *testing.T) {
    events := []es.Event{
        &OrderCreatedEvent{},
        &OrderCanceledEvent{},
        &OrderPaymetCreatedEvent{},
        &OrderPaymetCompletedEvent{},
        &OrderPaymetFailedEvent{},
    }
    assert.NotPanics(t, func(){ NewOrder().ApplyEvents(events) }, "Check Process All Event")
}

func TestOrderCommand(t *testing.T){
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
    guid := es.NewGuid()
    
    tests := []struct{
        order  *Order
        command  es.Command
        event  es.Event
    }{
        {
            &Order{},
            &CreateOrderCommand{ WithGuid:es.WithGuid{guid}, Items:orderItems },
            &OrderCreatedEvent{ WithGuid:es.WithGuid{guid}, Items:orderItems },
        },
        {
            &Order{ state:OrderCreated },
            &CancelOrderCommand{ WithGuid:es.WithGuid{Guid:guid} },
            &OrderCanceledEvent{ WithGuid:es.WithGuid{Guid:guid} },
        },
    }
    
    for _, v := range tests {
        assert.Equal(t, []es.Event{v.event}, v.order.ProcessCommand(v.command) )
    }
}

func TestOrderCommand_Panic(t *testing.T){
    tests := []struct{
        order  *Order
        command  es.Command
    }{
        {
            &Order{},
            &struct{es.WithGuid}{},
        },
        {
            &Order{},
            &CreateOrderCommand{ Items:[]OrderItem{} },
        },
        {
            &Order{},
            &CancelOrderCommand{},
        },
        {
            &Order{},
            &CreateOrderPaymetCommand{},
        },
        {
            &Order{},
            &CompleteOrderPaymetCommand{},
        },
        {
            &Order{},
            &FailOrderPaymetCommand{},
        },
    }
    
    for _, v := range tests {
        assert.Panics(t, func(){v.order.ProcessCommand(v.command)}, fmt.Sprintf("test panics error: command:%v", v.command))
    }
}

func TestOrderPaymetCommand(t *testing.T){
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
    guid := es.NewGuid()
    user := es.NewGuid()
    
    tests := []struct{
        order  *Order
        command  es.Command
        event  es.Event
    }{
        {
            &Order{ state:OrderCreated, items:orderItems, price:195 },
            &CreateOrderPaymetCommand{ WithGuid:es.WithGuid{guid}, User:user, UserAccount:"95588333", ManagedAccount:"93388388" },
            &OrderPaymetCreatedEvent{ WithGuid:es.WithGuid{guid}, User:user, UserAccount:"95588333", ManagedAccount:"93388388", Price:195 },
        },
        {
            &Order{ state:OrderPaymetCreated, items:orderItems, price:195  },
            &CompleteOrderPaymetCommand{ WithGuid:es.WithGuid{guid}, User:user, UserAccount:"95588333", ManagedAccount:"93388388", Price:195 },
            &OrderPaymetCompletedEvent{ WithGuid:es.WithGuid{guid}, OrderItems:orderItems, User:user },
        },
        {
            &Order{ state:OrderPaymetCreated, items:orderItems, price:195  },
            &FailOrderPaymetCommand{ WithGuid:es.WithGuid{guid}, User:user, UserAccount:"95588333", ManagedAccount:"93388388", Price:195 },
            &OrderPaymetFailedEvent{ WithGuid:es.WithGuid{guid}, OrderItems:orderItems, User:user },
        },
    }
    
    for _, v := range tests {
        assert.Equal(t, []es.Event{v.event}, v.order.ProcessCommand(v.command) )
    }
}
