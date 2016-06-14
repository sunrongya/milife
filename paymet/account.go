package paymet

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type Account struct {
	es.BaseAggregate
	name    string
	card    BankCard
	balance Money
}

var _ es.Aggregate = (*Account)(nil)

func (a *Account) ApplyEvents(events []es.Event) {
	for _, e := range events {
		switch event := e.(type) {
		case *AccountOpenedEvent:
			a.name, a.card, a.balance = event.Name, event.Card, event.Balance
		case *AccountCreditedEvent:
			a.balance += event.Amount
		case *AccountDebitedEvent:
			a.balance -= event.Amount
		case *AccountDebitFailedEvent:
		case *AccountDebitedBecauseOfTransferEvent:
			a.balance -= event.Amount
		case *AccountDebitedBecauseOfTransferFailedEvent:
		case *AccountCreditedBecauseOfTransferEvent:
			a.balance += event.Amount
		default:
			panic(fmt.Errorf("Unknown event %#v", event))
		}
	}
	a.SetVersion(len(events))
}

func (a *Account) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *OpenAccountCommand:
		event = a.proccessOpenAccountCommand(c)
	case *CreditAccountCommand:
		event = a.proccessCreditAccountCommand(c)
	case *DebitAccountCommand:
		event = a.proccessDebitAccountCommand(c)
	case *DebitAccountBecauseOfTransferCommand:
		event = a.proccessDebitAccountBecauseOfTransferCommand(c)
	case *CreditAccountBecauseOfTransferCommand:
		event = a.proccessCreditAccountBecauseOfTransferCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (a *Account) proccessOpenAccountCommand(command *OpenAccountCommand) es.Event {
	return &AccountOpenedEvent{
		Name:    command.Name,
		Card:    command.Card,
		Balance: command.Balance,
	}
}

func (a *Account) proccessCreditAccountCommand(command *CreditAccountCommand) es.Event {
	return &AccountCreditedEvent{Amount: command.Amount}
}

func (a *Account) proccessDebitAccountCommand(command *DebitAccountCommand) es.Event {
	if a.balance < command.Amount {
		return &AccountDebitFailedEvent{}
	}
	return &AccountDebitedEvent{Amount: command.Amount}
}

func (a *Account) proccessDebitAccountBecauseOfTransferCommand(command *DebitAccountBecauseOfTransferCommand) es.Event {
	if a.balance < command.Amount {
		return &AccountDebitedBecauseOfTransferFailedEvent{mTDetails: command.mTDetails}
	}
	return &AccountDebitedBecauseOfTransferEvent{mTDetails: command.mTDetails}
}

func (a *Account) proccessCreditAccountBecauseOfTransferCommand(command *CreditAccountBecauseOfTransferCommand) es.Event {
	return &AccountCreditedBecauseOfTransferEvent{mTDetails: command.mTDetails}
}

func NewAccount() es.Aggregate {
	return &Account{}
}
