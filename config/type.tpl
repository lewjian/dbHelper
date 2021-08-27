{{if ne .PackageName ""}}package {{.PackageName}}
{{end}}
// {{.ModelName}} -> table: {{.TableName}}
{{if ne .Comment ""}}// comment：{{.Comment}}{{end}}
type {{.ModelName}} struct {
	{{range $i,$v := .Columns}}{{$v.ColumnName}} {{$v.ColumnType}} `gorm:"{{$v.Tag}}"` {{if ne $v.Comment ""}}// {{$v.Comment}} {{end}}
	{{end}}
}

func ({{.ModelName}}) TableName() string {
    return "{{.TableName}}"
}
