package sqlgen

import "fmt"

func TableAlias(table, alias string) tableAliasExpr {
	return tableAliasExpr{table, alias}
}

func (e tableAliasExpr) String() string {
	return fmt.Sprintf("%s AS %s", e.table, e.alias)
}

func (e tableAliasExpr) Colume(column string) string {
	return fmt.Sprintf("%s.%s", e.alias, column)
}

/*
func (e tableAliasExpr) ToSql() (sql string, args []interface{}, err error) {
	sql = e.String()
	return
}
*/
