// A small library that lets you use Pongo2 with Beego
//
// When Render is called, it will populate the render context with Beego's flash messages.
// You can also use {% urlfor "MyController.Action" ":key" "value" %} in your templates, and
// it'll work just like `urlfor` would with `html/template`. It takes one controller argument and
// zero or more key/value pairs to fill the URL.
//
package pongo2

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"path"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	p2 "gopkg.in/flosch/pongo2.v3"

	_ "github.com/flosch/pongo2-addons"
)

//const (
//	templateDir = "templates"
//)

type Context map[string]interface{}

var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

var devMode bool
var Pango2TemplatesPath string

// Render takes a Beego context, template name and a Context (map[string]interface{}).
// The template is parsed and cached, and gets executed into beegoCtx's ResponseWriter.
//
// Templates are looked up in `templates/` instead of Beego's default `views/` so that
// Beego doesn't attempt to load and parse our templates with `html/template`.
func Render(beegoCtx *context.Context, tmpl string, ctx Context) error {
	var template *p2.Template
	var err error
	devMode = beego.RunMode == "dev"
	if devMode {
		beego.Warn("p2.FromFile_Mode:", beego.RunMode)
		beego.Warn(path.Join(Pango2TemplatesPath, tmpl))
		template, err = p2.FromFile(path.Join(Pango2TemplatesPath, tmpl))
	} else {
		beego.Info("p2.FromCache_Mode:", beego.RunMode)
		template, err = p2.FromCache(path.Join(Pango2TemplatesPath, tmpl))
	}

	if err != nil {
		panic(err)
	}

	var pCtx p2.Context
	if ctx == nil {
		pCtx = p2.Context{}
	} else {
		pCtx = p2.Context(ctx)
	}

	if xsrf, ok := beegoCtx.GetSecureCookie(beego.XSRFKEY, "_xsrf"); ok {
		pCtx["_xsrf"] = xsrf
	}

	// Only override "flash" if it hasn't already been set in Context
	if _, ok := ctx["flash"]; !ok {
		if ctx == nil {
			ctx = Context{}
		}
		ctx["flash"] = readFlash(beegoCtx)
	}

	// FIXME 当是api的时候，直接返回页面参数
	if devMode && strings.HasPrefix(beegoCtx.Input.Request.URL.Path, "/api/") {
		var pCtx p2.Context
		if ctx == nil {
			pCtx = p2.Context{}
		} else {
			pCtx = p2.Context(ctx)
		}
		content, _ := json.MarshalIndent(pCtx, "", "  ")
		_, err := beegoCtx.ResponseWriter.Write(content)
		return err
	}

	return template.ExecuteWriter(pCtx, beegoCtx.ResponseWriter)
}

// Same as Render() but returns a string
func RenderString(tmpl string, ctx Context) (string, error) {
	var template *p2.Template
	var err error
	devMode = beego.RunMode == "dev"

	if devMode {
		beego.Warn("p2.FromFile_Mode:", beego.RunMode)
		template, err = p2.FromFile(path.Join(Pango2TemplatesPath, tmpl))
	} else {
		beego.Info("p2.FromCache_Mode:", beego.RunMode)
		template, err = p2.FromCache(path.Join(Pango2TemplatesPath, tmpl))
	}

	if err != nil {
		panic(err)
	}

	var pCtx p2.Context
	if ctx == nil {
		pCtx = p2.Context{}
	} else {
		pCtx = p2.Context(ctx)
	}

	// str, _ := template.Execute(pCtx)
	// return str
	return template.Execute(pCtx)
}

// readFlash is similar to beego.ReadFromRequest except that it takes a *context.Context instead
// of a *beego.Controller, and returns a map[string]string directly instead of a Beego.FlashData
// (which only has a Data field anyway).
func readFlash(ctx *context.Context) map[string]string {
	data := map[string]string{}
	if cookie, err := ctx.Request.Cookie(beego.FlashName); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, "\x23"+beego.FlashSeperator+"\x23")
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}
		}
		// read one time then delete it
		ctx.SetCookie(beego.FlashName, "", -1, "/")
	}
	return data
}

func init() {
	// FIXME 不需要设置,直接用beego.RunMode 判断
	//devMode = beego.AppConfig.String("runmode") == "dev"
	//beego.Error("beego-pango2.v3_run_mode:", beego.AppConfig.String("run_mode"))
	//beego.Error("Pango2TemplatesPath:", beego.AppConfig.String("Pango2TemplatesPath"))
	//devMode = beego.RunMode == "dev"
	beego.AutoRender = false
}
