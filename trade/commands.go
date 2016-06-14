package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Goods Commands
// --------------

type PublishGoodsCommand struct {
	es.WithGuid
	Name     string
	Price    Money
	Quantity Quantity
	SN       SN
}

type AuditGoodsCommand struct {
	es.WithGuid
	IsPass bool
}

type OnlineGoodsCommand struct {
	es.WithGuid
}

type OfflineGoodsCommand struct {
	es.WithGuid
}

type PurchaseGoodsBecauseOfPurchaseCommand struct {
	es.WithGuid
	PurchaseDetails
}

type CommentGoodsBecauseOfCommentCommand struct {
	es.WithGuid
	CommentDetails
}

type CompletePaymetGoodsBecauseOfOrderCommand struct {
	es.WithGuid
	User     es.Guid
	Order    es.Guid
	Purchase es.Guid
	Quantity Quantity
}

type FailPaymetGoodsBecauseOfOrderCommand struct {
	es.WithGuid
	User     es.Guid
	Order    es.Guid
	Purchase es.Guid
	Quantity Quantity
}

// ----------------------
// GoodsPurchase Commands
// ----------------------

type CreateGoodsPurchaseCommand struct {
	es.WithGuid
	PurchaseDetails
}

type CompleteGoodsPurchaseCommand struct {
	es.WithGuid
	PurchaseDetails
}

type FailGoodsPurchaseCommand struct {
	es.WithGuid
	PurchaseDetails
}

// ----------------------
// GoodsComment Commands
// ----------------------

type CreateGoodsCommentCommand struct {
	es.WithGuid
	CommentDetails
}

type CompleteGoodsCommentCommand struct {
	es.WithGuid
	CommentDetails
}

type FailGoodsCommentCommand struct {
	es.WithGuid
	CommentDetails
}

// ----------------------
// Order Commands
// ----------------------

type CreateOrderCommand struct {
	es.WithGuid
	Items []OrderItem
}

type CancelOrderCommand struct {
	es.WithGuid
}

type CreateOrderPaymetCommand struct {
	es.WithGuid
	User           es.Guid
	UserAccount    BankAccount
	ManagedAccount BankAccount
}

type CompleteOrderPaymetCommand struct {
	es.WithGuid
	Price          Money
	User           es.Guid
	UserAccount    BankAccount
	ManagedAccount BankAccount
}

type FailOrderPaymetCommand struct {
	es.WithGuid
	Price          Money
	User           es.Guid
	UserAccount    BankAccount
	ManagedAccount BankAccount
}
