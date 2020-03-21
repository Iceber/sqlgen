package sqlgen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type InsertBuilder struct {
	intoTable string
	columns   []string
	values    [][]interface{}
}

func (b *InsertBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(b.intoTable) == 0 {
		err = fmt.Errorf("insert statements must specify a table")
	}

	if len(b.values) == 0 {
		err = fmt.Errorf("insert statements must have at least one set of values")
	}

	sql := &bytes.Buffer{}
	sql.WriteString("INSERT ")

	sql.WriteString("INTO ")
	sql.WriteString(b.intoTable)
	sql.WriteString(" ")

	if len(b.columns) > 0 {
		sql.WriteString("(")
		sql.WriteString(strings.Join(b.columns, ","))
		sql.WriteString(")")
	}

	args, err = b.appendValuesToSQL(sql, args)
	if err != nil {
		return
	}
	return
}

func (b *InsertBuilder) appendValuesToSQL(w io.Writer, args []interface{}) ([]interface{}, error) {
	if len(b.values) == 0 {
		return args, errors.New("values for insert statements are not set")
	}

	io.WriteString(w, "VALUES ")
	valuesStrings := make([]string, len(b.values))
	for r, row := range b.values {
		valueStrings := make([]string, len(row))
		for v, val := range row {
			switch typeVal := val.(type) {
			case Sqlizer:
				valSql, valArgs, err := typeVal.ToSql()
				if err != nil {
					return nil, err
				}
				valueStrings[v] = valSql
				args = append(args, valArgs...)
			default:
				valueStrings[v] = "?"
				args = append(args, val)
			}
		}
		valuesStrings[r] = fmt.Sprintf("(%s)", strings.Join(valueStrings, ","))
	}
	io.WriteString(w, strings.Join(valuesStrings, ","))
	return args, nil
}

func (b *InsertBuilder) Into(table string) *InsertBuilder {
	b.intoTable = table
	return b
}

func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	b.columns = append(b.columns, columns...)
	return b
}

func (b *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	b.values = append(b.values, values)
	return b
}

func (b *InsertBuilder) SetMap(clauses map[string]interface{}) *InsertBuilder {
	cols := make([]string, 0, len(clauses))
	vals := make([]interface{}, 0, len(clauses))
	for col, val := range clauses {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	b.columns = cols
	b.values = [][]interface{}{vals}
	return b
}
