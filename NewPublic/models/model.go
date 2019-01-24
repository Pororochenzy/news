package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	Name string
	Pwd  string
	Articles []*Article `orm:"rel(m2m)"` //设置多对多关系
}

//beego规定当没有设置主键时，以字段名为Id，类型为int的字段当默认主键
type Article struct {
	Id2 int  `orm:"pk;auto"`
	Title string `orm:"size(40)"`
	Content string `orm:"size(100)"`
	ReadCount int 	`orm:"default(0)"`
	Time time.Time	`orm:"type(datetime);auto_now_add"`
	Img string	`orm:"null"`
	//
	ArticleType *ArticleType  `orm:"rel(fk)"` //设置一对多关系
	Users []*User `orm:"reverse(many)"`  //设置多对多的反向关系
}

type ArticleType struct {
	Id int
	TypeName string
	Articles []*Article `orm:"reverse(many)"`//设置一对多的反向关系
}

func init (){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsPublic?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}