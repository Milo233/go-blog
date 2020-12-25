package controllers

import (
	"github.com/jicg/go-blog/syserrors"
	"fmt"
	"errors"
	"time"
)

type IndexController struct {
	BaseController
}


// @router /appinfo [get]
func (c *IndexController) Info() {
	db_t:=c.Dao.GetDBTime()
	var db_ts =""
	if db_t!=nil{
		db_ts=db_t.Format("2006-01-02 15:04:05")
	}
	c.JSONOkH("ok",H{
		"app_time":time.Now().Format("2006-01-02 15:04:05"),
		"db_time":db_ts,
	})
}

// @router / [get]
func (c *IndexController) Get() {
	userId := c.User.ID
	limit := 10;
	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		page = 1;
	}
	title := c.GetString("title", "")
	noteType := c.GetString("type", "")
	if c.Dao==nil{
		c.Abort500(errors.New("数据库初始化失败！"))
	}
	ns, err := c.Dao.QueryNotesByPage(page, limit, title, userId,noteType)
	if err != nil {
		c.Abort500(err)
	}
	if ns != nil {
		c.Data["notes"] = ns
	}
	var totpage int = 0;
	totcnt, _ := c.Dao.QueryNotesCount(title,userId,noteType)
	if totcnt%limit == 0 {
		totpage = totcnt / limit
	} else {
		totpage = totcnt/limit + 1
	}
	c.Data["totpage"] = totpage
	c.Data["page"] = page
	c.Data["title"] = title
	c.Data["type"] = noteType
	c.TplName = "index.html"
}

// @router /details/:key [get]
func (c *IndexController) GetDetail() {
	key := c.Ctx.Input.Param(":key")
	note, err := c.Dao.QueryNoteByKeyWithAuthLimit(key,int(c.User.ID))
	if err != nil {
		c.Abort500(syserrors.NewError("文章不存在 or 没有查看权限！", err))
	}
	go c.Dao.AllVisitCount(key)
	c.Data["praise"] = false
	//praise, err := c.Dao.QueryPraiseLog(key, int(c.User.ID), "note")
	//if err == nil {
	//	c.Data["praise"] = praise.Flag
	//}
	messages, _ := c.Dao.QueryMessageForNote(note.Key)
	c.Data["messages"] = messages
	c.Data["note"] = note
	c.TplName = "details.html"
}

// @router /comment/:key [get]
func (c *IndexController) GetComment() {

	key := c.Ctx.Input.Param(":key")
	note, err := c.Dao.QueryNoteByKey(key)
	if err != nil {
		c.Abort500(syserrors.NewError("文章不存在 or 没有查看权限！", err))
	}

	c.Data["note"] = note
	c.TplName = "comment.html"
}

// @router /setting [get]
func (c *IndexController) GetSetting() {
	c.TplName = "setting.html"
}

// @router /user [get]
func (c *IndexController) GetUser() {
	c.TplName = "user.html"
}

// @router /reg [get]
func (c *IndexController) GetReg() {
	c.TplName = "reg.html"
}

// @router /message [get]
func (c *IndexController) GetMessage() {
	messages, err := c.Dao.QueryMessageForNote("")
	if err != nil {
		c.Abort500(err)
	}
	fmt.Printf("%v \n", messages)
	c.Data["messages"] = messages
	c.TplName = "message.html"
}

// @router /about [get]
func (c *IndexController) GetAbout() {
	c.TplName = "about.html"
}
