package trade

import(
    "time"
    es "github.com/sunrongya/eventsourcing"
)

type GoodsCommentService struct{
    es.Service
}

func NewGoodsCommentService(store es.EventStore) *GoodsCommentService{
    service := &GoodsCommentService{
        Service: es.NewService(store, NewGoodsComment),
    }
    return service
}

func (g *GoodsCommentService) CreateGoodsComment(user es.Guid, goods es.Guid, content string, Time time.Time) es.Guid {
    guid := es.NewGuid()
    c := &CreateGoodsCommentCommand{
        WithGuid: es.WithGuid{guid},
        CommentDetails:CommentDetails{
            User:    user,
            Goods:   goods,
            Comment: guid,
            Content: content,
            Time:    Time,
        },
    }
    g.PublishCommand(c)
    return guid
}
