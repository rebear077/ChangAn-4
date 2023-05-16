package server

import (
	receive "github.com/rebear077/changan/connApi"
	"github.com/sirupsen/logrus"
)

// 打包加密贸易数据-发票信息，其中owner字段置为空
func (s *Server) PackedTradeData_InvoiceInfo(invoices map[string]*receive.InvoiceInformation) []packedInvoiceMessage {
	packedInvoices := make([]packedInvoiceMessage, 0)
	//加密打包
	for uuid, invoiceList := range invoices {
		for _, invoice := range invoiceList.Invoiceinfos {
			packedInvoice := packedInvoiceMessage{}
			tempStr := invoiceList.Certificateid + "," + invoiceList.Customerid + "," + invoiceList.Corpname + "," + invoiceList.Certificatetype + "," + invoiceList.Intercustomerid + "," + invoice.Invoicenotaxamt + "," + invoice.Invoiceccy + "," + invoice.Sellername + "," + invoice.Invoicetype + "," + invoice.Buyername + "," + invoice.Buyerusccode + "," + invoice.Invoicedate + "," + invoice.Sellerusccode + "," + invoice.Invoicecode + "," + invoice.Invoicenum + "," + invoice.Checkcode + "," + invoice.Invoiceamt
			id := invoiceList.Customerid + ":" + invoice.Invoicedate
			cipher, encryptionKey, signed, err := s.DataEncryption([]byte(tempStr))
			checkError(err)
			packedInvoice.Uuid = uuid
			packedInvoice.Header = id
			packedInvoice.Cipher = cipher
			packedInvoice.EncryptionKey = encryptionKey
			packedInvoice.Params = invoice.Invoicedate + "," + invoice.Invoicetype + "," + invoice.Invoicenum + "," + string(signed) + "," + ""
			packedInvoices = append(packedInvoices, packedInvoice)
		}
	}
	return packedInvoices
}

// 打包加密贸易数据-历史交易信息
// 加密是按照basestr+use.Tradeyearmonth + "," + use.Usedamount + "," + use.Ccy的形式
func (s *Server) PackedTradeData_HistoricalInfo(historicalInfos map[string]*receive.TransactionHistory) ([]packedHistoricalMessage, []packedHistoricalMessage, []packedHistoricalMessage, []packedHistoricalMessage) {
	historicalUsed := make([]packedHistoricalMessage, 0)
	historicalSettle := make([]packedHistoricalMessage, 0)
	historicalOrder := make([]packedHistoricalMessage, 0)
	historicalReceivable := make([]packedHistoricalMessage, 0)
	for uuid, historicalInfo := range historicalInfos {
		//ID
		header := historicalInfo.Customerid
		//common
		base := historicalInfo.Customergrade + "," + historicalInfo.Certificatetype + "," + historicalInfo.Intercustomerid + "," + historicalInfo.Corpname + "," + historicalInfo.Financeid + "," + historicalInfo.Certificateid + "," + historicalInfo.Customerid
		usedInfos := packedHistoricalMessage{}
		for _, use := range historicalInfo.Usedinfos {
			if use.Tradeyearmonth != "" {
				usedStr := use.Tradeyearmonth + "," + use.Usedamount + "," + use.Ccy
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + usedStr))
				checkError(err)
				usedInfos.Cipher = cipher
				usedInfos.Uuid = uuid
				usedInfos.EncryptionKey = encryptionKey
				usedInfos.Header = header
				usedInfos.Params = use.Tradeyearmonth + "," + historicalInfo.Financeid + "," + string(signed) + "," + ""
				historicalUsed = append(historicalUsed, usedInfos)
			}

		}
		settleInfos := packedHistoricalMessage{}
		for _, settle := range historicalInfo.Settleinfos {
			if settle.Tradeyearmonth != "" {
				settleStr := settle.Tradeyearmonth + "," + settle.Settleamount + "," + settle.Ccy
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + settleStr))
				checkError(err)
				settleInfos.Cipher = cipher
				settleInfos.Uuid = uuid
				settleInfos.Header = header
				settleInfos.EncryptionKey = encryptionKey
				settleInfos.Params = settle.Tradeyearmonth + "," + historicalInfo.Financeid + "," + string(signed) + "," + ""
				historicalSettle = append(historicalSettle, settleInfos)
			}

		}
		orderInfos := packedHistoricalMessage{}
		for _, order := range historicalInfo.Orderinfos {
			if order.Tradeyearmonth != "" {
				orderStr := order.Tradeyearmonth + "," + order.Orderamount + "," + order.Ccy
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + orderStr))
				checkError(err)
				orderInfos.Cipher = cipher
				orderInfos.Uuid = uuid
				orderInfos.Header = header
				orderInfos.EncryptionKey = encryptionKey
				orderInfos.Params = order.Tradeyearmonth + "," + historicalInfo.Financeid + "," + string(signed) + "," + ""
				historicalOrder = append(historicalOrder, orderInfos)
			}

		}
		receivableInfos := packedHistoricalMessage{}
		for _, receivable := range historicalInfo.Receivableinfos {
			if receivable.Tradeyearmonth != "" {
				receivableStr := receivable.Tradeyearmonth + "," + receivable.Receivableamount + "," + receivable.Ccy
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + receivableStr))
				checkError(err)
				receivableInfos.Cipher = cipher
				receivableInfos.Uuid = uuid
				receivableInfos.Header = header
				receivableInfos.EncryptionKey = encryptionKey
				receivableInfos.Params = receivable.Tradeyearmonth + "," + historicalInfo.Financeid + "," + string(signed) + "," + ""
				historicalReceivable = append(historicalReceivable, receivableInfos)
			}

		}
	}
	return historicalUsed, historicalSettle, historicalOrder, historicalReceivable

}

