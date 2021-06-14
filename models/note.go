package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Note struct {
	//("created_at","updated_at","deleted_at","key","user_id","title","summary","content","source","editor","files"
	Model
	Key     string `gorm:"unique_index;not null;"`
	UserID  int
	User    User
	Title   string
	Summary string `gorm:"type:text"`
	Content string `gorm:"type:text"`
	Source  string `gorm:"type:text" json:"source"`
	Editor  string `gorm:"varchar(40)'" `
	Files   string `gorm:"type:text"`
	Visit   int    `gorm:"default:0"`
	Praise  int    `gorm:"default:0"`
	Open  int
	Type string
//	给字段设置default的值的时候，如果插入的时候，字段值和上面给的默认值相等，就不会插入这个字段，让数据库自动写入这个数据 跪了。。
}

func (db *DB) QueryNoteByKeyAndUserId(key string, userId int) (note Note, err error) {
	return note, db.db.Model(&Note{}).Where("`key` = ? and user_id = ?", key, userId).Take(&note).Error
}

func (db *DB) QueryNoteByKey(key string) (note Note, err error) {
	return note, db.db.Model(&Note{}).Where("`key` = ?", key ).Take(&note).Error
}

// 详情页展示，增加权限控制
func (db *DB) QueryNoteByKeyWithAuthLimit(key string, userId int) (note Note, err error) {
	return note, db.db.Model(&Note{}).Where("`key` = ? and (open = 1 or user_id = ?)", key,userId).Take(&note).Error
}

func (db *DB) AllVisitCount(key string) error {
	return db.db.Model(&Note{}).Where("`key` = ?", key).UpdateColumn("visit", gorm.Expr("visit + 1")).Error
}

func (db *DB) DelNoteByKey(key string, userId int) (error) {
	return db.db.Delete(Note{}, "`key` = ? and user_id = ? ", key, userId).Error
}
func (db *DB) QueryNotesByPage(page, limit int, title string, userId uint, notetype string) (note []*Note, err error) {
	return note, db.db.Model(&Note{}).Where("(title like ? or content like ? ) and (open = 1 or user_id = ?) and `type` like ? ",
		fmt.Sprintf("%%%s%%", title),fmt.Sprintf("%%%s%%", title), userId,fmt.Sprintf("%%%s%%", notetype), ).Offset((page - 1) * limit).Limit(limit).Order("updated_at DESC").Find(&note).Error
}
func (db *DB) QueryNotesCount(title string, userId uint,notetype string) (cnt int, err error) {
	return cnt, db.db.Model(&Note{}).Where("(title like ? or content like ? ) and (open = 1 or user_id = ?) and `type` like ?",
		fmt.Sprintf("%%%s%%", title),fmt.Sprintf("%%%s%%", title), userId,fmt.Sprintf("%%%s%%", notetype)).Offset(-1).Limit(-1).Count(&cnt).Error
}

func (db *DB) SaveNote(n *Note) error {
	fmt.Println("open before saving is ")
	fmt.Println(n.Open)
	return db.db.Save(n).Error
}
