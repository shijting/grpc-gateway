package camera_messages

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/models"
	"github.com/showiot/camera/pkg/v1/cameras"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/camera_messages_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc/codes"
	"time"
)

type impl struct{}

func NewCameraMessagesImpl() *impl {
	return &impl{}
}

/**
# AddCameraMessage 新增设备消息
*/
func (*impl) AddCameraMessage(ctx context.Context, req *camera_messages_pb.AddCameraMessageRequest) (*empty.Empty, error) {
	var (
		cameraMessage *models.CameraMessage
		log           = logger.GetLogger()
		db            = psql.GetDB()
		camera        models.Camera
		err           error
	)

	if err = db.Model(&camera).Where("id = ?", req.GetCameraId()).First(); err != nil {
		if err == pg.ErrNoRows {
			return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_CAMERA_NOT_FOUND)
		}
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	// 设备没有被用户绑定
	if camera.UserID == 0 {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_CAMERA_NOT_BEING_BINDED)
	}
	cameraMessage = &models.CameraMessage{
		CameraId: req.GetCameraId(),
		VideoUrl: req.GetVideoUrl(),
		ImageUrl: req.GetImageUrl(),
		Title:    req.GetTitle(),
	}
	result, err := db.Model(cameraMessage).Insert()
	if err != nil {
		log.WithError(err).Error("insert camera message error")
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	if result.RowsAffected() == 0 {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}

	return &empty.Empty{}, nil
}

/**
# ListCameraMessage 获取设备消息
*/
func (*impl) ListCameraMessage(ctx context.Context, _ *empty.Empty) (*camera_messages_pb.ListCameraMessagesResponse, error) {
	var (
		err               error
		tokenVal          *utils.UserTokenValue
		db                = psql.GetDB()
		cameraMessageInfo []models.CameraMessageInfo
		log               = logger.GetLogger()
	)
	// 从token中获取用户Id
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	_, err = db.Query(&cameraMessageInfo, `SELECT 
			DISTINCT ON ("cm"."camera_id") camera_id, 
			"cm"."id", 
			"cm"."video_url", 
			"cm"."image_url", 
			"cm"."title", 
			"cm"."created_at", 
			"cm"."updated_at",
			"c"."name",
			"c"."model"
		FROM "camera_messages" AS "cm" 
		LEFT JOIN cameras as "c"
		ON "c"."id" = "cm"."camera_id"
		WHERE ("c"."user_id" = ?)
		ORDER BY "cm"."camera_id" ASC,"cm"."id" DESC `, tokenVal.UserID)
	if err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	cameraMessages := make([]*camera_messages_pb.CameraMessage, 0)
	for _, val := range cameraMessageInfo {
		cameraName := val.Name
		if cameraName == "" {
			cameraName = val.Model
		}
		var cameraMessage camera_messages_pb.CameraMessage
		cameraMessage.CameraName = cameraName
		copier.Copy(&cameraMessage, &val)
		cameraMessages = append(cameraMessages, &cameraMessage)
	}
	return &camera_messages_pb.ListCameraMessagesResponse{CameraMessages: cameraMessages}, nil
}

/**
# ShowCameraMessage 设备消息
*/
func (*impl) ShowCameraMessage(ctx context.Context, req *camera_messages_pb.CameraMessagesRequest) (*camera_messages_pb.CameraMessagesResponse, error) {
	var (
		err            error
		tokenVal       *utils.UserTokenValue
		cameraMessages []models.CameraMessage
		db             = psql.GetDB()
		log            = logger.GetLogger()
	)
	// 获取登录用户信息
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	// 权限判断
	if ok, err := cameras.NewCamerasImpl().IsCameraAdmin(req.GetCameraId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	err = db.Model(&cameraMessages).
		Where("camera_id = ?", req.GetCameraId()).
		Order("id desc").
		Offset(int(req.GetOffset())).
		Limit(int(req.GetLimit())).
		Select()
	if err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INTERNAL_SERVER)
	}
	respCameraMessages := make([]*camera_messages_pb.CameraMessage, 0)

	for _, val := range cameraMessages {
		cameraMessage := new(camera_messages_pb.CameraMessage)
		copier.Copy(cameraMessage, val)
		respCameraMessages = append(respCameraMessages, cameraMessage)
	}
	return &camera_messages_pb.CameraMessagesResponse{CameraMessages: respCameraMessages}, nil
}

/**
# DeleteCameraMessages 删除设备消息
# 只有管理员有操作权限
*/
func (*impl) DeleteCameraMessages(ctx context.Context, req *camera_messages_pb.DeleteCameraMessagesRequest) (*empty.Empty, error) {
	var (
		err            error
		tokenVal       *utils.UserTokenValue
		cameraMessages models.CameraMessage
		db             = psql.GetDB()
		log            = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	// 判断权限
	if ok, err := cameras.NewCamerasImpl().IsCameraAdmin(req.GetCameraId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	if _, err = db.Model(&cameraMessages).Where("camera_id = ?", req.GetCameraId()).Delete(); err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &empty.Empty{}, nil
}

/**
# DeleteCameraMessage 删除一条消息
# 只有管理员有操作权限
*/
func (*impl) DeleteCameraMessage(ctx context.Context, req *camera_messages_pb.DeleteCameraMessageRequest) (*empty.Empty, error) {
	var (
		err            error
		tokenVal       *utils.UserTokenValue
		cameraMessages models.CameraMessage
		db             = psql.GetDB()
		log            = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	cameraMessages.Id = req.GetMessageId()
	if err =db.Model(&cameraMessages).WherePK().Select(); err !=nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	// 判断权限
	if ok, err := cameras.NewCamerasImpl().IsCameraAdmin(cameraMessages.CameraId, tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}

	if _, err = db.Model(&cameraMessages).WherePK().Delete(); err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &empty.Empty{}, nil
}
func (*impl) ReadUserCameraMessage(ctx context.Context, req *camera_messages_pb.ReadUserCameraMessageRequest) (*empty.Empty, error) {
	var (
		err            error
		tokenVal       *utils.UserTokenValue
		cameraMessages models.CameraMessage
		db             = psql.GetDB()
		log            = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	if ok, err := cameras.NewCamerasImpl().IsCameraAdmin(req.GetCameraId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	cameraMessages.IsRead = true
	cameraMessages.UpdatedAt = time.Now()
	if _, err = db.Model(&cameraMessages).Where("camera_id = ?", req.GetCameraId()).Column("is_read", "updated_at").Update(); err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}

	return &empty.Empty{}, nil
}
