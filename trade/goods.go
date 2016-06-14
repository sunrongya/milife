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

func (g *Goods) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *GoodsPublishedEvent:
			g.name, g.price, g.quantity, g.sn, g.state = e.Name, e.Price, e.Quantity, e.SN, Published
		case *GoodsAuditedPassEvent:
			g.state = AuditedPass
		case *GoodsAuditedNoPassEvent:
			g.state = AuditedNoPass
		case *GoodsOnlinedEvent:
			g.state = Onlined
		case *GoodsOfflinedEvent:
			g.state = Offlined
		case *GoodsPurchaseSuccessedEvent:
			g.quantity -= e.Quantity
		case *GoodsPurchaseFailuredEvent:
		case *GoodsCommentSuccessedEvent:
			delete(g.comments, e.Purchase)
		case *GoodsCommentFailuredEvent:
		case *PaymetGoodsCompletedBecauseOfOrderEvent:
			g.comments[e.Purchase] = e.User
		case *PaymetGoodsFailedBecauseOfOrderEvent:
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	g.SetVersion(len(events))
}

func (g *Goods) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *PublishGoodsCommand:
		event = g.processPublishGoodsCommand(c)
	case *AuditGoodsCommand:
		event = g.processAuditGoodsCommand(c)
	case *OnlineGoodsCommand:
		event = g.processOnlineGoodsCommand(c)
	case *OfflineGoodsCommand:
		event = g.processOfflineGoodsCommand(c)
	case *PurchaseGoodsBecauseOfPurchaseCommand:
		event = g.processPurchaseGoodsBecauseOfPurchaseCommand(c)
	case *CompletePaymetGoodsBecauseOfOrderCommand:
		event = g.processCompletePaymetGoodsBecauseOfOrderCommand(c)
	case *FailPaymetGoodsBecauseOfOrderCommand:
		event = g.processFailPaymetGoodsBecauseOfOrderCommand(c)
	case *CommentGoodsBecauseOfCommentCommand:
		event = g.processCommentGoodsBecauseOfCommentCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (g *Goods) processPublishGoodsCommand(command *PublishGoodsCommand) es.Event {
	return &GoodsPublishedEvent{
		Name:     command.Name,
		Price:    command.Price,
		Quantity: command.Quantity,
		SN:       command.SN,
	}
}

func (g *Goods) processAuditGoodsCommand(command *AuditGoodsCommand) es.Event {
	if g.state != Published {
		panic(fmt.Errorf("Can't process AuditGoodsCommand of state:%s", g.state))
	}
	if command.IsPass {
		return &GoodsAuditedPassEvent{}
	}
	return &GoodsAuditedNoPassEvent{}
}

func (g *Goods) processOnlineGoodsCommand(command *OnlineGoodsCommand) es.Event {
	if g.state != AuditedPass {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", g.state))
	}
	return &GoodsOnlinedEvent{}
}

func (g *Goods) processOfflineGoodsCommand(command *OfflineGoodsCommand) es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process OnlineGoodsCommand of state:%s", g.state))
	}
	return &GoodsOfflinedEvent{}
}

func (g *Goods) processPurchaseGoodsBecauseOfPurchaseCommand(command *PurchaseGoodsBecauseOfPurchaseCommand) es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process PurchaseGoodsBecauseOfPurchaseCommand of state:%s", g.state))
	}
	if g.checkPurchaseQuantity(command) {
		return &GoodsPurchaseSuccessedEvent{PurchaseDetails: command.PurchaseDetails}
	}
	return &GoodsPurchaseFailuredEvent{PurchaseDetails: command.PurchaseDetails}
}

func (g *Goods) checkPurchaseQuantity(command *PurchaseGoodsBecauseOfPurchaseCommand) bool {
	userQuantity := g.purchaseLimit[command.User] + command.Quantity
	return g.quantity >= userQuantity && userQuantity <= 3
}

func (g *Goods) processCompletePaymetGoodsBecauseOfOrderCommand(command *CompletePaymetGoodsBecauseOfOrderCommand) es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process CompletePaymetGoodsBecauseOfOrderCommand of state:%s", g.state))
	}
	return &PaymetGoodsCompletedBecauseOfOrderEvent{
		User:     command.User,
		Order:    command.Order,
		Purchase: command.Purchase,
		Quantity: command.Quantity,
	}
}

func (g *Goods) processFailPaymetGoodsBecauseOfOrderCommand(command *FailPaymetGoodsBecauseOfOrderCommand) es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process FailPaymetGoodsBecauseOfOrderCommand of state:%s", g.state))
	}
	return &PaymetGoodsFailedBecauseOfOrderEvent{
		User:     command.User,
		Order:    command.Order,
		Purchase: command.Purchase,
		Quantity: command.Quantity,
	}
}

func (g *Goods) processCommentGoodsBecauseOfCommentCommand(command *CommentGoodsBecauseOfCommentCommand) es.Event {
	if g.state != Onlined {
		panic(fmt.Errorf("Can't process CommentGoodsBecauseOfCommentCommand of state:%s", g.state))
	}
	if _, ok := g.comments[command.Purchase]; ok {
		return &GoodsCommentSuccessedEvent{CommentDetails: command.CommentDetails}
	}
	return &GoodsCommentFailuredEvent{CommentDetails: command.CommentDetails}
}

func NewGoods() es.Aggregate {
	return &Goods{
		purchaseLimit: make(map[es.Guid]Quantity),
		comments:      make(map[es.Guid]es.Guid),
	}
}
