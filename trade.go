package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	ES "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/estore"
	"github.com/sunrongya/milife/trade"
	"github.com/xyproto/simplebolt"
)

func Trade() {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	creator := simplebolt.NewCreator(db)
	eventFactory := ES.NewEventFactory()
	eventFactory.RegisterAggregate(trade.NewGoods(), trade.NewGoodsComment(), trade.NewGoodsPurchase(), trade.NewOrder())
	store := estore.NewXyprotoEStore(creator, estore.NewEncoder(eventFactory), estore.NewDecoder(eventFactory))

	//var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 初始化
	goodsService := trade.NewGoodsService(store)
	purchaseService := trade.NewGoodsPurchaseService(store)
	orderService := trade.NewOrderService(store)
	commentService := trade.NewGoodsCommentService(store)
	eventbus := ES.NewInternalEventBus(store)

	purchaseHandler := trade.NewPurchaseEventHandler(goodsService.CommandChannel(), purchaseService.CommandChannel())
	commentHandler := trade.NewCommentEventHandler(goodsService.CommandChannel(), commentService.CommandChannel())
	paymetHandler := trade.NewPaymetEventHandler(new(NullPaymetService), goodsService.CommandChannel(), orderService.CommandChannel())
	readRepository := ES.NewMemoryReadRepository()
	goodsProjector := trade.NewGoodsProjector(readRepository)

	eventbus.RegisterHandlers(purchaseHandler)
	eventbus.RegisterHandlers(commentHandler)
	eventbus.RegisterHandlers(paymetHandler)
	eventbus.RegisterHandlers(goodsProjector)

	go eventbus.HandleEvents()
	go goodsService.HandleCommands()
	go purchaseService.HandleCommands()
	go orderService.HandleCommands()
	go commentService.HandleCommands()

	// 执行命令
	// 商品发布-审核-上线
	fmt.Printf("- 发布商品[name:mangos, price:25, quantity:50, sn:SN12345]\tOK\n")
	mangos := goodsService.PublishGoods("mangos", 25, 50, "SN12345")
	fmt.Printf("- 发布商品[name:apple, price:10, quantity:100, sn:SN12345]\tOK\n")
	apple := goodsService.PublishGoods("apple", 10, 100, "SN12346")
	fmt.Printf("- 商品通过验证[name:mangos] \tOK\n")
	goodsService.AuditGoods(mangos, true)
	fmt.Printf("- 商品通过验证[name:apple] \tOK\n")
	goodsService.AuditGoods(apple, true)
	fmt.Printf("- 商品上线[name:mangos] \tOK\n")
	goodsService.OnlineGoods(mangos)
	fmt.Printf("- 商品上线[name:apple] \tOK\n")
	goodsService.OnlineGoods(apple)
	fmt.Printf("- 商品下线[name:apple] \tOK\n")
	goodsService.OfflineGoods(apple)
	// 抢购商品
	user1, user2 := ES.NewGuid(), ES.NewGuid()
	fmt.Printf("- user1 抢购商品[name:mangos] 2个\tOK\n")
	purchase1 := purchaseService.CreateGoodsPurchase(mangos, user1, 2)
	fmt.Printf("- user2 抢购商品[name:mangos] 3个\tOK\n")
	purchase2 := purchaseService.CreateGoodsPurchase(mangos, user2, 3)
	// 订单生成及支付
	orderItems1 := []trade.OrderItem{
		{
			Goods:    mangos,
			Purchase: purchase1,
			Name:     "mangos",
			Price:    25,
			Quantity: 2,
		},
	}
	orderItems2 := []trade.OrderItem{
		{
			Goods:    mangos,
			Purchase: purchase2,
			Name:     "mangos",
			Price:    25,
			Quantity: 3,
		},
	}
	fmt.Printf("- 生成订单1 \tOK\n")
	order1 := orderService.CreateOrder(orderItems1)
	fmt.Printf("- 生成订单2 \tOK\n")
	order2 := orderService.CreateOrder(orderItems2)
	fmt.Printf("- 支付订单1 \tOK\n")
	orderService.PaymetOrder(order1, user1, "user1Account", "managedAccount")
	fmt.Printf("- 支付订单2 \tOK\n")
	orderService.PaymetOrder(order2, user2, "user2Account", "managedAccount")

	// 查看结果 wait and print
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(mangos))
		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(apple))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase1))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase2))
		fmt.Printf("%v\n------------------\n", orderService.RestoreAggregate(order1))
		fmt.Printf("%v\n------------------\n", orderService.RestoreAggregate(order2))

		fmt.Printf("-----------------\nRead Model:\n\n")
		if goodsMangos, err := readRepository.Find(mangos); err == nil {
			goods := goodsMangos.(*trade.RGoods)
			fmt.Printf("商品 %v\n------------------\n", goods.Name)
			fmt.Printf("SN:%v, Name:%v, Price:%v, Quantity:%v, Purchases:%v, State:%v \n",
				goods.SN, goods.Name, goods.Price, goods.Quantity, goods.Purchases, goods.State)
			fmt.Println("######## SuccessedPurchaseRecords #######")
			for i, v := range goods.SuccessedPurchaseRecords {
				fmt.Printf("success %d: %v, %v, %v, %v\n", i, v.User, v.Purchase, v.Quantity, v.OrderState)
			}
			fmt.Println("######## FailuredPurchaseRecords #######")
			for i, v := range goods.FailuredPurchaseRecords {
				fmt.Printf("success %d: %v, %v, %v, %v\n", i, v.User, v.Purchase, v.Quantity, v.OrderState)
			}
		}

		if goodsApple, err := readRepository.Find(apple); err == nil {
			fmt.Printf("商品2: %v\n------------------\n", goodsApple)
		}

		wg.Done()
	}()

	wg.Wait()
}
