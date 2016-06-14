package trade

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestGoodsReadModel(t *testing.T) {
	readRepository := es.NewMemoryReadRepository()
	goodsProjector := NewGoodsProjector(readRepository)

	// 商品发布
	publishedEvents := []*GoodsPublishedEvent{
		&GoodsPublishedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Name: "mangos", Price: 50, Quantity: 100, SN: "SN1234"},
		&GoodsPublishedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Name: "apple", Price: 20, Quantity: 60, SN: "SN1235"},
		&GoodsPublishedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Name: "banana", Price: 15, Quantity: 80, SN: "SN1236"},
	}

	for _, event := range publishedEvents {
		goodsProjector.HandleGoodsPublishedEvent(event)
	}

	// 商品发布验证
	for _, event := range publishedEvents {
		i, err := readRepository.Find(event.GetGuid())
		assert.NoError(t, err, fmt.Sprintf("读取已发布商品[%s]信息错误", event.Name))
		goods := i.(*RGoods)

		assert.Equal(t, event.GetGuid(), goods.Id, "ID 不相等")
		assert.Equal(t, event.Name, goods.Name, "Name 不相等")
		assert.Equal(t, event.Price, goods.Price, "Price 不相等")
		assert.Equal(t, event.Quantity, goods.Quantity, "Quantity 不相等")
		assert.Equal(t, event.SN, goods.SN, "SN 不相等")
		assert.Equal(t, Published, goods.State, "应该是已发布状态")
	}

	// 商品审核
	goodsProjector.HandleGoodsAuditedPassEvent(&GoodsAuditedPassEvent{WithGuid: es.WithGuid{publishedEvents[0].GetGuid()}})
	goodsProjector.HandleGoodsAuditedPassEvent(&GoodsAuditedPassEvent{WithGuid: es.WithGuid{publishedEvents[1].GetGuid()}})
	goodsProjector.HandleGoodsAuditedNoPassEvent(&GoodsAuditedNoPassEvent{WithGuid: es.WithGuid{publishedEvents[2].GetGuid()}})

	// 商品审核验证
	auditedStates := []State{AuditedPass, AuditedPass, AuditedNoPass}
	for i, event := range publishedEvents {
		model, _ := readRepository.Find(event.GetGuid())
		goods := model.(*RGoods)
		assert.Equal(t, auditedStates[i], goods.State, "状态错误")
	}

	// 商品上线
	goodsProjector.HandleGoodsOnlinedEvent(&GoodsOnlinedEvent{WithGuid: es.WithGuid{publishedEvents[0].GetGuid()}})
	goodsProjector.HandleGoodsOnlinedEvent(&GoodsOnlinedEvent{WithGuid: es.WithGuid{publishedEvents[1].GetGuid()}})

	// 商品上线验证
	for i := 0; i < 2; i++ {
		model, _ := readRepository.Find(publishedEvents[i].GetGuid())
		goods := model.(*RGoods)
		assert.Equal(t, Onlined, goods.State, "应该是已上线状态")
	}

	// 商品下线
	goodsProjector.HandleGoodsOfflinedEvent(&GoodsOfflinedEvent{WithGuid: es.WithGuid{publishedEvents[1].GetGuid()}})

	// 商品下线验证
	model, _ := readRepository.Find(publishedEvents[1].GetGuid())
	goods := model.(*RGoods)
	assert.Equal(t, Offlined, goods.State, "应该是已下线状态")
}

func TestPurchaseOfGoodsReadModel(t *testing.T) {
	readRepository := es.NewMemoryReadRepository()
	goodsProjector := NewGoodsProjector(readRepository)

	publishedEvent := &GoodsPublishedEvent{
		WithGuid: es.WithGuid{es.NewGuid()},
		Name:     "mangos",
		Price:    50,
		Quantity: 100,
		SN:       "SN1234",
	}

	goodsProjector.HandleGoodsPublishedEvent(publishedEvent)
	goodsProjector.HandleGoodsAuditedPassEvent(&GoodsAuditedPassEvent{WithGuid: es.WithGuid{publishedEvent.GetGuid()}})
	goodsProjector.HandleGoodsOnlinedEvent(&GoodsOnlinedEvent{WithGuid: es.WithGuid{publishedEvent.GetGuid()}})

	// 抢购成功
	purchaseSuccessedEvents := []*GoodsPurchaseSuccessedEvent{
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 1),
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 2),
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 3),
	}
	for _, event := range purchaseSuccessedEvents {
		goodsProjector.HandleGoodsPurchaseSuccessedEvent(event)
	}

	// 抢购成功验证
	model, _ := readRepository.Find(publishedEvent.GetGuid())
	goods := model.(*RGoods)

	var successRecords []*PurchaseRecord
	for i, event := range purchaseSuccessedEvents {
		successRecords = append(successRecords, &PurchaseRecord{
			Quantity: Quantity(i + 1),
			User:     event.User,
			Purchase: event.Purchase,
		})
	}
	assert.Equal(t, Quantity(6), goods.Purchases, "成功抢购数错误")
	assert.Equal(t, successRecords, goods.SuccessedPurchaseRecords, "验证成功抢购数据失败")

	// 抢购失败
	purchaseFailuredEvents := []*GoodsPurchaseFailuredEvent{
		getPurchaseFailuredEvent(publishedEvent.GetGuid(), 2),
		getPurchaseFailuredEvent(publishedEvent.GetGuid(), 3),
		getPurchaseFailuredEvent(publishedEvent.GetGuid(), 4),
		getPurchaseFailuredEvent(publishedEvent.GetGuid(), 5),
	}
	for _, event := range purchaseFailuredEvents {
		goodsProjector.HandleGoodsPurchaseFailuredEvent(event)
	}

	// 抢购失败验证
	model, _ = readRepository.Find(publishedEvent.GetGuid())
	goods = model.(*RGoods)

	assert.Equal(t, Quantity(6), goods.Purchases, "成功抢购数错误")
	assert.Equal(t, successRecords, goods.SuccessedPurchaseRecords, "验证成功抢购记录失败")

	var failureRecords []*PurchaseRecord
	for i, event := range purchaseFailuredEvents {
		failureRecords = append(failureRecords, &PurchaseRecord{
			Quantity: Quantity(i + 2),
			User:     event.User,
			Purchase: event.Purchase,
		})
	}
	assert.Equal(t, failureRecords, goods.FailuredPurchaseRecords, "验证失败抢购记录失败")
}

