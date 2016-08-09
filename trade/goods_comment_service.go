package trade

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type GoodsCommentService struct {
	es.Service
}

func NewGoodsCommentService(store es.EventStore) *GoodsCommentService {
	service := &GoodsCommentService{
		Service: es.NewService(store, NewGoodsComment),
	}
	return service
}

func (this *GoodsCommentService) CreateGoodsComment(user es.Guid, goods es.Guid, content string, Time time.Time) es.Guid {
	guid := es.NewGuid()
	c := &CreateGoodsCommentCommand{
		WithGuid: es.WithGuid{guid},
		CommentDetails: CommentDetails{
			User:    user,
			Goods:   goods,
			Comment: guid,
			Content: content,
			Time:    Time,
		},
	}
	this.PublishCommand(c)
	return guid
}
