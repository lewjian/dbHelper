package config

import _ "embed"

const DefaultDateTimeFormatTpl = "2006-01-02 15:04:05"

//go:embed type.tpl
var StructTpl string
