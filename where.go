package sqlgen

import "fmt"

type wherePart part

func newWherePart(pred interface{}, args ...interface{}) Sqlizer {
	return &wherePart{pred: pred, args: args}
}

func (p wherePart) ToSql() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
	case Sqlizer:
		return pred.ToSql()
	case map[string]interface{}:
		return Eq(pred).ToSql()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}