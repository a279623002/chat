package models

type User struct {
	Id         int    `orm:"column(user_id);auto"`
	Name       string `orm:"column(name);size(20)"`
	Password   string `orm:"column(password)"`
	Headimgurl string `orm:"column(headimgurl)"`
	Addtime    int    `orm:"column(addtime)"`
	Status     int8   `orm:"column(status)"`
}
