package trade

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestGoodsRestore(t *testing.T) {
	guid := es.NewGuid()
	goods := &Goods{}
	goods.ApplyEvents([]es.Event{
		&GoodsPublishedEvent{WithGuid: es.WithGuid{guid}, Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"},
		&GoodsAuditedPassEvent{WithGuid: es.WithGuid{guid}},
		&GoodsOnlinedEvent{WithGuid: es.WithGuid{guid}},
	})
	assert.Equal(t, 3, goods.Version(), "version error")
	assert.Equal(t, "mmm", goods.name, "name error")
	assert.Equal(t, Money(50), goods.price, "price error")
	assert.Equal(t, Quantity(100), goods.quantity, "quantity error")
	assert.Equal(t, SN("20160601333"), goods.sn, "sn error")
	assert.Equal(t, Onlined, goods.state, "state error")
}

func TestGoodsRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewGoods().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckGoodsApplyEvents(t *testing.T) {
	events := []es.Event{
		&GoodsPublishedEvent{},
		&GoodsAuditedPassEvent{},
		&GoodsAuditedNoPassEvent{},
		&GoodsOnlinedEvent{},
		&GoodsOfflinedEvent{},
		&GoodsPurchaseSuccessedEvent{},
		&GoodsPurchaseFailuredEvent{},
		&GoodsCommentSuccessedEvent{},
		&GoodsCommentFailuredEvent{},
		&PaymetGoodsCompletedBecauseOfOrderEvent{},
		&PaymetGoodsFailedBecauseOfOrderEvent{},
	}
	assert.NotPanics(t, func() { NewGoods().ApplyEvents(events) }, "Check Process All Event")
}

func TestGoodsCommand(t *testing.T) {
	guid := es.NewGuid()

	tests := []struct {
		goods   *Goods
		command es.Command
		event   es.Event
	}{
		{
			&Goods{},
			&PublishGoodsCommand{WithGuid: es.WithGuid{Guid: guid}, Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"},
			&GoodsPublishedEvent{WithGuid: es.WithGuid{Guid: guid}, Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"},
		},
		{
			&Goods{state: Published},
			&AuditGoodsCommand{WithGuid: es.WithGuid{Guid: guid}, IsPass: true},
			&GoodsAuditedPassEvent{WithGuid: es.WithGuid{Guid: guid}},
		},
		{
			&Goods{state: Published},
			&AuditGoodsCommand{WithGuid: es.WithGuid{Guid: guid}, IsPass: false},
			&GoodsAuditedNoPassEvent{WithGuid: es.WithGuid{Guid: guid}},
		},
		{
			&Goods{state: AuditedPass},
			&OnlineGoodsCommand{WithGuid: es.WithGuid{Guid: guid}},
			&GoodsOnlinedEvent{WithGuid: es.WithGuid{Guid: guid}},
		},
		{
			&Goods{state: Onlined},
			&OfflineGoodsCommand{WithGuid: es.WithGuid{Guid: guid}},
			&GoodsOfflinedEvent{WithGuid: es.WithGuid{Guid: guid}},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.goods.ProcessCommand(v.command))
	}
}

func TestGoodsCommand_Panic(t *testing.T) {
	tests := []struct {
		goods   *Goods
		command es.Command
	}{
		{
			&Goods{},
			&struct{ es.WithGuid }{},
		},
		{
			&Goods{},
			&AuditGoodsCommand{},
		},
		{
			&Goods{state: AuditedPass},
			&AuditGoodsCommand{},
		},
		{
			&Goods{state: Onlined},
			&AuditGoodsCommand{},
		},
		{
			&Goods{state: Offlined},
			&AuditGoodsCommand{},
		},
		{
			&Goods{},
			&OnlineGoodsCommand{},
		},
		{
			&Goods{state: Published},
			&OnlineGoodsCommand{},
		},
		{
			&Goods{state: Onlined},
			&OnlineGoodsCommand{},
		},
		{
			&Goods{state: Offlined},
			&OnlineGoodsCommand{},
		},
		{
			&Goods{},
			&OfflineGoodsCommand{},
		},
		{
			&Goods{state: Published},
			&OfflineGoodsCommand{},
		},
		{
			&Goods{state: AuditedPass},
			&OfflineGoodsCommand{},
		},
		{
			&Goods{state: AuditedNoPass},
			&OfflineGoodsCommand{},
		},
		{
			&Goods{state: Offlined},
			&OfflineGoodsCommand{},
		},
		{
			&Goods{state: Offlined},
			&PurchaseGoodsBecauseOfPurchaseCommand{},
		},
		{
			&Goods{state: Published},
			&PurchaseGoodsBecauseOfPurchaseCommand{},
		},
		{
			&Goods{state: AuditedPass},
			&PurchaseGoodsBecauseOfPurchaseCommand{},
		},
		{
			&Goods{state: AuditedNoPass},
			&PurchaseGoodsBecauseOfPurchaseCommand{},
		},
		{
			&Goods{},
			&CompletePaymetGoodsBecauseOfOrderCommand{},
		},
		{
			&Goods{},
			&FailPaymetGoodsBecauseOfOrderCommand{},
		},
		{
			&Goods{},
			&CommentGoodsBecauseOfCommentCommand{},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.goods.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}

func TestPurchaseGoodsCommand(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 3,
	}

	outRangeDetails := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 4,
	}

	tests := []struct {
		goods   *Goods
		command es.Command
		event   es.Event
	}{
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 23, sn: "sn1234"},
			&PurchaseGoodsBecauseOfPurchaseCommand{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
			&GoodsPurchaseSuccessedEvent{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
		},
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 2, sn: "sn1234"},
			&PurchaseGoodsBecauseOfPurchaseCommand{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
			&GoodsPurchaseFailuredEvent{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
		},
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 20, sn: "sn1234"},
			&PurchaseGoodsBecauseOfPurchaseCommand{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: outRangeDetails},
			&GoodsPurchaseFailuredEvent{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: outRangeDetails},
		},
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 20, sn: "sn1234", purchaseLimit: map[es.Guid]Quantity{details.User: 1}},
			&PurchaseGoodsBecauseOfPurchaseCommand{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
			&GoodsPurchaseFailuredEvent{WithGuid: es.WithGuid{Guid: details.Goods}, PurchaseDetails: details},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.goods.ProcessCommand(v.command))
	}
}

