package sqlgen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type PlaceholderFormat interface {
	ReplacePlaceholders(sql string) (string, error)
}

var (
	Question = questionFormat{}
	Dollar   = dollarFormat{}
)

type questionFormat struct{}

func (questionFormat) ReplacePlaceholders(sql string) (string, error) {
	return sql, nil
}

type dollarFormat struct{}

func (dollarFormat) ReplacePlaceholders(sql string) (string, error) {
	return replacePlaceholders(sql, func(w io.Writer, i int) error {
		io.WriteString(w, fmt.Sprintf("$%d", i))
		return nil
	})
}

func replacePlaceholders(sql string, replace func(w io.Writer, i int) error) (string, error) {
	var i int
	buf := &bytes.Buffer{}
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		i++
		buf.WriteString(sql[:p])
		if err := replace(buf, i); err != nil {
			return "", err
		}
		sql = sql[p+1:]
	}
	buf.WriteString(sql)
	return buf.String(), nil
}

func Placeholders(count int) string {
	if count < 1 {
		return ""
	}
	return strings.Repeat(",?", count)[1:]
}
