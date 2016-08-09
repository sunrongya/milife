package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

type PurchaseEventHandler struct {
	_goodsChan    chan<- es.Command
	_purchaseChan chan<- es.Command
}

func NewPurchaseEventHandler(goodsChan, purchaseChan chan<- es.Command) *PurchaseEventHandler {
	return &PurchaseEventHandler{
		_goodsChan:    goodsChan,
		_purchaseChan: purchaseChan,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseCreatedEvent(event *GoodsPurchaseCreatedEvent) {
	this._goodsChan <- &PurchaseGoodsBecauseOfPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Goods},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseSuccessedEvent(event *GoodsPurchaseSuccessedEvent) {
	this._purchaseChan <- &CompleteGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

func (this *PurchaseEventHandler) HandleGoodsPurchaseFailuredEvent(event *GoodsPurchaseFailuredEvent) {
	this._purchaseChan <- &FailGoodsPurchaseCommand{
		WithGuid:        es.WithGuid{Guid: event.Purchase},
		PurchaseDetails: event.PurchaseDetails,
	}
}

type CommentEventHandler struct {
	_goodsChan   chan<- es.Command
	_commentChan chan<- es.Command
}

func NewCommentEventHandler(goodsChan, commentChan chan<- es.Command) *CommentEventHandler {
	return &CommentEventHandler{
		_goodsChan:   goodsChan,
		_commentChan: commentChan,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentCreatedEvent(event *GoodsCommentCreatedEvent) {
	this._goodsChan <- &CommentGoodsBecauseOfCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Goods},
		CommentDetails: event.CommentDetails,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentCompletedEvent(event *GoodsCommentCompletedEvent) {
	this._commentChan <- &CompleteGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

func (this *CommentEventHandler) HandleGoodsCommentFailedEvent(event *GoodsCommentFailedEvent) {
	this._commentChan <- &FailGoodsCommentCommand{
		WithGuid:       es.WithGuid{Guid: event.Comment},
		CommentDetails: event.CommentDetails,
	}
}

type PaymetEventHandler struct {
	PaymetService
	_goodsChan chan<- es.Command
	_orderChan chan<- es.Command
}

func NewPaymetEventHandler(paymetService PaymetService, goodsCh, orderCh chan<- es.Command) *PaymetEventHandler {
	return &PaymetEventHandler{PaymetService: paymetService, _goodsChan: goodsCh, _orderChan: orderCh}
}

type PaymetService interface {
	Transfer(Money, BankAccount, BankAccount, func(bool))
}

func (this *PaymetEventHandler) HandleOrderPaymetCreatedEvent(event *OrderPaymetCreatedEvent) {
	go this.Transfer(event.Price, event.UserAccount, event.ManagedAccount, func(isOk bool) {
		if isOk {
			this._orderChan <- &CompleteOrderPaymetCommand{
				WithGuid:       es.WithGuid{Guid: event.GetGuid()},
				Price:          event.Price,
				User:           event.User,
				UserAccount:    event.UserAccount,
				ManagedAccount: event.ManagedAccount,
			}
		} else {
			this._orderChan <- &FailOrderPaymetCommand{
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
		this._goodsChan <- &CompletePaymetGoodsBecauseOfOrderCommand{
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
		this._goodsChan <- &FailPaymetGoodsBecauseOfOrderCommand{
			WithGuid: es.WithGuid{Guid: v.Goods},
			User:     event.User,
			Order:    event.GetGuid(),
			Purchase: v.Purchase,
			Quantity: v.Quantity,
		}
	}
}
