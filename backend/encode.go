package server

import (
	receive "github.com/rebear077/changan/connApi"
)

// 接收勾选的数据作为融资的凭证
//
//	func PackInfo(message receive.SelectedInfoToApplication) (map[int]map[string]string, map[int]map[string]string, map[int]map[string]string) {
//		invoice := encodeInvoiceInfo(message.Invoice)
//		history := encodeTransactionHistory(message.HistoryInfo)
//		pool := encodeEnterpoolData(message.PoolInfo)
//		return invoice, history, pool
//	}
func HandleFinancingIntentionAndSelectedInfos(message []*receive.SelectedInfosAndFinancingApplication) map[int]map[string]string {

}

// func encodeEnterpoolData(list []receive.EnterpoolData) map[int]map[string]string {
// 	mapping := make(map[int]map[string]string)
// 	for index, l := range list {
// 		mapping[index] = make(map[string]string)
// 		header := l.Customerid
// 		baseStr := l.Datetimepoint + "," + l.Ccy + "," + l.Customerid + "," + l.Intercustomerid + "," + l.Receivablebalance
// 		var planinfos string
// 		planinfos = "["
// 		for n, p := range l.Planinfos {
// 			planinfos += p.Tradeyearmonth + "," + p.Planamount + "," + p.Currency
// 			if n != len(l.Planinfos)-1 {
// 				planinfos += "|"
// 			} else {
// 				planinfos += "]"
// 			}
// 		}
// 		var usedinfos string
// 		usedinfos = "["
// 		for n, p := range l.Providerusedinfos {
// 			usedinfos += p.Tradeyearmonth + "," + p.Usedamount + "," + p.Currency
// 			if n != len(l.Providerusedinfos)-1 {
// 				usedinfos += "|"
// 			} else {
// 				usedinfos += "]"
// 			}
// 		}
// 		tempStr := baseStr + "," + planinfos + "," + usedinfos
// 		mapping[index][header] = tempStr
// 	}
// 	return mapping
// }
// func encodeTransactionHistory(list []receive.TransactionHistory) map[int]map[string]string {
// 	mapping := make(map[int]map[string]string)
// 	for index, l := range list {
// 		mapping[index] = make(map[string]string)
// 		header := l.Customerid
// 		baseStr := l.Customergrade + "," + l.Certificatetype + "," + l.Intercustomerid + "," + l.Corpname + "," + l.Financeid + "," + l.Certificateid + "," + l.Customerid
// 		var usedinfos string
// 		usedinfos = "["
// 		for n, u := range l.Usedinfos {
// 			usedinfos += u.Tradeyearmonth + "," + u.Usedamount + "," + u.Ccy
// 			if n != len(l.Usedinfos)-1 {
// 				usedinfos += "|"
// 			} else {
// 				usedinfos += "]"
// 			}
// 		}
// 		var settleinfos string
// 		settleinfos = "["
// 		for n, s := range l.Settleinfos {
// 			settleinfos += s.Tradeyearmonth + "," + s.Settleamount + "," + s.Ccy
// 			if n != len(l.Settleinfos)-1 {
// 				settleinfos += "|"
// 			} else {
// 				settleinfos += "]"
// 			}
// 		}
// 		var orderinfos string
// 		orderinfos = "["
// 		for n, o := range l.Orderinfos {
// 			orderinfos += o.Tradeyearmonth + "," + o.Orderamount + "," + o.Ccy
// 			if n != len(l.Orderinfos)-1 {
// 				orderinfos += "|"
// 			} else {
// 				orderinfos += "]"
// 			}
// 		}
// 		var receivableinfos string
// 		receivableinfos = "["
// 		for n, r := range l.Receivableinfos {
// 			receivableinfos += r.Tradeyearmonth + "," + r.Receivableamount + "," + r.Ccy
// 			if n != len(l.Receivableinfos)-1 {
// 				receivableinfos += "|"
// 			} else {
// 				receivableinfos += "]"
// 			}
// 		}
// 		tempStr := baseStr + "," + usedinfos + "," + settleinfos + "," + orderinfos + "," + receivableinfos
// 		mapping[index][header] = tempStr
// 	}
// 	return mapping
// }

func encodeInvoiceInfo(list []*receive.InvoiceInfo) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	guide := 0
	for _, info := range list {
		guide += 1
		header := info.CustomerID
		tempStr := info.CertificateID + "," + info.CustomerID + "," + info.CorpName + "," + info.CertificateType + "," + info.InterCustomerID + "," + info.InvoiceNotaxAmt + "," + info.InvoiceCcy + "," + info.SellerName + "," + info.InvoiceType + "," + info.BuyerName + "," + info.BuyerUsccode + "," + info.InvoiceDate + "," + info.SellerUsccode + "," + info.InvoiceCode + "," + info.InvoiceNum + "," + info.CheckCode + "," + info.InvoiceAmt
		mapping[guide] = make(map[string]string)
		mapping[guide][header] = tempStr
	}
	return mapping
}

