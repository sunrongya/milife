package trade

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Goods Events
// --------------

type GoodsPublishedEvent struct {
	es.WithGuid
	Name     string
	Price    Money
	Quantity Quantity
	SN       SN
}

type GoodsAuditedPassEvent struct {
	es.WithGuid
}

type GoodsAuditedNoPassEvent struct {
	es.WithGuid
}

type GoodsOnlinedEvent struct {
	es.WithGuid
}

type GoodsOfflinedEvent struct {
	es.WithGuid
}

type GoodsPurchaseSuccessedEvent struct {
	es.WithGuid
	PurchaseDetails
}

type GoodsPurchaseFailuredEvent struct {
	es.WithGuid
	PurchaseDetails
}

type GoodsCommentSuccessedEvent struct {
	es.WithGuid
	CommentDetails
}

type GoodsCommentFailuredEvent struct {
	es.WithGuid
	CommentDetails
}

type PaymetGoodsCompletedBecauseOfOrderEvent struct {
	es.WithGuid
	User     es.Guid
	Order    es.Guid
	Purchase es.Guid
	Quantity Quantity
}

type PaymetGoodsFailedBecauseOfOrderEvent struct {
	es.WithGuid
	User     es.Guid
	Order    es.Guid
	Purchase es.Guid
	Quantity Quantity
}

// --------------------
// GoodsPurchase Events
// --------------------

type GoodsPurchaseCreatedEvent struct {
	es.WithGuid
	PurchaseDetails
}

type GoodsPurchaseCompletedEvent struct {
	es.WithGuid
	PurchaseDetails
}

type GoodsPurchaseFailedEvent struct {
	es.WithGuid
	PurchaseDetails
}

// --------------------
// GoodsComment Events
// --------------------

type GoodsCommentCreatedEvent struct {
	es.WithGuid
	CommentDetails
}

type GoodsCommentCompletedEvent struct {
	es.WithGuid
	CommentDetails
}

type GoodsCommentFailedEvent struct {
	es.WithGuid
	CommentDetails
}

// --------------------
// Order Events
// --------------------

type OrderCreatedEvent struct {
	es.WithGuid
	Items []OrderItem
}

type OrderCanceledEvent struct {
	es.WithGuid
}

type OrderPaymetCreatedEvent struct {
	es.WithGuid
	Price          Money
	User           es.Guid
	UserAccount    BankAccount
	ManagedAccount BankAccount
}

type OrderPaymetCompletedEvent struct {
	es.WithGuid
	User       es.Guid
	OrderItems []OrderItem
}

type OrderPaymetFailedEvent struct {
	es.WithGuid
	User       es.Guid
	OrderItems []OrderItem
}