func TestPaymetGoodsCommand(t *testing.T) {
	guid, userId, orderId := es.NewGuid(), es.NewGuid(), es.NewGuid()

	tests := []struct {
		goods   *Goods
		command es.Command
		event   es.Event
	}{
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 23, sn: "sn1234"},
			&CompletePaymetGoodsBecauseOfOrderCommand{WithGuid: es.WithGuid{Guid: guid}, User: userId, Order: orderId, Quantity: 2},
			&PaymetGoodsCompletedBecauseOfOrderEvent{WithGuid: es.WithGuid{Guid: guid}, User: userId, Order: orderId, Quantity: 2},
		},
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 23, sn: "sn1234"},
			&FailPaymetGoodsBecauseOfOrderCommand{WithGuid: es.WithGuid{Guid: guid}, User: userId, Order: orderId, Quantity: 2},
			&PaymetGoodsFailedBecauseOfOrderEvent{WithGuid: es.WithGuid{Guid: guid}, User: userId, Order: orderId, Quantity: 2},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.goods.ProcessCommand(v.command))
	}
}

func TestCommentGoodsCommand(t *testing.T) {
	details := CommentDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Comment:  es.NewGuid(),
		Purchase: es.NewGuid(),
		Content:  "dadfds",
		Time:     time.Now(),
	}

	tests := []struct {
		goods   *Goods
		command es.Command
		event   es.Event
	}{
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 23, sn: "sn1234", comments: map[es.Guid]es.Guid{details.Purchase: details.User}},
			&CommentGoodsBecauseOfCommentCommand{WithGuid: es.WithGuid{Guid: details.Goods}, CommentDetails: details},
			&GoodsCommentSuccessedEvent{WithGuid: es.WithGuid{Guid: details.Goods}, CommentDetails: details},
		},
		{
			&Goods{state: Onlined, name: "goods1", price: 30, quantity: 23, sn: "sn1234", comments: map[es.Guid]es.Guid{}},
			&CommentGoodsBecauseOfCommentCommand{WithGuid: es.WithGuid{Guid: details.Goods}, CommentDetails: details},
			&GoodsCommentFailuredEvent{WithGuid: es.WithGuid{Guid: details.Goods}, CommentDetails: details},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.goods.ProcessCommand(v.command))
	}
}
