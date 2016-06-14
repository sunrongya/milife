package paymet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestMoneyTransferRestore(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	moneyTransfer := &MoneyTransfer{}
	moneyTransfer.ApplyEvents([]es.Event{
		&TransferCreatedEvent{mTDetails: details},
		&TransferDebitedEvent{mTDetails: details},
		&TransferCompletedEvent{mTDetails: details},
	})
	assert.Equal(t, Completed, moneyTransfer.state, "")
	assert.Equal(t, details, moneyTransfer.mTDetails, "")
	assert.Equal(t, 3, moneyTransfer.Version())
}

func TestMoneyTransferForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewMoneyTransfer().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckMoneyTransferApplyEvents(t *testing.T) {
	events := []es.Event{
		&TransferCreatedEvent{},
		&TransferDebitedEvent{},
		&TransferCompletedEvent{},
		&TransferFailedEvent{},
	}
	assert.NotPanics(t, func() { NewMoneyTransfer().ApplyEvents(events) }, "Check Process All Event")
}

func TestMoneyTransferProcessCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}

	tests := []struct {
		transfer *MoneyTransfer
		command  es.Command
		event    es.Event
	}{
		{
			&MoneyTransfer{},
			&CreateTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
			&TransferCreatedEvent{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
		},
		{
			&MoneyTransfer{state: Created},
			&DebitedTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
			&TransferDebitedEvent{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
		},
		{
			&MoneyTransfer{state: Debited},
			&CompletedTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
			&TransferCompletedEvent{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
		},
		{
			&MoneyTransfer{state: Debited},
			&FailedTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
			&TransferFailedEvent{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
		},
		{
			&MoneyTransfer{},
			&FailedTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
			&TransferFailedEvent{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.transfer.ProcessCommand(v.command))
	}
}

func TestMoneyTransferCommand_Panic(t *testing.T) {
	tests := []struct {
		transfer *MoneyTransfer
		command  es.Command
	}{
		{
			&MoneyTransfer{},
			&struct{ es.WithGuid }{},
		},
		{
			&MoneyTransfer{},
			&DebitedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Debited},
			&DebitedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Completed},
			&DebitedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Failed},
			&DebitedTransferCommand{},
		},

		{
			&MoneyTransfer{},
			&CompletedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Created},
			&CompletedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Completed},
			&CompletedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Failed},
			&CompletedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Created},
			&FailedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Completed},
			&FailedTransferCommand{},
		},
		{
			&MoneyTransfer{state: Failed},
			&FailedTransferCommand{},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.transfer.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
