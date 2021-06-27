package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"
	"github.com/Milo233/go-blog/models"
	_ "github.com/Milo233/go-blog/models"
	_ "github.com/Milo233/go-blog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	initLog()
	initSession()
	initTemplate()
	openBrowser()
	beego.Run()
}
func initLog() {
	if err := os.MkdirAll("data/logs", 0777); err != nil {
		beego.Error(err)
		return
	}
	logs.SetLogger("file", `{"filename":"data/logs/lyblog.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.Async(1e3)
}

func initSession() {
	gob.Register(models.User{})
	//https://beego.me/docs/mvc/controller/session.md
	beego.SetStaticPath("assert", "assert")
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "go-blog-key"
	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24 * 7 // 七天过期

	beego.BConfig.WebConfig.Session.SessionProviderConfig = "data/session"
}
func initTemplate() {
	beego.AddFuncMap("equrl", func(x, y string) bool {
		s1 := strings.Trim(x, "/")
		s2 := strings.Trim(y, "/")
		return strings.Compare(s1, s2) == 0
	})
	beego.AddFuncMap("eq2", func(x, y interface{}) bool {
		s1 := fmt.Sprintf("%v", x)
		s2 := fmt.Sprintf("%v", y)
		return strings.Compare(s1, s2) == 0
	})
	beego.AddFuncMap("add", func(x, y int) int {
		return x + y
	})
	beego.AddFuncMap("json", func(obj interface{}) string {
		bs, err := json.Marshal(obj)
		if err != nil {
			return "{id:0}"
		}
		return string(bs)
	})

}

func openBrowser() {
	fmt.Println("start ... ")
	time.Sleep(1 * time.Second)
	url := "http://" + GetOutboundIP() + ":6010"
	fmt.Println("Running at " + url)
	err := open(url)
	if err != nil {
		fmt.Println("failed to open browser!")
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}