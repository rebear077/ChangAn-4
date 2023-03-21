package server

import (
	"github.com/rebear077/changan/jsonsplit"
)

func HistoricalInformationSlice(key string, value string, baselength int) ([]string, []string, []string, []string) {
	//fmt.Println("key: ", key)
	//fmt.Println("value: ", value)
	splitvalue := jsonsplit.SplitTransactionHistory(value)
	// fmt.Println("splitvalue: ", splitvalue)
	//最小拆分单位
	// fmt.Println(splitvalue.Orderinfos[0])
	//头部信息
	hisheader := splitvalue.Customergrade + "," + splitvalue.Certificatetype + "," + splitvalue.Intercustomerid + "," + splitvalue.Corpname + "," + splitvalue.Financeid + "," + splitvalue.Certificateid + "," + splitvalue.Customerid + ","
	//fmt.Println("hisheader: ", hisheader)
	var temp string
	//入库信息
	temp = "["
	var usedinfos []string
	for i := 0; i < len(splitvalue.Usedinfos); i++ {
		temp = temp + splitvalue.Usedinfos[i].Tradeyearmonth + "," + splitvalue.Usedinfos[i].Usedamount + "," + splitvalue.Usedinfos[i].Ccy
		if i == len(splitvalue.Usedinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			usedinfos = append(usedinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			usedinfos = append(usedinfos, temp)
			temp = "["
		}
	}

	//结算信息
	temp = "["
	var settleinfos []string
	for i := 0; i < len(splitvalue.Settleinfos); i++ {
		temp = temp + splitvalue.Settleinfos[i].Tradeyearmonth + "," + splitvalue.Settleinfos[i].Settleamount + "," + splitvalue.Settleinfos[i].Ccy
		if i == len(splitvalue.Settleinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			settleinfos = append(settleinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			settleinfos = append(settleinfos, temp)
			temp = "["
		}
	}
	//订单信息

	temp = "["
	var orderinfos []string
	for i := 0; i < len(splitvalue.Orderinfos); i++ {
		temp = temp + splitvalue.Orderinfos[i].Tradeyearmonth + "," + splitvalue.Orderinfos[i].Orderamount + "," + splitvalue.Orderinfos[i].Ccy
		if i == len(splitvalue.Orderinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			orderinfos = append(orderinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			orderinfos = append(orderinfos, temp)
			temp = "["
		}
	}
	// fmt.Println("orderinfos: ", orderinfos)

	//应收账款信息
	temp = "["
	var receivableinfos []string
	for i := 0; i < len(splitvalue.Receivableinfos); i++ {
		temp = temp + splitvalue.Receivableinfos[i].Tradeyearmonth + "," + splitvalue.Receivableinfos[i].Receivableamount + "," + splitvalue.Receivableinfos[i].Ccy
		if i == len(splitvalue.Receivableinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			receivableinfos = append(receivableinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			receivableinfos = append(receivableinfos, temp)
			temp = "["
		}
	}
	// fmt.Println("receivableinfos: ", receivableinfos)

	return usedinfos, settleinfos, orderinfos, receivableinfos
}

func PoolInformationSlice(key string, value string, baselength int) ([]string, []string) {
	splitvalue := jsonsplit.SplitEnterpoolDataPool(value)
	//头部信息
	hisheader := splitvalue.Datetimepoint + "," + splitvalue.Ccy + "," + splitvalue.Customerid + "," + splitvalue.Intercustomerid + "," + splitvalue.Receivablebalance + ","
	var temp string
	//入库信息
	temp = "["
	var planinfos []string
	for i := 0; i < len(splitvalue.Planinfos); i++ {
		temp = temp + splitvalue.Planinfos[i].Tradeyearmonth + "," + splitvalue.Planinfos[i].Planamount + "," + splitvalue.Planinfos[i].Currency
		if i == len(splitvalue.Planinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			planinfos = append(planinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			planinfos = append(planinfos, temp)
			temp = "["
		}
	}
	// fmt.Println("planinfos: ", planinfos)

	//结算信息
	temp = "["
	var providerusedinfos []string
	for i := 0; i < len(splitvalue.Providerusedinfos); i++ {
		temp = temp + splitvalue.Providerusedinfos[i].Tradeyearmonth + "," + splitvalue.Providerusedinfos[i].Usedamount + "," + splitvalue.Providerusedinfos[i].Currency
		if i == len(splitvalue.Providerusedinfos)-1 {
			// 如果遍历到最后
			temp = temp + "]"
			temp = hisheader + temp
			providerusedinfos = append(providerusedinfos, temp)
		} else if i%baselength != baselength-1 {
			temp = temp + "|"
		} else if i%baselength == baselength-1 {
			// 如果达到最小切片长度
			temp = temp + "]"
			temp = hisheader + temp
			providerusedinfos = append(providerusedinfos, temp)
			temp = "["
		}
	}
	// fmt.Println("providerusedinfos: ", providerusedinfos)

	return planinfos, providerusedinfos
}
