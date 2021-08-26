package cmd

import (
	"fmt"
	"git.iglou.eu/Imported/go-wildcard"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"github.com/lewjian/dbHelper/config"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "列出数据库或者表列表",
		Long:  "将数据库的所有名称或者指定数据库的所有表列举出来",
	}

	dbCmd = &cobra.Command{
		Use:   "database",
		Short: "列出数据库",
		Long:  "列出数据库",
		Run: func(cmd *cobra.Command, args []string) {
			ms, err := ConnectDB(cmd)
			if err != nil {
				color.Warnln("连接数据库失败", err)
				return
			}
			dbs, err := ms.ListDatabases()
			if err != nil {
				color.Warnln("列出数据库失败", err)
				return
			}
			dbLen := len(dbs)
			var out string
			for i := 0; i < dbLen ; i++ {
				if i != 0 && i % 4 == 0 {
					color.Success.Println(out)
					out = ""
				}
				out = fmt.Sprintf("%s%-40s", out, dbs[i])
			}
			if out != "" {
				color.Success.Println(out)
			}

			color.Warnln("\n合计数据库数量：", dbLen)

		},
	}

	tableCmd = &cobra.Command{
		Use:   "table",
		Short: "列出所有table",
		Long:  "列出所有table",
		Run: func(cmd *cobra.Command, args []string) {
			VerifyFlagRequired(cmd, []string{"database"})
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

			color.Yellow.Println(fmt.Sprintf("%-50s|%-16s|%-16s|%-16s|%-16s|%-16s|%-24s|%s",
				"TableName", "TableEngine", "Rows", "DataLength", "IndexLength", "AutoIncrement", "CreateTime", "Comment"))
			tableLen := len(tables)
			matchNum := 0
			tablePattern, _ := cmd.Flags().GetString("table")
			if tableLen > 0 {
				for _, item := range tables {
					if tablePattern != "" && !wildcard.MatchSimple(tablePattern, item.Name) {
						continue
					}
					color.Success.Println(fmt.Sprintf("%-50s|%-16s|%-16d|%-16s|%-16s|%-16d|%-24s|%s",
						item.Name, item.Engine, item.Rows, humanize.Bytes(uint64(item.DataLength)),
						humanize.Bytes(uint64(item.IndexLength)), item.AutoIncrement, item.CreateTime.Format(config.DefaultDateTimeFormatTpl),
						item.Comment))
					matchNum++
				}
			}

			color.Warnln("\n合计表数量：", matchNum)
		},
	}
)

func init() {
	tableCmd.Flags().StringP("table", "t", "", "列出指定的表，支持*通配符")
	tableCmd.Flags().StringP("database", "d", "", "指定数据库")
	tableCmd.MarkFlagRequired("database")
	listCmd.AddCommand(dbCmd, tableCmd)
}
