package psql

import (
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/go-pg/pg/v10/orm"
	"github.com/showiot/camera/inits/config"
	"log"
	"sync"
)
var db *pg.DB
func Init() error  {
	db = pg.Connect(&pg.Options{
		Addr:     config.Conf.DatabaseConfig.Addr,
		User:     config.Conf.DatabaseConfig.User,
		Password: config.Conf.DatabaseConfig.Passwrod,
		Database: config.Conf.DatabaseConfig.DB,
	})
	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		return err
	}
	return nil
}
func Reload() (err error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	if db != nil {
		err = db.Close()
		if err != nil {
			return
		}
		db = nil

	}
	err = Init()
	log.Println("psql 重载成功...")
	return
}
func GetDB() *pg.DB {
	return db
}
func Close() (err error)  {
	if db !=nil {
		err = db.Close()
	}
	return
}