package paymet

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestTransferServicePublishCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}

	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := TransferService{Service: service}
		details.Transaction = as.Transfer(details.Amount, details.From, details.To)
		return &CreateTransferCommand{WithGuid: es.WithGuid{Guid: details.Transaction}, mTDetails: details}
	})
}
