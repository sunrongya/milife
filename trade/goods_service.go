package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type GoodsService struct {
	es.Service
}

func NewGoodsService(store es.EventStore) *GoodsService {
	service := &GoodsService{
		Service: es.NewService(store, NewGoods),
	}
	return service
}

func (g *GoodsService) PublishGoods(name string, price Money, quantity Quantity, sn SN) es.Guid {
	guid := es.NewGuid()
	c := &PublishGoodsCommand{
		WithGuid: es.WithGuid{guid},
		Name:     name,
		Price:    price,
		Quantity: quantity,
		SN:       sn,
	}
	g.PublishCommand(c)
	return guid
}

func (g *GoodsService) AuditGoods(guid es.Guid, isPass bool) {
	c := &AuditGoodsCommand{
		WithGuid: es.WithGuid{guid},
		IsPass:   isPass,
	}
	g.PublishCommand(c)
}

func (g *GoodsService) OnlineGoods(guid es.Guid) {
	c := &OnlineGoodsCommand{
		WithGuid: es.WithGuid{guid},
	}
	g.PublishCommand(c)
}

func (g *GoodsService) OfflineGoods(guid es.Guid) {
	c := &OfflineGoodsCommand{
		WithGuid: es.WithGuid{guid},
	}
	g.PublishCommand(c)
}
