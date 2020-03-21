package sqlgen

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type SelectBuilder struct {
	distinct   bool
	columns    []Sqlizer
	from       []string
	whereParts []Sqlizer
	groupBys   []string
	orderBys   []string

	limit  uint64
	offset uint64
}

type CountBuilder struct {
	*SelectBuilder
}

func (b *CountBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}
	sql.WriteString("SELECT COUNT(*) ")
	sql.WriteString(" FROM ")
	sql.WriteString(strings.Join(b.from, ", "))
	if len(b.whereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(sql, args, b.whereParts, ", ")
	}
	if len(b.groupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupBys, ", "))
	}
	return
}

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

func (b *SelectBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(b.columns) == 0 {
		err = fmt.Errorf("select statements must have at least one result columns")
		return
	}
	sql := &bytes.Buffer{}

	sql.WriteString("SELECT ")
	if b.distinct {
		sql.WriteString("DISTINCT ")
	}

	args, err = appendToSql(sql, args, b.columns, ", ")
	if err != nil {
		return
	}

	sql.WriteString(" FROM ")
	sql.WriteString(strings.Join(b.from, ", "))

	if len(b.whereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(sql, args, b.whereParts, ", ")
	}
	if len(b.groupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupBys, ", "))
	}

	if len(b.orderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderBys, ", "))
	}

	if b.limit >= 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.FormatUint(b.limit, 10))
	}

	if b.offset >= 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.FormatUint(b.offset, 10))
	}
	return
}

func (b *SelectBuilder) Count() CountBuilder {
	return CountBuilder{b}
}

func (b *SelectBuilder) Columns(columns ...string) *SelectBuilder {
	for _, col := range columns {
		b.columns = append(b.columns, newPart(col))
	}
	return b
}

func (b *SelectBuilder) Column(column interface{}, args ...interface{}) *SelectBuilder {
	b.columns = append(b.columns, newPart(column, args...))
	return b
}

func (b *SelectBuilder) From(tables ...string) *SelectBuilder {
	b.from = tables
	return b
}

func (b *SelectBuilder) Where(pred interface{}, args ...interface{}) *SelectBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

func (b *SelectBuilder) GroupBy(groupBys ...string) *SelectBuilder {
	b.groupBys = append(b.groupBys, groupBys...)
	return b
}

func (b *SelectBuilder) OrderBy(orderBys ...string) *SelectBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

func (b *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	b.limit = limit
	return b
}

func (b *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	b.offset = offset
	return b
}
