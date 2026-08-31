package compile

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/kalo-build/go-util/strcase"
)

var pluralizeClient = pluralize.NewClient()

func GetTableNameFromModel(modelName string) string {
	return Pluralize(strcase.ToSnakeCaseLower(modelName))
}

func GetColumnNameFromField(fieldName string) string {
	return strcase.ToSnakeCaseLower(fieldName)
}

func GetForeignKeyColumnName(relatedModelName, relatedFieldName string) string {
	return fmt.Sprintf("%s_%s",
		strcase.ToSnakeCaseLower(relatedModelName),
		strcase.ToSnakeCaseLower(relatedFieldName))
}

func Pluralize(word string) string {
	return pluralizeClient.Plural(word)
}

func FormatSQLValue(value any) string {
	if value == nil {
		return "NULL"
	}
	switch v := value.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
