package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbname := beego.AppConfig.String("dbname")
	dbhost := beego.AppConfig.String("dbhost")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbport := beego.AppConfig.String("dbport")
	dbprefix := beego.AppConfig.String("dbprefix")

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpassword, dbhost, dbport, dbname))
	// 需要在init中注册定义的model
	orm.RegisterModelWithPrefix(dbprefix, new(User))
}
