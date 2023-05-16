package receive

// 发票信息推送接口
type InvoiceInformation struct {
	UUID            string         `json:"-"`
	Certificateid   string         `json:"certificateId"`
	Customerid      string         `json:"customerId"`
	Corpname        string         `json:"corpName"`
	Certificatetype string         `json:"certificateType"`
	Intercustomerid string         `json:"interCustomerId"`
	Invoiceinfos    []Invoiceinfos `json:"invoiceInfos"`
}

type Invoiceinfos struct {
	Invoicenotaxamt string `json:"invoiceNotaxAmt"`
	Invoiceccy      string `json:"invoiceCcy"`
	Sellername      string `json:"sellerName"`
	Invoicetype     string `json:"invoiceType"`
	Buyername       string `json:"buyerName"`
	Buyerusccode    string `json:"buyerUsccode"`
	Invoicedate     string `json:"invoiceDate"`
	Sellerusccode   string `json:"sellerUsccode"`
	Invoicecode     string `json:"invoiceCode"`
	Invoicenum      string `json:"invoiceNum"`
	Checkcode       string `json:"checkCode"`
	Invoiceamt      string `json:"invoiceAmt"`
}

// 推送历史交易信息接口
type TransactionHistory struct {
	UUID            string            `json:"-"`
	Customergrade   string            `json:"customerGrade"`
	Certificatetype string            `json:"certificateType"`
	Intercustomerid string            `json:"interCustomerId"`
	Corpname        string            `json:"corpName"`
	Financeid       string            `json:"financeId"`
	Certificateid   string            `json:"certificateId"`
	Customerid      string            `json:"customerId"`
	Usedinfos       []Usedinfos       `json:"usedInfos"`
	Settleinfos     []Settleinfos     `json:"settleInfos"`
	Orderinfos      []Orderinfos      `json:"orderInfos"`
	Receivableinfos []Receivableinfos `json:"receivableInfos"`
}

type Usedinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Usedamount     string `json:"usedAmount"`
	Ccy            string `json:"ccy"`
}
type Settleinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Settleamount   string `json:"settleAmount"`
	Ccy            string `json:"ccy"`
}
type Orderinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Orderamount    string `json:"orderAmount"`
	Ccy            string `json:"ccy"`
}
type Receivableinfos struct {
	Tradeyearmonth   string `json:"tradeYearMonth"`
	Receivableamount string `json:"receivableAmount"`
	Ccy              string `json:"ccy"`
}

// 推送入池数据接口
type EnterpoolData struct {
	UUID              string      `json:"-"`
	Datetimepoint     string      `json:"dateTimePoint"`
	Ccy               string      `json:"ccy"`
	Customerid        string      `json:"customerId"`
	Intercustomerid   string      `json:"interCustomerId"`
	Receivablebalance string      `json:"receivableBalance"`
	Planinfos         []Planinfos `json:"planInfos"`
	UsedInfos         []UsedInfos `json:"usedInfos"`
}

type Planinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Planamount     string `json:"planAmount"`
	Currency       string `json:"currency"`
}
type UsedInfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Usedamount     string `json:"usedAmount"`
	Currency       string `json:"currency"`
}

// 提交融资意向接口
type FinancingIntention struct {
	Custcdlinkposition string `json:"custcdLinkPosition"`
	Custcdlinkname     string `json:"custcdLinkName"`
	Certificateid      string `json:"certificateId"`
	Corpname           string `json:"corpName"`
	Remark             string `json:"remark"`
	Bankcontact        string `json:"bankContact"`
	Banklinkname       string `json:"bankLinkName"`
	Custcdcontact      string `json:"custcdContact"`
	Customerid         string `json:"customerId"`
	Financeid          string `json:"financeId"`
	Cooperationyears   string `json:"cooperationYears"`
	Certificatetype    string `json:"certificateType"`
	Intercustomerid    string `json:"interCustomerId"`
}

// 推送回款账户接口
type RawAccount struct {
	Backaccount     string `json:"backAccount"`
	Certificateid   string `json:"certificateId"`
	Customerid      string `json:"customerId"`
	Corpname        string `json:"corpName"`
	Lockremark      string `json:"lockRemark"`
	Certificatetype string `json:"certificateType"`
	Intercustomerid string `json:"interCustomerId"`
}
type UpdateCollectionAccount struct {
	UUID       string     `json:"-"`
	OldAccount RawAccount `json:"oldAccount"`
	NewAccount RawAccount `json:"newAccount"`
	FinanceId  string     `json:"financeId"`
}
type LockAccount struct {
	UUID            string `json:"-"`
	FinanceId       string `json:"financeId"`
	Backaccount     string `json:"backAccount"`
	Certificateid   string `json:"certificateId"`
	Customerid      string `json:"customerId"`
	Corpname        string `json:"corpName"`
	Lockremark      string `json:"lockRemark"`
	Certificatetype string `json:"certificateType"`
	Intercustomerid string `json:"interCustomerId"`
}

type SelectedInfosAndFinancingApplication struct {
	UUID                 string             `json:"-"`
	FinancingApplication FinancingIntention `json:"financingApplication"`
	Invoice              []InvoiceInfo      `json:"invoice"`
}
type InvoiceInfo struct {
	CertificateID   string `json:"certificateId"`
	CustomerID      string `json:"customerId"`
	CorpName        string `json:"corpName"`
	CertificateType string `json:"certificateType"`
	InterCustomerID string `json:"interCustomerId"`
	InvoiceNotaxAmt string `json:"invoiceNotaxAmt"`
	InvoiceCcy      string `json:"invoiceCcy"`
	SellerName      string `json:"sellerName"`
	InvoiceType     string `json:"invoiceType"`
	BuyerName       string `json:"buyerName"`
	BuyerUsccode    string `json:"buyerUsccode"`
	InvoiceDate     string `json:"invoiceDate"`
	SellerUsccode   string `json:"sellerUsccode"`
	InvoiceCode     string `json:"invoiceCode"`
	InvoiceNum      string `json:"invoiceNum"`
	CheckCode       string `json:"checkCode"`
	InvoiceAmt      string `json:"invoiceAmt"`
	FinancingID     string `json:"-"`
}
