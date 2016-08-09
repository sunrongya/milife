package main

import (
	"time"

	ES "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/milife/paymet"
	"github.com/sunrongya/milife/trade"
)

type NullPaymetService struct {
}

func (this *NullPaymetService) Transfer(amount trade.Money, user, managed trade.BankAccount, completeFn func(bool)) {
	completeFn(true)
}

type PaymentTransferService struct {
	_accounts        map[trade.BankAccount]ES.Guid
	_transferService *paymet.TransferService
}

func NewPaymentTransferService(transferService *paymet.TransferService) *PaymentTransferService {
	return &PaymentTransferService{
		_accounts:        make(map[trade.BankAccount]ES.Guid),
		_transferService: transferService,
	}
}

func (this *PaymentTransferService) Transfer(amount trade.Money, userAccount, managedAccount trade.BankAccount, completeFn func(bool)) {
	userGuid, ok := this._accounts[userAccount]
	if !ok {
		completeFn(false)
		return
	}
	managedGuid, ok := this._accounts[managedAccount]
	if !ok {
		completeFn(false)
		return
	}
	trans := this._transferService.Transfer(paymet.Money(amount), userGuid, managedGuid)
	// 这里的做法很简单：等待200毫秒后判断支付聚合根是否支付完成
	time.Sleep(200 * time.Millisecond)
	moneyTransfer := this._transferService.RestoreAggregate(trans).(*paymet.MoneyTransfer)
	if moneyTransfer.State() != paymet.Completed {
		completeFn(false)
		return
	}
	completeFn(true)
}

func (this *PaymentTransferService) Register(account string, accountId ES.Guid) {
	this._accounts[trade.BankAccount(account)] = accountId
}
