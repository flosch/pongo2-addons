package pongo2addons

import (
	"errors"
	"time"

	"github.com/flosch/pongo2"

	"github.com/extemporalgenome/slug"
	"github.com/russross/blackfriday"
	"github.com/flosch/go-humanize"
)

func init() {
	pongo2.RegisterFilter("markdown", filterMarkdown)
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("filesizeformat", filterFilesizeformat)
	pongo2.RegisterFilter("timeuntil", filterTimeuntilTimesince)
	pongo2.RegisterFilter("timesince", filterTimeuntilTimesince)
	pongo2.RegisterFilter("naturaltime", filterTimeuntilTimesince)
	pongo2.RegisterFilter("intcomma", filterIntcomma)
}

func filterMarkdown(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(string(blackfriday.MarkdownCommon([]byte(in.String())))), nil
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(slug.Slug(in.String())), nil
}

func filterFilesizeformat(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(humanize.IBytes(uint64(in.Integer()))), nil
}

func filterTimeuntilTimesince(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	basetime, is_time := in.Interface().(time.Time)
	if !is_time {
		return nil, errors.New("timeuntil-value is not a time.Time-instance.")
	}
	var paramtime time.Time
	if !param.IsNil() {
		paramtime, is_time = param.Interface().(time.Time)
		if !is_time {
			return nil, errors.New("timeuntil-parameter is not a time.Time-instance.")
		}
	} else {
		paramtime = time.Now()
	}

	return pongo2.AsValue(humanize.TimeDuration(basetime.Sub(paramtime))), nil
}

func filterIntcomma(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(humanize.Comma(int64(in.Integer()))), nil
}
