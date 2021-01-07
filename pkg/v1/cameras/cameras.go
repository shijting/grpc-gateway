package cameras

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/models"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/cameras_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc/codes"
	"sync"
	"time"
)

type impl struct {
	sync.Mutex
}

func NewCamerasImpl() *impl {
	return &impl{}
}

// ShowCamera 查询设备
func (this *impl) ShowCamera(ctx context.Context, req *cameras_pb.ShowCameraRequest) (resp *cameras_pb.Camera, err error) {
	data := &models.Camera{Id: req.GetCameraId()}
	err = psql.GetDB().Model(data).
		WherePK().
		Relation("UserCamera").
		First()
	if err != nil && err != pg.ErrNoRows {
		// 设备不存在
		if err == pg.ErrNoRows {
			return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_CAMERA_NOT_FOUND)
		}
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	resp = new(cameras_pb.Camera)
	copier.Copy(resp, data)
	return
}

// UpdateCamera 更新设备
func (this *impl) UpdateCamera(ctx context.Context, req *cameras_pb.UpdateCameraRequest) (*empty.Empty, error) {
	var (
		db       = psql.GetDB()
		err      error
		tokenVal *utils.UserTokenValue
	)

	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	// 不是设备管理员，无权操作
	if ok, err := this.IsCameraAdmin(req.GetCameraId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	data := &models.Camera{Id: req.CameraId}
	column := make([]string, 0)
	switch req.GetUpdateType() {
	case cameras_pb.UpdateCameraRequest_IS_ALARM:
		column = append(column, "is_alarm")
		data.IsAlarm = req.GetIsAlarm()
	case cameras_pb.UpdateCameraRequest_PASSWORD:
		column = append(column, "password")
		password, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
		}
		data.Password = password
	case cameras_pb.UpdateCameraRequest_NAME:
		column = append(column, "name")
		data.Name = req.GetName()
	}
	_, err = db.Model(data).
		WherePK().
		Column(column...).
		Update()
	if err != nil {
		logger.GetLogger().Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &empty.Empty{}, nil
}

// ResetCamera 重置设备
func (this *impl) ResetCamera(ctx context.Context, req *cameras_pb.ResetCameraRequest) (*empty.Empty, error) {
	data := &models.Camera{
		Id:        req.CameraId,
		Password:  "",
		IsAlarm:   false,
		UserID:    0,
		UpdatedAt: time.Now(),
		Name:      "",
	}

	tx, err := psql.GetDB().Begin()
	if err != nil {
		logger.GetLogger().Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	_, err = tx.Model(data).WherePK().Column("password", "is_alarm", "user_id", "name", "updated_at").Update()
	if err != nil {
		logger.GetLogger().Error(err)
		tx.Rollback()
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	// 删除 user_camera 相关数据
	_, err = tx.Model(&models.UserCamera{}).Where("camera_id = ?", req.GetCameraId()).Delete()
	if err != nil {
		logger.GetLogger().Error(err)
		tx.Rollback()
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	if err = tx.Commit(); err != nil {
		logger.GetLogger().Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &empty.Empty{}, nil
}

// BindCamera 绑定设备
func (*impl) BindCamera(ctx context.Context, req *cameras_pb.BindCameraRequest) (*cameras_pb.BindCameraResponse, error) {
	var (
		camera   *models.Camera
		tx       *pg.Tx
		err      error
		user     *models.User
		tokenVal *utils.UserTokenValue
		db       = psql.GetDB()
		log      = logger.GetLogger()
	)
	// 从token中获取用户Id
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	if tx, err = db.Begin(); err != nil {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	camera = &models.Camera{}
	if err = tx.Model(camera).Where("no = ? ", req.GetNo()).First(); err != nil {
		// 设备不存在
		if err == pg.ErrNoRows {
			return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_CAMERA_NOT_FOUND)
		}
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}

	user = &models.User{Id: tokenVal.UserID}
	if err = tx.Model(user).WherePK().First(); err != nil && err != pg.ErrNoRows {
		logger.GetLogger().WithField("userID", camera.UserID).WithError(err).Error("get user failed")
	}
	// 设备已被绑定
	if camera.UserID > 0 {
		return &cameras_pb.BindCameraResponse{UserId: camera.UserID, IsSuccess: 0, PhoneNumber: utils.MaskedMobile(user.PhoneNumber)}, nil
	}
	camera = &models.Camera{Id: camera.Id, UserID: user.Id, UpdatedAt: time.Now()}
	updateResult, err := tx.Model(camera).WherePK().Column("user_id", "updated_at").Update()
	if err != nil || updateResult.RowsAffected() == 0 {
		tx.Rollback()
		log.WithError(err).Error("updating camera table error")
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	userCamera := &models.UserCamera{
		CameraId: camera.Id,
		UserId:   user.Id,
		IsAdmin:  true,
	}
	userResult, err := tx.Model(userCamera).Insert()
	if err != nil || userResult.RowsAffected() == 0 {
		tx.Rollback()
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}
	if err = tx.Commit(); err != nil {
		log.WithError(err).Error("commit transaction error")
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &cameras_pb.BindCameraResponse{UserId: user.Id, IsSuccess: 1, PhoneNumber: utils.MaskedMobile(user.PhoneNumber)}, nil
}

// ListCameras 获取用户设备列表
func (*impl) ListUserCamera(ctx context.Context, _ *empty.Empty) (*cameras_pb.ListUserCameraResponse, error) {
	var (
		userCameras []models.UserCamera
		db          = psql.GetDB()
		tokenVal    *utils.UserTokenValue
		err         error
	)
	// 从token中获取用户Id
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}

	err = db.Model(&userCameras).
		Where("user_camera.user_id = ?", tokenVal.UserID).
		Relation("Camera").
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return &cameras_pb.ListUserCameraResponse{
				UserCameras: nil,
			}, nil
		}
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	respUserCamera := make([]*cameras_pb.UserCamera, 0)
	for _, val := range userCameras {
		userCamera := new(cameras_pb.UserCamera)
		copier.Copy(userCamera, val)
		respUserCamera = append(respUserCamera, userCamera)
	}

	return &cameras_pb.ListUserCameraResponse{
		UserCameras: respUserCamera,
	}, nil
}

/**
# UpdateSharePermission 设置分享权限
# 必须是设备管理员才有操作权限
# 必须是已分享给该用户的设备
# SharaId 设备id
# UserId 被分享人id
*/
func (this *impl) UpdateSharePermission(ctx context.Context, req *cameras_pb.UpdateSharePermissionRequest) (*empty.Empty, error) {
	var (
		db       = psql.GetDB()
		tokenVal *utils.UserTokenValue
		err      error
		log      = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		log.WithError(err).Error("")
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}

	// 不是设备管理员，无权操作
	if ok, err := this.IsCameraAdmin(req.GetShareId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	data := &models.UserCamera{Permissions: req.GetPermission(), UpdatedAt: time.Now()}
	result, err := db.Model(data).Where("camera_id = ? and user_id = ?", req.GetShareId(), req.GetUserId()).Column("permissions", "updated_at").Update()
	if err != nil || result.RowsAffected() == 0 {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_OPERATION_FAILED)
	}
	return &empty.Empty{}, nil
}

/**
# AddShare 添加分享
# 选择用户 + 选择权限 = 分享
# 必须是设备管理员分享给其他用户
# 管理员不能分享给自己
# 被分享的用户必须已注册
*/
func (this *impl) AddShare(ctx context.Context, req *cameras_pb.AddShareRequest) (*empty.Empty, error) {
	var (
		userCamera models.UserCamera
		db         = psql.GetDB()
		tokenVal   *utils.UserTokenValue
		err        error
		user       models.User
		log        = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}

	// 不是设备管理员，无权操作
	if ok, err := this.IsCameraAdmin(req.GetCameraId(), tokenVal.UserID); err != nil || !ok {
		return nil, proto.Errorg(codes.PermissionDenied, proto.Error_ERR_UNAUTHORIZED)
	}
	// 不能分享给自己
	if tokenVal.UserID == req.UserId {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_UNAUTHORIZED)
	}
	if err = db.Model(&user).Where("id = ? ", req.UserId).First(); err != nil {
		// 被分享用户不存在
		if err == pg.ErrNoRows {
			return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_SHARED_USER_NOT_EXIST)
		}
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}

	//
	if err = db.Model(&userCamera).Where("user_id = ? and camera_id = ? ", req.UserId, req.CameraId).First(); err != nil && err != pg.ErrNoRows {
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	// 重复的分享
	if userCamera.Id != 0 {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_SHARE_REPEATED)
	}
	userCamera = models.UserCamera{
		CameraId:    req.CameraId,
		UserId:      req.UserId,
		Permissions: req.Permission,
	}
	result, err := db.Model(&userCamera).Insert()
	if err != nil || result.RowsAffected() == 0 {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}
	return &empty.Empty{}, nil
}

/**
# CancelShare 取消分享
# ShareId 设备id
# 设备管理员取消分享
*/
func (*impl) CancelShare(ctx context.Context, req *cameras_pb.CancelShareRequest) (*empty.Empty, error) {
	var (
		camera     models.Camera
		userCamera models.UserCamera
		db         = psql.GetDB()
		tokenVal   *utils.UserTokenValue
		err        error
		log        = logger.GetLogger()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	if err = db.Model(&camera).
		Where("id = ? and user_id = ?", req.ShareId, tokenVal.UserID).
		First(); err != nil {
		if err == pg.ErrNoRows {
			return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_CAMERA_NOT_FOUND)
		}
		log.Error(err)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	// TODO: 优化
	if _, err = db.Model(&userCamera).Where("camera_id = ? and user_id not in (?) ", req.ShareId, tokenVal.UserID).Delete(); err != nil {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}

	return &empty.Empty{}, nil
}

/**
# DeleteShare 删除分享
# ShareId 设备id
# 被分享用户主动删除分享
*/
func (*impl) DeleteShare(ctx context.Context, req *cameras_pb.DeleteShareRequest) (*empty.Empty, error) {
	tokenVal, err := utils.GetUserInfoFromToken(ctx)
	if err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	result, err := psql.GetDB().Model(&models.UserCamera{}).Where("camera_id = ? and user_id =? ", req.ShareId, tokenVal.UserID).Delete()
	if err != nil || result.RowsAffected() == 0 {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_OPERATION_FAILED)
	}
	return &empty.Empty{}, nil
}

// 是否是设备管理员
func (*impl) IsCameraAdmin(cameraId, userId uint32) (bool, error) {
	var (
		err    error
		camera models.Camera
		db     = psql.GetDB()
		log    = logger.GetLogger()
	)
	camera = models.Camera{Id: cameraId}
	if err = db.Model(&camera).WherePK().First(); err != nil {
		if err == pg.ErrNoRows {
			return false, fmt.Errorf("camera:%d is not rows", cameraId)
		}
		log.Error(err)
		return false, err
	}
	if camera.UserID != userId {
		return false, fmt.Errorf("is not camera admin")
	}
	return true, nil
}
