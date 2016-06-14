package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Account Commands
// --------------

type OpenAccountCommand struct {
	es.WithGuid
	Name    string
	Card    BankCard
	Balance Money
}

type CreditAccountCommand struct {
	es.WithGuid
	Amount Money
}

type DebitAccountCommand struct {
	es.WithGuid
	Amount Money
}

type DebitAccountBecauseOfTransferCommand struct {
	es.WithGuid
	mTDetails
}

type CreditAccountBecauseOfTransferCommand struct {
	es.WithGuid
	mTDetails
}

//-----------------------
//Money Transfer Commands
//-----------------------

type CreateTransferCommand struct {
	es.WithGuid
	mTDetails
}

type DebitedTransferCommand struct {
	es.WithGuid
	mTDetails
}

type CompletedTransferCommand struct {
	es.WithGuid
	mTDetails
}

type FailedTransferCommand struct {
	es.WithGuid
	mTDetails
}
