package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type EventHandler struct {
	goodsChan    chan<- es.Command
	purchaseChan chan<- es.Command
	commentChan  chan<- es.Command
}

func (this *EventHandler) HandleGoodsPurchaseCreatedEvent(event *GoodsPurchaseCreatedEvent) {
	this.goodsChan <- &PurchaseGoodsBecauseOfPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Goods},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *EventHandler) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	this.purchaseChan <- &CompleteGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *EventHandler) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
	this.purchaseChan <- &FailGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *EventHandler) HandleGoodsCommentCreatedEvent(event *GoodsCommentCreatedEvent) {
	this.goodsChan <- &CommentGoodsBecauseOfCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Goods},
		CommentDetails: event.CommentDetails,
	}
}

func (this *EventHandler) HandleGoodsCommentCompletedEvent(event *GoodsCommentCompletedEvent) {
	this.commentChan <- &CompleteGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

func (this *EventHandler) HandleGoodsCommentFailedEvent(event *GoodsCommentFailedEvent) {
	this.commentChan <- &FailGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

func NewEventHandler(goodsChan, purchaseChan, commentChan chan<- es.Command) *EventHandler {
	return &EventHandler{
		goodsChan:    goodsChan,
		purchaseChan: purchaseChan,
		commentChan:  commentChan,
	}
}

type PaymetEventHandler struct {
	PaymetService
	goodsChan chan<- es.Command
	orderChan chan<- es.Command
}

func (this *PaymetEventHandler) HandleOrderPaymetCreatedEvent(event *OrderPaymetCreatedEvent) {
	go this.Transfer(event.Price, event.UserAccount, event.ManagedAccount, func(isOk bool) {
		if isOk {
			this.orderChan <- &CompleteOrderPaymetCommand{
				WithGuid:       es.WithGuid{Guid: event.GetGuid()},
				Price:          event.Price,
				User:           event.User,
				UserAccount:    event.UserAccount,
				ManagedAccount: event.ManagedAccount,
			}
		} else {
			this.orderChan <- &FailOrderPaymetCommand{
				WithGuid:       es.WithGuid{Guid: event.GetGuid()},
				Price:          event.Price,
				User:           event.User,
				UserAccount:    event.UserAccount,
				ManagedAccount: event.ManagedAccount,
			}
		}
	})
}

func (this *PaymetEventHandler) HandleOrderPaymetCompletedEvent(event *OrderPaymetCompletedEvent) {
	for _, v := range event.OrderItems {
		this.goodsChan <- &CompletePaymetGoodsBecauseOfOrderCommand{
			WithGuid: es.WithGuid{Guid: v.Goods},
			User:     event.User,
			Order:    event.GetGuid(),
			Purchase: v.Purchase,
			Quantity: v.Quantity,
		}
	}
}

func (this *PaymetEventHandler) HandleOrderPaymetFailedEvent(event *OrderPaymetFailedEvent) {
	for _, v := range event.OrderItems {
		this.goodsChan <- &FailPaymetGoodsBecauseOfOrderCommand{
			WithGuid: es.WithGuid{Guid: v.Goods},
			User:     event.User,
			Order:    event.GetGuid(),
			Purchase: v.Purchase,
			Quantity: v.Quantity,
		}
	}
}

func NewPaymetEventHandler(paymetService PaymetService, goodsCh, orderCh chan<- es.Command) *PaymetEventHandler {
	return &PaymetEventHandler{PaymetService: paymetService, goodsChan: goodsCh, orderChan: orderCh}
}

type PaymetService interface {
	Transfer(Money, BankAccount, BankAccount, func(bool))
}

type NullPaymetService struct {
}

func (n *NullPaymetService) Transfer(amount Money, userAccount, managedAccount BankAccount, completeFn func(bool)) {
	completeFn(true)
}
