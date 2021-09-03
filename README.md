# dbHelper
一个用于自动将数据库表信息转为go的struct的工具，目前仅支持mysql

# 用法
> 列出所有数据库
```shell
dbHelper list database -u root -H 127.0.0.1 --password xxx
# 输出
information_schema                      mysql                                   test                                  japi

```
> 列出所有表
```shell
dbHelper list table -u root -H 127.0.0.1 --password xxx -d goods
# 输出
TableName                                         |TableEngine     |Rows            |DataLength      |IndexLength     |AutoIncrement   |CreateTime              |Comment
article_list                                      |MyISAM          |28              |6.1 kB          |2.0 kB          |29              |2017-08-23 22:46:38     |单页面内容表
article_orders                                    |InnoDB          |22              |16 kB           |0 B             |41              |2021-04-23 18:49:58     |知识付费订单
```
> 生成go struct
```shell
dbHelper gen -u root -H 127.0.0.1 --password xxx -d goods 
```
- 结果为
```go
package goods

import "time"

// ArticleList -> table: article_list
// comment：单页面内容表
type ArticleList struct {
	ID         uint      `gorm:"column:id;primaryKey;not null;default:null"`
	Title      string    `gorm:"column:title;default:null"`
	Content    string    `gorm:"column:content;default:null"`
	CreateTime int       `gorm:"column:createTime;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:updateTime;autoUpdateTime"`
}

func (ArticleList) TableName() string {
	return "article_list"
}

// ArticleOrder -> table: article_orders
// comment：知识付费订单
type ArticleOrder struct {
	ID            int       `gorm:"column:id;primaryKey;not null;default:null"`
	CustomerID    string    `gorm:"column:customer_id"`                // 会员id
	WxPlat        int       `gorm:"column:wx_plat"`                    // 支付到微信平台
	FromType      string    `gorm:"column:from_type"`                  // 来源
	ProductID     string    `gorm:"column:product_id"`                 // 产品id
	TransactionID string    `gorm:"column:transaction_id"`             // 交易单号
	OrderNo       string    `gorm:"column:order_no"`                   // 订单编号
	OutTradeNo    string    `gorm:"column:out_trade_no"`               // 商户订单
	TotalAmount   float64   `gorm:"column:total_amount;default:null"`  // 支付金额
	PayFinishTm   time.Time `gorm:"column:pay_finish_tm;default:null"` // 支付完成时间
	Detail        string    `gorm:"column:detail;default:null"`        // 支付内容描述
	PayMethod     int8      `gorm:"column:pay_method"`                 // 支付方式：1-微信支付
	TradeType     string    `gorm:"column:trade_type"`                 // 采用的支付方式类型
	PayStatus     int8      `gorm:"column:pay_status"`                 // 1-待支付；2-已支付；
	CreateTm      time.Time `gorm:"column:create_tm;autoCreateTime"`   // 创建时间
	UpdateTm      time.Time `gorm:"column:update_tm;not null;autoUpdateTime"`
}

func (ArticleOrder) TableName() string {
	return "article_orders"
}
```
> 帮助
```shell
dbHelper -h
根据数据库表生成go的struct模板，目前支持MySQL数据库，对应的ORM框架为gorm

Usage:
  dbHelper [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  gen         生成go struct文件
  help        Help about any command
  list        列出数据库或者表列表

Flags:
  -c, --charset string    数据库连接编码 (default "utf8mb4")
  -h, --help              help for dbHelper
  -H, --host string       数据库连接host
      --password string   数据库连接密码
  -p, --port int          数据库连接端口 (default 3306)
  -u, --user string       数据库连接用户名

Use "dbHelper [command] --help" for more information about a command.
```
