package pongo2addons

import (
	"github.com/flosch/pongo2"

	"github.com/extemporalgenome/slug"
	"github.com/russross/blackfriday"
)

func init() {
	pongo2.RegisterFilter("markdown", filterMarkdown)
	pongo2.RegisterFilter("slugify", filterSlugify)
}

func filterMarkdown(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(string(blackfriday.MarkdownCommon([]byte(in.String())))), nil
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(slug.Slug(in.String())), nil
}
