package main

import (
	"fmt"
	ES "github.com/sunrongya/eventsourcing"
	"sync"
	"milife/paymet"
	"milife/trade"
	"time"
)

func testTrade() {
	var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 初始化
	goodsService := trade.NewGoodsService(store)
	purchaseService := trade.NewGoodsPurchaseService(store)
	orderService := trade.NewOrderService(store)
	commentService := trade.NewGoodsCommentService(store)
	eventbus := ES.NewInternalEventBus(store)

	eh := trade.NewEventHandler(goodsService.CommandChannel(), purchaseService.CommandChannel(), commentService.CommandChannel())
	paymetHandler := trade.NewPaymetEventHandler(new(trade.NullPaymetService), goodsService.CommandChannel(), orderService.CommandChannel())
	readRepository := ES.NewMemoryReadRepository()
	goodsProjector := trade.NewGoodsProjector(readRepository)

	eventbus.RegisterHandlers(eh)
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
		printEvents(store.GetEvents(0, 100))
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(mangos))
		fmt.Printf("%v\n------------------\n", goodsService.RestoreAggregate(apple))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase1))
		fmt.Printf("%v\n------------------\n", purchaseService.RestoreAggregate(purchase2))

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

func testPaymet() {
	var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	as := paymet.NewAccountService(store)
	mt := paymet.NewTransferService(store)
	eventbus := ES.NewInternalEventBus(store)

	// 注册EventHandler/读模型Handler
	eh := paymet.NewEventHandler(as.CommandChannel(), mt.CommandChannel())
	readRepository := ES.NewMemoryReadRepository()
	accProjector := paymet.NewRAccountProjector(readRepository)
	rateProjector := paymet.NewMoneyFlowRateProjector(readRepository, ES.NewGuid())
	eventbus.RegisterHandlers(eh)
	eventbus.RegisterHandlers(accProjector)
	eventbus.RegisterHandlers(rateProjector)

	go eventbus.HandleEvents()
	go as.HandleCommands()
	go mt.HandleCommands()

	// 执行命令
	fmt.Printf("- Open account 1 with balance 10:\tOK\n")
	acc1 := as.OpenAccount("sry", "95588", 10) // acc1: balance=10
	fmt.Printf("- Open account 2 with balance 10:\tOK\n")
	acc2 := as.OpenAccount("jjj", "95533", 10) // acc2: balance=10
	fmt.Printf("- Credit account 1 with amount 190:\tOK\n")
	as.CreditAccount(acc1, 190) // acc1: balance=200
	fmt.Printf("- Debit account 1 with amount 100:\tOK\n")
	as.DebitAccount(acc1, 100) // acc1: balance=100
	fmt.Printf("- Debit account 2 with amount 100:\tFAIL\n")
	as.DebitAccount(acc2, 100) // Will fail -> no change
	fmt.Printf("- Transfer 10 from account 1 to account 2:\tOK\n")
	trans1 := mt.Transfer(10, acc1, acc2)
	fmt.Printf("- Transfer 100 from account 2 to account 1:\tFAIL\n")
	trans2 := mt.Transfer(100, acc2, acc1)

	// 验证
	//wait and print
	go func() {
		time.Sleep(200 * time.Millisecond)
		printEvents(store.GetEvents(0, 100))
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", as.RestoreAggregate(acc1))
		fmt.Printf("%v\n------------------\n", as.RestoreAggregate(acc2))
		fmt.Printf("%v\n------------------\n", mt.RestoreAggregate(trans1))
		fmt.Printf("%v\n------------------\n", mt.RestoreAggregate(trans2))

		fmt.Printf("-----------------\nRead Model:\n\n")
		if account1, err := readRepository.Find(acc1); err == nil {
			fmt.Printf("账户1: %v\n------------------\n", account1)
		}
		if account2, err := readRepository.Find(acc2); err == nil {
			fmt.Printf("账户2: %v\n------------------\n", account2)
		}
		if flowRate, err := readRepository.Find(rateProjector.Id); err == nil {
			fmt.Printf("交易统计: %v\n------------------\n", flowRate)
		}

		wg.Done()
	}()

	wg.Wait()
}

func testTrade2Paymet() {
}

func printEvents(events []ES.Event) {
	fmt.Printf("-----------------\nEvents after all operations:\n\n")
	for i, e := range events {
		fmt.Printf("%v: %T\n", i, e)
	}
}

func main() {
	testTrade()
	//testPaymet()
}
