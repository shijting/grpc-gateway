package psql

import (
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/go-pg/pg/v10/orm"
	"github.com/showiot/camera/inits/config"
	"log"
	"sync"
	"time"
)
var db *pg.DB
func Init() error  {
	db = pg.Connect(&pg.Options{
		Addr:     config.Conf.DatabaseConfig.Addr,
		User:     config.Conf.DatabaseConfig.User,
		Password: config.Conf.DatabaseConfig.Passwrod,
		Database: config.Conf.DatabaseConfig.DB,
	})
	//db.AddQueryHook(dbLogger{})
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
type dbLogger struct{}
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	c2 := context.WithValue(c, "start", time.Now())
	return c2, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	start := c.Value("start").(time.Time)
	sql, err := q.FormattedQuery()
	spent := time.Since(start)
	//if err == nil {
	//	if viper.GetBool("db.sql.print") && len(sql) < viper.GetInt("db.sql.length") {
	//		println("spent: ", int(spent/time.Millisecond), " sql: ", string(sql))
	//	}
	//}
	println("spent: ", int(spent/time.Millisecond), " sql: ", string(sql))
	//if spent > viper.GetDuration("db.sql.slow_time") && len(sql) < viper.GetInt("db.sql.length") {
	//	println("spent: ", int(spent/time.Millisecond), " slow query: ", string(sql))
	//}
	return err
}