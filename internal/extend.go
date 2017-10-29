package internal

import (
	"fmt"
	"strings"

	"github.com/knq/snaker"
)

func FixRelkindWithComment(args *ArgType, typeMap map[string]*Type) error {
	for _, typeTpl := range typeMap {
		for _, field := range typeTpl.Fields {
			// load column type from comment
			err := LoadColumnType(args, field)
			if err != nil {
				return err
			}
		}

		// load extra fields
		err := LoadTableExtraFields(args, typeTpl)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadTableExtraFields(args *ArgType, typeTpl *Type) error {
	arr := strings.Split(typeTpl.Table.TableComment, "\n")
	if len(arr) == 0 {
		return nil
	}

	for _, line := range arr {
		data, err := parseData(line)
		if err != nil || len(data) == 0 {
			return err
		}

		var extraField *ExtraField
		if ref, ok := data["ref"]; ok {
			selfKey, otherTable, otherKey := parseExtraRef(ref)
			index, err := getIndex(otherTable, otherKey)
			if err != nil {
				return err
			}

			t := typeFromIndex(index)
			extraField = &ExtraField{
				Field: &Field{
					Name:    data["name"],
					Type:    t,
					Comment: line,

					Ref: &Ref{
						Type:          t,
						RefTableName:  snaker.SnakeToCamelIdentifier(otherTable),
						RefColumnName: snaker.CamelToSnakeIdentifier(otherKey),
						RefKeyName:    snaker.SnakeToCamelIdentifier(otherKey),
						SelfKeyName:   snaker.SnakeToCamelIdentifier(selfKey),
						FuncName:      index.FuncName,
						IsUnique:      index.Index.IsUnique,
					},
				},
				JsonName: snaker.CamelToSnake(data["name"]),
			}
		} else {
			extraField = &ExtraField{
				Field: &Field{
					Name:    data["name"],
					Type:    data["type"],
					Comment: line,

					Ref: nil,
				},
				JsonName: snaker.CamelToSnake(data["name"]),
			}
		}

		typeTpl.ExtraFields = append(typeTpl.ExtraFields, extraField)
	}

	return nil
}

// ref should be "selfKey#otherTable.otherKey"
func parseExtraRef(ref string) (string, string, string) {
	arr := strings.Split(ref, "#")
	selfKey := arr[0]
	arr = strings.Split(arr[1], ".")

	//selfKey, otherTable, otherKey
	return selfKey, arr[0], arr[1]
}

func parseData(line string) (map[string]string, error) {
	start := strings.Index(line, "`") + 1
	end := strings.LastIndex(line, "`")
	//not extra data format
	if start < 0 || end < 0 {
		return nil, nil
	}
	line = line[start:end]

	//not extra data format
	if !strings.HasPrefix(line, "xo:") {
		return nil, nil
	}
	line = line[3:]
	arr := strings.Split(line, ",")
	result := map[string]string{}

	for _, item := range arr {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("wrong param: %s", line)
		}

		result[kv[0]] = kv[1]
	}

	return result, nil
}

func getIndex(table string, column string) (*Index, error) {
	for _, index := range s_indexMap {
		if index.Type.Table.TableName != table {
			continue
		}

		if len(index.Fields) != 1 {
			continue
		}

		if index.Fields[0].Col.ColumnName != column {
			continue
		}

		return index, nil
	}

	return nil, fmt.Errorf("index `%v` in table `%v` not exists", column, table)
}

func typeFromIndex(index *Index) string {
	if index.Index.IsUnique {
		return "*" + index.Type.Name
	}

	return "[]*" + index.Type.Name
}

func LoadColumnType(args *ArgType, field *Field) error {
	data, err := parseData(field.Col.ColumnComment)
	if len(data) == 0 || err != nil {
		return err
	}

	// ref
	if ref, ok := data["ref"]; ok {
		tk := strings.Split(ref, ".")
		if len(tk) != 2 {
			return fmt.Errorf("wrong param: %s", field.Col.ColumnComment)
		}

		index, err := getIndex(tk[0], tk[1])
		if err != nil {
			return err
		}

		field.Ref = &Ref{
			Type:          typeFromIndex(index),
			RefTableName:  snaker.SnakeToCamelIdentifier(tk[0]),
			RefColumnName: snaker.CamelToSnakeIdentifier(tk[1]),
			RefKeyName:    snaker.SnakeToCamelIdentifier(tk[1]),
			SelfKeyName:   field.Name,
			FuncName:      index.FuncName,
			IsUnique:      index.Index.IsUnique,
		}

		field.Name = data["name"]
		field.Type = field.Ref.Type

		// spew.Dump(field)

		return nil
	}

	// conv
	if conv, ok := data["conv"]; ok && conv == "json" {
		field.Type = data["type"]
		field.Conv = &Conv{
			JsFieldName: "js" + field.Name,
		}

		return nil
	}

	return nil
}
