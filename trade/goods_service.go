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

func (this *GoodsService) PublishGoods(name string, price Money, quantity Quantity, sn SN) es.Guid {
	guid := es.NewGuid()
	c := &PublishGoodsCommand{
		WithGuid: es.WithGuid{guid},
		Name:     name,
		Price:    price,
		Quantity: quantity,
		SN:       sn,
	}
	this.PublishCommand(c)
	return guid
}

func (this *GoodsService) AuditGoods(guid es.Guid, isPass bool) {
	c := &AuditGoodsCommand{
		WithGuid: es.WithGuid{guid},
		IsPass:   isPass,
	}
	this.PublishCommand(c)
}

func (this *GoodsService) OnlineGoods(guid es.Guid) {
	c := &OnlineGoodsCommand{
		WithGuid: es.WithGuid{guid},
	}
	this.PublishCommand(c)
}

func (this *GoodsService) OfflineGoods(guid es.Guid) {
	c := &OfflineGoodsCommand{
		WithGuid: es.WithGuid{guid},
	}
	this.PublishCommand(c)
}
