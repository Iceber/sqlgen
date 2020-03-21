package sqlgen

type statementBuilder struct {
	placeholderFormat PlaceholderFormat
}

var defaultStatementBuilder = statementBuilder{}

func NewStatementBuilder(placeholderFormat PlaceholderFormat) statementBuilder {
	return statementBuilder{
		placeholderFormat: placeholderFormat,
	}
}

func Select(columns ...string) *SelectBuilder {
	return defaultStatementBuilder.Select(columns...)
}

func Update(table string) *UpdateBuilder {
	return defaultStatementBuilder.Update(table)
}

func Delete(table string) *DeleteBuilder {
	return defaultStatementBuilder.Delete(table)
}

func Insert(table string) *InsertBuilder {
	return defaultStatementBuilder.Insert(table)
}

func (b statementBuilder) Select(columns ...string) *SelectBuilder {
	return NewSelectBuilder(b).Columns(columns...)
}

func (b statementBuilder) Insert(table string) *InsertBuilder {
	return NewInsertBuilder(b).Table(table)
}

func (b statementBuilder) Update(table string) *UpdateBuilder {
	return NewUpdateBuilder(b).Table(table)
}

func (b statementBuilder) Delete(table string) *DeleteBuilder {
	return NewDeleteBuilder(b).Table(table)
}
