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

const (
	Issue  = "issue"
	Modify = "modify"
)

type FrontEnd struct {
	InvoicePool                             map[string]*InvoiceInformation
	TransactionHistoryPool                  map[string]*TransactionHistory
	EnterpoolDataPool                       map[string]*EnterpoolData
	FinancingIntentionWithSelectedInfosPool map[string]*SelectedInfosAndFinancingApplication
	UpdateCollectionAccountPool             map[string]*UpdateCollectionAccount
	LockAccountPool                         map[string]*LockAccount
	ModifyFinancingWithSelectedInfosPool    map[string]*SelectedInfosAndFinancingApplication

	IssueInvoicemutex                         sync.RWMutex
	TransactionHistorymutex                   sync.RWMutex
	EnterpoolDatamutex                        sync.RWMutex
	FinancingIntentionWithSelectedInfosMutex  sync.RWMutex
	ModifyFinancingWithSelectedInfosPoolMutex sync.RWMutex
	UpdateCollectionAccountPoolMutex          sync.RWMutex
	LockAccountPoolMutex                      sync.RWMutex

	IssueInvoiceOKChan                  chan interface{}
	IssueHistoryUsedInfoOKChan          chan interface{}
	IssueHistoricalOrderInfoOKChan      chan interface{}
	IssueHistoricalSettleInfoOKChan     chan interface{}
	IssueHistoricalReceivableInfoOKChan chan interface{}
	IssueEnterPoolPlanOKChan            chan interface{}
	IssueEnterPoolUsedOKChan            chan interface{}
	UpdateAndLockAccountOKChan          chan interface{}
	LockAccountOKChan                   chan interface{}
	//提交融资意向时使用
	FinancingIntentionIssueOKChan    chan interface{}
	ModifyFinancingOKChan            chan interface{}
	ModifyInvoiceOKChan              chan interface{}
	ModifyInvoiceWhenFinancingOKChan chan interface{}
}
type PackedResponse struct {
	Success map[string]uptoChain.ResponseMessage
	Fail    map[string]uptoChain.ResponseMessage
}

func NewPackedResponse() *PackedResponse {
	return &PackedResponse{
		Success: make(map[string]uptoChain.ResponseMessage),
		Fail:    make(map[string]uptoChain.ResponseMessage),
	}
}
func NewFrontEnd() *FrontEnd {
	return &FrontEnd{
		InvoicePool:                             make(map[string]*InvoiceInformation, 0),
		TransactionHistoryPool:                  make(map[string]*TransactionHistory, 0),
		EnterpoolDataPool:                       make(map[string]*EnterpoolData, 0),
		FinancingIntentionWithSelectedInfosPool: make(map[string]*SelectedInfosAndFinancingApplication, 0),
		UpdateCollectionAccountPool:             make(map[string]*UpdateCollectionAccount, 0),
		LockAccountPool:                         make(map[string]*LockAccount, 0),
		ModifyFinancingWithSelectedInfosPool:    make(map[string]*SelectedInfosAndFinancingApplication, 0),
		IssueInvoiceOKChan:                      make(chan interface{}),
		IssueHistoryUsedInfoOKChan:              make(chan interface{}),
		IssueHistoricalOrderInfoOKChan:          make(chan interface{}),
		IssueHistoricalSettleInfoOKChan:         make(chan interface{}),
		IssueHistoricalReceivableInfoOKChan:     make(chan interface{}),
		IssueEnterPoolPlanOKChan:                make(chan interface{}),
		IssueEnterPoolUsedOKChan:                make(chan interface{}),
		UpdateAndLockAccountOKChan:              make(chan interface{}),
		LockAccountOKChan:                       make(chan interface{}),
		FinancingIntentionIssueOKChan:           make(chan interface{}),
		ModifyFinancingOKChan:                   make(chan interface{}),
		ModifyInvoiceOKChan:                     make(chan interface{}),
	}
}