func getPurchaseSuccessedEvent(guid es.Guid, quantity Quantity) *GoodsPurchaseSuccessedEvent {
	return &GoodsPurchaseSuccessedEvent{WithGuid: es.WithGuid{guid}, PurchaseDetails: PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    guid,
		Purchase: es.NewGuid(),
		Quantity: quantity,
	}}
}

func getPurchaseFailuredEvent(guid es.Guid, quantity Quantity) *GoodsPurchaseFailuredEvent {
	return &GoodsPurchaseFailuredEvent{WithGuid: es.WithGuid{guid}, PurchaseDetails: PurchaseDetails{
		User:     es.NewGuid(),
		Goods:    guid,
		Purchase: es.NewGuid(),
		Quantity: quantity,
	}}
}

func TestPaymetOfGoodsReadModel(t *testing.T) {
	readRepository := es.NewMemoryReadRepository()
	goodsProjector := NewGoodsProjector(readRepository)

	publishedEvent := &GoodsPublishedEvent{
		WithGuid: es.WithGuid{es.NewGuid()},
		Name:     "mangos",
		Price:    50,
		Quantity: 100,
		SN:       "SN1234",
	}

	goodsProjector.HandleGoodsPublishedEvent(publishedEvent)
	goodsProjector.HandleGoodsAuditedPassEvent(&GoodsAuditedPassEvent{WithGuid: es.WithGuid{publishedEvent.GetGuid()}})
	goodsProjector.HandleGoodsOnlinedEvent(&GoodsOnlinedEvent{WithGuid: es.WithGuid{publishedEvent.GetGuid()}})

	// 抢购成功
	purchaseSuccessedEvents := []*GoodsPurchaseSuccessedEvent{
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 1),
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 2),
		getPurchaseSuccessedEvent(publishedEvent.GetGuid(), 3),
	}
	for _, event := range purchaseSuccessedEvents {
		goodsProjector.HandleGoodsPurchaseSuccessedEvent(event)
	}

	// 支付
	paymetSuccessedEvents := []*PaymetGoodsCompletedBecauseOfOrderEvent{
		&PaymetGoodsCompletedBecauseOfOrderEvent{
			WithGuid: es.WithGuid{publishedEvent.GetGuid()},
			User:     purchaseSuccessedEvents[0].User,
			Order:    es.NewGuid(),
			Purchase: purchaseSuccessedEvents[0].Purchase,
			Quantity: purchaseSuccessedEvents[0].Quantity,
		},
		&PaymetGoodsCompletedBecauseOfOrderEvent{
			WithGuid: es.WithGuid{publishedEvent.GetGuid()},
			User:     purchaseSuccessedEvents[1].User,
			Order:    es.NewGuid(),
			Purchase: purchaseSuccessedEvents[1].Purchase,
			Quantity: purchaseSuccessedEvents[1].Quantity,
		},
	}

	paymetFailedEvent := &PaymetGoodsFailedBecauseOfOrderEvent{
		WithGuid: es.WithGuid{publishedEvent.GetGuid()},
		User:     purchaseSuccessedEvents[2].User,
		Order:    es.NewGuid(),
		Purchase: purchaseSuccessedEvents[2].Purchase,
		Quantity: purchaseSuccessedEvents[2].Quantity,
	}

	for _, event := range paymetSuccessedEvents {
		goodsProjector.HandlePaymetGoodsCompletedBecauseOfOrderEvent(event)
	}
	goodsProjector.HandlePaymetGoodsFailedBecauseOfOrderEvent(paymetFailedEvent)

	// 支付验证
	model, _ := readRepository.Find(publishedEvent.GetGuid())
	goods := model.(*RGoods)

	assert.Equal(t, 3, len(goods.SuccessedPurchaseRecords), "成功抢购数错误")
	assert.Equal(t, 0, len(goods.FailuredPurchaseRecords), "失败抢购数错误")
	assert.Equal(t, Quantity(3), goods.Purchases, "已抢购量错误")
	assert.Equal(t, OrderPaymetCompleted, goods.SuccessedPurchaseRecords[0].OrderState, "支付状态错误")
	assert.Equal(t, OrderPaymetCompleted, goods.SuccessedPurchaseRecords[1].OrderState, "支付状态错误")
	assert.Equal(t, OrderPaymetFailed, goods.SuccessedPurchaseRecords[2].OrderState, "支付状态错误")
}

func TestCommentOfGoodsReadModel(t *testing.T) {
}
