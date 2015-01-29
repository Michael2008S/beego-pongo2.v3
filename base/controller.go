package base

import (
	"github.com/astaxie/beego"
	renderer "github.com/danegigi/beego-pongo2.v3"
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) Render(tplName string, tplData renderer.Context) error {
	return renderer.Render(b.Ctx, tplName, tplData)
}

func (b *BaseController) RenderString(tplName string, tplData renderer.Context) string {
	toReturn, err := renderer.RenderString(tplName, tplData)
	if err != nil {
		panic(err)
	}
	return toReturn
}
