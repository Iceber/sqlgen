package sqlgen

import (
	"fmt"
	"io"
)

type part struct {
	pred interface{}
	args []interface{}
}

func newPart(pred interface{}, args ...interface{}) Sqlizer {
	return &part{pred, args}
}

func (p part) ToSql() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
	case Sqlizer:
		sql, args, err = pred.ToSql()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expcted string or Sqlizer, not %T", pred)
	}
	return
}

func appendToSql(w io.Writer, args []interface{}, parts []Sqlizer, sep string) ([]interface{}, error) {
	for i, part := range parts {
		partSql, partArgs, err := part.ToSql()
		if err != nil {
			return nil, err
		}
		if len(partSql) == 0 {
			continue
		}

		if i > 0 {
			if _, err := io.WriteString(w, sep); err != nil {
				return nil, err
			}
		}

		if _, err = io.WriteString(w, partSql); err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
