package internal

import (
	"fmt"
	"strings"

	"github.com/knq/snaker"
)

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

		typeTpl.ExtraFields = append(typeTpl.ExtraFields, &ExtraField{
			Name:     data["name"],
			Type:     data["type"],
			JsonName: snaker.CamelToSnake(data["name"]),
			Comment:  line,
		})
	}

	return nil
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

func LoadColumnType(args *ArgType, field *Field) error {
	data, err := parseData(field.Col.ColumnComment)
	if len(data) == 0 || err != nil {
		return err
	}

	if ref, ok := data["ref"]; ok {
		tk := strings.Split(ref, ".")
		if len(tk) != 2 {
			return fmt.Errorf("wrong param: %s", field.Col.ColumnComment)
		}

		field.Ref = &Ref{
			Type:       "*" + SingularizeTableName(tk[0], args.KeepTablePrefix),
			TableName:  snaker.SnakeToCamelIdentifier(tk[0]),
			ColumnName: snaker.CamelToSnakeIdentifier(tk[1]),
			KeyName:    snaker.SnakeToCamelIdentifier(tk[1]),
		}
		field.Ref.FuncName = field.Ref.TableName + "By" + field.Ref.KeyName

		field.Name = data["name"]
		field.Type = field.Ref.Type

		// spew.Dump(field)

		return nil
	}

	//TODO: jennal conv
	// if conv, ok := data["conv"]; ok {
	// 	return nil
	// }

	return nil
}
