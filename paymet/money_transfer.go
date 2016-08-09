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
	_mTDetails mTDetails
	_state     State
}

var _ es.Aggregate = (*MoneyTransfer)(nil)

func NewMoneyTransfer() es.Aggregate {
	return &MoneyTransfer{}
}

func (this *MoneyTransfer) ProcessCreateTransferCommand(command *CreateTransferCommand) []es.Event {
	return []es.Event{&TransferCreatedEvent{mTDetails: command.mTDetails}}
}

func (this *MoneyTransfer) ProcessDebitedTransferCommand(command *DebitedTransferCommand) []es.Event {
	if this._state != Created {
		panic(fmt.Errorf("Can't process DebitedTransferCommand of state:%s", this._state))
	}
	return []es.Event{&TransferDebitedEvent{mTDetails: command.mTDetails}}
}

func (this *MoneyTransfer) ProcessCompletedTransferCommand(command *CompletedTransferCommand) []es.Event {
	if this._state != Debited {
		panic(fmt.Errorf("Can't process CompletedTransferCommand of state:%s", this._state))
	}
	return []es.Event{&TransferCompletedEvent{mTDetails: command.mTDetails}}
}

func (this *MoneyTransfer) ProcessFailedTransferCommand(command *FailedTransferCommand) []es.Event {
	if this._state == Completed || this._state == Failed {
		panic(fmt.Errorf("Can't process FailedTransferCommand of state:%s", this._state))
	}
	return []es.Event{&TransferFailedEvent{mTDetails: command.mTDetails}}
}

func (this *MoneyTransfer) HandleTransferCreatedEvent(event *TransferCreatedEvent) {
	this._mTDetails, this._state = event.mTDetails, Created
}

func (this *MoneyTransfer) HandleTransferDebitedEvent(event *TransferDebitedEvent) {
	this._state = Debited
}

func (this *MoneyTransfer) HandleTransferCompletedEvent(event *TransferCompletedEvent) {
	this._state = Completed
}

func (this *MoneyTransfer) HandleTransferFailedEvent(event *TransferFailedEvent) {
	this._state = Failed
}

func (this *MoneyTransfer) State() State {
	return this._state
}
