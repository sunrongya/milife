package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type GoodsPurchaseService struct {
	es.Service
}

func NewGoodsPurchaseService(store es.EventStore) *GoodsPurchaseService {
	service := &GoodsPurchaseService{
		Service: es.NewService(store, NewGoodsPurchase),
	}
	return service
}

func (g *GoodsPurchaseService) CreateGoodsPurchase(goods, user es.Guid, quantity Quantity) es.Guid {
	guid := es.NewGuid()
	c := &CreateGoodsPurchaseCommand{
		WithGuid: es.WithGuid{guid},
		PurchaseDetails: PurchaseDetails{
			User:     user,
			Goods:    goods,
			Purchase: guid,
			Quantity: quantity,
		},
	}
	g.PublishCommand(c)
	return guid
}
