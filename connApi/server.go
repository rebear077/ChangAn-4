package receive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	logloader "github.com/rebear077/changan/logs"
)

var logs = logloader.NewLog()

type FrontEnd struct {
	InvoicePool             []*InvoiceInformation
	TransactionHistoryPool  []*TransactionHistory
	EnterpoolDataPool       []*EnterpoolData
	FinancingIntentionPool  []*FinancingIntention
	CollectionAccountPool   []*CollectionAccount
	Invoicemutex            sync.RWMutex
	TransactionHistorymutex sync.RWMutex
	EnterpoolDatamutex      sync.RWMutex
	FinancingIntentionmutex sync.RWMutex
	CollectionAccountmutex  sync.RWMutex
}

func NewFrontEnd() *FrontEnd {
	return &FrontEnd{
		InvoicePool:            make([]*InvoiceInformation, 0),
		TransactionHistoryPool: make([]*TransactionHistory, 0),
		EnterpoolDataPool:      make([]*EnterpoolData, 0),
		FinancingIntentionPool: make([]*FinancingIntention, 0),
		CollectionAccountPool:  make([]*CollectionAccount, 0),
	}
}
func (f *FrontEnd) HandleInvoiceInformation(writer http.ResponseWriter, request *http.Request) {
	pubKey, err := ioutil.ReadFile("./connApi/confs/public.pem")
	if err != nil {
		logs.Info(err)
	}
	request.Header.Set("Connection", "close")
	if request.Header.Get("verify") == "SHA256withRSAVerify" {
		cipertext := request.Header.Get("apisign")
		appid := request.Header.Get("appid")
		//时间戳处理
		timestamp := request.Header.Get("timestamp")
		formatTimeStr := convertimeStamp(timestamp)
		sign := request.Header.Get("sign")
		sourcedata := appid + "&" + timestamp + "&" + sign
		res, err := rsaVerySignWithSha256([]byte(sourcedata), cipertext, pubKey)
		if err != nil {
			logs.Info(err)
		}
		if res {
			if checkTimeStamp(formatTimeStr) {
				var message InvoiceInformation
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					jsonData := sucessCode()
					f.Invoicemutex.Lock()
					// fmt.Println(message)
					f.InvoicePool = append(f.InvoicePool, &message)
					f.Invoicemutex.Unlock()
					fmt.Fprint(writer, jsonData)
				}
			} else {
				jsonData := timeExceeded()
				fmt.Fprint(writer, jsonData)
			}
		} else {
			jsonData := verySignatureFailed()
			fmt.Fprint(writer, jsonData)
		}
	} else {
		jsonData := wrongVerifyMethod()
		fmt.Fprint(writer, jsonData)
	}
}

// 推送历史交易信息接口
func (f *FrontEnd) HandleTransactionHistory(writer http.ResponseWriter, request *http.Request) {
	pubKey, err := ioutil.ReadFile("./connApi/confs/public.pem")
	if err != nil {
		logs.Info(err)
	}
	request.Header.Set("Connection", "close")
	if request.Header.Get("verify") == "SHA256withRSAVerify" {
		cipertext := request.Header.Get("apisign")
		appid := request.Header.Get("appid")
		//时间戳处理
		timestamp := request.Header.Get("timestamp")
		formatTimeStr := convertimeStamp(timestamp)
		sign := request.Header.Get("sign")
		sourcedata := appid + "&" + timestamp + "&" + sign
		res, err := rsaVerySignWithSha256([]byte(sourcedata), cipertext, pubKey)
		if err != nil {
			logs.Info(err)
		}
		if res {
			// fmt.Println("签名信息验证成功！！")
			if checkTimeStamp(formatTimeStr) {
				var message TransactionHistory
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					jsonData := sucessCode()
					f.TransactionHistorymutex.Lock()
					// fmt.Println(message)
					f.TransactionHistoryPool = append(f.TransactionHistoryPool, &message)
					f.TransactionHistorymutex.Unlock()
					fmt.Fprint(writer, jsonData)
				}
			} else {
				jsonData := timeExceeded()
				fmt.Fprint(writer, jsonData)
			}
		} else {
			jsonData := verySignatureFailed()
			fmt.Fprint(writer, jsonData)
		}
	} else {
		jsonData := wrongVerifyMethod()
		fmt.Fprint(writer, jsonData)
	}
}

