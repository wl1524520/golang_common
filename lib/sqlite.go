package lib

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqliteDBPool(path string) error {
	//普通的db方式
	DbConfMap := &SqliteMapConf{}
	err := ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty sqlite config.")
	}

	GORMSqliteMapPool = map[string]*gorm.DB{}
	for confName, DbConf := range DbConfMap.List {
		//gorm连接方式
		dbgorm, err := gorm.Open(sqlite.Open(DbConf.DataSourceName), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, err := dbgorm.DB()
		// dbgorm.SingularTable(true)
		err = sqlDB.Ping()
		if err != nil {
			return err
		}
		// dbgorm.LogMode(true)
		// dbgorm.LogCtx(true)
		// dbgorm.SetLogger(&MysqlGormLogger{Trace: NewTrace()})
		sqlDB.SetMaxIdleConns(DbConf.MaxIdleConn)
		sqlDB.SetMaxOpenConns(DbConf.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		GORMSqliteMapPool[confName] = dbgorm
	}

	if dbpool, err := GetSqliteGormPool("default"); err == nil {
		GORMSqliteDefaultPool = dbpool
	}
	return nil
}

func GetSqliteGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMSqliteMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func CloseSqliteDB() error {
	for _, dbpool := range GORMSqliteMapPool {
		sqlDB, _ := dbpool.DB()
		sqlDB.Close()
	}
	return nil
}
