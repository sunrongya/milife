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
	name          string
	price         Money
	quantity      Quantity
	sn            SN
	state         State
	purchaseLimit map[es.Guid]Quantity
	comments      map[es.Guid]es.Guid
}

var _ es.Aggregate = (*Goods)(nil)

func NewGoods() es.Aggregate {
	return &Goods{
		purchaseLimit: make(map[es.Guid]Quantity),
		comments:      make(map[es.Guid]es.Guid),
	}
}

func (g *Goods) ProcessPublishGoodsCommand(command *PublishGoodsCommand) []es.Event {
	return []es.Event{
		&GoodsPublishedEvent{
			Name:     command.Name,
			Price:    command.Price,
			Quantity: command.Quantity,
			SN:       command.SN,
		},
	}
}

func (g *Goods) ProcessAuditGoodsCommand(command *AuditGoodsCommand) []es.Event {
	if g.state != Published {
		panic(fmt.Errorf("Can't process AuditGoodsCommand of state:%s", g.state))
	}
	if command.IsPass {
		return []es.Event{&GoodsAuditedPassEvent{}}
	}
	return []es.Event{&GoodsAuditedNoPassEvent{}}
}

func (g *Goods) ProcessOnlineGoodsCommand(command *OnlineGoodsCommand) []es.Event {
	if g.state != AuditedPass {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsOnlinedEvent{}}
}

func (g *Goods) ProcessOfflineGoodsCommand(command *OfflineGoodsCommand) []es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsOfflinedEvent{}}
}

func (g *Goods) ProcessPurchaseGoodsBecauseOfPurchaseCommand(command *PurchaseGoodsBecauseOfPurchaseCommand) []es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process PurchaseGoodsBecauseOfPurchaseCommand of state:%s", g.state))
	}
	if g.checkPurchaseQuantity(command) {
		return []es.Event{&GoodsPurchaseSuccessedEvent{PurchaseDetails: command.PurchaseDetails}}
	}
	return []es.Event{&GoodsPurchaseFailuredEvent{PurchaseDetails: command.PurchaseDetails}}
}

func (g *Goods) checkPurchaseQuantity(command *PurchaseGoodsBecauseOfPurchaseCommand) bool {
	userQuantity := g.purchaseLimit[command.User] + command.Quantity
	return g.quantity >= userQuantity && userQuantity <= 3
}

func (g *Goods) ProcessCompletePaymetGoodsBecauseOfOrderCommand(command *CompletePaymetGoodsBecauseOfOrderCommand) []es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process CompletePaymetGoodsBecauseOfOrderCommand of state:%s", g.state))
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

func (g *Goods) ProcessFailPaymetGoodsBecauseOfOrderCommand(command *FailPaymetGoodsBecauseOfOrderCommand) []es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process FailPaymetGoodsBecauseOfOrderCommand of state:%s", g.state))
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

func (g *Goods) ProcessCommentGoodsBecauseOfCommentCommand(command *CommentGoodsBecauseOfCommentCommand) []es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process CommentGoodsBecauseOfCommentCommand of state:%s", g.state))
	}
	if _, ok := g.comments[command.Purchase]; ok {
		return []es.Event{&GoodsCommentSuccessedEvent{CommentDetails: command.CommentDetails}}
	}
	return []es.Event{&GoodsCommentFailuredEvent{CommentDetails: command.CommentDetails}}
}

func (g *Goods) HandleGoodsPublishedEvent(event *GoodsPublishedEvent) {
	g.name, g.price, g.quantity, g.sn, g.state = event.Name, event.Price, event.Quantity, event.SN, Published
}

func (g *Goods) HandleGoodsAuditedPassEvent(event *GoodsAuditedPassEvent) {
	g.state = AuditedPass
}

func (g *Goods) HandleGoodsAuditedNoPassEvent(event *GoodsAuditedNoPassEvent) {
	g.state = AuditedNoPass
}

func (g *Goods) HandleGoodsOnlinedEvent(event *GoodsOnlinedEvent) {
	g.state = Onlined
}

func (g *Goods) HandleGoodsOfflinedEvent(event *GoodsOfflinedEvent) {
	g.state = Offlined
}

func (g *Goods) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	g.quantity -= event.Quantity
}

func (g *Goods) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
}

func (g *Goods) HandleGoodsCommentSuccessedEvent(event *GoodsCommentSuccessedEvent) {
	delete(g.comments, event.Purchase)
}

func (g *Goods) HandleGoodsCommentFailuredEvent(event *GoodsCommentFailuredEvent) {
}

func (g *Goods) HandlePaymetGoodsCompletedBecauseOfOrderEvent(event *PaymetGoodsCompletedBecauseOfOrderEvent) {
	g.comments[event.Purchase] = event.User
}

func (g *Goods) HandlePaymetGoodsFailedBecauseOfOrderEvent(event *PaymetGoodsFailedBecauseOfOrderEvent) {
}
