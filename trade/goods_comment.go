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

type CommentDetails struct {
	User     es.Guid
	Goods    es.Guid
	Purchase es.Guid
	Comment  es.Guid
	Content  string
	Time     time.Time
}

func (g *GoodsComment) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *GoodsCommentCreatedEvent:
			g.CommentDetails, g.state = e.CommentDetails, CommentStarted
		case *GoodsCommentCompletedEvent:
			g.state = CommentCompleted
		case *GoodsCommentFailedEvent:
			g.state = CommentFailed
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	g.SetVersion(len(events))
}

func (g *GoodsComment) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateGoodsCommentCommand:
		event = g.processCreateGoodsCommentCommand(c)
	case *CompleteGoodsCommentCommand:
		event = g.processCompleteGoodsCommentCommand(c)
	case *FailGoodsCommentCommand:
		event = g.processFailGoodsCommentCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (g *GoodsComment) processCreateGoodsCommentCommand(command *CreateGoodsCommentCommand) es.Event {
	return &GoodsCommentCreatedEvent{CommentDetails: command.CommentDetails}
}

func (g *GoodsComment) processCompleteGoodsCommentCommand(command *CompleteGoodsCommentCommand) es.Event {
	if g.state != CommentStarted {
		panic(fmt.Errorf("Can't process CompleteGoodsCommentCommand of state:%s", g.state))
	}
	return &GoodsCommentCompletedEvent{CommentDetails: command.CommentDetails}
}

func (g *GoodsComment) processFailGoodsCommentCommand(command *FailGoodsCommentCommand) es.Event {
	if g.state != CommentStarted {
		panic(fmt.Errorf("Can't process FailGoodsCommentCommand of state:%s", g.state))
	}
	return &GoodsCommentFailedEvent{CommentDetails: command.CommentDetails}
}

func NewGoodsComment() es.Aggregate {
	return &GoodsComment{}
}
