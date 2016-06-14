package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type PurchaseRecord struct {
	Quantity   Quantity
	User       es.Guid
	Purchase   es.Guid
	OrderState OrderState
}

type RGoods struct {
	Id                       es.Guid
	SN                       SN
	Name                     string
	Price                    Money
	Quantity                 Quantity
	Purchases                Quantity
	State                    State
	SuccessedPurchaseRecords []*PurchaseRecord
	FailuredPurchaseRecords  []*PurchaseRecord
}

type GoodsProjector struct {
	repository es.ReadRepository
}

func NewGoodsProjector(repository es.ReadRepository) *GoodsProjector {
	return &GoodsProjector{repository: repository}
}

func (g *GoodsProjector) HandleGoodsPublishedEvent(event *GoodsPublishedEvent) {
	goods := &RGoods{
		Id:       event.GetGuid(),
		Name:     event.Name,
		Price:    event.Price,
		Quantity: event.Quantity,
		SN:       event.SN,
		State:    Published,
	}
	g.repository.Save(goods.Id, goods)
}

func (g *GoodsProjector) HandleGoodsAuditedPassEvent(event *GoodsAuditedPassEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = AuditedPass
	})
}

func (g *GoodsProjector) HandleGoodsAuditedNoPassEvent(event *GoodsAuditedNoPassEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = AuditedNoPass
	})
}

func (g *GoodsProjector) HandleGoodsOnlinedEvent(event *GoodsOnlinedEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = Onlined
	})
}

func (g *GoodsProjector) HandleGoodsOfflinedEvent(event *GoodsOfflinedEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = Offlined
	})
}

func (g *GoodsProjector) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.Purchases += event.Quantity
		goods.SuccessedPurchaseRecords = append(goods.SuccessedPurchaseRecords, &PurchaseRecord{
			Quantity: event.Quantity,
			User:     event.User,
			Purchase: event.Purchase,
		})
	})
}

func (g *GoodsProjector) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		goods.FailuredPurchaseRecords = append(goods.FailuredPurchaseRecords, &PurchaseRecord{
			Quantity: event.Quantity,
			User:     event.User,
			Purchase: event.Purchase,
		})
	})
}

func (g *GoodsProjector) HandleGoodsCommentSuccessedEvent(event *GoodsCommentSuccessedEvent) {
}

func (g *GoodsProjector) HandleGoodsCommentFailuredEvent(event *GoodsCommentFailuredEvent) {
}

func (g *GoodsProjector) HandlePaymetGoodsCompletedBecauseOfOrderEvent(event *PaymetGoodsCompletedBecauseOfOrderEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		for _, record := range goods.SuccessedPurchaseRecords {
			if record.Purchase == event.Purchase {
				record.OrderState = OrderPaymetCompleted
			}
		}
	})
}

func (g *GoodsProjector) HandlePaymetGoodsFailedBecauseOfOrderEvent(event *PaymetGoodsFailedBecauseOfOrderEvent) {
	g.do(event.GetGuid(), func(goods *RGoods) {
		for _, record := range goods.SuccessedPurchaseRecords {
			if record.Purchase == event.Purchase {
				record.OrderState = OrderPaymetFailed
				goods.Purchases -= event.Quantity
			}
		}
	})
}

func (g *GoodsProjector) do(id es.Guid, assignRGoodsFn func(*RGoods)) {
	i, err := g.repository.Find(id)
	if err != nil {
		return
	}
	goods := i.(*RGoods)
	assignRGoodsFn(goods)
	g.repository.Save(id, goods)
}