// 打包加密贸易数据-入池数据
func (s *Server) PackedTradeData_EnterPoolInfo(enterPool map[string]*receive.EnterpoolData) ([]packedPoolMessage, []packedPoolMessage) {
	poolPlan := make([]packedPoolMessage, 0)
	poolUsed := make([]packedPoolMessage, 0)
	for uuid, poolInfos := range enterPool {
		header := poolInfos.Customerid
		base := poolInfos.Datetimepoint + "," + poolInfos.Ccy + "," + poolInfos.Customerid + "," + poolInfos.Intercustomerid + "," + poolInfos.Receivablebalance
		planInfo := packedPoolMessage{}
		for _, plan := range poolInfos.Planinfos {
			if plan.Tradeyearmonth != "" {
				planStr := plan.Tradeyearmonth + "," + plan.Planamount + "," + plan.Currency
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + planStr))
				checkError(err)
				planInfo.Cipher = cipher
				planInfo.Uuid = uuid
				planInfo.EncryptionKey = encryptionKey
				planInfo.Header = header
				planInfo.Params = plan.Tradeyearmonth + "," + string(signed) + "," + ""
				poolPlan = append(poolPlan, planInfo)
			}

		}
		usedInfo := packedPoolMessage{}
		for _, used := range poolInfos.UsedInfos {
			if used.Tradeyearmonth != "" {
				usedStr := used.Tradeyearmonth + "," + used.Usedamount + "," + used.Currency
				cipher, encryptionKey, signed, err := s.DataEncryption([]byte(base + "," + usedStr))
				checkError(err)
				usedInfo.Cipher = cipher
				usedInfo.Uuid = uuid
				usedInfo.EncryptionKey = encryptionKey
				usedInfo.Header = header
				usedInfo.Params = used.Tradeyearmonth + "," + string(signed) + "," + ""
				poolUsed = append(poolUsed, usedInfo)
			}

		}

	}
	return poolPlan, poolUsed
}

// 打包加密贸易数据-回款账户信息
func (s *Server) PackedTradeData_UpdateAccountInfo(accounts map[string]*receive.UpdateCollectionAccount) []packedUpdateAccountMessage {
	collectionAccounts := make([]packedUpdateAccountMessage, 0)
	for uuid, account := range accounts {
		collectionAccount := packedUpdateAccountMessage{}
		collectionAccount.Header = account.NewAccount.Customerid
		collectionAccount.FinanceID = account.FinanceId
		tempNewStr := account.NewAccount.Backaccount + "," + account.NewAccount.Certificateid + "," + account.NewAccount.Customerid + "," + account.NewAccount.Corpname + "," + account.NewAccount.Lockremark + "," + account.NewAccount.Certificatetype + "," + account.NewAccount.Intercustomerid
		cipher, encryptionKey, newHash, err := s.DataEncryption([]byte(tempNewStr))
		checkError(err)
		collectionAccount.Cipher = cipher
		collectionAccount.EncryptionKey = encryptionKey
		collectionAccount.NewHash = newHash
		collectionAccount.Uuid = uuid
		tempOldStr := account.OldAccount.Backaccount + "," + account.OldAccount.Certificateid + "," + account.OldAccount.Customerid + "," + account.OldAccount.Corpname + "," + account.OldAccount.Lockremark + "," + account.OldAccount.Certificatetype + "," + account.OldAccount.Intercustomerid
		_, _, oldHash, err := s.DataEncryption([]byte(tempOldStr))
		checkError(err)
		collectionAccount.OldHash = oldHash
		collectionAccounts = append(collectionAccounts, collectionAccount)
	}
	return collectionAccounts
}
func (s *Server) PackedTradeData_LockAccountInfo(accounts map[string]*receive.LockAccount) []packedLockAccountMessage {
	collectionAccounts := make([]packedLockAccountMessage, 0)
	for uuid, account := range accounts {
		collectionAccount := packedLockAccountMessage{}
		collectionAccount.Header = account.Customerid
		collectionAccount.FinanceID = account.FinanceId
		tempStr := account.Backaccount + "," + account.Certificateid + "," + account.Customerid + "," + account.Corpname + "," + account.Lockremark + "," + account.Certificatetype + "," + account.Intercustomerid
		cipher, encryptionKey, signed, err := s.DataEncryption([]byte(tempStr))
		checkError(err)
		collectionAccount.Cipher = cipher
		collectionAccount.EncryptionKey = encryptionKey
		collectionAccount.Signed = signed
		collectionAccount.Uuid = uuid
		collectionAccounts = append(collectionAccounts, collectionAccount)
	}
	return collectionAccounts
}

