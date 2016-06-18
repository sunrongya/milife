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
	state CommentState
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

func (g *GoodsComment) ProcessCreateGoodsCommentCommand(command *CreateGoodsCommentCommand) []es.Event {
	return []es.Event{&GoodsCommentCreatedEvent{CommentDetails: command.CommentDetails}}
}

func (g *GoodsComment) ProcessCompleteGoodsCommentCommand(command *CompleteGoodsCommentCommand) []es.Event {
	if g.state != CommentStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsCommentCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsCommentCompletedEvent{CommentDetails: command.CommentDetails}}
}

func (g *GoodsComment) ProcessFailGoodsCommentCommand(command *FailGoodsCommentCommand) []es.Event {
	if g.state != CommentStarted {
		panic(fmt.Errorf("Can't process FailGoodsCommentCommand of state:%s", g.state))
	}
	return []es.Event{&GoodsCommentFailedEvent{CommentDetails: command.CommentDetails}}
}

func (g *GoodsComment) HandleGoodsCommentCreatedEvent(event *GoodsCommentCreatedEvent) {
	g.CommentDetails, g.state = event.CommentDetails, CommentStarted
}

func (g *GoodsComment) HandleGoodsCommentCompletedEvent(event *GoodsCommentCompletedEvent) {
	g.state = CommentCompleted
}

func (g *GoodsComment) HandleGoodsCommentFailedEvent(event *GoodsCommentFailedEvent) {
	g.state = CommentFailed
}
