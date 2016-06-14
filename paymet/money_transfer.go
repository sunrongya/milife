package paymet

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type State string

const (
	Created   = State("Created")
	Debited   = State("Debited")
	Completed = State("Completed")
	Failed    = State("Failed")
)

type mTDetails struct {
	From        es.Guid
	To          es.Guid
	Amount      Money
	Transaction es.Guid
}

type MoneyTransfer struct {
	es.BaseAggregate
	mTDetails mTDetails
	state     State
}

var _ es.Aggregate = (*MoneyTransfer)(nil)

func (a *MoneyTransfer) ApplyEvents(events []es.Event) {
	for _, e := range events {
		switch event := e.(type) {
		case *TransferCreatedEvent:
			a.mTDetails, a.state = event.mTDetails, Created
		case *TransferDebitedEvent:
			a.state = Debited
		case *TransferCompletedEvent:
			a.state = Completed
		case *TransferFailedEvent:
			a.state = Failed
		default:
			panic(fmt.Errorf("Unknown event %#v", event))
		}
	}
	a.SetVersion(len(events))
}

func (a *MoneyTransfer) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateTransferCommand:
		event = a.processCreateTransferCommand(c)
	case *DebitedTransferCommand:
		event = a.processDebitedTransferCommand(c)
	case *CompletedTransferCommand:
		event = a.processCompletedTransferCommand(c)
	case *FailedTransferCommand:
		event = a.processFailedTransferCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (a *MoneyTransfer) processCreateTransferCommand(command *CreateTransferCommand) es.Event {
	return &TransferCreatedEvent{mTDetails: command.mTDetails}
}

func (a *MoneyTransfer) processDebitedTransferCommand(command *DebitedTransferCommand) es.Event {
	if a.state != Created {
		panic(fmt.Errorf("Can't process DebitedTransferCommand of state:%s", a.state))
	}
	return &TransferDebitedEvent{mTDetails: command.mTDetails}
}

func (a *MoneyTransfer) processCompletedTransferCommand(command *CompletedTransferCommand) es.Event {
	if a.state != Debited {
		panic(fmt.Errorf("Can't process CompletedTransferCommand of state:%s", a.state))
	}
	return &TransferCompletedEvent{mTDetails: command.mTDetails}
}

func (a *MoneyTransfer) processFailedTransferCommand(command *FailedTransferCommand) es.Event {
	if a.state == Created || a.state == Completed || a.state == Failed {
		panic(fmt.Errorf("Can't process FailedTransferCommand of state:%s", a.state))
	}
	return &TransferFailedEvent{mTDetails: command.mTDetails}
}

func NewMoneyTransfer() es.Aggregate {
	return &MoneyTransfer{}
}
