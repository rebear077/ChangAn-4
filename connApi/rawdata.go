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
	Invoicenotaxamt string `json:"InvoiceNotaxAmt"`
	Invoiceccy      string `json:"InvoiceCcy"`
	Sellername      string `json:"SellerName"`
	Invoicetype     string `json:"InvoiceType"`
	Buyername       string `json:"BuyerName"`
	Buyerusccode    string `json:"BuyerUsccode"`
	Invoicedate     string `json:"InvoiceDate"`
	Sellerusccode   string `json:"SellerUsccode"`
	Invoicecode     string `json:"InvoiceCode"`
	Invoicenum      string `json:"InvoiceNum"`
	Checkcode       string `json:"CheckCode"`
	Invoiceamt      string `json:"InvoiceAmt"`
}

// 推送历史交易信息接口
type TransactionHistory struct {
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
	Tradeyearmonth string `json:"TradeYearMonth"`
	Usedamount     string `json:"UsedAmount"`
	Ccy            string `json:"Ccy"`
}
type Settleinfos struct {
	Tradeyearmonth string `json:"TradeYearMonth"`
	Settleamount   string `json:"SettleAmount"`
	Ccy            string `json:"Ccy"`
}
type Orderinfos struct {
	Tradeyearmonth string `json:"TradeYearMonth"`
	Orderamount    string `json:"OrderAmount"`
	Ccy            string `json:"Ccy"`
}
type Receivableinfos struct {
	Tradeyearmonth   string `json:"TradeYearMonth"`
	Receivableamount string `json:"ReceivableAmount"`
	Ccy              string `json:"Ccy"`
}

// 推送入池数据接口
type EnterpoolData struct {
	Datetimepoint     string              `json:"dateTimePoint"`
	Ccy               string              `json:"ccy"`
	Customerid        string              `json:"customerId"`
	Intercustomerid   string              `json:"interCustomerId"`
	Receivablebalance string              `json:"receivableBalance"`
	Planinfos         []Planinfos         `json:"planInfos"`
	Providerusedinfos []Providerusedinfos `json:"ProviderUsedInfos"`
}

type Planinfos struct {
	Tradeyearmonth string `json:"TradeYearMonth"`
	Planamount     string `json:"PlanAmount"`
	Currency       string `json:"Currency"`
}
type Providerusedinfos struct {
	Tradeyearmonth string `json:"TradeYearMonth"`
	Usedamount     string `json:"UsedAmount"`
	Currency       string `json:"Currency"`
}

// 提交融资意向接口
type FinancingIntention struct {
	Custcdlinkposition string `json:"CustcdLinkPosition"`
	Custcdlinkname     string `json:"CustcdLinkName"`
	Certificateid      string `json:"CertificateId"`
	Corpname           string `json:"CorpName"`
	Remark             string `json:"Remark"`
	Bankcontact        string `json:"BankContact"`
	Banklinkname       string `json:"BankLinkName"`
	Custcdcontact      string `json:"CustcdContact"`
	Customerid         string `json:"CustomerId"`
	Financeid          string `json:"FinanceId"`
	Cooperationyears   string `json:"CooperationYears"`
	Certificatetype    string `json:"CertificateType"`
	Intercustomerid    string `json:"InterCustomerId"`
}

// type FinancingIntention []struct {
// 	Custcdlinkposition string `json:"CustcdLinkPosition"`
// 	Custcdlinkname     string `json:"CustcdLinkName"`
// 	Certificateid      string `json:"CertificateId"`
// 	Corpname           string `json:"CorpName"`
// 	Remark             string `json:"Remark"`
// 	Bankcontact        string `json:"BankContact"`
// 	Banklinkname       string `json:"BankLinkName"`
// 	Custcdcontact      string `json:"CustcdContact"`
// 	Customerid         string `json:"CustomerId"`
// 	Financeid          string `json:"FinanceId"`
// 	Cooperationyears   string `json:"CooperationYears"`
// 	Certificatetype    string `json:"CertificateType"`
// 	Intercustomerid    string `json:"InterCustomerId"`
// }

// 推送回款账户接口
type CollectionAccount struct {
	Backaccount     string `json:"BackAccount"`
	Certificateid   string `json:"CertificateId"`
	Customerid      string `json:"CustomerId"`
	Corpname        string `json:"CorpName"`
	Lockremark      string `json:"LockRemark"`
	Certificatetype string `json:"CertificateType"`
	Intercustomerid string `json:"InterCustomerId"`
}

type SelectedInfosAndFinancingApplication struct {
	FinancingApplication FinancingIntention `json:"FinancingApplication"`
	Invoice              []InvoiceInfo      `json:"invoice"`
}

type SelectedInfoToApplication struct {
	Invoice     []InvoiceInfo        `json:"invoice"`
	HistoryInfo []TransactionHistory `json:"historyInfo"`
	PoolInfo    []EnterpoolData      `json:"poolInfo"`
}

type InvoiceInfo struct {
	CertificateID   string `json:"certificateId"`
	CustomerID      string `json:"customerId"`
	CorpName        string `json:"corpName"`
	CertificateType string `json:"certificateType"`
	InterCustomerID string `json:"interCustomerId"`
	InvoiceNotaxAmt string `json:"InvoiceNotaxAmt"`
	InvoiceCcy      string `json:"InvoiceCcy"`
	SellerName      string `json:"SellerName"`
	InvoiceType     string `json:"InvoiceType"`
	BuyerName       string `json:"BuyerName"`
	BuyerUsccode    string `json:"BuyerUsccode"`
	InvoiceDate     string `json:"InvoiceDate"`
	SellerUsccode   string `json:"SellerUsccode"`
	InvoiceCode     string `json:"InvoiceCode"`
	InvoiceNum      string `json:"InvoiceNum"`
	CheckCode       string `json:"CheckCode"`
	InvoiceAmt      string `json:"InvoiceAmt"`
}

//	type TransactionHistoryData struct {
//		CustomerGrade   string            `json:"customerGrade"`
//		CertificateType string            `json:"certificateType"`
//		InterCustomerID string            `json:"interCustomerId"`
//		CorpName        string            `json:"corpName"`
//		FinanceID       string            `json:"financeId"`
//		CertificateID   string            `json:"certificateId"`
//		CustomerID      string            `json:"customerId"`
//		UsedInfos       []Usedinfos       `json:"usedInfos"`
//		SettleInfos     []Settleinfos     `json:"settleInfos"`
//		OrderInfos      []Orderinfos      `json:"orderInfos"`
//		ReceivableInfos []Receivableinfos `json:"receivableInfos"`
//	}
// type Enterpool struct {
// 	DateTimePoint     string              `json:"dateTimePoint"`
// 	Ccy               string              `json:"ccy"`
// 	CustomerID        string              `json:"customerId"`
// 	InterCustomerID   string              `json:"interCustomerId"`
// 	ReceivableBalance string              `json:"receivableBalance"`
// 	PlanInfos         []Planinfos         `json:"planInfos"`
// 	ProviderUsedInfos []Providerusedinfos `json:"ProviderUsedInfos"`
// }
