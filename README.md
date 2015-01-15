beego-pongo2.v2
============

A tiny little helper for using Pongo2 (v3) with Beego.

Status: With little tests. IF YOU FIND BUGS, PLS LET ME KNOW.

Documentation: http://godoc.org/github.com/danegigi/beego-pongo2.v3

Based on [https://github.com/ipfans/beego-pongo2.v2](https://github.com/ipfans/beego-pongo2.v2)

## Usage

```go
package controllers

import (
    "github.com/astaxie/beego"
    "github.com/danegigi/beego-pongo2.v3"
)

type MainController struct {
    beego.Controller
}

func (this *MainController) Page1() {
    err := pongo2.Render(this.Ctx, "page.html", pongo2.Context{
        "ints": []int{1, 2, 3, 4, 5},
    })

    if err != nil{
      this.Abort("500")
    }
}

func (this *MainController) Page2(){
  str, err := pongo2.RenderString("page2.html", pongo2.Context{
    "name": "My Name"
  })

  if err != nil{
    this.Abort("500")
  }

  this.Ctx.WriteString(str)
}
```