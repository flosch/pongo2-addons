package pongo2addons

import (
	"testing"
	"time"

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
	c.Assert(pongo2.RenderTemplateString("{{ \"**test**\"|markdown }}", nil), Equals, "<p><strong>test</strong></p>\n")

	// Slugify
	c.Assert(pongo2.RenderTemplateString("{{ \"this is ä test!\"|slugify }}", nil), Equals, "this-is-a-test")

	// Filesizeformat
	c.Assert(pongo2.RenderTemplateString("{{ 123456789|filesizeformat }}", nil), Equals, "118MiB")

	// Timesince/timeuntil
	base_date := time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	future_date := base_date.Add(24*7*4*time.Hour + 2*time.Hour)
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

	// Natural time
	base_date = time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	future_date = base_date.Add(4 * time.Second)
	c.Assert(pongo2.RenderTemplateString("{{ base_date|naturaltime:future_date }}",
		pongo2.Context{"base_date": base_date, "future_date": future_date}), Equals, "4 seconds ago")

	// Naturalday
	today := time.Date(2014, time.February, 1, 8, 30, 00, 00, time.UTC)
	yesterday := today.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	today_plus_3 := today.Add(3 * 24 * time.Hour)
	c.Assert(pongo2.RenderTemplateString("{{ date|naturalday:today }}",
		pongo2.Context{"date": today, "today": today}), Equals, "today")
	c.Assert(pongo2.RenderTemplateString("{{ date|naturalday:today }}",
		pongo2.Context{"date": yesterday, "today": today}), Equals, "yesterday")
	c.Assert(pongo2.RenderTemplateString("{{ date|naturalday:today }}",
		pongo2.Context{"date": tomorrow, "today": today}), Equals, "tomorrow")
	c.Assert(pongo2.RenderTemplateString("{{ date|naturalday:today }}",
		pongo2.Context{"date": today_plus_3, "today": today}), Equals, "3 days from now")

	// Intcomma
	c.Assert(pongo2.RenderTemplateString("{{ 123456789|intcomma }}", nil), Equals, "123,456,789")

	// Ordinal
	c.Assert(pongo2.RenderTemplateString("{{ 1|ordinal }} {{ 2|ordinal }} {{ 3|ordinal }} {{ 18241|ordinal }}", nil),
		Equals, "1st 2nd 3rd 18241st")

	// Truncatesentences
	c.Assert(pongo2.RenderTemplateString("{{ text|truncatesentences:3|safe }}", pongo2.Context{
		"text": `This is a first sentence with a 4.50 number. The second one is even more fun! Isn't it? Last sentence, okay.`}),
		Equals, "This is a first sentence with a 4.50 number. The second one is even more fun! Isn't it?")

	// Truncatesentences_html
	c.Assert(pongo2.RenderTemplateString("{{ text|truncatesentences_html:2|safe }}", pongo2.Context{
		"text": `<div class="test"><ul><li>This is a first sentence with a 4.50 number.</li><li>The second one is even more fun! Isn't it?</li><li>Last sentence, okay.</li></ul></div>`}),
		Equals, `<div class="test"><ul><li>This is a first sentence with a 4.50 number.</li><li>The second one is even more fun!</li></ul></div>`)
	c.Assert(pongo2.RenderTemplateString("{{ text|truncatesentences_html:3|safe }}", pongo2.Context{
		"text": `<div class="test"><ul><li>This is a first sentence with a 4.50 number.</li><li>The second one is even more fun! Isn't it?</li><li>Last sentence, okay.</li></ul></div>`}),
		Equals, `<div class="test"><ul><li>This is a first sentence with a 4.50 number.</li><li>The second one is even more fun! Isn't it?</li></ul></div>`)
}
