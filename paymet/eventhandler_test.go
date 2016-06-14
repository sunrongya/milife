package paymet

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func testHandleEvent(t *testing.T, methodName string, doAccountHandle func(chan es.Command, mTDetails) es.Command) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	ch := make(chan es.Command)
	command := doAccountHandle(ch, details)
	select {
	case c := <-ch:
		assert.Equal(t, c, command, methodName)
	case <-time.After(1 * time.Second):
		t.Error(methodName)
	}
}

func TestHandleTransferCreatedEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleTransferCreatedEvent", func(accountCH chan es.Command, details mTDetails) es.Command {
		handler := NewEventHandler(accountCH, nil)
		go handler.HandleTransferCreatedEvent(&TransferCreatedEvent{WithGuid: es.WithGuid{details.Transaction}, mTDetails: details})
		return &DebitAccountBecauseOfTransferCommand{WithGuid: es.WithGuid{details.From}, mTDetails: details}
	})
}

func TestHandleTransferDebitedEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleTransferDebitedEvent", func(accountCH chan es.Command, details mTDetails) es.Command {
		handler := NewEventHandler(accountCH, nil)
		go handler.HandleTransferDebitedEvent(&TransferDebitedEvent{WithGuid: es.WithGuid{details.Transaction}, mTDetails: details})
		return &CreditAccountBecauseOfTransferCommand{WithGuid: es.WithGuid{details.To}, mTDetails: details}
	})
}

func TestHandleAccountDebitedBecauseOfTransferEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleAccountDebitedBecauseOfTransferEvent", func(transferCH chan es.Command, details mTDetails) es.Command {
		handler := NewEventHandler(nil, transferCH)
		go handler.HandleAccountDebitedBecauseOfTransferEvent(&AccountDebitedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.From}, mTDetails: details})
		return &DebitedTransferCommand{WithGuid: es.WithGuid{details.Transaction}, mTDetails: details}
	})
}

func TestHandleAccountDebitedBecauseOfTransferFailedEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleAccountDebitedBecauseOfTransferFailedEvent", func(transferCH chan es.Command, details mTDetails) es.Command {
		handler := NewEventHandler(nil, transferCH)
		go handler.HandleAccountDebitedBecauseOfTransferFailedEvent(&AccountDebitedBecauseOfTransferFailedEvent{WithGuid: es.WithGuid{details.From}, mTDetails: details})
		return &FailedTransferCommand{WithGuid: es.WithGuid{details.Transaction}, mTDetails: details}
	})
}

func TestHandleAccountCreditedBecauseOfTransferEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleAccountCreditedBecauseOfTransferEvent", func(transferCH chan es.Command, details mTDetails) es.Command {
		handler := NewEventHandler(nil, transferCH)
		go handler.HandleAccountCreditedBecauseOfTransferEvent(&AccountCreditedBecauseOfTransferEvent{WithGuid: es.WithGuid{details.To}, mTDetails: details})
		return &CompletedTransferCommand{WithGuid: es.WithGuid{details.Transaction}, mTDetails: details}
	})
}