// 推送发票信息接口
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
				var messages *InvoiceInformation
				if json.NewDecoder(request.Body).Decode(&messages) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					messages.UUID = id.String()
					f.IssueInvoicemutex.Lock()
					f.InvoicePool[id.String()] = messages
					f.IssueInvoicemutex.Unlock()
					<-f.IssueInvoiceOKChan
					jsonData := NewPackedResponse()
					uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.InvoiceMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.InvoiceMapLock.Unlock()
							}
							uptoChain.InvoiceMap.Delete(uuid)
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
			if checkTimeStamp(formatTimeStr) {
				var messages *TransactionHistory
				if json.NewDecoder(request.Body).Decode(&messages) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					messages.UUID = id.String()
					f.TransactionHistorymutex.Lock()
					f.TransactionHistoryPool[id.String()] = messages
					f.TransactionHistorymutex.Unlock()
					fmt.Println(messages)
					<-f.IssueHistoryUsedInfoOKChan
					<-f.IssueHistoricalOrderInfoOKChan
					<-f.IssueHistoricalReceivableInfoOKChan
					<-f.IssueHistoricalSettleInfoOKChan
					jsonData := NewPackedResponse()
					uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.HistoricalOrderMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("HistoricalOrder:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.HistoricalOrderMapLock.Unlock()
								uptoChain.HistoricalOrderMap.Delete(uuid)
							}
						}
						return true
					})
					uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.HistoricalSettleMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("HistoricalSettle:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.HistoricalSettleMapLock.Unlock()
								uptoChain.HistoricalSettleMap.Delete(uuid)
							}
						}
						return true
					})
					uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.HistoricalUsedMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("HistoricalUsed:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.HistoricalUsedMapLock.Unlock()
								uptoChain.HistoricalUsedMap.Delete(uuid)

							}
						}
						return true
					})
					uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.HistoricalReceivableMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("HistoricalReceivable:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.HistoricalReceivableMapLock.Unlock()
								uptoChain.HistoricalReceivableMap.Delete(uuid)
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
				var messages *EnterpoolData
				if json.NewDecoder(request.Body).Decode(&messages) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					messages.UUID = id.String()
					fmt.Println(".....", messages)
					f.EnterpoolDatamutex.Lock()
					f.EnterpoolDataPool[id.String()] = messages
					f.EnterpoolDatamutex.Unlock()
					<-f.IssueEnterPoolPlanOKChan
					<-f.IssueEnterPoolUsedOKChan
					jsonData := NewPackedResponse()
					uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.PoolPlanMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("PoolPlan:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.PoolPlanMapLock.Unlock()
								uptoChain.PoolPlanMap.Delete(uuid)
							}
						}
						return true
					})
					uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.PoolUsedMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("PoolUsed:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.PoolUsedMapLock.Unlock()
								uptoChain.PoolUsedMap.Delete(uuid)
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
				//处理发布情况
				var message *SelectedInfosAndFinancingApplication
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else if !VerifyInvoice(*message) {
					jsonData := wrongVerifyInvoice()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					message.UUID = id.String()
					for index := range message.Invoice {
						message.Invoice[index].FinancingID = message.FinancingApplication.Financeid
					}
					f.FinancingIntentionWithSelectedInfosMutex.Lock()
					f.FinancingIntentionWithSelectedInfosPool[id.String()] = message
					f.FinancingIntentionWithSelectedInfosMutex.Unlock()
					<-f.FinancingIntentionIssueOKChan
					<-f.ModifyInvoiceOKChan
					jsonData := NewPackedResponse()
					uptoChain.FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.FinancingApplicationIssueMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("FinancingApplication:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.FinancingApplicationIssueMapLock.Unlock()
								uptoChain.FinancingApplicationIssueMap.Delete(uuid)
							}
						}
						return true
					})
					uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.ModifyInvoiceMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("ModifyInvoice:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.ModifyInvoiceMapLock.Unlock()
								uptoChain.ModifyInvoiceMap.Delete(uuid)
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

// 修改融资意向申请接口，与所选的发票数据一同接收
func (f *FrontEnd) HandleModifyFinancingIntentionWithSelectedInfos(writer http.ResponseWriter, request *http.Request) {
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
				var message *SelectedInfosAndFinancingApplication
				if json.NewDecoder(request.Body).Decode(&message) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else if !VerifyInvoice(*message) {
					jsonData := wrongVerifyInvoice()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					message.UUID = id.String()
					for index := range message.Invoice {
						message.Invoice[index].FinancingID = message.FinancingApplication.Financeid
					}
					f.ModifyFinancingWithSelectedInfosPoolMutex.Lock()
					f.ModifyFinancingWithSelectedInfosPool[id.String()] = message
					f.ModifyFinancingWithSelectedInfosPoolMutex.Unlock()
					<-f.ModifyFinancingOKChan
					<-f.ModifyInvoiceWhenFinancingOKChan
					jsonData := NewPackedResponse()
					uptoChain.ModifyFinancingMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.ModifyFinancingMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("ModifyFinancingApplication:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.ModifyFinancingMapLock.Unlock()
								uptoChain.ModifyFinancingMap.Delete(uuid)
							}
						}
						return true
					})
					uptoChain.ModifyInvoiceWhenMFAMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.ModifyInvoiceWhenMFAMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									message.AddMessage("ModifyFinancingAndInvoice:")
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.ModifyInvoiceWhenMFAMapLock.Unlock()
								uptoChain.ModifyInvoiceWhenMFAMap.Delete(uuid)
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

