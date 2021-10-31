package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/wl1524520/golang_common/lib"
)

type Test2 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Test2) Table() string {
	return "test2"
}

var (
	createTableSQL2 = "CREATE TABLE `test2` (`id`	INTEGER," +
		" `name` varchar(20) NOT NULL DEFAULT '', " +
		"`created_at` datetime NOT NULL DEFAULT '', " +
		"PRIMARY KEY(`id` AUTOINCREMENT))"
	insertSQL2    = "INSERT INTO `test2` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2018-08-29 11:01:43');"
	dropTableSQL2 = "DROP TABLE `test2`"
	beginSQL2     = "start transaction;"
	commitSQL2    = "commit;"
	rollbackSQL2  = "rollback;"
)

func Test_Sqlite_GORM(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetGormPool("sqlite")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	// traceCtx := lib.NewTrace()

	//设置trace信息
	// db = db.SetCtx(traceCtx)
	if err := db.Exec(createTableSQL2).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t2 := &Test2{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t2).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test2{}
	if err := db.Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)

	// 删除表数据
	// if err := db.Exec(dropTableSQL2).Error; err != nil {
	// 	db.Rollback()
	// 	t.Fatal(err)
	// }
	db.Commit()
	TearDown()
}
