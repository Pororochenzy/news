package controllers


import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"NewPublic/models"
	"math"
)

type ArticleController struct {
	beego.Controller
}

func(this*ArticleController)ShowIndex(){
	pageIndex,err:=this.GetInt("pageIndex")
		if err !=nil {
			pageIndex=1
		}
    typeName:= this.GetString("select")

	//查找文章
	o:=orm.NewOrm()
	//获取所有文章    select * from article;  queryseter
	qs:=o.QueryTable("Article")
	//分页查找
	//获取总记录数和总页数
	count,_:=qs.RelatedSel("ArticleType").Count() //总记录数
	pageSize := 2 // 一页显示多岁条
	pageCount:=math.Ceil(float64(count)/float64(pageSize))

	var articles []models.Article
	//limit(pageSize,start) -->pageSize获得几条数据,start从第几条开始
	_,err=qs.Limit(pageSize,(pageIndex-1)*pageSize).RelatedSel("ArticleType").All(&articles)

	if err!=nil {
		beego.Info("find_err:",err)
	}

	 this.Data["articles"]=articles
	 this.Data["pageCount"]=pageCount
	 this.Data["count"]=count
	this.Data["pageIndex"]=pageIndex

	//指定视图
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAdd(){
	//查询文章类型
	o2:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o2.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	//指定视图
	this.TplName = "add.html"
}

func (this*ArticleController)HandleAdd(){



	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	file,head,err:=this.GetFile("uploadname")

	//获取数据
	if articleName==""||content=="" ||err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return
	}
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return
	}

	defer file.Close()

	beego.Info(head.Filename)
	ext:=path.Ext(head.Filename)
	if ext != ".jpg"&&ext != ".JPG"  && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return
	}
	//防止重名
	beego.Info(time.Now().Format("2006-01-02 15:04:05"))
	fileName := time.Now().Format("20060102150405")

	this.SaveToFile("uploadname","static/img/"+fileName+ext)




	//把数据插入到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext

	//获得新闻类型
	typeName:=this.GetString("select")
	var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Read(&articleType,"TypeName")  //先根据类型名字查出来 类型表的信息的所有信息，然后放回对象中
	article.ArticleType=&articleType   //注意，artcle的articleType成员是个对象指针

	//
	//插入
	o.Insert(&article)

	//返回数据
	this.Redirect("/index",302)
}
//查看文章内容
func (this*ArticleController)ShowArticle(){
	id,err:=this.GetInt("id")
	if err!=nil {
		beego.Error("id error")
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id2=id

	o.QueryTable(&article)
	//阅读次数加一

	//更新操作
	article.ReadCount += 1
	_,err=o.Update(&article)
	if err !=nil {
		beego.Error("insert fail",err)//insert fail Error 1048: Column 'time' cannot be null
	}

	//返回数据
	this.Data["article"] = article

	//指定视图
	this.TplName = "content.html"
}

//展示编辑页面
func (this*ArticleController)ShowEdit(){
	id,err:=this.GetInt("id")
	//数据校验
	if err != nil{
		beego.Error("获取数据错误",err)
		this.TplName = "index.html"
		return
	}
	//处理数据
	//查询
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//指定查询条件
	article.Id2 = id
	//查询
	o.Read(&article)
	this.TplName="update.html"

	//返回数据
	this.Data["article"] = article

	//指定视图
	this.TplName = "update.html"
}
//编辑文章内容
func (this*ArticleController)HandleEdit(){
	id,err:=this.GetInt("id")
	if err !=nil {
		beego.Error("id wrong:",err)
		return
	}
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filePath:=UploadFunc(this,"uploadname")
	//校验数据
	if err != nil || articleName == "" || content == "" || filePath == ""{
		beego.Error("获取数据错误")
		this.TplName = "update.html"
		return
	}



	//更新
	//获取orm对象
	o:=orm.NewOrm()
	//获取更新对象
	var article models.Article
	//给更新条件赋值
	article.Id2=id
	//先read一下，判断要更新的数据
	err =o.Read(&article)
	if err != nil{
		beego.Error("更新数据不存在")
		this.TplName = "update.html"
		return
	}
	//更新
	article.Title = articleName
	article.Content = content
	article.Img = filePath
	o.Update(&article)

	this.Redirect("index.html",302)
}
//图片处理
func  UploadFunc(this*ArticleController,fileName string) string{

	file,head,err:=this.GetFile(fileName)
	if err!=nil {
		beego.Error("图片上传失败",err)
		this.TplName = "add.html"
		return ""
	}
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return ""
	}

	defer file.Close()
	beego.Info(head.Filename)
	ext:=path.Ext(head.Filename)
	if ext != ".jpg"&&ext != ".JPG"  && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return ""
	}
	//防止重名
	beego.Info(time.Now().Format("2006-01-02 15:04:05"))
	filePath := time.Now().Format("20060102150405")

	this.SaveToFile("uploadname","static/img/"+filePath+ext)
	return "/static/img/"+filePath+ext
}

//删除内容
func (this*ArticleController)HandleDelete(){
	id,err:=this.GetInt("id")
	if err!=nil {
		beego.Error("删除请求数据错误")
		this.TplName="index.html"
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id2=id
	_,err=o.Delete(&article)
	if err!=nil {
		beego.Error("删除失败")
		this.TplName="index.html"
		return
	}

	//返回数据
	this.Redirect("/index",302)
}

//添加文章类型
func (this*ArticleController)ShowAddType(){
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("article_type").All(&articleTypes)

	this.Data["articleTypes"]=articleTypes
	this.TplName="addType.html"
}
func (this*ArticleController)HandleAddType(){
	typeName:=this.GetString("typeName")
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	_,err:=o.Insert(&articleType)
	if err!=nil {
		beego.Error(err)
		return
	}
	this.Redirect("/addType",302)

}