// 打包融资意向申请和需要修改的发票信息
func (s *Server) PackedApplicationAndModifyInvoiceInfos(applications map[string]*receive.SelectedInfosAndFinancingApplication, state string) ([]packedFinancingMessage, []packedModifyInvoiceMessage) {
	financingApplications := make([]packedFinancingMessage, 0)
	modifyInvoices := make([]packedModifyInvoiceMessage, 0)
	for uuid, application := range applications {
		financingApplication := packedFinancingMessage{}
		tempApplicationStr := application.FinancingApplication.Custcdlinkposition + "," + application.FinancingApplication.Custcdlinkname + "," + application.FinancingApplication.Certificateid + "," + application.FinancingApplication.Corpname + "," + application.FinancingApplication.Remark + "," + application.FinancingApplication.Bankcontact + "," + application.FinancingApplication.Banklinkname + "," + application.FinancingApplication.Custcdcontact + "," + application.FinancingApplication.Customerid + "," + application.FinancingApplication.Financeid + "," + application.FinancingApplication.Cooperationyears + "," + application.FinancingApplication.Certificatetype + "," + application.FinancingApplication.Intercustomerid
		cipher, encryptionKey, signed, err := s.DataEncryption([]byte(tempApplicationStr))
		checkError(err)
		financingApplication.CustomerID = application.FinancingApplication.Customerid
		financingApplication.Cipher = cipher
		financingApplication.EncryptionKey = encryptionKey
		financingApplication.Uuid = uuid
		financingApplication.Header = application.FinancingApplication.Financeid
		financingApplication.Signed = signed
		financingApplication.Financingid = application.FinancingApplication.Financeid
		financingApplication.State = state
		financingApplications = append(financingApplications, financingApplication)
		for _, invoice := range application.Invoice {
			invoiceInfo := packedModifyInvoiceMessage{}
			invoiceInfo.Header = invoice.CustomerID + ":" + invoice.InvoiceDate
			tempInvoiceStr := invoice.CertificateID + "," + invoice.CustomerID + "," + invoice.CorpName + "," + invoice.CertificateType + "," + invoice.InterCustomerID + "," + invoice.InvoiceNotaxAmt + "," + invoice.InvoiceCcy + "," + invoice.SellerName + "," + invoice.InvoiceType + "," + invoice.BuyerName + "," + invoice.BuyerUsccode + "," + invoice.InvoiceDate + "," + invoice.SellerUsccode + "," + invoice.InvoiceCode + "," + invoice.InvoiceNum + "," + invoice.CheckCode + "," + invoice.InvoiceAmt
			cipher, encryptionKey, signed, err := s.DataEncryption([]byte(tempInvoiceStr))
			checkError(err)
			invoiceInfo.Cipher = cipher
			invoiceInfo.EncryptionKey = encryptionKey
			invoiceInfo.FinancingID = invoice.FinancingID
			invoiceInfo.Uuid = uuid
			invoiceInfo.Sign = string(signed)
			modifyInvoices = append(modifyInvoices, invoiceInfo)
		}
	}
	return financingApplications, modifyInvoices
}

func checkError(err error) {
	if err != nil {
		logrus.Fatalln("数据加密失败 失败信息为:", err)
	}
}
