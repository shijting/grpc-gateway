package models

import "time"

type User struct {
	Id          uint32 `pg:"id,type:serial"`
	PhoneNumber string `pg:",unique,notnull"`
	Nickname    string `pg:"nickname"`
	Status      int8   `pg:"type:smallint"`
	LastLoginAt time.Time
	LastLoginIp string
	Password    string
	CreatedAt   time.Time
	Avatar      string
	UpdatedAt   time.Time
}
