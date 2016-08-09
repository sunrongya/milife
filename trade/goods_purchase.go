package trade

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type PurchaseState string

const (
	PurchaseStarted   = PurchaseState("PurchaseStarted")
	PurchaseCompleted = PurchaseState("PurchaseCompleted")
	PurchaseFailed    = PurchaseState("PurchaseFailed")
)

type GoodsPurchase struct {
	es.BaseAggregate
	PurchaseDetails
	_state PurchaseState
}

type PurchaseDetails struct {
	Quantity
	User     es.Guid
	Goods    es.Guid
	Purchase es.Guid
}

var _ es.Aggregate = (*GoodsPurchase)(nil)

func NewGoodsPurchase() es.Aggregate {
	return &GoodsPurchase{}
}

func (this *GoodsPurchase) ProcessCreateGoodsPurchaseCommand(command *CreateGoodsPurchaseCommand) []es.Event {
	return []es.Event{&GoodsPurchaseCreatedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (this *GoodsPurchase) ProcessCompleteGoodsPurchaseCommand(command *CompleteGoodsPurchaseCommand) []es.Event {
	if this._state != PurchaseStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsPurchaseCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsPurchaseCompletedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (this *GoodsPurchase) ProcessFailGoodsPurchaseCommand(command *FailGoodsPurchaseCommand) []es.Event {
	if this._state != PurchaseStarted {
		panic(fmt.Errorf("Can't process FailGoodsPurchaseCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsPurchaseFailedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (this *GoodsPurchase) HandleGoodsPurchaseCreatedEvent(event *GoodsPurchaseCreatedEvent) {
	this.PurchaseDetails, this._state = event.PurchaseDetails, PurchaseStarted
}

func (this *GoodsPurchase) HandleGoodsPurchaseCompletedEvent(event *GoodsPurchaseCompletedEvent) {
	this._state = PurchaseCompleted
}

func (this *GoodsPurchase) HandleGoodsPurchaseFailedEvent(event *GoodsPurchaseFailedEvent) {
	this._state = PurchaseFailed
}
