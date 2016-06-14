package paymet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestAccountRestore(t *testing.T) {
	acc := &Account{}
	acc.ApplyEvents([]es.Event{
		&AccountOpenedEvent{Name: "sry", Card: "955884334444", Balance: 100},
		&AccountCreditedEvent{Amount: 100},
		&AccountDebitedEvent{Amount: 50},
		&AccountDebitFailedEvent{},
	})
	assert.Equal(t, "sry", acc.name)
	assert.Equal(t, BankCard("955884334444"), acc.card)
	assert.Equal(t, Money(150), acc.balance)
	assert.Equal(t, 4, acc.Version())
}

func TestAccountRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewAccount().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckApplyEvents(t *testing.T) {
	events := []es.Event{
		&AccountOpenedEvent{},
		&AccountCreditedEvent{},
		&AccountDebitedEvent{},
		&AccountDebitFailedEvent{},
		&AccountDebitedBecauseOfTransferEvent{},
		&AccountDebitedBecauseOfTransferFailedEvent{},
		&AccountCreditedBecauseOfTransferEvent{},
	}
	assert.NotPanics(t, func() { NewAccount().ApplyEvents(events) }, "Check Process All Event")
}

func TestAccountCommand(t *testing.T) {
	details1 := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	details2 := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      50,
		Transaction: es.NewGuid(),
	}

	tests := []struct {
		account *Account
		command es.Command
		event   es.Event
	}{
		{
			&Account{},
			&OpenAccountCommand{WithGuid: es.WithGuid{Guid: "1234"}, Name: "sry", Card: "955884334444", Balance: 100},
			&AccountOpenedEvent{WithGuid: es.WithGuid{Guid: "1234"}, Name: "sry", Card: "955884334444", Balance: 100},
		},
		{
			&Account{},
			&CreditAccountCommand{Amount: 100},
			&AccountCreditedEvent{Amount: 100},
		},
		{
			&Account{balance: 50},
			&DebitAccountCommand{Amount: 50},
			&AccountDebitedEvent{Amount: 50},
		},
		{
			&Account{balance: 100},
			&DebitAccountCommand{Amount: 50},
			&AccountDebitedEvent{Amount: 50},
		},
		{
			&Account{balance: 100},
			&DebitAccountCommand{Amount: 101},
			&AccountDebitFailedEvent{},
		},
		{
			&Account{balance: 45},
			&DebitAccountBecauseOfTransferCommand{WithGuid: es.WithGuid{details1.From}, mTDetails: details1},
			&AccountDebitedBecauseOfTransferEvent{WithGuid: es.WithGuid{details1.From}, mTDetails: details1},
		},
		{
			&Account{balance: 40},
			&DebitAccountBecauseOfTransferCommand{WithGuid: es.WithGuid{details2.From}, mTDetails: details2},
			&AccountDebitedBecauseOfTransferFailedEvent{WithGuid: es.WithGuid{details2.From}, mTDetails: details2},
		},
		{
			&Account{balance: 100},
			&CreditAccountBecauseOfTransferCommand{WithGuid: es.WithGuid{details1.To}, mTDetails: details1},
			&AccountCreditedBecauseOfTransferEvent{WithGuid: es.WithGuid{details1.To}, mTDetails: details1},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.account.ProcessCommand(v.command))
	}
}

func TestAccountCommand_Panic(t *testing.T) {
	tests := []struct {
		account *Account
		command es.Command
	}{
		{
			&Account{},
			&struct{ es.WithGuid }{},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.account.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
