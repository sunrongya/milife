package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type Account struct {
	es.BaseAggregate
	_name    string
	_card    BankCard
	_balance Money
}

var _ es.Aggregate = (*Account)(nil)

func NewAccount() es.Aggregate {
	return &Account{}
}

func (this *Account) ProcessOpenAccountCommand(command *OpenAccountCommand) []es.Event {
	return []es.Event{
		&AccountOpenedEvent{
			Name:    command.Name,
			Card:    command.Card,
			Balance: command.Balance,
		},
	}
}

func (this *Account) ProcessCreditAccountCommand(command *CreditAccountCommand) []es.Event {
	return []es.Event{&AccountCreditedEvent{Amount: command.Amount}}
}

func (this *Account) ProcessDebitAccountCommand(command *DebitAccountCommand) []es.Event {
	if this._balance < command.Amount {
		return []es.Event{&AccountDebitFailedEvent{}}
	}
	return []es.Event{&AccountDebitedEvent{Amount: command.Amount}}
}

func (this *Account) ProcessDebitAccountBecauseOfTransferCommand(command *DebitAccountBecauseOfTransferCommand) []es.Event {
	if this._balance < command.Amount {
		return []es.Event{&AccountDebitedBecauseOfTransferFailedEvent{mTDetails: command.mTDetails}}
	}
	return []es.Event{&AccountDebitedBecauseOfTransferEvent{mTDetails: command.mTDetails}}
}

func (this *Account) ProcessCreditAccountBecauseOfTransferCommand(command *CreditAccountBecauseOfTransferCommand) []es.Event {
	return []es.Event{&AccountCreditedBecauseOfTransferEvent{mTDetails: command.mTDetails}}
}

func (this *Account) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	this._name, this._card, this._balance = event.Name, event.Card, event.Balance
}

func (this *Account) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	this._balance += event.Amount
}

func (this *Account) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	this._balance -= event.Amount
}

func (this *Account) HandleAccountDebitFailedEvent(event *AccountDebitFailedEvent) {
}

func (this *Account) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	this._balance -= event.Amount
}

func (this *Account) HandleAccountDebitedBecauseOfTransferFailedEvent(event *AccountDebitedBecauseOfTransferFailedEvent) {
}

func (this *Account) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	this._balance += event.Amount
}
