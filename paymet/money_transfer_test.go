package paymet

import (
	//"fmt"
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
	moneyTransfer.HandleTransferCreatedEvent(&TransferCreatedEvent{mTDetails: details})
	moneyTransfer.HandleTransferDebitedEvent(&TransferDebitedEvent{mTDetails: details})
	moneyTransfer.HandleTransferCompletedEvent(&TransferCompletedEvent{mTDetails: details})

	assert.Equal(t, Completed, moneyTransfer.state, "")
	assert.Equal(t, details, moneyTransfer.mTDetails, "")
}

func TestCreateTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	command := &CreateTransferCommand{mTDetails: details}
	events := []es.Event{&TransferCreatedEvent{mTDetails: details}}

	assert.Equal(t, events, new(MoneyTransfer).ProcessCreateTransferCommand(command), "执行CreateTransferCommand产生的事件有误")
}

func TestDebitedTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	moneyTransfer := &MoneyTransfer{state: Created}
	command := &DebitedTransferCommand{mTDetails: details}
	events := []es.Event{&TransferDebitedEvent{mTDetails: details}}

	assert.Equal(t, events, moneyTransfer.ProcessDebitedTransferCommand(command), "执行DebitedTransferCommand产生的事件有误")
}

func TestCompletedTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	moneyTransfer := &MoneyTransfer{state: Debited}
	command := &CompletedTransferCommand{mTDetails: details}
	events := []es.Event{&TransferCompletedEvent{mTDetails: details}}

	assert.Equal(t, events, moneyTransfer.ProcessCompletedTransferCommand(command), "执行CompletedTransferCommand产生的事件有误")
}

func TestFailedTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	moneyTransfers := []*MoneyTransfer{&MoneyTransfer{state: Debited}, &MoneyTransfer{}}
	command := &FailedTransferCommand{mTDetails: details}
	events := []es.Event{&TransferFailedEvent{mTDetails: details}}

	for _, moneyTransfer := range moneyTransfers {
		assert.Equal(t, events, moneyTransfer.ProcessFailedTransferCommand(command), "执行FailedTransferCommand产生的事件有误")
	}
}

func TestDebitedTransferCommand_Panic(t *testing.T) {
	moneyTransfers := []*MoneyTransfer{
		&MoneyTransfer{},
		&MoneyTransfer{state: Debited},
		&MoneyTransfer{state: Completed},
		&MoneyTransfer{state: Failed},
	}
	for _, moneyTransfer := range moneyTransfers {
		assert.Panics(t, func() {
			moneyTransfer.ProcessDebitedTransferCommand(&DebitedTransferCommand{})
		}, "执行DebitedTransferCommand命令应该抛出异常")
	}
}

func TestCompletedTransferCommand_Panic(t *testing.T) {
	moneyTransfers := []*MoneyTransfer{
		&MoneyTransfer{},
		&MoneyTransfer{state: Created},
		&MoneyTransfer{state: Completed},
		&MoneyTransfer{state: Failed},
	}
	for _, moneyTransfer := range moneyTransfers {
		assert.Panics(t, func() {
			moneyTransfer.ProcessCompletedTransferCommand(&CompletedTransferCommand{})
		}, "执行CompletedTransferCommand命令应该抛出异常")
	}
}

func TestFailedTransferCommand_Panic(t *testing.T) {
	moneyTransfers := []*MoneyTransfer{
		&MoneyTransfer{state: Created},
		&MoneyTransfer{state: Completed},
		&MoneyTransfer{state: Failed},
	}
	for _, moneyTransfer := range moneyTransfers {
		assert.Panics(t, func() {
			moneyTransfer.ProcessFailedTransferCommand(&FailedTransferCommand{})
		}, "执行FailedTransferCommand命令应该抛出异常")
	}
}
