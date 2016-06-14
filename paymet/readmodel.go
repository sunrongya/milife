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
	repository es.ReadRepository
}

func NewRAccountProjector(repository es.ReadRepository) *RAccountProjector {
	return &RAccountProjector{repository: repository}
}

func (r *RAccountProjector) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	account := &RAccount{
		Id:      event.GetGuid(),
		Name:    event.Name,
		Card:    string(event.Card),
		Balance: event.Balance,
	}
	r.repository.Save(account.Id, account)
}

func (r *RAccountProjector) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	r.do(event.GetGuid(), func(account *RAccount) {
		account.Balance += event.Amount
	})
}

func (r *RAccountProjector) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	r.do(event.GetGuid(), func(account *RAccount) {
		account.Balance -= event.Amount
	})
}

func (r *RAccountProjector) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	r.do(event.GetGuid(), func(account *RAccount) {
		account.Balance -= event.Amount
	})
}

func (r *RAccountProjector) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	r.do(event.GetGuid(), func(account *RAccount) {
		account.Balance += event.Amount
	})
}

func (r *RAccountProjector) do(id es.Guid, assignRAccountFn func(*RAccount)) {
	i, err := r.repository.Find(id)
	if err != nil {
		return
	}
	account := i.(*RAccount)
	assignRAccountFn(account)
	r.repository.Save(id, account)
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
	repository es.ReadRepository
	Id         es.Guid
}

func NewMoneyFlowRateProjector(repository es.ReadRepository, id es.Guid) *MoneyFlowRateProjector {
	return &MoneyFlowRateProjector{repository: repository, Id: id}
}

func (r *MoneyFlowRateProjector) HandleAccountOpenedEvent(event *AccountOpenedEvent) {
	r.do(func(rate *MoneyFlowRate) {
		rate.NumOpened += 1
		rate.Amount += event.Balance
	})
}

func (r *MoneyFlowRateProjector) HandleAccountCreditedEvent(event *AccountCreditedEvent) {
	r.do(func(rate *MoneyFlowRate) {
		rate.NumCredited += 1
		rate.Amount += event.Amount
	})
}

func (r *MoneyFlowRateProjector) HandleAccountDebitedEvent(event *AccountDebitedEvent) {
	r.do(func(rate *MoneyFlowRate) {
		rate.NumDebited += 1
		rate.Amount += event.Amount
	})
}

func (r *MoneyFlowRateProjector) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	r.do(func(rate *MoneyFlowRate) {
		rate.NumDebitedBecauseOfTransfer += 1
		rate.Amount += event.Amount
	})
}

func (r *MoneyFlowRateProjector) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	r.do(func(rate *MoneyFlowRate) {
		rate.NumCreditedBecauseOfTransfer += 1
		rate.Amount += event.Amount
	})
}

func (r *MoneyFlowRateProjector) do(assignRateFn func(*MoneyFlowRate)) {
	i, _ := r.repository.Find(r.Id)
	if i == nil {
		i = &MoneyFlowRate{}
	}
	rate := i.(*MoneyFlowRate)
	assignRateFn(rate)
	r.repository.Save(r.Id, rate)
}
