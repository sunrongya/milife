package trade

import(
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
    es "github.com/sunrongya/eventsourcing"
)

func TestGoodsPurchaseRestore(t *testing.T) {
    details := PurchaseDetails{
        User:     es.NewGuid(),
        Goods:    es.NewGuid(),
        Purchase: es.NewGuid(),
        Quantity: 5,
    }
    purchase := &GoodsPurchase{}
    purchase.ApplyEvents([]es.Event{
        &GoodsPurchaseCreatedEvent{ WithGuid:es.WithGuid{details.Purchase}, PurchaseDetails:details },
        &GoodsPurchaseCompletedEvent{ WithGuid:es.WithGuid{details.Purchase}, PurchaseDetails:details },
    })
    assert.Equal(t, 2, purchase.Version(), "version error")
    assert.Equal(t, details, purchase.PurchaseDetails, "PurchaseDetails error")
    assert.Equal(t, PurchaseCompleted, purchase.state, "state error")
}

func TestGoodsPurchaseRestoreForErrorEvent(t *testing.T){
    assert.Panics(t, func(){ 
        NewGoodsPurchase().ApplyEvents([]es.Event{ &struct{es.WithGuid}{} }) 
    }, "restore error event must panic error")
}

func TestCheckGoodsPurchaseApplyEvents(t *testing.T) {
    events := []es.Event{
        &GoodsPurchaseCreatedEvent{},
        &GoodsPurchaseCompletedEvent{},
        &GoodsPurchaseFailedEvent{},
    }
    assert.NotPanics(t, func(){ NewGoodsPurchase().ApplyEvents(events) }, "Check Process All Event")
}

func TestGoodsPurchaseCommand(t *testing.T){
    details := PurchaseDetails{
        User:     es.NewGuid(),
        Goods:    es.NewGuid(),
        Purchase: es.NewGuid(),
        Quantity: 5,
    }
    
    tests := []struct{
        goodsPurchase  *GoodsPurchase
        command  es.Command
        event  es.Event
    }{
        {
            &GoodsPurchase{},
            &CreateGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
            &GoodsPurchaseCreatedEvent{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseStarted },
            &CompleteGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
            &GoodsPurchaseCompletedEvent{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseStarted },
            &FailGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
            &GoodsPurchaseFailedEvent{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
    }
    
    for _, v := range tests {
        assert.Equal(t, []es.Event{v.event}, v.goodsPurchase.ProcessCommand(v.command) )
    }
}

func TestGoodsPurchaseCommand_Panic(t *testing.T){
    details := PurchaseDetails{
        User:     es.NewGuid(),
        Goods:    es.NewGuid(),
        Purchase: es.NewGuid(),
        Quantity: 5,
    }
    
    tests := []struct{
        goodsPurchase  *GoodsPurchase
        command  es.Command
    }{
        {
            &GoodsPurchase{},
            &struct{es.WithGuid}{ WithGuid:es.WithGuid{Guid:details.Purchase} },
        },
        {
            &GoodsPurchase{},
            &CompleteGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseCompleted },
            &CompleteGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseFailed },
            &CompleteGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{},
            &FailGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseCompleted },
            &FailGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
        {
            &GoodsPurchase{ state:PurchaseFailed },
            &FailGoodsPurchaseCommand{WithGuid:es.WithGuid{Guid:details.Purchase}, PurchaseDetails:details},
        },
    }
    
    for _, v := range tests {
        assert.Panics(t, func(){v.goodsPurchase.ProcessCommand(v.command)}, fmt.Sprintf("test panics error: command:%v", v.command))
    }
}

