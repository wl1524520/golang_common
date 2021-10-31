package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/wl1524520/golang_common/lib"
)

type Test1 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Test1) Table() string {
	return "test1"
}

var (
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	dropTableSQL = "DROP TABLE `test1`"
)

func Test_MySQL_GORM(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetGormPool("default")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	// traceCtx := lib.NewTrace()

	//设置trace信息
	// db = db.SetCtx(traceCtx)
	if err := db.Exec(createTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t1 := &Test1{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t1).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test1{}
	if err := db.Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)

	//删除表数据
	if err := db.Exec(dropTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()
	TearDown()
}
