package models

import "time"

type CameraMessage struct {
	Id                uint32 `pg:",pk"`
	CameraId          uint32
	VideoUrl          string             `pg:"video_url,use_zero"`
	ImageUrl          string             `pg:"image_url,use_zero"`
	Title             string             `pg:"title,use_zero"`
	IsRead            bool               `pg:",use_zero"`
	Camera            *Camera            `pg:"rel:has-one"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CameraMessageInfo struct {
	tableName struct{} `pg:"camera_messages"`
	Id        uint32   `pg:",pk"`
	CameraId  uint32
	VideoUrl  string `pg:"video_url,use_zero"`
	ImageUrl  string `pg:"image_url,use_zero"`
	Title     string `pg:"title,use_zero"`
	IsRead    bool   `pg:",use_zero"`
	Name      string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
