package server

// 特别针对发票信息的Message结构体
type packedInvoiceMessage struct {
	Uuid          string
	Header        string
	Params        string
	Cipher        []byte
	EncryptionKey []byte
}
type packedHistoricalMessage struct {
	Uuid          string
	Header        string
	Params        string
	Cipher        []byte
	EncryptionKey []byte
}
type packedPoolMessage struct {
	Uuid          string
	Header        string
	Params        string
	Cipher        []byte
	EncryptionKey []byte
}
type packedUpdateAccountMessage struct {
	Uuid          string
	FinanceID     string
	Header        string
	Cipher        []byte
	EncryptionKey []byte
	OldHash       []byte
	NewHash       []byte
}
type packedLockAccountMessage struct {
	Uuid          string
	FinanceID     string
	Header        string
	Cipher        []byte
	EncryptionKey []byte
	Signed        []byte
}
type packedFinancingMessage struct {
	Uuid          string
	Header        string
	Financingid   string
	CustomerID    string
	State         string
	Cipher        []byte
	EncryptionKey []byte
	Signed        []byte
}

// 用于修改发票的owner字段
type packedModifyInvoiceMessage struct {
	FinancingID   string
	Uuid          string
	Header        string
	Sign          string
	Cipher        []byte
	EncryptionKey []byte
}
