package routers

import (
	"NewPublic/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/index", &controllers.ArticleController{},"get:ShowIndex")
	beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/add",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
	beego.Router("/article",&controllers.ArticleController{},"get:ShowArticle")
	beego.Router("/edit",&controllers.ArticleController{},"get:ShowEdit;post:HandleEdit")
	beego.Router("/delete",&controllers.ArticleController{},"get:HandleDelete")
	//添加文章类型
	beego.Router("/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
}
