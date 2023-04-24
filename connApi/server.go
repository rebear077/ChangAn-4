package receive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/uuid"
	logloader "github.com/rebear077/changan/logs"
	uptoChain "github.com/rebear077/changan/tochain"
	"github.com/sirupsen/logrus"
)

var logs = logloader.NewLog()

type FrontEnd struct {
	InvoicePool                             map[string][]*InvoiceInformation
	TransactionHistoryPool                  map[string]*TransactionHistory
	EnterpoolDataPool                       []*EnterpoolData
	FinancingIntentionWithSelectedInfosPool []*SelectedInfosAndFinancingApplication
	CollectionAccountPool                   []*CollectionAccount

	IssueInvoicemutex                        sync.RWMutex
	TransactionHistorymutex                  sync.RWMutex
	EnterpoolDatamutex                       sync.RWMutex
	FinancingIntentionWithSelectedInfosMutex sync.RWMutex
	CollectionAccountmutex                   sync.RWMutex

	IssueInvoiceOKChan     chan bool
	IssueHistoryInfoOKChan chan bool
}
type PackedResponse struct {
	Success map[string]*uptoChain.ResponseMessage
	Fail    map[string]*uptoChain.ResponseMessage
}

func NewPackedResponse() *PackedResponse {
	return &PackedResponse{
		Success: make(map[string]*uptoChain.ResponseMessage),
		Fail:    make(map[string]*uptoChain.ResponseMessage),
	}
}
func NewFrontEnd() *FrontEnd {
	return &FrontEnd{
		InvoicePool:                             make(map[string][]*InvoiceInformation, 0),
		TransactionHistoryPool:                  make(map[string]*TransactionHistory, 0),
		EnterpoolDataPool:                       make([]*EnterpoolData, 0),
		FinancingIntentionWithSelectedInfosPool: make([]*SelectedInfosAndFinancingApplication, 0),
		CollectionAccountPool:                   make([]*CollectionAccount, 0),
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
				var message []*InvoiceInformation
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					message.UUID = id.String()
					f.IssueInvoicemutex.Lock()
					f.InvoicePool[id.String()] = &message
					f.IssueInvoicemutex.Unlock()
					<-f.IssueInvoiceOKChan
					jsonData := NewPackedResponse()
					uptoChain.M.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = message
									} else {
										jsonData.Fail[txHash] = message
									}
								}
							}
							uptoChain.M.Delete(uuid)
						}
						return true
					})
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
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					message.UUID = id.String()
					f.TransactionHistorymutex.Lock()
					f.TransactionHistoryPool[id.String()] = &message
					f.TransactionHistorymutex.Unlock()
					<-f.IssueHistoryInfoOKChan
					jsonData := NewPackedResponse()
					uptoChain.M.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = message
									} else {
										jsonData.Fail[txHash] = message
									}
								}
							}
						}
						return true
					})
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

// 提交融资意向接口，与所勾选的发票数据一同接收
func (f *FrontEnd) HandleFinancingIntentionWithSelectedInfos(writer http.ResponseWriter, request *http.Request) {
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
				var message SelectedInfosAndFinancingApplication
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					jsonData := sucessCode()
					f.FinancingIntentionWithSelectedInfosMutex.Lock()
					f.FinancingIntentionWithSelectedInfosPool = append(f.FinancingIntentionWithSelectedInfosPool, &message)
					f.FinancingIntentionWithSelectedInfosMutex.Unlock()
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

// // 处理选取借贷的数据
// func (f *FrontEnd) HandleSelectedToApplication(writer http.ResponseWriter, request *http.Request) {
// 	// request.Header.Set("Connection", "close")
// 	var message SelectedInfoToApplication
// 	if json.NewDecoder(request.Body).Decode(&message) != nil {
// 		jsonData := wrongJsonType()
// 		fmt.Fprint(writer, jsonData)
// 	} else {
// 		//返回成功字段
// 		fmt.Println(message)
// 		f.SelectedInfoToApplicationMutex.Lock()
// 		f.SelectedInfoToApplicationData = append(f.SelectedInfoToApplicationData, &message)
// 		f.SelectedInfoToApplicationMutex.Unlock()
// 		fmt.Println(len(f.SelectedInfoToApplicationData))
// 		select {
// 		case res := <-f.Ok:
// 			if res {
// 				jsonData := sucessCode()
// 				fmt.Fprint(writer, jsonData)
// 				return
// 			} else {
// 				jsonData := failedCode()
// 				fmt.Fprintln(writer, jsonData)
// 				return
// 			}
// 		}

//		}
//	}
func check(err error) {
	if err != nil {
		logs.Fatalln(err)
	}
}
