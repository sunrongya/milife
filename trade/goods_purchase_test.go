package trade

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestGoodsPurchaseRestore(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 5,
	}
	purchase := &GoodsPurchase{}
	purchase.HandleGoodsPurchaseCreatedEvent(&GoodsPurchaseCreatedEvent{PurchaseDetails: details})
	purchase.HandleGoodsPurchaseCompletedEvent(&GoodsPurchaseCompletedEvent{PurchaseDetails: details})

	assert.Equal(t, details, purchase.PurchaseDetails, "PurchaseDetails error")
	assert.Equal(t, PurchaseCompleted, purchase._state, "state error")
}

func TestCreateGoodsPurchaseCommand(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 5,
	}
	command := &CreateGoodsPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseCreatedEvent{PurchaseDetails: details}}
	goodsPurchase := &GoodsPurchase{}

	assert.Equal(t, events, goodsPurchase.ProcessCreateGoodsPurchaseCommand(command))
}

func TestCompleteGoodsPurchaseCommand(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 5,
	}
	command := &CompleteGoodsPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseCompletedEvent{PurchaseDetails: details}}
	goodsPurchase := &GoodsPurchase{_state: PurchaseStarted}

	assert.Equal(t, events, goodsPurchase.ProcessCompleteGoodsPurchaseCommand(command))
}

func TestFailGoodsPurchaseCommand(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 5,
	}
	command := &FailGoodsPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseFailedEvent{PurchaseDetails: details}}
	goodsPurchase := &GoodsPurchase{_state: PurchaseStarted}

	assert.Equal(t, events, goodsPurchase.ProcessFailGoodsPurchaseCommand(command))
}

func TestCompleteGoodsPurchaseCommand_Panic(t *testing.T) {
	goodsPurchases := []*GoodsPurchase{
		&GoodsPurchase{},
		&GoodsPurchase{_state: PurchaseCompleted},
		&GoodsPurchase{_state: PurchaseFailed},
	}
	for _, goodsPurchase := range goodsPurchases {
		assert.Panics(t, func() {
			goodsPurchase.ProcessCompleteGoodsPurchaseCommand(&CompleteGoodsPurchaseCommand{})
		}, "执行命令CompleteGoodsPurchaseCommand应该抛出异常")
	}
}

func TestFailGoodsPurchaseCommand_Panic(t *testing.T) {
	goodsPurchases := []*GoodsPurchase{
		&GoodsPurchase{},
		&GoodsPurchase{_state: PurchaseCompleted},
		&GoodsPurchase{_state: PurchaseFailed},
	}
	for _, goodsPurchase := range goodsPurchases {
		assert.Panics(t, func() {
			goodsPurchase.ProcessFailGoodsPurchaseCommand(&FailGoodsPurchaseCommand{})
		}, "执行命令FailGoodsPurchaseCommand应该抛出异常")
	}
}
