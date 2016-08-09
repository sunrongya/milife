package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	ES "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/estore"
	"github.com/sunrongya/milife/paymet"
	"github.com/sunrongya/milife/trade"
	"github.com/xyproto/simplebolt"
)

func Trade2Paymet() {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	creator := simplebolt.NewCreator(db)
	eventFactory := ES.NewEventFactory()
	eventFactory.RegisterAggregate(
		paymet.NewAccount(),
		paymet.NewMoneyTransfer(),
		trade.NewGoods(),
		trade.NewGoodsComment(),
		trade.NewGoodsPurchase(),
		trade.NewOrder(),
	)
	store := estore.NewXyprotoEStore(creator, estore.NewEncoder(eventFactory), estore.NewDecoder(eventFactory))

	//var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 初始化
	accountService := paymet.NewAccountService(store)
	transferService := paymet.NewTransferService(store)

	goodsService := trade.NewGoodsService(store)
	purchaseService := trade.NewGoodsPurchaseService(store)
	orderService := trade.NewOrderService(store)
	commentService := trade.NewGoodsCommentService(store)
	eventbus := ES.NewInternalEventBus(store)

	paymentTransferService := NewPaymentTransferService(transferService)

	// 注册EventHandler/读模型Handler
	purchaseHandler := trade.NewPurchaseEventHandler(goodsService.CommandChannel(), purchaseService.CommandChannel())
	commentHandler := trade.NewCommentEventHandler(goodsService.CommandChannel(), commentService.CommandChannel())
	paymetHandler := trade.NewPaymetEventHandler(
		paymentTransferService,
		goodsService.CommandChannel(),
		orderService.CommandChannel(),
	)
	readRepository := ES.NewMemoryReadRepository()
	goodsProjector := trade.NewGoodsProjector(readRepository)

	paymetEventHandler := paymet.NewEventHandler(accountService.CommandChannel(), transferService.CommandChannel())
	accProjector := paymet.NewRAccountProjector(readRepository)
	rateProjector := paymet.NewMoneyFlowRateProjector(readRepository, ES.NewGuid())

	eventbus.RegisterHandlers(purchaseHandler)
	eventbus.RegisterHandlers(commentHandler)
	eventbus.RegisterHandlers(paymetHandler)
	eventbus.RegisterHandlers(goodsProjector)
	eventbus.RegisterHandlers(paymetEventHandler)
	eventbus.RegisterHandlers(accProjector)
	eventbus.RegisterHandlers(rateProjector)

	go eventbus.HandleEvents()
	go goodsService.HandleCommands()
	go purchaseService.HandleCommands()
	go orderService.HandleCommands()
	go commentService.HandleCommands()
	go accountService.HandleCommands()
	go transferService.HandleCommands()

	// 执行命令
	// 创建账户
	fmt.Printf("- Open managedAccount with balance 0:\tOK\n")
	managedAccount := accountService.OpenAccount("sry", "9558866", 0) // managedAccount: balance=0
	fmt.Printf("- Open account 1 with balance 1000:\tOK\n")
	acc1 := accountService.OpenAccount("sry", "95588", 100) // acc1: balance=100
	fmt.Printf("- Open account 2 with balance 10:\tOK\n")
	acc2 := accountService.OpenAccount("jjj", "95533", 10) // acc2: balance=10
	fmt.Printf("- Credit account 1 with amount 1900:\tOK\n")
	accountService.CreditAccount(acc1, 1900) // acc1: balance=2000

	paymentTransferService.Register("managedAccount", managedAccount)
	paymentTransferService.Register("user1Account", acc1)
	paymentTransferService.Register("user2Account", acc2)

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

	// 抢购商品
	user1, user2 := ES.NewGuid(), ES.NewGuid()
	fmt.Printf("- user1 抢购商品[name:mangos] 2个\tOK\n")
	purchase1 := purchaseService.CreateGoodsPurchase(mangos, user1, 2)
	fmt.Printf("- user2 抢购商品[name:mangos] 3个\tOK\n")
	purchase2 := purchaseService.CreateGoodsPurchase(mangos, user2, 3)
	fmt.Printf("- user2 抢购商品[name:apple] 2个\tOK\n")
	purchase3 := purchaseService.CreateGoodsPurchase(apple, user2, 2)

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
		{
			Goods:    apple,
			Purchase: purchase3,
			Name:     "apple",
			Price:    10,
			Quantity: 2,
		},
	}
	fmt.Printf("- 生成订单1 \tOK\n")
	order1 := orderService.CreateOrder(orderItems1)
	fmt.Printf("- 生成订单2 \tOK\n")
	order2 := orderService.CreateOrder(orderItems2)
	fmt.Printf("- 支付订单1 \tOK\n")
	orderService.PaymetOrder(order1, user1, "user1Account", "managedAccount")
	fmt.Printf("- 支付订单2 [支付失败,账户余额不够] \tOK\n")
	orderService.PaymetOrder(order2, user2, "user2Account", "managedAccount")

	// 验证
	//wait and print
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", accountService.RestoreAggregate(managedAccount))
		fmt.Printf("%v\n------------------\n", accountService.RestoreAggregate(acc1))
		fmt.Printf("%v\n------------------\n", accountService.RestoreAggregate(acc2))

		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(mangos))
		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(apple))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase1))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase2))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase3))
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
