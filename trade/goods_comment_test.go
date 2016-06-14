package trade

import(
    "fmt"
    "time"
    "testing"
    "github.com/stretchr/testify/assert"
    es "github.com/sunrongya/eventsourcing"
)

func TestGoodsCommentRestore(t *testing.T) {
    details := CommentDetails {
        User:    es.NewGuid(),
        Goods:   es.NewGuid(),
        Comment: es.NewGuid(),
        Content: "dadfds",
        Time:    time.Now(),
    }
    comment := &GoodsComment{}
    comment.ApplyEvents([]es.Event{
        &GoodsCommentCreatedEvent{WithGuid:es.WithGuid{details.Comment}, CommentDetails:details },
        &GoodsCommentCompletedEvent{ WithGuid:es.WithGuid{details.Comment}, CommentDetails:details },
    })
    assert.Equal(t, 2,                 comment.Version(), "version error")
    assert.Equal(t, details,           comment.CommentDetails, "details error")
    assert.Equal(t, CommentCompleted,  comment.state, "state error")
}

func TestGoodsCommentRestoreForErrorEvent(t *testing.T){
    assert.Panics(t, func(){ 
        NewGoodsComment().ApplyEvents([]es.Event{ &struct{es.WithGuid}{} }) 
    }, "restore error event must panic error")
}

func TestCheckGoodsCommentApplyEvents(t *testing.T) {
    events := []es.Event{
        &GoodsCommentCreatedEvent{},
        &GoodsCommentCompletedEvent{},
        &GoodsCommentFailedEvent{},
    }
    assert.NotPanics(t, func(){ NewGoodsComment().ApplyEvents(events) }, "Check Process All Event")
}

func TestGoodsCommentCommand(t *testing.T){
    details := CommentDetails {
        User:    es.NewGuid(),
        Goods:   es.NewGuid(),
        Comment: es.NewGuid(),
        Content: "dadfds",
        Time:    time.Now(),
    }
    
    tests := []struct{
        comment  *GoodsComment
        command  es.Command
        event  es.Event
    }{
        {
            &GoodsComment{},
            &CreateGoodsCommentCommand{WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
            &GoodsCommentCreatedEvent{WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
        },
        {
            &GoodsComment{ state:CommentStarted },
            &CompleteGoodsCommentCommand{ WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
            &GoodsCommentCompletedEvent{ WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
        },
        {
            &GoodsComment{ state:CommentStarted },
            &FailGoodsCommentCommand{ WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
            &GoodsCommentFailedEvent{ WithGuid:es.WithGuid{Guid:details.Comment}, CommentDetails:details },
        },
    }
    
    for _, v := range tests {
        assert.Equal(t, []es.Event{v.event}, v.comment.ProcessCommand(v.command) )
    }
}

func TestGoodsCommentCommand_Panic(t *testing.T){
    tests := []struct{
        comment  *GoodsComment
        command  es.Command
    }{
        {
            &GoodsComment{},
            &struct{es.WithGuid}{},
        },
        {
            &GoodsComment{},
            &CompleteGoodsCommentCommand{},
        },
        {
            &GoodsComment{ state:CommentCompleted },
            &CompleteGoodsCommentCommand{},
        },
        {
            &GoodsComment{ state:CommentFailed },
            &CompleteGoodsCommentCommand{},
        },
        {
            &GoodsComment{},
            &FailGoodsCommentCommand{},
        },
        {
            &GoodsComment{ state:CommentCompleted },
            &FailGoodsCommentCommand{},
        },
        {
            &GoodsComment{ state:CommentFailed },
            &FailGoodsCommentCommand{},
        },
    }
    
    for _, v := range tests {
        assert.Panics(t, func(){v.comment.ProcessCommand(v.command)}, fmt.Sprintf("test panics error: command:%v", v.command))
    }
}
