package trade

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type State string

const (
	Published     = State("Published")
	AuditedPass   = State("AuditedPass")
	AuditedNoPass = State("AuditedNoPass")
	Onlined       = State("Onlined")
	Offlined      = State("Offlined")
	Completed     = State("Completed")
)

// TODO 后面考虑添加商品SN不能重复，修改商品需求
type Goods struct {
	es.BaseAggregate
	_name          string
	_price         Money
	_quantity      Quantity
	_sn            SN
	_state         State
	_purchaseLimit map[es.Guid]Quantity
	_comments      map[es.Guid]es.Guid
}

var _ es.Aggregate = (*Goods)(nil)

func NewGoods() es.Aggregate {
	return &Goods{
		_purchaseLimit: make(map[es.Guid]Quantity),
		_comments:      make(map[es.Guid]es.Guid),
	}
}

func (this *Goods) ProcessPublishGoodsCommand(command *PublishGoodsCommand) []es.Event {
	return []es.Event{
		&GoodsPublishedEvent{
			Name:     command.Name,
			Price:    command.Price,
			Quantity: command.Quantity,
			SN:       command.SN,
		},
	}
}

func (this *Goods) ProcessAuditGoodsCommand(command *AuditGoodsCommand) []es.Event {
	if this._state != Published {
		panic(fmt.Errorf("Can't process AuditGoodsCommand of state:%s", this._state))
	}
	if command.IsPass {
		return []es.Event{&GoodsAuditedPassEvent{}}
	}
	return []es.Event{&GoodsAuditedNoPassEvent{}}
}

func (this *Goods) ProcessOnlineGoodsCommand(command *OnlineGoodsCommand) []es.Event {
	if this._state != AuditedPass {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsOnlinedEvent{}}
}

func (this *Goods) ProcessOfflineGoodsCommand(command *OfflineGoodsCommand) []es.Event {
	if this._state != Onlined {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsOfflinedEvent{}}
}

func (this *Goods) ProcessPurchaseGoodsBecauseOfPurchaseCommand(command *PurchaseGoodsBecauseOfPurchaseCommand) []es.Event {
	if this._state != Onlined {
		panic(fmt.Errorf("Can't process PurchaseGoodsBecauseOfPurchaseCommand of state:%s", this._state))
	}
	if this.checkPurchaseQuantity(command) {
		return []es.Event{&GoodsPurchaseSuccessedEvent{PurchaseDetails: command.PurchaseDetails}}
	}
	return []es.Event{&GoodsPurchaseFailuredEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (this *Goods) checkPurchaseQuantity(command *PurchaseGoodsBecauseOfPurchaseCommand) bool {
	userQuantity := this._purchaseLimit[command.User] + command.Quantity
	return this._quantity >= userQuantity && userQuantity <= 3
}

func (this *Goods) ProcessCompletePaymetGoodsBecauseOfOrderCommand(command *CompletePaymetGoodsBecauseOfOrderCommand) []es.Event {
	if this._state != Onlined {
		panic(fmt.Errorf("Can't process CompletePaymetGoodsBecauseOfOrderCommand of state:%s", this._state))
	}
	return []es.Event{
		&PaymetGoodsCompletedBecauseOfOrderEvent{
			User:     command.User,
			Order:    command.Order,
			Purchase: command.Purchase,
			Quantity: command.Quantity,
		},
	}
}

func (this *Goods) ProcessFailPaymetGoodsBecauseOfOrderCommand(command *FailPaymetGoodsBecauseOfOrderCommand) []es.Event {
	if this._state != Onlined {
		panic(fmt.Errorf("Can't process FailPaymetGoodsBecauseOfOrderCommand of state:%s", this._state))
	}
	return []es.Event{
		&PaymetGoodsFailedBecauseOfOrderEvent{
			User:     command.User,
			Order:    command.Order,
			Purchase: command.Purchase,
			Quantity: command.Quantity,
		},
	}
}

func (this *Goods) ProcessCommentGoodsBecauseOfCommentCommand(command *CommentGoodsBecauseOfCommentCommand) []es.Event {
	if this._state != Onlined {
		panic(fmt.Errorf("Can't process CommentGoodsBecauseOfCommentCommand of state:%s", this._state))
	}
	if _, ok := this._comments[command.Purchase]; ok {
		return []es.Event{&GoodsCommentSuccessedEvent{CommentDetails: command.CommentDetails}}
	}
	return []es.Event{&GoodsCommentFailuredEvent{CommentDetails: command.CommentDetails}}
}

func (this *Goods) HandleGoodsPublishedEvent(event *GoodsPublishedEvent) {
	this._name, this._price, this._quantity, this._sn, this._state = event.Name, event.Price, event.Quantity, event.SN, Published
}

func (this *Goods) HandleGoodsAuditedPassEvent(event *GoodsAuditedPassEvent) {
	this._state = AuditedPass
}

func (this *Goods) HandleGoodsAuditedNoPassEvent(event *GoodsAuditedNoPassEvent) {
	this._state = AuditedNoPass
}

func (this *Goods) HandleGoodsOnlinedEvent(event *GoodsOnlinedEvent) {
	this._state = Onlined
}

func (this *Goods) HandleGoodsOfflinedEvent(event *GoodsOfflinedEvent) {
	this._state = Offlined
}

func (this *Goods) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	this._quantity -= event.Quantity
}

func (this *Goods) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
}

func (this *Goods) HandleGoodsCommentSuccessedEvent(event *GoodsCommentSuccessedEvent) {
	delete(this._comments, event.Purchase)
}

func (this *Goods) HandleGoodsCommentFailuredEvent(event *GoodsCommentFailuredEvent) {
}

func (this *Goods) HandlePaymetGoodsCompletedBecauseOfOrderEvent(event *PaymetGoodsCompletedBecauseOfOrderEvent) {
	this._comments[event.Purchase] = event.User
}

func (this *Goods) HandlePaymetGoodsFailedBecauseOfOrderEvent(event *PaymetGoodsFailedBecauseOfOrderEvent) {
}
