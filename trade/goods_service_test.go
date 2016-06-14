package trade

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestGoodsServiceDoPublishGoods(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := GoodsService{Service: service}
		guid := gs.PublishGoods("mmm", 50, 100, "20160601333")
		return &PublishGoodsCommand{WithGuid: es.WithGuid{guid}, Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"}
	})
}

func TestGoodsServiceDoAuditGoods(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := GoodsService{Service: service}
		guid := es.NewGuid()
		as.AuditGoods(guid, true)
		return &AuditGoodsCommand{WithGuid: es.WithGuid{Guid: guid}, IsPass: true}
	})

	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := GoodsService{Service: service}
		guid := es.NewGuid()
		as.AuditGoods(guid, false)
		return &AuditGoodsCommand{WithGuid: es.WithGuid{Guid: guid}, IsPass: false}
	})
}

func TestGoodsServiceDoOnlineGoods(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := GoodsService{Service: service}
		guid := es.NewGuid()
		as.OnlineGoods(guid)
		return &OnlineGoodsCommand{WithGuid: es.WithGuid{Guid: guid}}
	})
}

func TestGoodsServiceDoOfflineGoods(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := GoodsService{Service: service}
		guid := es.NewGuid()
		as.OfflineGoods(guid)
		return &OfflineGoodsCommand{WithGuid: es.WithGuid{Guid: guid}}
	})
}
