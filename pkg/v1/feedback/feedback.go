package feedback

import (
	"context"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/models"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/feedback_pb"
	"google.golang.org/grpc/codes"
)

type impl struct{}

func NewFeedBackImpl() *impl {
	return &impl{}
}

func (f *impl) Create(ctx context.Context, req *feedback_pb.CreateRequest) (*feedback_pb.CreateResponse, error) {
	data := models.Feedback{
		Content:     req.GetContent(),
		PhoneNumber: req.GetPhoneNumber(),
	}
	_, err := psql.GetDB().Model(&data).Insert()
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}
	return &feedback_pb.CreateResponse{Id: data.Id}, nil
}
