package trade

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type CommentState string

const (
	CommentStarted   = CommentState("CommentStarted")
	CommentCompleted = CommentState("CommentCompleted")
	CommentFailed    = CommentState("CommentFailed")
)

type GoodsComment struct {
	es.BaseAggregate
	CommentDetails
	_state CommentState
}

var _ es.Aggregate = (*GoodsComment)(nil)

func NewGoodsComment() es.Aggregate {
	return &GoodsComment{}
}

type CommentDetails struct {
	User     es.Guid
	Goods    es.Guid
	Purchase es.Guid
	Comment  es.Guid
	Content  string
	Time     time.Time
}

func (this *GoodsComment) ProcessCreateGoodsCommentCommand(command *CreateGoodsCommentCommand) []es.Event {
	return []es.Event{&GoodsCommentCreatedEvent{CommentDetails: command.CommentDetails}}
}

func (this *GoodsComment) ProcessCompleteGoodsCommentCommand(command *CompleteGoodsCommentCommand) []es.Event {
	if this._state != CommentStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsCommentCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsCommentCompletedEvent{CommentDetails: command.CommentDetails}}
}

func (this *GoodsComment) ProcessFailGoodsCommentCommand(command *FailGoodsCommentCommand) []es.Event {
	if this._state != CommentStarted {
		panic(fmt.Errorf("Can't process FailGoodsCommentCommand of state:%s", this._state))
	}
	return []es.Event{&GoodsCommentFailedEvent{CommentDetails: command.CommentDetails}}
}

func (this *GoodsComment) HandleGoodsCommentCreatedEvent(event *GoodsCommentCreatedEvent) {
	this.CommentDetails, this._state = event.CommentDetails, CommentStarted
}

func (this *GoodsComment) HandleGoodsCommentCompletedEvent(event *GoodsCommentCompletedEvent) {
	this._state = CommentCompleted
}

func (this *GoodsComment) HandleGoodsCommentFailedEvent(event *GoodsCommentFailedEvent) {
	this._state = CommentFailed
}