// 更新并锁定回款账户接口
func (f *FrontEnd) HandleUpdateCollectionAccount(writer http.ResponseWriter, request *http.Request) {
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
				var messages *UpdateCollectionAccount
				if json.NewDecoder(request.Body).Decode(&messages) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					messages.UUID = id.String()
					f.UpdateCollectionAccountPoolMutex.Lock()
					f.UpdateCollectionAccountPool[id.String()] = messages
					f.UpdateCollectionAccountPoolMutex.Unlock()
					<-f.UpdateAndLockAccountOKChan
					jsonData := NewPackedResponse()
					uptoChain.UpdateAndLockAccountMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.UpdateAndLockAccountMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.UpdateAndLockAccountMapLock.Unlock()
							}
							uptoChain.UpdateAndLockAccountMap.Delete(uuid)
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

// 锁定回款账户接口
func (f *FrontEnd) HandleLockAccount(writer http.ResponseWriter, request *http.Request) {
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
				var messages *LockAccount
				if json.NewDecoder(request.Body).Decode(&messages) != nil {
					jsonData := wrongJsonType()
					fmt.Fprint(writer, jsonData)
				} else {
					id, err := uuid.NewUUID()
					if err != nil {
						logrus.Fatalf("newChannelMessage error: %v", err)
					}
					messages.UUID = id.String()
					f.LockAccountPoolMutex.Lock()
					f.LockAccountPool[id.String()] = messages
					f.LockAccountPoolMutex.Unlock()
					<-f.LockAccountOKChan
					jsonData := NewPackedResponse()
					uptoChain.LockAccountsMap.Range(func(key, value interface{}) bool {
						if uuid, ok := key.(string); ok {
							if uuid == id.String() {
								uptoChain.LockAccountsMapLock.Lock()
								mapping := value.(map[string]*uptoChain.ResponseMessage)
								for txHash, message := range mapping {
									if message.GetWhetherOK() {
										jsonData.Success[txHash] = *message
									} else {
										jsonData.Fail[txHash] = *message
									}
								}
								uptoChain.LockAccountsMapLock.Unlock()
							}
							uptoChain.LockAccountsMap.Delete(uuid)
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
