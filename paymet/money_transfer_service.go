package paymet

import (
	es "github.com/sunrongya/eventsourcing"
)

type TransferService struct {
	es.Service
}

func NewTransferService(store es.EventStore) *TransferService {
	acc := &TransferService{
		Service: es.NewService(store, NewMoneyTransfer),
	}
	return acc
}

func (this *TransferService) Transfer(amount Money, from, to es.Guid) es.Guid {
	guid := es.NewGuid()
	c := &CreateTransferCommand{
		WithGuid: es.WithGuid{guid},
		mTDetails: mTDetails{
			From:        from,
			To:          to,
			Amount:      amount,
			Transaction: guid,
		},
	}
	this.PublishCommand(c)
	return guid
}
