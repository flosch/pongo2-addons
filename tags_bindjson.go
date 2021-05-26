package pongo2addons

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/flosch/pongo2"
)

type tagBindJSONNode struct {
	file    string
	key     string
	wrapper *pongo2.NodeWrapper
}

func tagBindJSONParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	// verify that there are exactly 3 arguments
	if arguments.Count() != 3 {
		syntax := "{% bindjson \"filename\" as variable %} ... {% endbindjson %}"
		errsyntax := "Tag 'bindjson' uses the following syntax: " + syntax
		return nil, arguments.Error(errsyntax, nil)
	}

	// make sure there are no arguments on the end tag
	wrapper, endargs, err := doc.WrapUntilTag("endbindjson")
	if err != nil {
		return nil, err
	}
	if endargs.Count() > 0 {
		return nil, endargs.Error("Tag 'bindjson' has no end arguments", nil)
	}

	// parse filename
	relpath := arguments.MatchType(pongo2.TokenString)
	if relpath == nil {
		return nil, arguments.Error("Expected filename string.", nil)
	}
	// TODO use sandboxed loader
	filename := filepath.Join(filepath.Dir(start.Filename), relpath.Val)

	// parse "as" keyword
	if arguments.Match(pongo2.TokenKeyword, "as") == nil {
		return nil, arguments.Error("Expected 'as' keyword.", nil)
	}

	// parse variable identifier
	key := arguments.MatchType(pongo2.TokenIdentifier)
	if key == nil {
		return nil, arguments.Error("Expected variable identifier.", nil)
	}

	// create the Node
	node := &tagBindJSONNode{}
	node.file = filename
	node.key = key.Val
	node.wrapper = wrapper

	return node, nil
}

func (node *tagBindJSONNode) Execute(ctx *pongo2.ExecutionContext, w pongo2.TemplateWriter) *pongo2.Error {
	// read the json file
	// TODO make use of the sandboxed loader
	raw, err := ioutil.ReadFile(node.file)
	if err != nil {
		return ctx.Error(err.Error(), nil)
	}

	// parse the json
	var jsonObject interface{}
	err = json.Unmarshal(raw, &jsonObject)
	if err != nil {
		errmsg := "Invalid JSON in " + node.file + ": " + err.Error()
		return ctx.Error(errmsg, nil)
	}

	// create local context
	localContext := pongo2.NewChildExecutionContext(ctx)

	// bind to key in context
	localContext.Private.Update(pongo2.Context{
		node.key: jsonObject,
	})

	// render
	return node.wrapper.Execute(localContext, w)
}

func init() {
	pongo2.RegisterTag("bindjson", tagBindJSONParser)
}
