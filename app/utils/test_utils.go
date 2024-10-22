package utils

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"strings"
	"unicode"
)

func GenerateSqlMockRow(data any) *sqlmock.Rows {
	columns, values := structToColumnsAndValues(data)
	return sqlmock.NewRows(columns).AddRow(values...)
}

func GenerateSqlMockRows(data any) *sqlmock.Rows {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("expected a slice of pointers to structs")
	}

	if v.Len() == 0 {
		panic("empty slice provided")
	}

	firstElem := v.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}
	columns, _ := structToColumnsAndValues(firstElem.Interface())
	rows := sqlmock.NewRows(columns)

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		_, values := structToColumnsAndValues(elem)
		rows.AddRow(values...)
	}

	return rows
}

func structToColumnsAndValues(data any) ([]string, []driver.Value) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic("expected a pointer to a struct")
	}

	var columns []string
	var values []driver.Value

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if shouldSkipField(field) {
			continue
		}

		columns = append(columns, getColumnName(field))
		values = append(values, v.Field(i).Interface())
	}

	return columns, values
}

func getColumnName(field reflect.StructField) string {
	columnTag := field.Tag.Get("gorm")
	if columnTag == "" || !strings.HasPrefix(columnTag, "column:") {
		return toSnakeCase(field.Name)
	}

	return strings.TrimPrefix(columnTag, "column:")
}

func shouldSkipField(field reflect.StructField) bool {
	gormTag := field.Tag.Get("gorm")
	for _, tag := range excludeTags {
		if strings.Contains(gormTag, tag) {
			return true
		}
	}

	return false
}

var excludeTags = []string{
	"many2many", "hasMany", "belongsTo", "hasOne",
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
