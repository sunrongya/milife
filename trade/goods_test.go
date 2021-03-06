package trade

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestGoodsRestore(t *testing.T) {
	goods := &Goods{}
	goods.HandleGoodsPublishedEvent(&GoodsPublishedEvent{Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"})
	goods.HandleGoodsAuditedPassEvent(&GoodsAuditedPassEvent{})
	goods.HandleGoodsOnlinedEvent(&GoodsOnlinedEvent{})

	assert.Equal(t, "mmm", goods._name, "name error")
	assert.Equal(t, Money(50), goods._price, "price error")
	assert.Equal(t, Quantity(100), goods._quantity, "quantity error")
	assert.Equal(t, SN("20160601333"), goods._sn, "sn error")
	assert.Equal(t, Onlined, goods._state, "state error")
}

func TestPublishGoodsCommand(t *testing.T) {
	command := &PublishGoodsCommand{Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"}
	events := []es.Event{&GoodsPublishedEvent{Name: "mmm", Price: 50, Quantity: 100, SN: "20160601333"}}
	goods := &Goods{}

	assert.Equal(t, events, goods.ProcessPublishGoodsCommand(command))
}

func TestAuditGoodsCommandOfPass(t *testing.T) {
	command := &AuditGoodsCommand{IsPass: true}
	events := []es.Event{&GoodsAuditedPassEvent{}}
	goods := &Goods{_state: Published}

	assert.Equal(t, events, goods.ProcessAuditGoodsCommand(command))
}

func TestAuditGoodsCommandOfNoPass(t *testing.T) {
	command := &AuditGoodsCommand{IsPass: false}
	events := []es.Event{&GoodsAuditedNoPassEvent{}}
	goods := &Goods{_state: Published}

	assert.Equal(t, events, goods.ProcessAuditGoodsCommand(command))
}

func TestOnlineGoodsCommand(t *testing.T) {
	command := &OnlineGoodsCommand{}
	events := []es.Event{&GoodsOnlinedEvent{}}
	goods := &Goods{_state: AuditedPass}

	assert.Equal(t, events, goods.ProcessOnlineGoodsCommand(command))
}

func TestOfflineGoodsCommand(t *testing.T) {
	command := &OfflineGoodsCommand{}
	events := []es.Event{&GoodsOfflinedEvent{}}
	goods := &Goods{_state: Onlined}

	assert.Equal(t, events, goods.ProcessOfflineGoodsCommand(command))
}

func TestPurchaseGoodsBecauseOfPurchaseCommand2Successed(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 3,
	}
	command := &PurchaseGoodsBecauseOfPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseSuccessedEvent{PurchaseDetails: details}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 23, _sn: "sn1234"}

	assert.Equal(t, events, goods.ProcessPurchaseGoodsBecauseOfPurchaseCommand(command))
}

func TestPurchaseGoodsBecauseOfPurchaseCommand2Failured(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 3,
	}
	command := &PurchaseGoodsBecauseOfPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseFailuredEvent{PurchaseDetails: details}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 2, _sn: "sn1234"}

	assert.Equal(t, events, goods.ProcessPurchaseGoodsBecauseOfPurchaseCommand(command))
}

func TestPurchaseGoodsBecauseOfPurchaseCommand2FailuredOfLargeQuantity(t *testing.T) {
	outRangeDetails := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 4,
	}
	command := &PurchaseGoodsBecauseOfPurchaseCommand{PurchaseDetails: outRangeDetails}
	events := []es.Event{&GoodsPurchaseFailuredEvent{PurchaseDetails: outRangeDetails}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 20, _sn: "sn1234"}

	assert.Equal(t, events, goods.ProcessPurchaseGoodsBecauseOfPurchaseCommand(command))
}

func TestPurchaseGoodsBecauseOfPurchaseCommand2FailuredOfLargeQuantity2(t *testing.T) {
	details := PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Purchase: es.NewGuid(),
		Quantity: 3,
	}
	command := &PurchaseGoodsBecauseOfPurchaseCommand{PurchaseDetails: details}
	events := []es.Event{&GoodsPurchaseFailuredEvent{PurchaseDetails: details}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 20, _sn: "sn1234", _purchaseLimit: map[es.Guid]Quantity{details.User: 1}}

	assert.Equal(t, events, goods.ProcessPurchaseGoodsBecauseOfPurchaseCommand(command))
}

func TestCompletePaymetGoodsBecauseOfOrderCommand(t *testing.T) {
	userId, orderId := es.NewGuid(), es.NewGuid()
	command := &CompletePaymetGoodsBecauseOfOrderCommand{User: userId, Order: orderId, Quantity: 2}
	events := []es.Event{&PaymetGoodsCompletedBecauseOfOrderEvent{User: userId, Order: orderId, Quantity: 2}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 23, _sn: "sn1234"}

	assert.Equal(t, events, goods.ProcessCompletePaymetGoodsBecauseOfOrderCommand(command))
}

