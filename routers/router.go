package routers

import (
	"chat/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "*:Index")
	beego.Router("/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/register", &controllers.HomeController{}, "*:Register")
	beego.Router("/edit", &controllers.HomeController{}, "*:Editing")
	beego.Router("/logout", &controllers.HomeController{}, "*:Logout")

	// 消息
	beego.Router("/post", &controllers.HomeController{}, "post:Post")
	beego.Router("/fetch", &controllers.HomeController{}, "*:Fetch")
}
