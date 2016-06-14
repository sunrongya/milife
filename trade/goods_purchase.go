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

func (g *GoodsPurchase) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *GoodsPurchaseCreatedEvent:
			g.PurchaseDetails, g.state = e.PurchaseDetails, PurchaseStarted
		case *GoodsPurchaseCompletedEvent:
			g.state = PurchaseCompleted
		case *GoodsPurchaseFailedEvent:
			g.state = PurchaseFailed
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	g.SetVersion(len(events))
}

func (g *GoodsPurchase) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateGoodsPurchaseCommand:
		event = g.processCreateGoodsPurchaseCommand(c)
	case *CompleteGoodsPurchaseCommand:
		event = g.processCompleteGoodsPurchaseCommand(c)
	case *FailGoodsPurchaseCommand:
		event = g.processFailGoodsPurchaseCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (g *GoodsPurchase) processCreateGoodsPurchaseCommand(command *CreateGoodsPurchaseCommand) es.Event {
	return &GoodsPurchaseCreatedEvent{PurchaseDetails: command.PurchaseDetails}
}

func (g *GoodsPurchase) processCompleteGoodsPurchaseCommand(command *CompleteGoodsPurchaseCommand) es.Event {
	if g.state != PurchaseStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsPurchaseCommand of state:%s", g.state))
	}
	return &GoodsPurchaseCompletedEvent{PurchaseDetails: command.PurchaseDetails}
}

func (g *GoodsPurchase) processFailGoodsPurchaseCommand(command *FailGoodsPurchaseCommand) es.Event {
	if g.state != PurchaseStarted {
		panic(fmt.Errorf("Can't process FailGoodsPurchaseCommand of state:%s", g.state))
	}
	return &GoodsPurchaseFailedEvent{PurchaseDetails: command.PurchaseDetails}
}

func NewGoodsPurchase() es.Aggregate {
	return &GoodsPurchase{}
}
