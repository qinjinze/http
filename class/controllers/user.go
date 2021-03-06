package controllers

import (
	"class/models"
	"encoding/base64"

	//	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandlePost() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//beego.Info(userName,pwd)
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整，请重新注册"
		beego.Info("注册数据不完整，请重新注册")
		this.TplName = "register.html"
		return
	}
	//3.操作数据
	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	var user models.User
	//给插入对象赋值
	user.Name = userName
	user.PassWord = pwd
	//插入
	o.Insert(&user)
	//返回结果

	//4.返回页面
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login", 302)
	//this.TplName = "login.html"
}

//展示登陆页面
func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	data, _ := base64.StdEncoding.DecodeString(userName)
	beego.Info(userName)
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = string(data)
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}

//输入登陆账号登录
func (this *UserController) HandleLogin() {
	//1.获取数据
	fmt.Println("================")
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	fmt.Println("userName=", userName, "pwd:", pwd)
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整"
		this.TplName = "login.html"
		return
	}

	//3.操作数据
	//1。获取ORM对象
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		fmt.Println("err:", err)
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}
	fmt.Println("登陆成功")
	//4.返回页面
	//this.Ctx.WriteString("登录成功")
	data := this.GetString("remember")
	fmt.Println("data==", data)
	beego.Info(data)
	if data == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName", temp, 100)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.Data["userName"] = userName
	fmt.Println("userName=", userName)
	this.SetSession("remember", userName)
	Name := this.GetSession("remember")
	fmt.Println("Name=", Name)

	this.Redirect("/article/showArticleList", 302)
}

//退出登录
func (this *UserController) Logout() {
	//删除session
	this.DelSession("userName")
	//跳转登录页面
	this.Redirect("/login", 302)
}
