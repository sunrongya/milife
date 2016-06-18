package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type Account struct {
	es.BaseAggregate
	name    string
	card    BankCard
	balance Money
}

var _ es.Aggregate = (*Account)(nil)

func NewAccount() es.Aggregate {
	return &Account{}
}

func (a *Account) ProccessOpenAccountCommand(command *OpenAccountCommand) []es.Event {
	return []es.Event{
		&AccountOpenedEvent{
			Name:    command.Name,
			Card:    command.Card,
			Balance: command.Balance,
		},
	}
}

func (a *Account) ProccessCreditAccountCommand(command *CreditAccountCommand) []es.Event {
	return []es.Event{&AccountCreditedEvent{Amount: command.Amount}}
}

func (a *Account) ProccessDebitAccountCommand(command *DebitAccountCommand) []es.Event {
	if a.balance < command.Amount {
		return []es.Event{&AccountDebitFailedEvent{}}
	}
	return []es.Event{&AccountDebitedEvent{Amount: command.Amount}}
}

func (a *Account) ProccessDebitAccountBecauseOfTransferCommand(command *DebitAccountBecauseOfTransferCommand) []es.Event {
	if a.balance < command.Amount {
		return []es.Event{&AccountDebitedBecauseOfTransferFailedEvent{mTDetails: command.mTDetails}}
	}
	return []es.Event{&AccountDebitedBecauseOfTransferEvent{mTDetails: command.mTDetails}}
}

func (a *Account) ProccessCreditAccountBecauseOfTransferCommand(command *CreditAccountBecauseOfTransferCommand) []es.Event {
	return []es.Event{&AccountCreditedBecauseOfTransferEvent{mTDetails: command.mTDetails}}
}

func (a *Account) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	a.name, a.card, a.balance = event.Name, event.Card, event.Balance
}

func (a *Account) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	a.balance += event.Amount
}

func (a *Account) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	a.balance -= event.Amount
}

func (a *Account) HandleAccountDebitFailedEvent(event *AccountDebitFailedEvent) {
}

func (a *Account) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	a.balance -= event.Amount
}

func (a *Account) HandleAccountDebitedBecauseOfTransferFailedEvent(event *AccountDebitedBecauseOfTransferFailedEvent) {
}

func (a *Account) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	a.balance += event.Amount
}
