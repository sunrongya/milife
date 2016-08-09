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

func (this *AccountService) OpenAccount(name, card string, balance int) es.Guid {
	guid := es.NewGuid()
	c := &OpenAccountCommand{
		WithGuid: es.WithGuid{guid},
		Name:     name,
		Card:     BankCard(card),
		Balance:  Money(balance),
	}
	this.PublishCommand(c)
	return guid
}

func (this *AccountService) CreditAccount(guid es.Guid, amount int) {
	c := &CreditAccountCommand{
		WithGuid: es.WithGuid{guid},
		Amount:   Money(amount),
	}
	this.PublishCommand(c)
}

func (this *AccountService) DebitAccount(guid es.Guid, amount int) {
	c := &DebitAccountCommand{
		WithGuid: es.WithGuid{guid},
		Amount:   Money(amount),
	}
	this.PublishCommand(c)
}