func EncodeInvoiceInformation(list map[string][]*receive.InvoiceInformation) map[string]map[int]map[string]string {
	resMap := make(map[string]map[int]map[string]string)
	mapping := make(map[int]map[string]string)
	guide := 0
	for UUID, invoices := range list {
		for _, invoice := range invoices {
			header := invoice.Customerid
			for _, info := range invoice.Invoiceinfos {
				id := header + ":" + info.Invoicedate
				guide += 1
				mapping[guide] = make(map[string]string)
				tempStr := invoice.Certificateid + "," + invoice.Customerid + "," + invoice.Corpname + "," + invoice.Certificatetype + "," + invoice.Intercustomerid + "," + info.Invoicenotaxamt + "," + info.Invoiceccy + "," + info.Sellername + "," + info.Invoicetype + "," + info.Buyername + "," + info.Buyerusccode + "," + info.Invoicedate + "," + info.Sellerusccode + "," + info.Invoicecode + "," + info.Invoicenum + "," + info.Checkcode + "," + info.Invoiceamt
				mapping[guide][id] = tempStr
			}
		}
		resMap[UUID] = mapping
	}
	return resMap
}
func EncodeTransactionHistory(list map[string][]*receive.TransactionHistory) map[string]map[int]map[string]string {
	resMap := make(map[string]map[int]map[string]string)
	mapping := make(map[int]map[string]string)
	for UUID, historyInfos := range list {
		for index, l := range historyInfos {
			mapping[index] = make(map[string]string)
			header := l.Customerid
			baseStr := l.Customergrade + "," + l.Certificatetype + "," + l.Intercustomerid + "," + l.Corpname + "," + l.Financeid + "," + l.Certificateid + "," + l.Customerid
			var usedinfos string
			usedinfos = "["
			for n, u := range l.Usedinfos {
				usedinfos += u.Tradeyearmonth + "," + u.Usedamount + "," + u.Ccy
				if n != len(l.Usedinfos)-1 {
					usedinfos += "|"
				} else {
					usedinfos += "]"
				}
			}
			var settleinfos string
			settleinfos = "["
			for n, s := range l.Settleinfos {
				settleinfos += s.Tradeyearmonth + "," + s.Settleamount + "," + s.Ccy
				if n != len(l.Settleinfos)-1 {
					settleinfos += "|"
				} else {
					settleinfos += "]"
				}
			}
			var orderinfos string
			orderinfos = "["
			for n, o := range l.Orderinfos {
				orderinfos += o.Tradeyearmonth + "," + o.Orderamount + "," + o.Ccy
				if n != len(l.Orderinfos)-1 {
					orderinfos += "|"
				} else {
					orderinfos += "]"
				}
			}
			var receivableinfos string
			receivableinfos = "["
			for n, r := range l.Receivableinfos {
				receivableinfos += r.Tradeyearmonth + "," + r.Receivableamount + "," + r.Ccy
				if n != len(l.Receivableinfos)-1 {
					receivableinfos += "|"
				} else {
					receivableinfos += "]"
				}
			}
			tempStr := baseStr + "," + usedinfos + "," + settleinfos + "," + orderinfos + "," + receivableinfos
			mapping[index][header] = tempStr
		}
		resMap[UUID] = mapping
	}

	return resMap
}
func EncodeEnterpoolData(list []*receive.EnterpoolData) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	for index, l := range list {
		mapping[index] = make(map[string]string)
		header := l.Customerid
		baseStr := l.Datetimepoint + "," + l.Ccy + "," + l.Customerid + "," + l.Intercustomerid + "," + l.Receivablebalance
		var planinfos string
		planinfos = "["
		for n, p := range l.Planinfos {
			planinfos += p.Tradeyearmonth + "," + p.Planamount + "," + p.Currency
			if n != len(l.Planinfos)-1 {
				planinfos += "|"
			} else {
				planinfos += "]"
			}
		}
		var usedinfos string
		usedinfos = "["
		for n, p := range l.Providerusedinfos {
			usedinfos += p.Tradeyearmonth + "," + p.Usedamount + "," + p.Currency
			if n != len(l.Providerusedinfos)-1 {
				usedinfos += "|"
			} else {
				usedinfos += "]"
			}
		}
		tempStr := baseStr + "," + planinfos + "," + usedinfos
		mapping[index][header] = tempStr
	}
	return mapping
}

func EncodeFinancingIntention(l receive.FinancingIntention) map[int]map[string]string {
	mapping := make(map[int]map[string]string)

	mapping[0] = make(map[string]string)
	header := l.Customerid
	tempStr := l.Custcdlinkposition + "," + l.Custcdlinkname + "," + l.Certificateid + "," + l.Corpname + "," + l.Remark + "," + l.Bankcontact + "," + l.Banklinkname + "," + l.Custcdcontact + "," + l.Customerid + "," + l.Financeid + "," + l.Cooperationyears + "," + l.Certificatetype + "," + l.Intercustomerid
	mapping[0][header] = tempStr
	return mapping
}
func EncodeCollectionAccount(list []*receive.CollectionAccount) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	index := 0
	// for _, v := range list {
	for _, l := range list {
		index += 1
		mapping[index] = make(map[string]string)
		header := l.Customerid
		tempStr := l.Backaccount + "," + l.Certificateid + "," + l.Customerid + "," + l.Corpname + "," + l.Lockremark + "," + l.Certificatetype + "," + l.Intercustomerid
		mapping[index][header] = tempStr
	}
	// }
	return mapping
}
