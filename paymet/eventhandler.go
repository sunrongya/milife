package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type EventHandler struct {
	_accChan   chan<- es.Command
	_transChan chan<- es.Command
}

func (this *EventHandler) HandleTransferCreatedEvent(event *TransferCreatedEvent) {
	this._accChan <- &DebitAccountBecauseOfTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.From},
	}
}

func (this *EventHandler) HandleTransferDebitedEvent(event *TransferDebitedEvent) {
	this._accChan <- &CreditAccountBecauseOfTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.To},
	}
}

func (this *EventHandler) HandleAccountDebitedBecauseOfTransferEvent(event *AccountDebitedBecauseOfTransferEvent) {
	this._transChan <- &DebitedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func (this *EventHandler) HandleAccountDebitedBecauseOfTransferFailedEvent(event *AccountDebitedBecauseOfTransferFailedEvent) {
	this._transChan <- &FailedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func (this *EventHandler) HandleAccountCreditedBecauseOfTransferEvent(event *AccountCreditedBecauseOfTransferEvent) {
	this._transChan <- &CompletedTransferCommand{
		mTDetails: event.mTDetails,
		WithGuid:  es.WithGuid{Guid: event.Transaction},
	}
}

func NewEventHandler(accChan chan<- es.Command, transChan chan<- es.Command) *EventHandler {
	return &EventHandler{_accChan: accChan, _transChan: transChan}
}