func TestFailPaymetGoodsBecauseOfOrderCommand(t *testing.T) {
	userId, orderId := es.NewGuid(), es.NewGuid()
	command := &FailPaymetGoodsBecauseOfOrderCommand{User: userId, Order: orderId, Quantity: 2}
	events := []es.Event{&PaymetGoodsFailedBecauseOfOrderEvent{User: userId, Order: orderId, Quantity: 2}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 23, _sn: "sn1234"}

	assert.Equal(t, events, goods.ProcessFailPaymetGoodsBecauseOfOrderCommand(command))
}

func TestCommentGoodsBecauseOfCommentCommand2Successed(t *testing.T) {
	details := CommentDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Comment:  es.NewGuid(),
		Purchase: es.NewGuid(),
		Content:  "dadfds",
		Time:     time.Now(),
	}

	command := &CommentGoodsBecauseOfCommentCommand{CommentDetails: details}
	events := []es.Event{&GoodsCommentSuccessedEvent{CommentDetails: details}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 23, _sn: "sn1234", _comments: map[es.Guid]es.Guid{details.Purchase: details.User}}

	assert.Equal(t, events, goods.ProcessCommentGoodsBecauseOfCommentCommand(command))
}

func TestCommentGoodsBecauseOfCommentCommand2Failured(t *testing.T) {
	details := CommentDetails{
		User:     es.NewGuid(),
		Goods:    es.NewGuid(),
		Comment:  es.NewGuid(),
		Purchase: es.NewGuid(),
		Content:  "dadfds",
		Time:     time.Now(),
	}

	command := &CommentGoodsBecauseOfCommentCommand{CommentDetails: details}
	events := []es.Event{&GoodsCommentFailuredEvent{CommentDetails: details}}
	goods := &Goods{_state: Onlined, _name: "goods1", _price: 30, _quantity: 23, _sn: "sn1234", _comments: map[es.Guid]es.Guid{}}

	assert.Equal(t, events, goods.ProcessCommentGoodsBecauseOfCommentCommand(command))
}

func TestAuditGoodsCommand_Panic(t *testing.T) {
	goodses := []*Goods{
		&Goods{},
		&Goods{_state: AuditedPass},
		&Goods{_state: Onlined},
		&Goods{_state: Offlined},
	}

	for _, goods := range goodses {
		assert.Panics(t, func() {
			goods.ProcessAuditGoodsCommand(&AuditGoodsCommand{})
		})
	}
}

func TestOnlineGoodsCommand_Panic(t *testing.T) {
	goodses := []*Goods{
		&Goods{},
		&Goods{_state: Published},
		&Goods{_state: Onlined},
		&Goods{_state: Offlined},
	}

	for _, goods := range goodses {
		assert.Panics(t, func() {
			goods.ProcessOnlineGoodsCommand(&OnlineGoodsCommand{})
		})
	}
}

func TestOfflineGoodsCommand_Panic(t *testing.T) {
	goodses := []*Goods{
		&Goods{},
		&Goods{_state: Published},
		&Goods{_state: AuditedPass},
		&Goods{_state: AuditedNoPass},
		&Goods{_state: Offlined},
	}

	for _, goods := range goodses {
		assert.Panics(t, func() {
			goods.ProcessOfflineGoodsCommand(&OfflineGoodsCommand{})
		})
	}
}

func TestPurchaseGoodsBecauseOfPurchaseCommand_Panic(t *testing.T) {
	goodses := []*Goods{
		&Goods{_state: Offlined},
		&Goods{_state: Published},
		&Goods{_state: AuditedPass},
		&Goods{_state: AuditedNoPass},
	}

	for _, goods := range goodses {
		assert.Panics(t, func() {
			goods.ProcessPurchaseGoodsBecauseOfPurchaseCommand(&PurchaseGoodsBecauseOfPurchaseCommand{})
		})
	}
}

func TestCompletePaymetGoodsBecauseOfOrderCommand_Panic(t *testing.T) {
	assert.Panics(t, func() {
		new(Goods).ProcessCompletePaymetGoodsBecauseOfOrderCommand(&CompletePaymetGoodsBecauseOfOrderCommand{})
	})
}

func TestFailPaymetGoodsBecauseOfOrderCommand_Panic(t *testing.T) {
	assert.Panics(t, func() {
		new(Goods).ProcessFailPaymetGoodsBecauseOfOrderCommand(&FailPaymetGoodsBecauseOfOrderCommand{})
	})
}

func TestCommentGoodsBecauseOfCommentCommand_Panic(t *testing.T) {
	assert.Panics(t, func() {
		new(Goods).ProcessCommentGoodsBecauseOfCommentCommand(&CommentGoodsBecauseOfCommentCommand{})
	})
}
