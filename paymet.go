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
	"github.com/xyproto/simplebolt"
)

func Paymet() {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	creator := simplebolt.NewCreator(db)
	eventFactory := ES.NewEventFactory()
	eventFactory.RegisterAggregate(paymet.NewAccount(), paymet.NewMoneyTransfer())
	store := estore.NewXyprotoEStore(creator, estore.NewEncoder(eventFactory), estore.NewDecoder(eventFactory))

	//var store = ES.NewInMemStore()
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
