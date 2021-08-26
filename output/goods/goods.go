package goods

// ArticleList -> table: article_list
type ArticleList struct {
	ID         uint      `gorm:"column:id;primaryKey;not null"`
	Title      string    `gorm:"column:title"`
	Content    string    `gorm:"column:content"`
	CreateTime int       `gorm:"column:createTime"`
	UpdateTime time.Time `gorm:"column:updateTime"`
}

func (ArticleList) TableName() string {
	return "article_list"
}

// ArticleOrders -> table: article_orders
type ArticleOrders struct {
	ID            int       `gorm:"column:id;primaryKey;not null"`
	CustomerId    string    `gorm:"column:customer_id"`    // 会员id
	WxPlat        int       `gorm:"column:wx_plat"`        // 支付到微信平台
	FromType      string    `gorm:"column:from_type"`      // 来源
	ProductId     string    `gorm:"column:product_id"`     // 产品id
	TransactionId string    `gorm:"column:transaction_id"` // 交易单号
	OrderNo       string    `gorm:"column:order_no"`       // 订单编号
	OutTradeNo    string    `gorm:"column:out_trade_no"`   // 商户订单
	TotalAmount   float64   `gorm:"column:total_amount"`   // 支付金额
	PayFinishTm   time.Time `gorm:"column:pay_finish_tm"`  // 支付完成时间
	Detail        string    `gorm:"column:detail"`         // 支付内容描述
	PayMethod     int8      `gorm:"column:pay_method"`     // 支付方式：1-微信支付
	TradeType     string    `gorm:"column:trade_type"`     // 采用的支付方式类型
	PayStatus     int8      `gorm:"column:pay_status"`     // 1-待支付；2-已支付；
	CreateTm      time.Time `gorm:"column:create_tm"`      // 创建时间
	UpdateTm      time.Time `gorm:"column:update_tm;not null"`
}

func (ArticleOrders) TableName() string {
	return "article_orders"
}
