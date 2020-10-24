package models

import "time"

type User struct {
	Id int32 `pg:"type:serial"`
	PhoneNumber string `pg:"unique,notnull"`
	Password string `pg:"notnull"`
	Status int8 `pg:"type:smallint"`
	LastLoginDate *time.Time
	LastLoginIp string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}