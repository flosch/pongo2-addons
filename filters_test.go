package pongo2addons

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/flosch/pongo2"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type TestSuite1 struct{}

var _ = Suite(&TestSuite1{})

func (s *TestSuite1) TestFilters(c *C) {
	// Markdown
	c.Assert(pongo2.RenderTemplateString("{{ \"**test**\"|markdown|safe }}", nil), Equals, "<p><strong>test</strong></p>\n")

	// Slugify
	c.Assert(pongo2.RenderTemplateString("{{ \"this is ä test!\"|slugify }}", nil), Equals, "this-is-a-test")
}
