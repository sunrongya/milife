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
	state PurchaseState
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

func (g *GoodsPurchase) ProcessCreateGoodsPurchaseCommand(command *CreateGoodsPurchaseCommand) []es.Event {
	return []es.Event{&GoodsPurchaseCreatedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (g *GoodsPurchase) ProcessCompleteGoodsPurchaseCommand(command *CompleteGoodsPurchaseCommand) []es.Event {
	if g.state != PurchaseStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsPurchaseCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsPurchaseCompletedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (g *GoodsPurchase) ProcessFailGoodsPurchaseCommand(command *FailGoodsPurchaseCommand) []es.Event {
	if g.state != PurchaseStarted {
		panic(fmt.Errorf("Can't process FailGoodsPurchaseCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsPurchaseFailedEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (g *GoodsPurchase) HandleGoodsPurchaseCreatedEvent(event *GoodsPurchaseCreatedEvent) {
	g.PurchaseDetails, g.state = event.PurchaseDetails, PurchaseStarted
}

func (g *GoodsPurchase) HandleGoodsPurchaseCompletedEvent(event *GoodsPurchaseCompletedEvent) {
	g.state = PurchaseCompleted
}

func (g *GoodsPurchase) HandleGoodsPurchaseFailedEvent(event *GoodsPurchaseFailedEvent) {
	g.state = PurchaseFailed
}
