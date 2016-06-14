package paymet

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestAccountReadModel(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	readRepository := es.NewMemoryReadRepository()
	accProjector := NewRAccountProjector(readRepository)

	accProjector.HandleAccountOpenedEvent(&AccountOpenedEvent{WithGuid: es.WithGuid{details.From}, Name: "sry", Card: "955884334444", Balance: 100})
	accProjector.HandleAccountOpenedEvent(&AccountOpenedEvent{WithGuid: es.WithGuid{details.To}, Name: "managed account", Card: "955884334888", Balance: 1000})
	accProjector.HandleAccountCreditedEvent(&AccountCreditedEvent{WithGuid: es.WithGuid{details.From}, Amount: 100})
	accProjector.HandleAccountDebitedEvent(&AccountDebitedEvent{WithGuid: es.WithGuid{details.From}, Amount: 50})
	accProjector.HandleAccountDebitedBecauseOfTransferEvent(&AccountDebitedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.From}, mTDetails: details})
	accProjector.HandleAccountCreditedBecauseOfTransferEvent(&AccountCreditedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.To}, mTDetails: details})

	fromAccountTmp, err := readRepository.Find(details.From)
	assert.Nil(t, err, "fromAccount read repository find error")
	fromAccount := fromAccountTmp.(*RAccount)
	toAccountTmp, err := readRepository.Find(details.To)
	assert.Nil(t, err, "toAccount read repository find error")
	toAccount := toAccountTmp.(*RAccount)

	assert.Equal(t, details.From, fromAccount.Id, "")
	assert.Equal(t, "sry", fromAccount.Name, "")
	assert.Equal(t, "955884334444", fromAccount.Card, "")
	assert.Equal(t, Money(110), fromAccount.Balance, "")

	assert.Equal(t, details.To, toAccount.Id, "")
	assert.Equal(t, "managed account", toAccount.Name, "")
	assert.Equal(t, "955884334888", toAccount.Card, "")
	assert.Equal(t, Money(1040), toAccount.Balance, "")
}

func TestMoneyFlowRateReadModel(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	id := es.NewGuid()
	readRepository := es.NewMemoryReadRepository()
	rateProjector := NewMoneyFlowRateProjector(readRepository, id)

	rateProjector.HandleAccountOpenedEvent(&AccountOpenedEvent{WithGuid: es.WithGuid{details.From}, Name: "sry", Card: "955884334444", Balance: 100})
	rateProjector.HandleAccountOpenedEvent(&AccountOpenedEvent{WithGuid: es.WithGuid{details.To}, Name: "managed account", Card: "955884334888", Balance: 1000})
	rateProjector.HandleAccountCreditedEvent(&AccountCreditedEvent{WithGuid: es.WithGuid{details.From}, Amount: 100})
	rateProjector.HandleAccountCreditedEvent(&AccountCreditedEvent{WithGuid: es.WithGuid{details.From}, Amount: 80})
	rateProjector.HandleAccountCreditedEvent(&AccountCreditedEvent{WithGuid: es.WithGuid{details.From}, Amount: 45})
	rateProjector.HandleAccountDebitedEvent(&AccountDebitedEvent{WithGuid: es.WithGuid{details.From}, Amount: 50})
	rateProjector.HandleAccountDebitedBecauseOfTransferEvent(&AccountDebitedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.From}, mTDetails: details})
	rateProjector.HandleAccountCreditedBecauseOfTransferEvent(&AccountCreditedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.To}, mTDetails: details})

	moneyFlowRateTmp, err := readRepository.Find(id)
	assert.Nil(t, err, "MoneyFlowRate read repository find error")
	moneyFlowRate := moneyFlowRateTmp.(*MoneyFlowRate)

	assert.Equal(t, 2, moneyFlowRate.NumOpened, "")
	assert.Equal(t, 3, moneyFlowRate.NumCredited, "")
	assert.Equal(t, 1, moneyFlowRate.NumDebited, "")
	assert.Equal(t, Money(1455), moneyFlowRate.Amount, "")
	assert.Equal(t, 1, moneyFlowRate.NumCreditedBecauseOfTransfer, "")
	assert.Equal(t, 1, moneyFlowRate.NumDebitedBecauseOfTransfer, "")
}
