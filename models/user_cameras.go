package models

import "time"

type UserCamera struct {
	Id          uint32 `pg:"id,pk"`
	CameraId    uint32 `pg:"camera_id"`
	UserId      uint32 `pg:"user_id"`
	Permissions int32
	IsAdmin     bool
	Camera      *Camera `pg:"rel:has-one"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
