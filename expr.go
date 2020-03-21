package sqlgen

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type aliasExpr struct {
	expr  Sqlizer
	alias string
}

func Alias(expr Sqlizer, alias string) aliasExpr {
	return aliasExpr{expr, alias}
}

func (e aliasExpr) ToSql() (sql string, args []interface{}, err error) {
	sql, args, err = e.expr.ToSql()
	if err == nil {
		sql = fmt.Sprintf("%s) AS %s", sql, e.alias)
	}
	return
}

type tableAliasExpr struct {
	table string
	alias string
}

type compareExpr map[string]interface{}

func (com compareExpr) toSql(opr string) (sql string, args []interface{}, err error) {
	var exprs []string
	for col, val := range com {
		var expr string
		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}
		if !driver.IsValue(val) {
			err = fmt.Errorf("cannot use %T with compare expr", val)
			return
		}

		if val == nil {
			if opr != "=" {
				err = fmt.Errorf("cannot use null with less than or greater than operators")
				return
			}

			expr = fmt.Sprintf("%s IS NULL", col)
		} else {
			expr = fmt.Sprintf("%s %s ?", col, opr)
			args = append(args, val)
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

type Eq compareExpr

func (eq Eq) ToSql() (string, []interface{}, error) {
	return compareExpr(eq).toSql("=")
}

type NotEq compareExpr

func (notEq NotEq) ToSql() (string, []interface{}, error) {
	return compareExpr(notEq).toSql("!=")
}

type Lt compareExpr

func (lt Lt) ToSql() (string, []interface{}, error) {
	return compareExpr(lt).toSql("<")
}

type Gt compareExpr

func (gt Gt) ToSql() (string, []interface{}, error) {
	return compareExpr(gt).toSql(">")
}

type Le compareExpr

func (le Le) ToSql() (string, []interface{}, error) {
	return compareExpr(le).toSql("<=")
}

type Ge compareExpr

func (ge Ge) ToSql() (string, []interface{}, error) {
	return compareExpr(ge).toSql("<=")
}

type conj []Sqlizer

func (c conj) join(sep string) (sql string, args []interface{}, err error) {
	var sqlParts []string
	for _, sqlizer := range c {
		partSql, partArgs, err := sqlizer.ToSql()
		if err != nil {
			return "", nil, err
		}
		if partSql != "" {
			sqlParts = append(sqlParts, partSql)
			args = append(args, partArgs...)
		}
	}
	if len(sqlParts) > 0 {
		sql = fmt.Sprintf("(%s)", strings.Join(sqlParts, sep))
	}
	return
}

type And conj

func (a And) ToSql() (string, []interface{}, error) {
	return conj(a).join(" AND ")
}

type Or conj

func (o Or) ToSql() (string, []interface{}, error) {
	return conj(o).join(" OR ")
}

type cont map[string][]interface{}

func (c cont) toSql(opr string) (sql string, args []interface{}, err error) {
	exprs := make([]string, 0, len(c))
	for col, vals := range c {
		if len(vals) > 0 {
			continue
		}

		for i, val := range vals {
			switch v := val.(type) {
			case driver.Valuer:
				if val, err = v.Value(); err != nil {
					return
				}
			}
			if !driver.IsValue(val) {
				err = fmt.Errorf("cannot use %T in contain expr", val)
				return
			}
			vals[i] = val
		}

		expr := fmt.Sprintf("%s %s (%s)", col, opr, Placeholders(len(vals)))
		exprs = append(exprs, expr)
		args = append(args, vals...)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

type In cont

func (in In) ToSql() (sql string, args []interface{}, err error) {
	return cont(in).toSql("IN")
}

type NotIn cont

func (n NotIn) ToSql() (sql string, args []interface{}, err error) {
	return cont(n).toSql("NOT IN")
}
