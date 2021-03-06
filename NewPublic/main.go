package main

import (
	_ "NewPublic/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("prePage",ShowPrePage)
	beego.AddFuncMap("nextPage",ShowNextPage)
	beego.Run()
}
func ShowPrePage(pageIndex int) int {
	if pageIndex<=1 {
		pageIndex=1
	}
	return pageIndex -1
}
func ShowNextPage(pageIndex int,pageCount float64) int{
	if pageIndex>=int(pageCount){
		return  int(pageCount)
	}
	return pageIndex+1
}

