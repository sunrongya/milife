package trade

import(
    "testing"
    "time"
    es "github.com/sunrongya/eventsourcing"
    "github.com/sunrongya/eventsourcing/utiltest"
)

func TestGoodsCommentServiceDoPublishGoods(t *testing.T) {
    utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command{
        details := CommentDetails{
            User:     es.NewGuid(),
            Goods:    es.NewGuid(),
            Comment:  es.NewGuid(),
            Content:  "abc",
            Time:     time.Now(),
        }
        
        gs := GoodsCommentService{ Service: service}
        details.Comment = gs.CreateGoodsComment(details.User, details.Goods, details.Content, details.Time)
        return &CreateGoodsCommentCommand{WithGuid:es.WithGuid{details.Comment}, CommentDetails:details }
    })
}


