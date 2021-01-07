package models

import "time"

type Camera struct {
	Id         uint32 `pg:"id,pk"`
	No         string `pg:"no,notnull"`
	Name       string `pg:"name,use_zero"`
	Model      string `pg:"model,notnull"`
	Mac        string `pg:"mac,notnull"`
	Ip         string
	Port       uint32
	UserID     uint32 `pg:"user_id,use_zero"`
	Password   string `pg:"password,use_zero"`
	IsAlarm    bool   `pg:"is_alarm,use_zero"`
	Status     int8   `pg:"default:1,use_zero"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserCamera *UserCamera `pg:"rel:belongs-to,fk:camera_id"`
}
