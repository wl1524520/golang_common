package lib

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDBPool(path string) error {
	//普通的db方式
	DbConfMap := &DBMapConf{}
	err := ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}

	GORMMapPool = map[string]*gorm.DB{}
	for confName, DbConf := range DbConfMap.List {
		var dbgorm *gorm.DB
		var err error
		var dial gorm.Dialector
		//gorm连接方式
		if DbConf.DriverName == "mysql" {
			// mysql
			dial = mysql.Open(DbConf.DataSourceName)
		} else {
			// sqlite
			dial = sqlite.Open(DbConf.DataSourceName)
		}
		dbgorm, err = gorm.Open(dial, &gorm.Config{})
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
		GORMMapPool[confName] = dbgorm
	}

	// 手动配置连接
	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}
	return nil
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func CloseDB() error {
	for _, dbpool := range GORMMapPool {
		sqlDB, _ := dbpool.DB()
		sqlDB.Close()
	}
	return nil
}
