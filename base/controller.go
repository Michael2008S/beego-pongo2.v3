package base

import (
	"github.com/astaxie/beego"
	renderer "github.com/danegigi/beego-pongo2.v3"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Template(tplName string, tplData renderer.Context) error {
	return renderer.Render(this.Ctx, tplName, tplData)
}

func (this *BaseController) TemplateToString(tplName string, tplData renderer.Context) string {
	toReturn, err := renderer.RenderString(tplName, tplData)
	if err != nil {
		panic(err)
	}
	return toReturn
}
