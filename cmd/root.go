package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/lewjian/dbHelper/db/mysql"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dbHelper",
		Short: "根据数据库表生成go的struct模板",
		Long:  `根据数据库表生成go的struct模板，目前支持MySQL数据库，对应的ORM框架为gorm`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("user", "u", "", "数据库连接用户名")
	rootCmd.PersistentFlags().StringP("host", "H", "", "数据库连接host")
	rootCmd.PersistentFlags().StringP("charset", "c", "utf8mb4", "数据库连接编码")
	rootCmd.PersistentFlags().String("password", "", "数据库连接密码")
	rootCmd.PersistentFlags().IntP("port", "p", 3306, "数据库连接端口")
	rootCmd.MarkPersistentFlagRequired("user")
	rootCmd.MarkPersistentFlagRequired("host")

	//viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	//viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	//viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	//viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	//viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	rootCmd.AddCommand(listCmd)
}

func ConnectDB(cmd *cobra.Command) (*mysql.MySQL, error) {
	user, _ := cmd.Flags().GetString("user")
	host, _ := cmd.Flags().GetString("host")
	password, _ := cmd.Flags().GetString("password")
	port, _ := cmd.Flags().GetInt("port")
	charset, _ := cmd.Flags().GetString("charset")
	database, _ := cmd.Flags().GetString("database")
	if password == "" {
		password, _ = GetPassword()
	}
	ms := mysql.NewMySQL(user, host, password, port, charset, database)
	err := ms.Connect()
	if err != nil {
		return nil, err
	}
	return ms, nil
}

// GetPassword 获取用户命令行输入的密码
func GetPassword() (string, error) {
	color.Info.Printf("请输入密码：")
	var password string
	_, err := fmt.Scanln(&password)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}

func VerifyFlagRequired(cmd *cobra.Command, flags []string) {
	for _, flag := range flags {
		res, err := cmd.Flags().GetString(flag)
		if err != nil || res == "" {
			color.Error.Println(flag, "参数不可为空")
			cmd.Help()
			os.Exit(0)
		}
	}
}
