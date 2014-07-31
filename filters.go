package pongo2addons

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/flosch/pongo2"

	"github.com/extemporalgenome/slug"
	"github.com/flosch/go-humanize"
	"github.com/russross/blackfriday"
)

func init() {
	// Regulars
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("filesizeformat", filterFilesizeformat)
	pongo2.RegisterFilter("truncatesentences", filterTruncatesentences)
	pongo2.RegisterFilter("truncatesentences_html", filterTruncatesentencesHtml)

	// Markup
	pongo2.RegisterFilter("markdown", filterMarkdown)

	// Humanize
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

var filterTruncatesentencesRe = regexp.MustCompile(`(?U:.*[\w]{3,}.*([\d][\.!?][\D]|[\D][\.!?][\s]|[\n$]))`)

func filterTruncatesentences(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	count := param.Integer()
	if count <= 0 {
		return pongo2.AsValue(""), nil
	}
	sentencens := filterTruncatesentencesRe.FindAllString(strings.TrimSpace(in.String()), -1)
	return pongo2.AsValue(strings.TrimSpace(strings.Join(sentencens[:min(count, len(sentencens))], ""))), nil
}

// Taken from pongo2/filters_builtin.go
func filterTruncateHtmlHelper(value string, new_output *bytes.Buffer, cond func() bool, fn func(c rune, s int, idx int) int, finalize func()) {
	vLen := len(value)
	tag_stack := make([]string, 0)
	idx := 0

	for idx < vLen && !cond() {
		c, s := utf8.DecodeRuneInString(value[idx:])
		if c == utf8.RuneError {
			idx += s
			continue
		}

		if c == '<' {
			new_output.WriteRune(c)
			idx += s // consume "<"

			if idx+1 < vLen {
				if value[idx] == '/' {
					// Close tag

					new_output.WriteString("/")

					tag := ""
					idx += 1 // consume "/"

					for idx < vLen {
						c2, size2 := utf8.DecodeRuneInString(value[idx:])
						if c2 == utf8.RuneError {
							idx += size2
							continue
						}

						// End of tag found
						if c2 == '>' {
							idx++ // consume ">"
							break
						}
						tag += string(c2)
						idx += size2
					}

					if len(tag_stack) > 0 {
						// Ideally, the close tag is TOP of tag stack
						// In malformed HTML, it must not be, so iterate through the stack and remove the tag
						for i := len(tag_stack) - 1; i >= 0; i-- {
							if tag_stack[i] == tag {
								// Found the tag
								tag_stack[i] = tag_stack[len(tag_stack)-1]
								tag_stack = tag_stack[:len(tag_stack)-1]
								break
							}
						}
					}

					new_output.WriteString(tag)
					new_output.WriteString(">")
				} else {
					// Open tag

					tag := ""

					params := false
					for idx < vLen {
						c2, size2 := utf8.DecodeRuneInString(value[idx:])
						if c2 == utf8.RuneError {
							idx += size2
							continue
						}

						new_output.WriteRune(c2)

						// End of tag found
						if c2 == '>' {
							idx++ // consume ">"
							break
						}

						if !params {
							if c2 == ' ' {
								params = true
							} else {
								tag += string(c2)
							}
						}

						idx += size2
					}

					// Add tag to stack
					tag_stack = append(tag_stack, tag)
				}
			}
		} else {
			idx = fn(c, s, idx)
		}
	}

	finalize()

	for i := len(tag_stack) - 1; i >= 0; i-- {
		tag := tag_stack[i]
		// Close everything from the regular tag stack
		new_output.WriteString(fmt.Sprintf("</%s>", tag))
	}
}

func filterTruncatesentencesHtml(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, error) {
	count := param.Integer()
	if count <= 0 {
		return pongo2.AsValue(""), nil
	}

	value := in.String()
	newLen := max(param.Integer(), 0)

	new_output := bytes.NewBuffer(nil)

	sentencefilter := 0

	filterTruncateHtmlHelper(value, new_output, func() bool {
		return sentencefilter >= newLen
	}, func(_ rune, _ int, idx int) int {
		// Get next word
		word_found := false

		for idx < len(value) {
			c2, size2 := utf8.DecodeRuneInString(value[idx:])
			if c2 == utf8.RuneError {
				idx += size2
				continue
			}

			if c2 == '<' {
				// HTML tag start, don't consume it
				return idx
			}

			new_output.WriteRune(c2)
			idx += size2

			if (c2 == '.' && !(idx+1 < len(value) && value[idx+1] >= '0' && value[idx+1] <= '9')) ||
				c2 == '!' || c2 == '?' || c2 == '\n' {
				// Sentence ends here, stop capturing it now
				break
			} else {
				word_found = true
			}
		}

		if word_found {
			sentencefilter++
		}

		return idx
	}, func() {})

	return pongo2.AsValue(new_output.String()), nil
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