// 推送入池数据接口
func (f *FrontEnd) HandleEnterpoolData(writer http.ResponseWriter, request *http.Request) {
	pubKey, err := ioutil.ReadFile("./connApi/confs/public.pem")
	if err != nil {
		logs.Info(err)
	}
	request.Header.Set("Connection", "close")
	if request.Header.Get("verify") == "SHA256withRSAVerify" {
		cipertext := request.Header.Get("apisign")
		appid := request.Header.Get("appid")
		//时间戳处理
		timestamp := request.Header.Get("timestamp")
		formatTimeStr := convertimeStamp(timestamp)
		sign := request.Header.Get("sign")
		sourcedata := appid + "&" + timestamp + "&" + sign
		res, err := rsaVerySignWithSha256([]byte(sourcedata), cipertext, pubKey)
		if err != nil {
			logs.Info(err)
		}
		if res {
			if checkTimeStamp(formatTimeStr) {
				var message EnterpoolData
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					jsonData := sucessCode()
					f.EnterpoolDatamutex.Lock()
					// fmt.Println(message)
					f.EnterpoolDataPool = append(f.EnterpoolDataPool, &message)
					f.EnterpoolDatamutex.Unlock()
					fmt.Fprint(writer, jsonData)
				}
			} else {
				jsonData := timeExceeded()
				fmt.Fprint(writer, jsonData)
			}
		} else {
			jsonData := verySignatureFailed()
			fmt.Fprint(writer, jsonData)
		}
	} else {
		jsonData := wrongVerifyMethod()
		fmt.Fprint(writer, jsonData)
	}
}

// 提交融资意向接口
func (f *FrontEnd) HandleFinancingIntention(writer http.ResponseWriter, request *http.Request) {
	pubKey, err := ioutil.ReadFile("./connApi/confs/public.pem")
	if err != nil {
		logs.Info(err)
	}
	request.Header.Set("Connection", "close")
	if request.Header.Get("verify") == "SHA256withRSAVerify" {
		cipertext := request.Header.Get("apisign")
		appid := request.Header.Get("appid")
		//时间戳处理
		timestamp := request.Header.Get("timestamp")
		formatTimeStr := convertimeStamp(timestamp)
		sign := request.Header.Get("sign")
		sourcedata := appid + "&" + timestamp + "&" + sign
		res, err := rsaVerySignWithSha256([]byte(sourcedata), cipertext, pubKey)
		if err != nil {
			logs.Info(err)
		}
		if res {
			if checkTimeStamp(formatTimeStr) {
				var message FinancingIntention
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					jsonData := sucessCode()
					f.FinancingIntentionmutex.Lock()
					f.FinancingIntentionPool = append(f.FinancingIntentionPool, &message)
					// fmt.Println(message)
					f.FinancingIntentionmutex.Unlock()
					fmt.Fprint(writer, jsonData)
				}
			} else {
				jsonData := timeExceeded()
				fmt.Fprint(writer, jsonData)
			}
		} else {
			jsonData := verySignatureFailed()
			fmt.Fprint(writer, jsonData)
		}
	} else {
		jsonData := wrongVerifyMethod()
		fmt.Fprint(writer, jsonData)
	}
}

// 推送回款账户接口
func (f *FrontEnd) HandleCollectionAccount(writer http.ResponseWriter, request *http.Request) {
	pubKey, err := ioutil.ReadFile("./connApi/confs/public.pem")
	if err != nil {
		logs.Info(err)
	}
	request.Header.Set("Connection", "close")
	if request.Header.Get("verify") == "SHA256withRSAVerify" {
		cipertext := request.Header.Get("apisign")
		appid := request.Header.Get("appid")
		//时间戳处理
		timestamp := request.Header.Get("timestamp")
		formatTimeStr := convertimeStamp(timestamp)
		sign := request.Header.Get("sign")
		sourcedata := appid + "&" + timestamp + "&" + sign
		// fmt.Println(sourcedata)
		res, err := rsaVerySignWithSha256([]byte(sourcedata), cipertext, pubKey)
		if err != nil {
			logs.Info(err)
		}
		if res {
			// fmt.Println("签名信息验证成功！！")
			if checkTimeStamp(formatTimeStr) {
				var message CollectionAccount
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					//返回成功字段
					jsonData := sucessCode()
					f.CollectionAccountmutex.Lock()
					f.CollectionAccountPool = append(f.CollectionAccountPool, &message)
					// fmt.Println(message)
					f.CollectionAccountmutex.Unlock()
					fmt.Fprint(writer, jsonData)
				}
			} else {
				jsonData := timeExceeded()
				fmt.Fprint(writer, jsonData)
			}
		} else {
			jsonData := verySignatureFailed()
			fmt.Fprint(writer, jsonData)
		}
	} else {
		jsonData := wrongVerifyMethod()
		fmt.Fprint(writer, jsonData)
	}
}

func check(err error) {
	if err != nil {
		logs.Fatalln(err)
	}
}
