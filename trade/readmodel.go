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
	_repository es.ReadRepository
}

func NewGoodsProjector(repository es.ReadRepository) *GoodsProjector {
	return &GoodsProjector{_repository: repository}
}

func (this *GoodsProjector) HandleGoodsPublishedEvent(event *GoodsPublishedEvent) {
	goods := &RGoods{
		Id:       event.GetGuid(),
		Name:     event.Name,
		Price:    event.Price,
		Quantity: event.Quantity,
		SN:       event.SN,
		State:    Published,
	}
	this._repository.Save(goods.Id, goods)
}

func (this *GoodsProjector) HandleGoodsAuditedPassEvent(event *GoodsAuditedPassEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = AuditedPass
	})
}

func (this *GoodsProjector) HandleGoodsAuditedNoPassEvent(event *GoodsAuditedNoPassEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = AuditedNoPass
	})
}

func (this *GoodsProjector) HandleGoodsOnlinedEvent(event *GoodsOnlinedEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = Onlined
	})
}

func (this *GoodsProjector) HandleGoodsOfflinedEvent(event *GoodsOfflinedEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.State = Offlined
	})
}

func (this *GoodsProjector) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.Purchases += event.Quantity
		goods.SuccessedPurchaseRecords = append(goods.SuccessedPurchaseRecords, &PurchaseRecord{
			Quantity: event.Quantity,
			User:     event.User,
			Purchase: event.Purchase,
		})
	})
}

func (this *GoodsProjector) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		goods.FailuredPurchaseRecords = append(goods.FailuredPurchaseRecords, &PurchaseRecord{
			Quantity: event.Quantity,
			User:     event.User,
			Purchase: event.Purchase,
		})
	})
}

func (this *GoodsProjector) HandleGoodsCommentSuccessedEvent(event *GoodsCommentSuccessedEvent) {
}

func (this *GoodsProjector) HandleGoodsCommentFailuredEvent(event *GoodsCommentFailuredEvent) {
}

func (this *GoodsProjector) HandlePaymetGoodsCompletedBecauseOfOrderEvent(event *PaymetGoodsCompletedBecauseOfOrderEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		for _, record := range goods.SuccessedPurchaseRecords {
			if record.Purchase == event.Purchase {
				record.OrderState = OrderPaymetCompleted
			}
		}
	})
}

func (this *GoodsProjector) HandlePaymetGoodsFailedBecauseOfOrderEvent(event *PaymetGoodsFailedBecauseOfOrderEvent) {
	this.do(event.GetGuid(), func(goods *RGoods) {
		for _, record := range goods.SuccessedPurchaseRecords {
			if record.Purchase == event.Purchase {
				record.OrderState = OrderPaymetFailed
				goods.Purchases -= event.Quantity
			}
		}
	})
}

func (this *GoodsProjector) do(id es.Guid, assignRGoodsFn func(*RGoods)) {
	i, err := this._repository.Find(id)
	if err != nil {
		return
	}
	goods := i.(*RGoods)
	assignRGoodsFn(goods)
	this._repository.Save(id, goods)
}
