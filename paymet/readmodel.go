package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type RAccount struct {
	Id      es.Guid
	Name    string
	Card    string
	Balance Money
}

type RAccountProjector struct {
	_repository es.ReadRepository
}

func NewRAccountProjector(repository es.ReadRepository) *RAccountProjector {
	return &RAccountProjector{_repository: repository}
}

func (this *RAccountProjector) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	account := &RAccount{
		Id:      event.GetGuid(),
		Name:    event.Name,
		Card:    string(event.Card),
		Balance: event.Balance,
	}
	this._repository.Save(account.Id, account)
}

func (this *RAccountProjector) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	this.do(event.GetGuid(), func(account *RAccount) {
		account.Balance += event.Amount
	})
}

func (this *RAccountProjector) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	this.do(event.GetGuid(), func(account *RAccount) {
		account.Balance -= event.Amount
	})
}

func (this *RAccountProjector) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	this.do(event.GetGuid(), func(account *RAccount) {
		account.Balance -= event.Amount
	})
}

func (this *RAccountProjector) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	this.do(event.GetGuid(), func(account *RAccount) {
		account.Balance += event.Amount
	})
}

func (this *RAccountProjector) do(id es.Guid, assignRAccountFn func(*RAccount)) {
	i, err := this._repository.Find(id)
	if err != nil {
		return
	}
	account := i.(*RAccount)
	assignRAccountFn(account)
	this._repository.Save(id, account)
}

type MoneyFlowRate struct {
	NumOpened                    int
	NumCredited                  int
	NumDebited                   int
	Amount                       Money
	NumCreditedBecauseOfTransfer int
	NumDebitedBecauseOfTransfer  int
}

type MoneyFlowRateProjector struct {
	_repository es.ReadRepository
	Id          es.Guid
}

func NewMoneyFlowRateProjector(repository es.ReadRepository, id es.Guid) *MoneyFlowRateProjector {
	return &MoneyFlowRateProjector{_repository: repository, Id: id}
}

func (this *MoneyFlowRateProjector) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	this.do(func(rate *MoneyFlowRate) {
		rate.NumOpened += 1
		rate.Amount += event.Balance
	})
}

func (this *MoneyFlowRateProjector) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	this.do(func(rate *MoneyFlowRate) {
		rate.NumCredited += 1
		rate.Amount += event.Amount
	})
}

func (this *MoneyFlowRateProjector) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	this.do(func(rate *MoneyFlowRate) {
		rate.NumDebited += 1
		rate.Amount += event.Amount
	})
}

func (this *MoneyFlowRateProjector) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	this.do(func(rate *MoneyFlowRate) {
		rate.NumDebitedBecauseOfTransfer += 1
		rate.Amount += event.Amount
	})
}

func (this *MoneyFlowRateProjector) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	this.do(func(rate *MoneyFlowRate) {
		rate.NumCreditedBecauseOfTransfer += 1
		rate.Amount += event.Amount
	})
}

func (this *MoneyFlowRateProjector) do(assignRateFn func(*MoneyFlowRate)) {
	i, _ := this._repository.Find(this.Id)
	if i == nil {
		i = &MoneyFlowRate{}
	}
	rate := i.(*MoneyFlowRate)
	assignRateFn(rate)
	this._repository.Save(this.Id, rate)
}
