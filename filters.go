package pongo2addons

import (
	"errors"
	"time"

	"github.com/flosch/pongo2"

	"github.com/extemporalgenome/slug"
	"github.com/flosch/go-humanize"
	"github.com/russross/blackfriday"
)

func init() {
	pongo2.RegisterFilter("markdown", filterMarkdown)
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("filesizeformat", filterFilesizeformat)
	pongo2.RegisterFilter("timeuntil", filterTimeuntilTimesince)
	pongo2.RegisterFilter("timesince", filterTimeuntilTimesince)
	pongo2.RegisterFilter("naturaltime", filterTimeuntilTimesince)
	pongo2.RegisterFilter("naturalday", filterNaturalday)
	pongo2.RegisterFilter("intcomma", filterIntcomma)
	pongo2.RegisterFilter("ordinal", filterOrdinal)
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
		return nil, errors.New("time-value is not a time.Time-instance.")
	}
	var paramtime time.Time
	if !param.IsNil() {
		paramtime, is_time = param.Interface().(time.Time)
		if !is_time {
			return nil, errors.New("time-parameter is not a time.Time-instance.")
		}
	} else {
		paramtime = time.Now()
	}

	return pongo2.AsValue(humanize.TimeDuration(basetime.Sub(paramtime))), nil
}

func filterIntcomma(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(humanize.Comma(int64(in.Integer()))), nil
}

func filterOrdinal(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	return pongo2.AsValue(humanize.Ordinal(in.Integer())), nil
}

func filterNaturalday(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	basetime, is_time := in.Interface().(time.Time)
	if !is_time {
		return nil, errors.New("naturalday-value is not a time.Time-instance.")
	}

	var reference_time time.Time
	if !param.IsNil() {
		reference_time, is_time = param.Interface().(time.Time)
		if !is_time {
			return nil, errors.New("naturalday-parameter is not a time.Time-instance.")
		}
	} else {
		reference_time = time.Now()
	}

	d := reference_time.Sub(basetime) / time.Hour

	switch {
	case d >= 0 && d < 24:
		// Today
		return pongo2.AsValue("today"), nil
	case d >= 24:
		return pongo2.AsValue("yesterday"), nil
	case d < 0 && d >= -24:
		return pongo2.AsValue("tomorrow"), nil
	}

	// Default behaviour
	return pongo2.ApplyFilter("naturaltime", in, param)
}
