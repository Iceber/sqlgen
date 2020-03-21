package sqlgen

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type setClause struct {
	column string
	value  interface{}
}

func (set setClause) ToSql() (sqlStr string, args []interface{}, err error) {
	var valSql string
	switch typeVal := set.value.(type) {
	case Sqlizer:
		var valArgs []interface{}
		valSql, valArgs, err = typeVal.ToSql()
		if err != nil {
			return
		}
		args = append(args, valArgs...)
	default:
		valSql = "?"
		args = append(args, set.value)
	}
	sqlStr = fmt.Sprintf("%s=%s", set.column, valSql)
	return
}

type UpdateBuilder struct {
	table string

	setClauses []setClause
	where      []Sqlizer
	orderBy    []string

	limit  uint64
	offset uint64
}

func (b *UpdateBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(b.table) == 0 {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(b.setClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set caluse")
		return
	}

	sql := &bytes.Buffer{}

	sql.WriteString("UPDATE ")
	sql.WriteString(b.table)
	sql.WriteString(" SET ")

	setSqls := make([]string, len(b.setClauses))
	for i, setClause := range b.setClauses {
		var valSql string
		switch typeVal := setClause.value.(type) {
		case Sqlizer:
			var valArgs []interface{}
			valSql, valArgs, err = typeVal.ToSql()
			if err != nil {
				return
			}
			args = append(args, valArgs...)
		default:
			valSql = "?"
			args = append(args, typeVal)
		}
		setSqls[i] = fmt.Sprintf("%s=%s", setClause.column, valSql)
	}
	sql.WriteString(strings.Join(setSqls, ", "))

	if len(b.where) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(sql, args, b.where, " AND ")
		if err != nil {
			return
		}
	}

	if len(b.orderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.limit > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.FormatUint(b.limit, 10))
	}

	if b.offset > 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.FormatUint(b.limit, 10))
	}

	return
}

func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.table = table
	return b
}

func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	b.setClauses = append(b.setClauses, setClause{column: column, value: value})
	return b
}

func (b *UpdateBuilder) SetMap(clauses map[string]interface{}) *UpdateBuilder {
	keys := make([]string, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	for _, key := range keys {
		val, _ := clauses[key]
		b.Set(key, val)
	}
	return b
}

func (b *UpdateBuilder) Where(pred interface{}, args ...interface{}) *UpdateBuilder {
	b.where = append(b.where, newWherePart(pred, args...))
	return b
}

func (b *UpdateBuilder) OrderBy(orderBy ...string) *UpdateBuilder {
	b.orderBy = append(b.orderBy, orderBy...)
	return b
}

func (b *UpdateBuilder) Offset(offset uint64) *UpdateBuilder {
	b.offset = offset
	return b
}

func (b *UpdateBuilder) Limit(limit uint64) *UpdateBuilder {
	b.limit = limit
	return b
}
