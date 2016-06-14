package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type EventHandler struct {
	accChan   chan<- es.Command
	transChan chan<- es.Command
}

func (this *EventHandler) HandleTransferCreatedEvent(event *TransferCreatedEvent) {
	this.accChan <- &DebitAccountBecauseOfTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.From},
	}
}

func (this *EventHandler) HandleTransferDebitedEvent(event *TransferDebitedEvent) {
	this.accChan <- &CreditAccountBecauseOfTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.To},
	}
}

func (this *EventHandler) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	this.transChan <- &DebitedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func (this *EventHandler) HandleAccountDebitedBecauseOfTransferFailedEvent(event *AccountDebitedBecauseOfTransferFailedEvent) {
	this.transChan <- &FailedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func (this *EventHandler) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	this.transChan <- &CompletedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func NewEventHandler(accChan chan<- es.Command, transChan chan<- es.Command) *EventHandler {
	return &EventHandler{accChan: accChan, transChan: transChan}
}
