package cmd

import (
	"bytes"
	"fmt"
	"git.iglou.eu/Imported/go-wildcard"
	"github.com/gookit/color"
	"github.com/lewjian/dbHelper/config"
	"github.com/lewjian/dbHelper/db"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

type TplField struct {
	PackageName string
	ModelName   string
	TableName   string
	Columns     []TplCol
}
type TplCol struct {
	ColumnName string
	ColumnType string
	Tag        string
	Comment    string
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成go struct文件",
	Long:  "按照指定规则生成go struct文件，具体可以看帮助",
	Run: func(cmd *cobra.Command, args []string) {
		ms, err := ConnectDB(cmd)
		if err != nil {
			color.Warnln("连接数据库失败", err)
			return
		}

		tables, err := ms.ListTables()
		if err != nil {
			color.Warnln("列出table失败", err)
			return
		}
		// 创建、更新时间字段
		createFieldName, _ := cmd.Flags().GetString("create_tm")
		updateFieldName, _ := cmd.Flags().GetString("update_tm")
		// 是否使用sqlnull形式
		useSqlNull, _ := cmd.Flags().GetBool("useSqlNull")
		// 是否将内容输出为单独文件
		separate, _ := cmd.Flags().GetBool("separate")

		// 输出文件夹
		outputDir, _ := cmd.Flags().GetString("out_dir")

		matchNum := 0
		tablePattern, _ := cmd.Flags().GetString("table")
		t, err := template.New("struct").Parse(config.StructTpl)
		if err != nil {
			color.Errorf("解析模板失败: %s", err.Error())
			return
		}
		outputDir = path.Join(outputDir, ms.Database)
		err = MkdirIfNotExists(outputDir)
		if err != nil {
			color.Errorf("创建目录[%s]失败，原因为：%s", outputDir, err.Error())
			return
		}
		tableFilename := path.Join(outputDir, fmt.Sprintf("%s.go", ms.Database))
		tf, err := os.OpenFile(tableFilename, os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			color.Errorf("创建文件[%s]失败，原因为：%s", tableFilename, err.Error())
			return
		}
		if len(tables) > 0 {
			for _, item := range tables {
				if tablePattern != "" && !wildcard.MatchSimple(tablePattern, item.Name) {
					continue
				}
				columns, err := ms.ListColumns(item.Name)
				if err != nil {
					color.Errorf("获取表%s字段信息失败，原因为：%s", item.Name, err.Error())
					break
				}
				item.Columns = columns
				var cols []TplCol
				for _, col := range columns {
					colType := TransferDataType(col, useSqlNull)
					if col.Name == createFieldName || col.Name == updateFieldName {
						colType = "time.Time"
					}
					tag := fmt.Sprintf("column:%s", col.Name)
					if col.IsPrimary {
						tag = fmt.Sprintf("%s;primaryKey", tag)
					}
					if !col.IsNullable {
						tag = fmt.Sprintf("%s;not null", tag)
					}
					cols = append(cols, TplCol{
						ColumnName: Camelize(col.Name),
						ColumnType: colType,
						Tag:        tag,
						Comment:    col.Comment,
					})
				}
				packageName := ms.Database

				if !separate && matchNum > 0 {
					packageName = ""
				}
				tpl := TplField{
					PackageName: packageName,
					ModelName:   Camelize(item.Name),
					TableName:   item.Name,
					Columns:     cols,
				}
				if separate {
					filename := path.Join(outputDir, fmt.Sprintf("%s.go", item.Name))
					f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0777)
					if err != nil {
						color.Errorf("创建文件[%s]失败，原因为：%s", filename, err.Error())
						return
					} else {
						err = t.Execute(f, tpl)
						if err != nil {
							color.Errorf("写入文件[%s]失败，原因为：%s", filename, err.Error())
							return
						}
					}
					f.Close()
					formatOutput(item.Name, filename)
					Gofmt(filename)
				} else {
					t.Execute(tf, tpl)
				}
				matchNum++
			}
			tf.Close()
			if separate {
				os.Remove(tableFilename)
			} else {
				formatOutput(ms.Database, tableFilename)
				Gofmt(tableFilename)
			}
		}
		color.Warnln("\n共处理：", matchNum)
	},
}

func init() {
	genCmd.Flags().StringP("out_dir", "o", "./output", "生成文件的目录，默认当前目录下的output下")
	genCmd.Flags().String("create_tm", "create_tm", "创建时间的字段名，默认create_tm")
	genCmd.Flags().String("update_tm", "update_tm", "更新时间的字段名，默认update_tm")
	genCmd.Flags().Bool("useSqlNull", false, "针对某些允许为null的字段，是否设置为sql.NullString等格式")
	genCmd.Flags().Bool("separate", false, "是否将一张表生成一个文件（文件名为表名），还是所有的表生成为一个文件(数据库名.go)")
	genCmd.Flags().StringP("database", "d", "", "指定数据库")
	genCmd.Flags().StringP("table", "t", "", "列出指定的表，支持*通配符")
	genCmd.MarkFlagRequired("database")

	rootCmd.AddCommand(genCmd)
}

func Camelize(name string) string {
	if name == "id" {
		return "ID"
	}
	res := strings.Split(name, "_")
	var s strings.Builder
	for _, sub := range res {
		s.WriteString(UCFirst(sub))
	}
	return s.String()
}

func UCFirst(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[:1]), s[1:])
}

func TransferDataType(col db.ColumnInfo, useSqlNull bool) string {
	unsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")
	if useSqlNull && col.IsNullable {
		switch strings.ToUpper(col.DataType) {
		case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER":
			return "sql.NullInt32"
		case "BIGINT":
			return "sql.NullInt64"
		case "FLOAT", "DOUBLE", "DECIMAL":
			return "sql.NullFloat64"
		case "DATE", "DATETIME", "TIMESTAMP":
			return "sql.NullTime"
		case "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
			return "[]byte"
		default:
			return "sql.NullString"
		}
	} else {
		switch strings.ToUpper(col.DataType) {
		case "TINYINT":
			return If(unsigned, "uint8", "int8")
		case "SMALLINT":
			return If(unsigned, "uint16", "int16")
		case "MEDIUMINT", "INT", "INTEGER":
			return If(unsigned, "uint", "int")
		case "BIGINT":
			return If(unsigned, "uint64", "int64")
		case "FLOAT", "DOUBLE", "DECIMAL":
			return "float64"
		case "DATE", "DATETIME", "TIMESTAMP":
			return "time.Time"
		case "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
			return "[]byte"
		default:
			return "string"
		}
	}
}

func If(ok bool, a, b string) string {
	if ok {
		return a
	}
	return b
}

func MkdirIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, 0644)
	}
	return nil
}

func formatOutput(name, outputFilename string) {
	color.Success.Println(fmt.Sprintf("%-50s——>\t%s", name, outputFilename))
}

func Gofmt(file string) {
	ExecCommand(fmt.Sprintf("go fmt %s", file))
}

// ExecCommand 执行shell名，阻塞，需要输出返回值
func ExecCommand(cmdStr string) (output, errStr string, err error) {
	var cmd *exec.Cmd
	var stdout, stderr bytes.Buffer
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell")
	} else {
		cmd = exec.Command("/bin/bash")
	}
	stdin := bytes.NewBufferString(cmdStr)
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	if err = cmd.Wait(); err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), err
}
