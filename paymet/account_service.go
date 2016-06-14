package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type AccountService struct {
	es.Service
}

func NewAccountService(store es.EventStore) *AccountService {
	acc := &AccountService{
		Service: es.NewService(store, NewAccount),
	}
	return acc
}

func (a *AccountService) OpenAccount(name, card string, balance int) es.Guid {
	guid := es.NewGuid()
	c := &OpenAccountCommand{
		WithGuid: es.WithGuid{guid},
		Name:     name,
		Card:     BankCard(card),
		Balance:  Money(balance),
	}
	a.PublishCommand(c)
	return guid
}

func (a *AccountService) CreditAccount(guid es.Guid, amount int) {
	c := &CreditAccountCommand{
		WithGuid: es.WithGuid{guid},
		Amount:   Money(amount),
	}
	a.PublishCommand(c)
}

func (a *AccountService) DebitAccount(guid es.Guid, amount int) {
	c := &DebitAccountCommand{
		WithGuid: es.WithGuid{guid},
		Amount:   Money(amount),
	}
	a.PublishCommand(c)
}
