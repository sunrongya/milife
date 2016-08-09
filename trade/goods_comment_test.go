package trade

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestGoodsCommentRestore(t *testing.T) {
	details := CommentDetails{
		User:    es.NewGuid(),
		Goods:   es.NewGuid(),
		Comment: es.NewGuid(),
		Content: "dadfds",
		Time:    time.Now(),
	}
	comment := &GoodsComment{}
	comment.HandleGoodsCommentCreatedEvent(&GoodsCommentCreatedEvent{CommentDetails: details})
	comment.HandleGoodsCommentCompletedEvent(&GoodsCommentCompletedEvent{CommentDetails: details})

	assert.Equal(t, details, comment.CommentDetails, "details error")
	assert.Equal(t, CommentCompleted, comment._state, "state error")
}

func TestCreateGoodsCommentCommand(t *testing.T) {
	details := CommentDetails{
		User:    es.NewGuid(),
		Goods:   es.NewGuid(),
		Comment: es.NewGuid(),
		Content: "dadfds",
		Time:    time.Now(),
	}
	command := &CreateGoodsCommentCommand{CommentDetails: details}
	events := []es.Event{&GoodsCommentCreatedEvent{CommentDetails: details}}

	assert.Equal(t, events, new(GoodsComment).ProcessCreateGoodsCommentCommand(command), "")
}

func TestCompleteGoodsCommentCommand(t *testing.T) {
	details := CommentDetails{
		User:    es.NewGuid(),
		Goods:   es.NewGuid(),
		Comment: es.NewGuid(),
		Content: "dadfds",
		Time:    time.Now(),
	}
	command := &CompleteGoodsCommentCommand{CommentDetails: details}
	events := []es.Event{&GoodsCommentCompletedEvent{CommentDetails: details}}
	goodsComment := &GoodsComment{_state: CommentStarted}

	assert.Equal(t, events, goodsComment.ProcessCompleteGoodsCommentCommand(command), "")
}

func TestFailGoodsCommentCommand(t *testing.T) {
	details := CommentDetails{
		User:    es.NewGuid(),
		Goods:   es.NewGuid(),
		Comment: es.NewGuid(),
		Content: "dadfds",
		Time:    time.Now(),
	}
	command := &FailGoodsCommentCommand{CommentDetails: details}
	events := []es.Event{&GoodsCommentFailedEvent{CommentDetails: details}}
	goodsComment := &GoodsComment{_state: CommentStarted}

	assert.Equal(t, events, goodsComment.ProcessFailGoodsCommentCommand(command), "")
}

func TestCompleteGoodsCommentCommand_Panic(t *testing.T) {
	goodsComments := []*GoodsComment{
		&GoodsComment{},
		&GoodsComment{_state: CommentCompleted},
		&GoodsComment{_state: CommentFailed},
	}
	for _, goodsComment := range goodsComments {
		assert.Panics(t, func() {
			goodsComment.ProcessCompleteGoodsCommentCommand(&CompleteGoodsCommentCommand{})
		}, "执行命令CompleteGoodsCommentCommand应该抛出异常")
	}
}

func TestFailGoodsCommentCommand_Panic(t *testing.T) {
	goodsComments := []*GoodsComment{
		&GoodsComment{},
		&GoodsComment{_state: CommentCompleted},
		&GoodsComment{_state: CommentFailed},
	}
	for _, goodsComment := range goodsComments {
		assert.Panics(t, func() {
			goodsComment.ProcessFailGoodsCommentCommand(&FailGoodsCommentCommand{})
		}, "执行命令FailGoodsCommentCommand应该抛出异常")
	}
}
