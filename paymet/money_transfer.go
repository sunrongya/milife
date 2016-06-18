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

func NewMoneyTransfer() es.Aggregate {
	return &MoneyTransfer{}
}

func (a *MoneyTransfer) ProcessCreateTransferCommand(command *CreateTransferCommand) []es.Event {
	return []es.Event{&TransferCreatedEvent{mTDetails: command.mTDetails}}
}

func (a *MoneyTransfer) ProcessDebitedTransferCommand(command *DebitedTransferCommand) []es.Event {
	if a.state != Created {
		panic(fmt.Errorf("Can't process DebitedTransferCommand of state:%s", a.state))
	}
	return []es.Event{&TransferDebitedEvent{mTDetails: command.mTDetails}}
}

func (a *MoneyTransfer) ProcessCompletedTransferCommand(command *CompletedTransferCommand) []es.Event {
	if a.state != Debited {
		panic(fmt.Errorf("Can't process CompletedTransferCommand of state:%s", a.state))
	}
	return []es.Event{&TransferCompletedEvent{mTDetails: command.mTDetails}}
}

func (a *MoneyTransfer) ProcessFailedTransferCommand(command *FailedTransferCommand) []es.Event {
	if a.state == Created || a.state == Completed || a.state == Failed {
		panic(fmt.Errorf("Can't process FailedTransferCommand of state:%s", a.state))
	}
	return []es.Event{&TransferFailedEvent{mTDetails: command.mTDetails}}
}

func (a *MoneyTransfer) HandleTransferCreatedEvent(event *TransferCreatedEvent) {
	a.mTDetails, a.state = event.mTDetails, Created
}

func (a *MoneyTransfer) HandleTransferDebitedEvent(event *TransferDebitedEvent) {
	a.state = Debited
}

func (a *MoneyTransfer) HandleTransferCompletedEvent(event *TransferCompletedEvent) {
	a.state = Completed
}

func (a *MoneyTransfer) HandleTransferFailedEvent(event *TransferFailedEvent) {
	a.state = Failed
}
