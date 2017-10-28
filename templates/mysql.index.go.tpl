{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" "XOLog" .Fields) -}}
{{- $table := (schema .Schema .Type.Table.TableName) -}}
{{- $convlist := (convlist .Type) -}}
// {{ .FuncName }} retrieves a row from '{{ $table }}' as a {{ .Type.Name }}.
//
// Generated from index '{{ .Index.IndexName }}'.
func {{ .FuncName }}(db XODB{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}[]{{ end }}*{{ .Type.Name }}, error) {
	var err error

	// sql query
	const sqlstr = "SELECT " +
		"{{ colnames .Type.Fields }} " +
		"FROM {{ $table }} " +
		"WHERE {{ colnamesquery .Fields " AND " }}"

	// run query
	XOLog(sqlstr{{ goparamlist .Fields true false }})
{{- if .Index.IsUnique }}
	{{ $short }} := {{ .Type.Name }}{
	{{- if .Type.PrimaryKey }}
		_exists: true,
	{{ end -}}
	}

	// ref init
	{{ refvalinit .Type $short }}

	err = db.QueryRow(sqlstr{{ goparamlist .Fields true false }}).Scan({{ fieldnames .Type.Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}

	// ref load
	{{ reffillval .Type $short "db" }}

	{{- if $convlist }}
		// json fields
		{{- range $convlist }}
			{{ $short }}.{{ .Name }} = &{{ puretype .Type }}{}
			if len({{ $short }}.{{ .Conv.JsFieldName }}) > 0 {
				//no care about error
				json.Unmarshal([]byte({{ $short }}.{{ .Conv.JsFieldName }}), {{ $short }}.{{ .Name }})
			}

		{{- end }}
	{{- end }}

	return &{{ $short }}, nil
{{- else }}
	q, err := db.Query(sqlstr{{ goparamlist .Fields true false }})
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*{{ .Type.Name }}{}
	for q.Next() {
		{{ $short }} := {{ .Type.Name }}{
		{{- if .Type.PrimaryKey }}
			_exists: true,
		{{ end -}}
		}

		// ref init
		{{ refvalinit .Type $short }}

		// scan
		err = q.Scan({{ fieldnames .Type.Fields (print "&" $short) }})
		if err != nil {
			return nil, err
		}

		// ref load
		{{ reffillval .Type $short "db" }}

		{{- if $convlist }}
			// json fields
			{{- range $convlist }}
				{{ $short }}.{{ .Name }} = &{{ puretype .Type }}{}
				if len({{ $short }}.{{ .Conv.JsFieldName }}) > 0 {
					err = json.Unmarshal([]byte({{ $short }}.{{ .Conv.JsFieldName }}), {{ $short }}.{{ .Name }})
					if err != nil {
						return nil, err
					}
				}

			{{- end }}
		{{- end }}

		res = append(res, &{{ $short }})
	}

	return res, nil
{{- end }}
}

