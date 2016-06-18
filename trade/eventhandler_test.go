package trade

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func checkCommandForChannel(t *testing.T, ch chan es.Command, command es.Command, methodName string) {
	select {
	case c := <-ch:
		assert.Equal(t, c, command, methodName)
	case <-time.After(1 * time.Second):
		t.Error(methodName)
	}
}

func testPurchaseHandleEvent(t *testing.T, methodName string, doPurchaseHandle func(chan es.Command, PurchaseDetails) es.Command) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 5,
	}
	ch := make(chan es.Command)
	command := doPurchaseHandle(ch, details)
	checkCommandForChannel(t, ch, command, methodName)
}

func testCommentHandleEvent(t *testing.T, methodName string, doCommentHandle func(chan es.Command, CommentDetails) es.Command) {
	details := CommentDetails{
		User:    es.NewGuid(),
		Goods:   es.NewGuid(),
		Comment: es.NewGuid(),
		Content: "dadfds",
		Time:    time.Now(),
	}
	ch := make(chan es.Command)
	command := doCommentHandle(ch, details)
	checkCommandForChannel(t, ch, command, methodName)
}

// Goods Purchase
func TestHandleGoodsPurchaseCreatedEvent(t *testing.T) {
	testPurchaseHandleEvent(t, "TestHandleGoodsPurchaseCreatedEvent", func(goodsCH chan es.Command, details PurchaseDetails) es.Command {
		handler := NewPurchaseEventHandler(goodsCH, nil)
		go handler.HandleGoodsPurchaseCreatedEvent(&GoodsPurchaseCreatedEvent{WithGuid: es.WithGuid{details.Purchase}, PurchaseDetails: details})
		return &PurchaseGoodsBecauseOfPurchaseCommand{WithGuid: es.WithGuid{details.Goods}, PurchaseDetails: details}
	})
}

func TestHandleGoodsPurchaseSuccessedEvent(t *testing.T) {
	testPurchaseHandleEvent(t, "TestHandleGoodsPurchaseSuccessedEvent", func(purchaseCH chan es.Command, details PurchaseDetails) es.Command {
		handler := NewPurchaseEventHandler(nil, purchaseCH)
		go handler.HandleGoodsPurchaseSuccessedEvent(&GoodsPurchaseSuccessedEvent{WithGuid: es.WithGuid{details.Goods}, PurchaseDetails: details})
		return &CompleteGoodsPurchaseCommand{WithGuid: es.WithGuid{details.Purchase}, PurchaseDetails: details}
	})
}

func TestHandleGoodsPurchaseFailuredEvent(t *testing.T) {
	testPurchaseHandleEvent(t, "TestHandleGoodsPurchaseFailuredEvent", func(purchaseCH chan es.Command, details PurchaseDetails) es.Command {
		handler := NewPurchaseEventHandler(nil, purchaseCH)
		go handler.HandleGoodsPurchaseFailuredEvent(&GoodsPurchaseFailuredEvent{WithGuid: es.WithGuid{details.Goods}, PurchaseDetails: details})
		return &FailGoodsPurchaseCommand{WithGuid: es.WithGuid{details.Purchase}, PurchaseDetails: details}
	})
}

// Goods Comment
func TestHandleGoodsCommentCreatedEvent(t *testing.T) {
	testCommentHandleEvent(t, "TestHandleGoodsCommentCreatedEvent", func(goodsCH chan es.Command, details CommentDetails) es.Command {
		handler := NewCommentEventHandler(goodsCH, nil)
		go handler.HandleGoodsCommentCreatedEvent(&GoodsCommentCreatedEvent{WithGuid: es.WithGuid{details.Comment}, CommentDetails: details})
		return &CommentGoodsBecauseOfCommentCommand{WithGuid: es.WithGuid{details.Goods}, CommentDetails: details}
	})
}

func TestHandleGoodsCommentSuccessedEvent(t *testing.T) {
	testCommentHandleEvent(t, "TestHandleGoodsCommentSuccessedEvent", func(commentCH chan es.Command, details CommentDetails) es.Command {
		handler := NewCommentEventHandler(nil, commentCH)
		go handler.HandleGoodsCommentCompletedEvent(&GoodsCommentCompletedEvent{WithGuid: es.WithGuid{details.Goods}, CommentDetails: details})
		return &CompleteGoodsCommentCommand{WithGuid: es.WithGuid{details.Comment}, CommentDetails: details}
	})
}

func TestHandleGoodsCommentFailuredEvent(t *testing.T) {
	testCommentHandleEvent(t, "TestHandleGoodsCommentFailuredEvent", func(commentCH chan es.Command, details CommentDetails) es.Command {
		handler := NewCommentEventHandler(nil, commentCH)
		go handler.HandleGoodsCommentFailedEvent(&GoodsCommentFailedEvent{WithGuid: es.WithGuid{details.Goods}, CommentDetails: details})
		return &FailGoodsCommentCommand{WithGuid: es.WithGuid{details.Comment}, CommentDetails: details}
	})
}
