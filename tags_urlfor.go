package pongo2

import (
	"bytes"

	"github.com/astaxie/beego"
	p2 "gopkg.in/flosch/pongo2.v3"

	_ "github.com/flosch/pongo2-addons"
)

type tagURLForNode struct {
	objectEvaluators []p2.IEvaluator
}

func (node *tagURLForNode) Execute(ctx *p2.ExecutionContext, buffer *bytes.Buffer) *p2.Error {
	args := make([]string, len(node.objectEvaluators))
	for i, ev := range node.objectEvaluators {
		obj, err := ev.Evaluate(ctx)
		if err != nil {
			return err
		}
		args[i] = obj.String()
	}

	// fix : cannot use args[1:] (type []string) as type []interface {} in argument to beego.UrlFor
	//url := beego.URLFor(args[0], args[1:]...)

	length := len(args[1:])
	argsInf := make([]interface{}, length)
	for i, v := range args[1:] {
		argsInf[i] = interface{}(v)
	}
	url := beego.URLFor(args[0], argsInf...)

	buffer.WriteString(url)
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagURLForParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	evals := []p2.IEvaluator{}
	for arguments.Remaining() > 0 {
		expr, err := arguments.ParseExpression()
		evals = append(evals, expr)
		if err != nil {
			return nil, err
		}
	}

	if (len(evals)-1)%2 != 0 {
		return nil, arguments.Error("URL takes one argument for the controller and any number of optional pairs of key/value pairs.", nil)
	}

	return &tagURLForNode{evals}, nil
}

func init() {
	p2.RegisterTag("urlfor", tagURLForParser)
}
