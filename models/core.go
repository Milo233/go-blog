package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"github.com/astaxie/beego/logs"
	"time"
	"github.com/astaxie/beego"
)

type DB struct {
	db *gorm.DB
}

func (db *DB) Begin() {
	db.db = db.db.Begin()
}

func (db *DB) Rollback() {
	db.db = db.db.Rollback()
}

func (db *DB) Commit() {
	db.db = db.db.Commit()
}

var (
	db *gorm.DB
)

func NewDB() *DB {
	return &DB{db: db}
}

func init() {
	var err error
	// 创建data目录
	if err = os.MkdirAll("data", 0777); err != nil {
		panic("failed to connect database," + err.Error())
	}
	if err = initDB(); err != nil {
		panic("failed to connect database," + err.Error())
	}
	// 自动同步表结构
	db.SetLogger(logs.GetLogger("orm"))
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Note{}, &Message{}, &PraiseLog{})
	// Model(&User{})查询用户表, Count(&count) 将用户表的数据赋值给count字段。
	var count int
	if err := db.Model(&User{}).Count(&count).Error; err == nil && count == 0 {
		db.Create(&User{Name: "admin",
			//邮箱
			Email: "admin@qq.com",
			//密码
			Pwd: "123123",
			//头像地址
			Avatar: "/static/images/go.ico",
			//是否认证 例： lyblog 作者
			Role: 0,
		})
	}
}

func initDB() error {
	var err error
	dbconf, err := beego.AppConfig.GetSection("database");
	if err != nil {
		logs.Error(err)
		dbconf = map[string]string{
			"type": "sqlite3",
		}
	}
	switch dbconf["type"] {
	case "mysql":
		db, err = gorm.Open("mysql", dbconf["url"])
	default: // fixme read from env config file
		// 原来写死了读的data文件夹下的data.db 
		// db, err = gorm.Open("sqlite3", "data/data.db")
		db, err = gorm.Open("sqlite3", dbconf["url"])
	}
	if err != nil {
		return err
	}
	return nil
}

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createtime"`
	UpdatedAt time.Time  `json:"updatetime"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (db *DB) GetDBTime() *time.Time {
	var t *time.Time
	row, err := db.db.DB().Query("select NOW()")
	if err != nil {
		logs.Error(err)
		return nil
	}
	row.Scan(t)
	return t
}
