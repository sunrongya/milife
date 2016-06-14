package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Account Events
// --------------

type AccountOpenedEvent struct {
	es.WithGuid
	Name    string
	Card    BankCard
	Balance Money
}

type AccountCreditedEvent struct {
	es.WithGuid
	Amount Money
}

type AccountDebitedEvent struct {
	es.WithGuid
	Amount Money
}

type AccountDebitFailedEvent struct {
	es.WithGuid
}

type AccountDebitedBecauseOfTransferEvent struct {
	es.WithGuid
	mTDetails
}

type AccountDebitedBecauseOfTransferFailedEvent struct {
	es.WithGuid
	mTDetails
}

type AccountCreditedBecauseOfTransferEvent struct {
	es.WithGuid
	mTDetails
}

// ---------------------
// Money Transfer Events
// ---------------------

type TransferCreatedEvent struct {
	es.WithGuid
	mTDetails
}

type TransferDebitedEvent struct {
	es.WithGuid
	mTDetails
}

type TransferCompletedEvent struct {
	es.WithGuid
	mTDetails
}

type TransferFailedEvent struct {
	es.WithGuid
	mTDetails
}
