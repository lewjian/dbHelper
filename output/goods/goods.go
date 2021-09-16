package goods

import "time"

// ArticleList -> table: article_list
// comment：单页面内容表
type ArticleList struct {
	ID         uint      `gorm:"column:id;primaryKey;not null"`
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
	ID            int       `gorm:"column:id;primaryKey;not null"`
	CustomerID    string    `gorm:"column:customer_id;default:0"`      // 会员id
	WxPlat        int       `gorm:"column:wx_plat;default:0"`          // 支付到微信平台
	FromType      string    `gorm:"column:from_type"`                  // 来源
	ProductID     string    `gorm:"column:product_id"`                 // 产品id
	TransactionID string    `gorm:"column:transaction_id"`             // 交易单号
	OrderNo       string    `gorm:"column:order_no"`                   // 订单编号
	OutTradeNo    string    `gorm:"column:out_trade_no"`               // 商户订单
	TotalAmount   float64   `gorm:"column:total_amount;default:null"`  // 支付金额
	PayFinishTm   time.Time `gorm:"column:pay_finish_tm;default:null"` // 支付完成时间
	Detail        string    `gorm:"column:detail;default:null"`        // 支付内容描述
	PayMethod     int8      `gorm:"column:pay_method;default:0"`       // 支付方式：1-微信支付
	TradeType     string    `gorm:"column:trade_type"`                 // 采用的支付方式类型
	PayStatus     int8      `gorm:"column:pay_status;default:0"`       // 1-待支付；2-已支付；
	CreateTm      time.Time `gorm:"column:create_tm;autoCreateTime"`   // 创建时间
	UpdateTm      time.Time `gorm:"column:update_tm;not null;autoUpdateTime"`
}

func (ArticleOrder) TableName() string {
	return "article_orders"
}
