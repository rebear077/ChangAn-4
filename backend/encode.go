package server

import (
	receive "github.com/rebear077/changan/connApi"
)

func EncodeInvoiceInformation(list []*receive.InvoiceInformation) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	guide := 0
	for _, l := range list {
		header := l.Customerid
		for _, info := range l.Invoiceinfos {
			guide += 1
			mapping[guide] = make(map[string]string)
			tempStr := l.Certificateid + "," + l.Customerid + "," + l.Corpname + "," + l.Certificatetype + "," + l.Intercustomerid + "," + info.Invoicenotaxamt + "," + info.Invoiceccy + "," + info.Sellername + "," + info.Invoicetype + "," + info.Buyername + "," + info.Buyerusccode + "," + info.Invoicedate + "," + info.Sellerusccode + "," + info.Invoicecode + "," + info.Invoicenum + "," + info.Checkcode + "," + info.Invoiceamt
			mapping[guide][header] = tempStr
		}
	}
	return mapping
}
func EncodeTransactionHistory(list []*receive.TransactionHistory) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	for index, l := range list {
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
	return mapping
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

func EncodeFinancingIntention(list []*receive.FinancingIntention) map[int]map[string]string {
	mapping := make(map[int]map[string]string)
	index := 0
	// for _, v := range list { //发来的时候是单独一条一条, 不是以数组的形式
	for _, l := range list {
		index += 1
		mapping[index] = make(map[string]string)
		header := l.Customerid
		tempStr := l.Custcdlinkposition + "," + l.Custcdlinkname + "," + l.Certificateid + "," + l.Corpname + "," + l.Remark + "," + l.Bankcontact + "," + l.Banklinkname + "," + l.Custcdcontact + "," + l.Customerid + "," + l.Financeid + "," + l.Cooperationyears + "," + l.Certificatetype + "," + l.Intercustomerid
		mapping[index][header] = tempStr
	}

	// }
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
