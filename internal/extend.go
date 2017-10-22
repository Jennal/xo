package internal

import (
	"fmt"
	"strings"
)

func LoadExtraFields(args *ArgType, typeTpl *Type) error {
	arr := strings.Split(typeTpl.Table.TableComment, "\n")
	if len(arr) == 0 {
		return nil
	}

	for _, line := range arr {
		data, err := parseExtraData(line)
		if err != nil || len(data) == 0 {
			return err
		}

		typeTpl.ExtraFields = append(typeTpl.ExtraFields, &ExtraField{
			Name:     data["name"],
			Type:     data["type"],
			JsonName: strings.ToLower(data["name"]),
			Comment:  line,
		})
	}

	return nil
}

func parseExtraData(line string) (map[string]string, error) {
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
			return nil, fmt.Errorf("wrong param: %s", item)
		}

		result[kv[0]] = kv[1]
	}

	return result, nil
}
