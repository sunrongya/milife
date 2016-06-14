package trade

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestGoodsPurchaseServiceDoPublishGoods(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		details := PurchaseDetails{
			User:     es.NewGuid(),
			Goods:    es.NewGuid(),
			Purchase: es.NewGuid(),
			Quantity: 5,
		}

		gs := GoodsPurchaseService{Service: service}
		details.Purchase = gs.CreateGoodsPurchase(details.Goods, details.User, details.Quantity)
		return &CreateGoodsPurchaseCommand{WithGuid: es.WithGuid{details.Purchase}, PurchaseDetails: details}
	})
}
