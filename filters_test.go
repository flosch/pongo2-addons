package pongo2addons

import (
	"time"
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

	// Filesizeformat
	c.Assert(pongo2.RenderTemplateString("{{ 123456789|filesizeformat }}", nil), Equals, "118MiB")

	// Timesince/timeuntil
	base_date := time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	future_date := base_date.Add(24 * 7 * 4 * time.Hour + 2 * time.Hour)
	c.Assert(pongo2.RenderTemplateString("{{ future_date|timeuntil:base_date }}",
			pongo2.Context{"base_date": base_date, "future_date": future_date}), Equals, "4 weeks from now")

	base_date = time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	future_date = base_date.Add(2 * time.Hour)
	c.Assert(pongo2.RenderTemplateString("{{ future_date|timeuntil:base_date }}",
			pongo2.Context{"base_date": base_date, "future_date": future_date}), Equals, "2 hours from now")

	base_date = time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	future_date = base_date.Add(2 * time.Hour)
	c.Assert(pongo2.RenderTemplateString("{{ base_date|timesince:future_date }}",
			pongo2.Context{"base_date": base_date, "future_date": future_date}), Equals, "2 hours ago")
}
