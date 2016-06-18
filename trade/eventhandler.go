package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type PurchaseEventHandler struct {
	goodsChan    chan<- es.Command
	purchaseChan chan<- es.Command
}

func NewPurchaseEventHandler(goodsChan, purchaseChan chan<- es.Command) *PurchaseEventHandler {
	return &PurchaseEventHandler{
		goodsChan:    goodsChan,
		purchaseChan: purchaseChan,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseCreatedEvent(event *GoodsPurchaseCreatedEvent) {
	this.goodsChan <- &PurchaseGoodsBecauseOfPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Goods},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	this.purchaseChan <- &CompleteGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
	this.purchaseChan <- &FailGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

type CommentEventHandler struct {
	goodsChan   chan<- es.Command
	commentChan chan<- es.Command
}

func NewCommentEventHandler(goodsChan, commentChan chan<- es.Command) *CommentEventHandler {
	return &CommentEventHandler{
		goodsChan:   goodsChan,
		commentChan: commentChan,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentCreatedEvent(event *GoodsCommentCreatedEvent) {
	this.goodsChan <- &CommentGoodsBecauseOfCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Goods},
		CommentDetails: event.CommentDetails,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentCompletedEvent(event *GoodsCommentCompletedEvent) {
	this.commentChan <- &CompleteGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentFailedEvent(event *GoodsCommentFailedEvent) {
	this.commentChan <- &FailGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

type PaymetEventHandler struct {
	PaymetService
	goodsChan chan<- es.Command
	orderChan chan<- es.Command
}

func NewPaymetEventHandler(paymetService PaymetService, goodsCh, orderCh chan<- es.Command) *PaymetEventHandler {
	return &PaymetEventHandler{PaymetService: paymetService, goodsChan: goodsCh, orderChan: orderCh}
}

type PaymetService interface {
	Transfer(Money, BankAccount, BankAccount, func(bool))
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
