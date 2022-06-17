package controllers

import (
	"chat/models"
	"chat/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)
type HomeController struct {
	beego.Controller
}

// 获取用户结构体
func (this *HomeController) GetUserByID(ID int) models.User {
	user := models.User{Id: ID}

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	o.Read(&user)
	return user
}

// 首页
func (this *HomeController) Index() {
	userID := this.GetSession("userID")
	if userID == nil {
		this.Ctx.WriteString("<script>alert('未登录');window.location.href='login';</script>")
	}
	data := this.GetUserByID(userID.(int))
	// 加入聊天室
	Join(data.Name, nil)
	this.TplName = "home/index.html"
	this.Data["user_id"] = data.Id
}

// 登录
func (this *HomeController) Login() {
	userID := this.GetSession("userID")
	if userID != nil {
		this.Ctx.WriteString("<script>alert('已登录');window.location.href='/';</script>")
	}
	if this.Ctx.Request.Method == "POST" {
		username := this.GetString("username")
		password := this.GetString("password")
		user := models.User{Name: username}

		o := orm.NewOrm()
		o.Using("default") // 默认使用 default，你可以指定为其他数据库

		err := o.Read(&user, "Name")

		if user.Password == "" {
			mystruct := util.Msg{false, "账号不存在"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		if util.Md5(password) != strings.Trim(user.Password, " ") {
			mystruct := util.Msg{false, "密码错误"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		if err == orm.ErrNoRows {
			fmt.Println("查询不到")
		} else if err == orm.ErrMissPK {
			fmt.Println("找不到主键")
		} else {
			mystruct := util.MsgLogin{true, "登录成功", user}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			this.SetSession("userID", user.Id)
		}

	}
	this.TplName = "home/login.html"
}

// 注册
func (this *HomeController) Register() {
	userID := this.GetSession("userID")
	if userID != nil {
		this.Ctx.WriteString("<script>alert('已登录');window.location.href='/';</script>")
	}
	if this.Ctx.Request.Method == "POST" {

		username := this.GetString("username")
		password := this.GetString("password")
		password_c := this.GetString("password_c")

		o := orm.NewOrm()
		o.Using("default") // 默认使用 default，你可以指定为其他数据库
		user := models.User{Name: username}

		err := o.Read(&user, "Name")

		err, state, src := this.UploadHead()
		if state {
			if err != nil {
				mystruct := util.MsgHead{false, "", err}
				this.Data["json"] = &mystruct
				this.ServeJSON()
			}
			user.Headimgurl = src

		}

		if username == "" {
			mystruct := util.Msg{false, "用户名不得为空"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		if password == "" {
			mystruct := util.Msg{false, "密码不得为空"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		if password != password_c {
			mystruct := util.Msg{false, "密码不一致"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		if user.Password != "" {
			mystruct := util.Msg{false, "用户名已存在"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

		user.Name = username
		user.Password = util.Md5(password)
		user.Addtime = int(time.Now().Unix())

		_, err = o.Insert(&user)
		if err == nil {
			mystruct := util.Msg{true, "注册成功"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
		} else {
			mystruct := util.Msg{false, "注册失败"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

	}
	this.TplName = "home/register.html"
}

// 修改
func (this *HomeController) Editing() {
	userID := this.GetSession("userID")
	if userID == nil {
		this.Ctx.WriteString("<script>alert('未登录');window.location.href='login';</script>")
	}
	data := this.GetUserByID(userID.(int))

	if this.Ctx.Request.Method == "POST" {
		o := orm.NewOrm()
		o.Using("default") // 默认使用 default，你可以指定为其他数据库

		password := this.GetString("password")
		password_c := this.GetString("password_c")

		if password != "" {
			if password != password_c {
				mystruct := util.Msg{false, "密码不一致"}
				this.Data["json"] = &mystruct
				this.ServeJSON()
			}
			data.Password = util.Md5(password)
		}

		err, state, src := this.UploadHead()
		if state {
			if err != nil {
				mystruct := util.MsgHead{false, "", err}
				this.Data["json"] = &mystruct
				this.ServeJSON()
			}
			data.Headimgurl = src

		}

		_, err = o.Update(&data)
		if err == nil {
			this.DelSession("user")
			this.SetSession("user", data)

			mystruct := util.Msg{true, "修改成功"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
		} else {
			mystruct := util.Msg{false, "修改失败"}
			this.Data["json"] = &mystruct
			this.ServeJSON()
			return
		}

	}
	this.TplName = "home/edit.html"
	this.Data["user"] = data
}

// 上传头像
func (this *HomeController) UploadHead() (err error, state bool, src string) {
	f, h, err := this.GetFile("headimgurl")
	if f != nil {
		if err != nil {
			// 文件上传错误
			state = false
			return
		}
		exStrArr := strings.Split(h.Filename, ".")
		exStr := strings.ToLower(exStrArr[len(exStrArr)-1])
		if exStr != "jpg" && exStr != "png" && exStr != "gif" {
			// "上传只能.jpg 或者png格式"
			state = false
			return
		}

		src = "static/upload/headimg/" + strconv.FormatInt(time.Now().Unix(), 10) + h.Filename
		defer f.Close()
		err = this.SaveToFile("headimgurl", src) // 保存位置在 static/upload, 没有文件夹要先创建
		if err != nil {
			// 文件上传保存错误
			state = false
			return
		} else {
			state = true
		}
	} else {
		// 未上传文件
		state = false
	}
	return
}

// 退出登录
func (this *HomeController) Logout() {
	userID := this.GetSession("userID")
	if userID != nil {
		// user is set and can be deleted
		data := this.GetUserByID(userID.(int))
		this.DelSession("userID")
		// 退出聊天室
		Leave(data.Name)
	}
	this.Ctx.WriteString("<script>alert('退出成功');window.location.href='login';</script>")

}

// 用户聊天
func (this *HomeController) Post() {

	content := this.GetString("content")
	userID := this.GetSession("userID")
	if userID == nil {
		mystruct := util.Msg{false, "用户未登录"}
		this.Data["json"] = &mystruct
		this.ServeJSON()
		return
	}

	if len(content) == 0 {
		mystruct := util.Msg{false, "消息不得为空"}
		this.Data["json"] = &mystruct
		this.ServeJSON()
		return
	}
	data := this.GetUserByID(userID.(int))

	// 把聊天事件填入发布消息管道
	publish <- newEvent(models.EVENT_MESSAGE, data.Id, data.Name, content, data.Headimgurl)
	mystruct := util.Msg{true, ""}
	this.Data["json"] = &mystruct
	this.ServeJSON()
	return
}

// 获取事件
func (this *HomeController) Fetch() {
	lastReceived, err := this.GetInt("lastReceived")
	if err != nil {
		mystruct := util.Msg{true, "参数错误"}
		this.Data["json"] = &mystruct
		this.ServeJSON()
		return
	}

    // 获取事件
	events := models.GetEvents(int(lastReceived))
    // 如果获取到事件，即可返回
	if len(events) > 0 {
		this.Data["json"] = events
		this.ServeJSON()
		return
	}

    // 如果获取事件为空，即等待消息，此处造成阻塞关键点
	// Wait for new message(s).
	ch := make(chan bool)
	waitingList.PushBack(ch)
	<-ch

    // 一有事件即刻返回
	this.Data["json"] = models.GetEvents(int(lastReceived))
	this.ServeJSON()
	return
}
