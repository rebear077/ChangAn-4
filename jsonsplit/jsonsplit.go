package jsonsplit

import (
	"fmt"
	"strings"
)

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

func SplitTransactionHistory(str string) TransactionHistory {
	flag := 0
	header := ""
	usedinfos := ""
	settleinfos := ""
	orderinfos := ""
	receivableinfos := ""
	for index, val := range str {
		if index+1 >= len(str) {
			break
		}
		if flag == 0 {
			if str[index] == ',' && str[index+1] == '[' {
				flag = 1
			} else {
				header = header + string(val)
			}
		} else if flag == 1 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 2
			} else if str[index] == ']' {
				flag = 2
			} else if str[index] != '[' && str[index] != ']' {
				usedinfos = usedinfos + string(val)
			}
		} else if flag == 2 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 3
			} else if str[index] == ']' && len(settleinfos) != 0 {
				flag = 3
			} else if str[index] != '[' && str[index] != ']' {
				if len(settleinfos) == 0 && str[index] == ',' {
					continue
				} else {
					settleinfos = settleinfos + string(val)
				}
			}
		} else if flag == 3 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 4
			} else if str[index] == ']' && len(orderinfos) != 0 {
				flag = 4
			} else if str[index] != '[' && str[index] != ']' {
				if len(orderinfos) == 0 && str[index] == ',' {
					continue
				} else {
					orderinfos = orderinfos + string(val)
				}
			}
		} else if flag == 4 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 5
			} else if str[index] == ']' && len(receivableinfos) != 0 {
				flag = 5
			} else if str[index] != '[' && str[index] != ']' {
				if len(receivableinfos) == 0 && str[index] == ',' {
					continue
				} else {
					receivableinfos = receivableinfos + string(val)
				}
			}
		}
	}
	header_split := strings.Split(header, ",")
	var UsedInfos []Usedinfos
	var SettleInfos []Settleinfos
	var OrderInfos []Orderinfos
	var ReceivableInfos []Receivableinfos

	usedinfos_split := strings.Split(usedinfos, "|")
	if usedinfos_split[0] != "" {
		for i := 0; i < len(usedinfos_split); i++ {
			us := strings.Split(usedinfos_split[i], ",")
			UIfo := Usedinfos{
				us[0],
				us[1],
				us[2],
			}
			UsedInfos = append(UsedInfos, UIfo)
		}
	}

	settleinfos_split := strings.Split(settleinfos, "|")
	if settleinfos_split[0] != "" {
		for i := 0; i < len(settleinfos_split); i++ {
			st := strings.Split(settleinfos_split[i], ",")
			SIfo := Settleinfos{
				st[0],
				st[1],
				st[2],
			}
			SettleInfos = append(SettleInfos, SIfo)
		}
	}

	orderinfos_split := strings.Split(orderinfos, "|")
	if orderinfos_split[0] != "" {
		for i := 0; i < len(orderinfos_split); i++ {
			od := strings.Split(orderinfos_split[i], ",")
			OIfo := Orderinfos{
				od[0],
				od[1],
				od[2],
			}
			OrderInfos = append(OrderInfos, OIfo)
		}
	}

	receivableinfos_split := strings.Split(receivableinfos, "|")
	if receivableinfos_split[0] != "" {
		for i := 0; i < len(receivableinfos_split); i++ {
			rc := strings.Split(receivableinfos_split[i], ",")
			RIfo := Receivableinfos{
				rc[0],
				rc[1],
				rc[2],
			}
			ReceivableInfos = append(ReceivableInfos, RIfo)
		}
	}

	trsh := TransactionHistory{
		header_split[0],
		header_split[1],
		header_split[2],
		header_split[3],
		header_split[4],
		header_split[5],
		header_split[6],
		UsedInfos,
		SettleInfos,
		OrderInfos,
		ReceivableInfos,
	}
	// fmt.Println(trsh)
	//fmt.Println(HIS)
	return trsh
}

func SplitEnterpoolDataPool(str string) EnterpoolData {
	fmt.Println(str)
	flag := 0
	header := ""
	planinfos := ""
	providerusedinfos := ""
	for index, val := range str {
		if index+1 >= len(str) {
			break
		}
		if flag == 0 {
			if str[index] == ',' && str[index+1] == '[' {
				flag = 1
			} else {
				header = header + string(val)
			}
		} else if flag == 1 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 2
			} else if str[index] == ']' {
				flag = 2
			} else if str[index] != '[' && str[index] != ']' {
				planinfos = planinfos + string(val)
			}
		} else if flag == 2 {
			if str[index] == '[' && str[index+1] == ',' {
				flag = 3
			} else if str[index] == ']' && len(providerusedinfos) != 0 {
				flag = 3
			} else if str[index] != '[' && str[index] != ']' {
				if len(providerusedinfos) == 0 && str[index] == ',' {
					continue
				} else {
					providerusedinfos = providerusedinfos + string(val)
				}
			}
		}
	}
	header_split := strings.Split(header, ",")
	var PlanInfos []Planinfos
	planinfos_split := strings.Split(planinfos, "|")
	if planinfos_split[0] != "" {
		for i := 0; i < len(planinfos_split); i++ {
			pl := strings.Split(planinfos_split[i], ",")
			PLfo := Planinfos{
				pl[0],
				pl[1],
				pl[2],
			}
			PlanInfos = append(PlanInfos, PLfo)
		}
	}
	var ProviderusedInfos []Providerusedinfos
	providerusedinfos_split := strings.Split(providerusedinfos, "|")
	if providerusedinfos_split[0] != "" {
		for i := 0; i < len(providerusedinfos_split); i++ {
			pr := strings.Split(providerusedinfos_split[i], ",")
			PRfo := Providerusedinfos{
				pr[0],
				pr[1],
				pr[2],
			}
			ProviderusedInfos = append(ProviderusedInfos, PRfo)
		}
	}

	epdt := EnterpoolData{
		header_split[0],
		header_split[1],
		header_split[2],
		header_split[3],
		header_split[4],
		PlanInfos,
		ProviderusedInfos,
	}

	return epdt
}
