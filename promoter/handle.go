package promote

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	server "github.com/rebear077/changan/backend"
	chainloader "github.com/rebear077/changan/chaininfos"
	receive "github.com/rebear077/changan/connApi"
	logloader "github.com/rebear077/changan/logs"
	uptoChain "github.com/rebear077/changan/tochain"
)

const (
	PoolPlanInfos             = "poolPlan"
	PoolUsedInfos             = "poolUsed"
	HistoricalOrderInfos      = "hisOrder"
	HistoricalReceivableInfos = "hisReceivable"
	HistoricalSettleInfos     = "hisSettle"
	HistoricalUsedInfos       = "hisUsed"
)

var logs = logloader.NewLog()

type Promoter struct {
	server        *server.Server
	DataApi       *receive.FrontEnd
	monitor       *server.Monitor
	encryptedPool *Pools
	loader        *logloader.Loader
	chaininfo     *chainloader.ChainInfo
}

func NewPromoter() *Promoter {
	ser := server.NewServer()
	api := receive.NewFrontEnd()
	monitor := server.NewMonitor()
	pool := NewPools()
	lder := logloader.NewLoader()
	chainld := chainloader.NewChainInfo()
	return &Promoter{
		server:        ser,
		DataApi:       api,
		monitor:       monitor,
		encryptedPool: pool,
		loader:        lder,
		chaininfo:     chainld,
	}
}

func (p *Promoter) Start() {
	// logrus.Infoln("开始运行")
	// go p.loader.Start()
	// go p.chaininfo.Start()
	logs.Infoln("开始运行")
	go p.monitor.Start()
	for {
		if p.monitor.VerifyChainStatus() {
			p.InvoiceInfoHandler()
			p.SupplierFinancingApplicationInfoWithSelectedInfosHandler()
			p.HistoricalInfoHandler()
			p.PushPaymentAccountsInfoHandler()
			p.PoolInfoHandler()
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func (p *Promoter) InvoiceInfoHandler() {
	if len(p.DataApi.InvoicePool) != 0 {
		logs.Infoln("开始同步发票信息")
		var wg sync.WaitGroup
		invoices := make(map[string]*receive.InvoiceInformation, 0)
		p.DataApi.IssueInvoicemutex.Lock()
		for uuid := range p.DataApi.InvoicePool {
			invoices[uuid] = p.DataApi.InvoicePool[uuid]
			delete(p.DataApi.InvoicePool, uuid)
		}
		p.DataApi.IssueInvoicemutex.Unlock()
		invoiceMapping := server.EncodeInvoiceInformation(invoices)
		for uuid, invoices := range invoiceMapping {
			for _, invoice := range invoices {
				for id, info := range invoice {
					wg.Add(1)
					tempheader := id
					tempinfo := info
					UUID := uuid
					go func(UUID, tempheader, tempinfo string) {
						p.packInvoiceInfo(UUID, tempheader, tempinfo, "fast", "invoice")
						wg.Done()
					}(UUID, tempheader, tempinfo)
				}
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("invoice", "fast")
		for _, message := range messages {
			temp, _ := message.(packedInvoiceMessage)
			err := p.server.IssueInvoiceInformation(temp.uuid, temp.header, temp.params, temp.cipher, temp.encryptionKey)
			if err != nil {
				logs.Errorln("发票信息上链失败:", temp.header, "失败信息为:", err)
			}
		}
		for {
			counter := 0
			uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
				uptoChain.InvoiceMapLock.Lock()
				mapping := value.(map[string]*uptoChain.ResponseMessage)
				counter += len(mapping)
				for _, message := range mapping {
					if !message.GetWhetherOK() {
						counter = 0
						break
					}
				}
				uptoChain.InvoiceMapLock.Unlock()
				return true
			})
			// fmt.Println(counter)
			if counter == len(messages) {
				p.DataApi.IssueInvoiceOKChan <- struct{}{}
				for {
					flag := 0
					uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
						if key != nil {
							flag++
							return false
						}
						return true
					})
					if flag == 0 {
						break
					}
				}
				break
			}
		}
	}
}

func (p *Promoter) HistoricalInfoHandler() {
	if len(p.DataApi.TransactionHistoryPool) != 0 {
		logs.Infoln("开始历史交易信息")
		var wg sync.WaitGroup
		hisinfos := make(map[string]*receive.TransactionHistory, 0)
		p.DataApi.TransactionHistorymutex.Lock()
		for uuid := range p.DataApi.TransactionHistoryPool {
			hisinfos[uuid] = p.DataApi.TransactionHistoryPool[uuid]
			delete(p.DataApi.TransactionHistoryPool, uuid)
		}
		p.DataApi.TransactionHistorymutex.Unlock()
		mapping := server.EncodeTransactionHistory(hisinfos)
		for UUID, historyInfos := range mapping {
			for index := range historyInfos {
				for header, info := range historyInfos[index] {
					tempheader := header
					tempinfo := info
					wg.Add(1)
					go func(UUID, tempheader, tempinfo string) {
						usedvalue, settlevalue, ordervalue, receivablevalue := server.HistoricalInformationSlice(tempheader, tempinfo, 1)
						p.packHistoricalInfos(UUID, tempheader, usedvalue, "fast", "historicalUsed")
						p.packHistoricalInfos(UUID, tempheader, settlevalue, "fast", "historicalSettle")
						p.packHistoricalInfos(UUID, tempheader, ordervalue, "fast", "historicalOrder")
						p.packHistoricalInfos(UUID, tempheader, receivablevalue, "fast", "historicalReceivable")
						wg.Done()
					}(UUID, tempheader, tempinfo)
				}
			}
		}
		wg.Wait()
		hisUsedMessage := p.encryptedPool.QueryMessages("historicalUsed", "fast")
		hisSettleMessage := p.encryptedPool.QueryMessages("historicalSettle", "fast")
		hisOrderMessage := p.encryptedPool.QueryMessages("historicalOrder", "fast")
		hisReceivableMessage := p.encryptedPool.QueryMessages("historicalReceivable", "fast")
		for _, message := range hisUsedMessage {
			tempUsed, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalUsedInformation(tempUsed.uuid, tempUsed.header, tempUsed.params, tempUsed.cipher, tempUsed.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}
		}
		for _, message := range hisSettleMessage {
			tempSettle, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalSettleInformation(tempSettle.uuid, tempSettle.header, tempSettle.params, tempSettle.cipher, tempSettle.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempSettle.header, "失败信息为:", err)
			}
		}
		for _, message := range hisOrderMessage {
			tempOrder, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalOrderInformation(tempOrder.uuid, tempOrder.header, tempOrder.params, tempOrder.cipher, tempOrder.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempOrder.header, "失败信息为:", err)
			}
		}
		for _, message := range hisReceivableMessage {
			tempReceivable, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalReceivableInformation(tempReceivable.uuid, tempReceivable.header, tempReceivable.params, tempReceivable.cipher, tempReceivable.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempReceivable.header, "失败信息为:", err)
			}
		}
		wg.Add(4)
		//order
		go func() {
			for {
				counter := 0
				uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
					uptoChain.HistoricalOrderMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.HistoricalOrderMapLock.Unlock()
					return true
				})
				if counter == len(hisOrderMessage) {
					p.DataApi.IssueHistoricalOrderInfoOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		//settle
		go func() {
			for {
				counter := 0
				uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
					uptoChain.HistoricalSettleMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.HistoricalSettleMapLock.Unlock()
					return true
				})
				if counter == len(hisOrderMessage) {
					p.DataApi.IssueHistoricalSettleInfoOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		//used
		go func() {
			for {
				counter := 0
				uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
					uptoChain.HistoricalUsedMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.HistoricalUsedMapLock.Unlock()
					return true
				})
				if counter == len(hisOrderMessage) {
					p.DataApi.IssueHistoryUsedInfoOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		//receivable
		go func() {
			for {
				counter := 0
				uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
					uptoChain.HistoricalReceivableMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.HistoricalReceivableMapLock.Unlock()
					return true
				})
				if counter == len(hisOrderMessage) {
					p.DataApi.IssueHistoricalReceivableInfoOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		logs.Println("退出")
	}
}
func (p *Promoter) PoolInfoHandler() {
	if len(p.DataApi.EnterpoolDataPool) != 0 {
		// logrus.Infoln("开始入池数据信息")
		logs.Infoln("开始入池数据信息")
		var wg sync.WaitGroup
		poolInfos := make(map[string]*receive.EnterpoolData, 0)
		p.DataApi.EnterpoolDatamutex.Lock()
		for uuid := range p.DataApi.EnterpoolDataPool {
			poolInfos[uuid] = p.DataApi.EnterpoolDataPool[uuid]
			delete(p.DataApi.EnterpoolDataPool, uuid)
		}
		p.DataApi.EnterpoolDatamutex.Unlock()
		mapping := server.EncodeEnterpoolData(poolInfos)
		for UUID, poolInfos := range mapping {
			for index := range poolInfos {
				for header, info := range poolInfos[index] {
					tempheader := header
					tempinfo := info
					wg.Add(1)
					go func(UUID, tempheader, tempinfo string) {
						var wwg sync.WaitGroup
						planvalue, providerusedvalue := server.PoolInformationSlice(tempheader, tempinfo, 1)
						wwg.Add(2)
						go func(UUID, tempheader string, planvalue []string) {
							p.packPoolInfos(UUID, tempheader, planvalue, "fast", "poolPlan")
							wwg.Done()
						}(UUID, tempheader, planvalue)
						go func(UUID, tempheader string, providerusedvalue []string) {
							p.packPoolInfos(UUID, tempheader, providerusedvalue, "fast", "poolUsed")
							wwg.Done()
						}(UUID, tempheader, providerusedvalue)
						wwg.Wait()
						wg.Done()
					}(UUID, tempheader, tempinfo)
				}
			}
		}
		wg.Wait()
		planMessages := p.encryptedPool.QueryMessages("poolPlan", "fast")
		usedMessages := p.encryptedPool.QueryMessages("poolUsed", "fast")
		for _, message := range planMessages {
			tempPlan, _ := message.(packedPoolMessage)
			err := p.server.IssuePoolPlanInformation(tempPlan.uuid, tempPlan.header, tempPlan.params, tempPlan.cipher, tempPlan.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempPlan.header, "失败信息为:", err)
			}
		}
		for _, message := range usedMessages {
			tempUsed, _ := message.(packedPoolMessage)
			err := p.server.IssuePoolUsedInformation(tempUsed.uuid, tempUsed.header, tempUsed.params, tempUsed.cipher, tempUsed.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}

		}
		wg.Add(2)
		go func() {
			for {
				counter := 0
				uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
					uptoChain.PoolPlanMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.PoolPlanMapLock.Unlock()
					return true
				})
				if counter == len(planMessages) {
					p.DataApi.IssueEnterPoolPlanOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for {
				counter := 0
				uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
					uptoChain.PoolUsedMapLock.Lock()
					mapping := value.(map[string]*uptoChain.ResponseMessage)
					counter += len(mapping)
					for _, message := range mapping {
						if !message.GetWhetherOK() {
							counter = 0
							break
						}
					}
					uptoChain.PoolUsedMapLock.Unlock()
					return true
				})
				if counter == len(planMessages) {
					p.DataApi.IssueEnterPoolUsedOKChan <- struct{}{}
					for {
						flag := 0
						uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
							if key != nil {
								flag++
								return false
							}
							return true
						})
						if flag == 0 {
							break
						}
					}
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		logs.Println("退出")
	}
}
func (p *Promoter) ModifyInvoiceInfoHandler(invoices map[string]map[string]map[int]map[string]string) {
	var wg sync.WaitGroup
	for uuid, invoicewithID := range invoices {
		for financingID, infos := range invoicewithID {
			for index := range infos {
				for id, info := range infos[index] {
					wg.Add(1)
					tempheader := id
					tempinfo := info
					UUID := uuid
					go func(financingID, UUID, tempheader, tempinfo string) {
						p.packModifyInvoiceInfo(financingID, UUID, tempheader, tempinfo, "fast", "modifyinvoice")
						wg.Done()
					}(financingID, UUID, tempheader, tempinfo)
				}
			}
		}
	}
	wg.Wait()
	messages := p.encryptedPool.QueryMessages("modifyinvoice", "fast")
	for _, message := range messages {
		temp, _ := message.(packedModifyInvoiceMessage)
		err := p.server.VerifyAndUpdateInvoiceInformation(temp.uuid, temp.header, temp.sign, temp.financingID)
		if err != nil {
			logs.Errorln("发票信息上链失败:", temp.header, "失败信息为:", err)
		}
	}
	for {
		counter := 0
		uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyInvoiceMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if !message.GetWhetherOK() {
					counter = 0
					break
				}
			}
			uptoChain.ModifyInvoiceMapLock.Unlock()
			return true
		})
		if counter == len(messages) {
			p.DataApi.ModifyInvoiceOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
}
func (p *Promoter) SupplierFinancingApplicationInfoWithSelectedInfosHandler() {
	if len(p.DataApi.FinancingIntentionWithSelectedInfosPool) != 0 {
		logs.Infoln("开始同步融资意向请求信息")
		var wg sync.WaitGroup
		finintensWithSelectedInfos := make(map[string]*receive.SelectedInfosAndFinancingApplication, 0)
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Lock()
		for uuid := range p.DataApi.FinancingIntentionWithSelectedInfosPool {
			finintensWithSelectedInfos[uuid] = p.DataApi.FinancingIntentionWithSelectedInfosPool[uuid]
			delete(p.DataApi.FinancingIntentionWithSelectedInfosPool, uuid)
		}
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Unlock()
		financingInfo, Invoices := server.HandleFinancingIntentionAndSelectedInfos(finintensWithSelectedInfos)
		go p.ModifyInvoiceInfoHandler(Invoices)
		for UUID := range financingInfo {
			for header, info := range financingInfo[UUID] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(UUID, tempheader, tempinfo string) {
					p.packFinancingInfo(UUID, tempheader, tempinfo, "fast", "application")
					wg.Done()
				}(UUID, tempheader, tempinfo)
			}
		}

		wg.Wait()
		messages := p.encryptedPool.QueryMessages("application", "fast")
		for _, message := range messages {
			temp, _ := message.(packedFinancingMessage)
			err := p.server.IssueSupplierFinancingApplication(temp.uuid, temp.header, temp.financingid, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				logs.Errorln("融资意向请求上链失败,", "失败信息为:", err)
			}
		}
		for {
			counter := 0
			uptoChain.FinancingApplicationMap.Range(func(key, value interface{}) bool {
				uptoChain.FinancingApplicationMapLock.Lock()
				mapping := value.(map[string]*uptoChain.ResponseMessage)
				counter += len(mapping)
				for _, message := range mapping {
					if !message.GetWhetherOK() {
						counter = 0
						break
					}
				}
				uptoChain.FinancingApplicationMapLock.Unlock()
				return true
			})
			if counter == len(messages) {
				p.DataApi.FinancingIntentionOKChan <- struct{}{}
				for {
					flag := 0
					uptoChain.FinancingApplicationMap.Range(func(key, value interface{}) bool {
						if key != nil {
							flag++
							return false
						}
						return true
					})
					if flag == 0 {
						break
					}
				}
				break
			}
		}
	}
}

func (p *Promoter) PushPaymentAccountsInfoHandler() {
	if len(p.DataApi.CollectionAccountPool) != 0 {
		logs.Infoln("开始同步回款信息")
		var wg sync.WaitGroup
		payinfos := make(map[string]*receive.CollectionAccount, 0)
		p.DataApi.CollectionAccountmutex.Lock()
		for uuid := range p.DataApi.CollectionAccountPool {
			payinfos[uuid] = p.DataApi.CollectionAccountPool[uuid]
			delete(p.DataApi.EnterpoolDataPool, uuid)
		}
		p.DataApi.CollectionAccountmutex.Unlock()
		mapping := server.EncodeCollectionAccount(payinfos)
		for UUID, accounts := range mapping {
			for index := range accounts {
				for header, info := range accounts[index] {
					wg.Add(1)
					tempheader := header
					tempinfo := info
					go func(UUID, tempheader, tempinfo string) {
						p.packInfo(UUID, tempheader, tempinfo, "fast", "payment")
						wg.Done()
					}(UUID, tempheader, tempinfo)
				}
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("payment", "fast")
		for _, message := range messages {
			temp, ok := message.(packedMessage)
			if !ok {
				fmt.Println("errorerror")
			}
			err := p.server.UpdatePushPaymentAccount(temp.uuid, temp.header, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		for {
			counter := 0
			uptoChain.CollectionAccountMap.Range(func(key, value interface{}) bool {
				uptoChain.CollectionAccountMapLock.Lock()
				mapping := value.(map[string]*uptoChain.ResponseMessage)
				counter += len(mapping)
				for _, message := range mapping {
					if !message.GetWhetherOK() {
						counter = 0
						break
					}
				}
				uptoChain.CollectionAccountMapLock.Unlock()
				return true
			})
			if counter == len(messages) {
				p.DataApi.ModifyAccountOKChan <- struct{}{}
				for {
					flag := 0
					uptoChain.CollectionAccountMap.Range(func(key, value interface{}) bool {
						if key != nil {
							flag++
							return false
						}
						return true
					})
					if flag == 0 {
						break
					}
				}
				break
			}
		}
	}
}

func (p *Promoter) packInfo(uuid, header, info, poolType, method string) {
	cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(info))
	if err != nil {
		// logrus.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
		logs.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
	}
	temp := packedMessage{}
	temp.cipher = cipher
	temp.encryptionKey = encryptionKey
	temp.signed = signed
	temp.header = header
	temp.uuid = uuid
	p.encryptedPool.Insert(temp, method, poolType)
}

// 针对发票信息的packInfo
// 加密后存入缓存池
func (p *Promoter) packInvoiceInfo(UUID string, header string, info string, poolType string, method string) {
	cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(info))
	if err != nil {
		// logrus.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
		logs.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
	}
	//info是发票信息的字符串形式，各个参数之间用逗号分割
	fields := strings.Split(info, ",")
	temp := packedInvoiceMessage{}
	//参数11是开票日期，参数8是发票类型，参数14是发票号码
	temp.uuid = UUID
	temp.params = fields[11] + "," + fields[8] + "," + fields[14] + "," + string(signed) + "," + ""
	temp.cipher = cipher
	temp.encryptionKey = encryptionKey
	temp.header = header
	p.encryptedPool.InsertInvoice(temp, method, poolType)
}
func (p *Promoter) packModifyInvoiceInfo(finangcingID, UUID string, header string, info string, poolType string, method string) {
	cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(info))
	if err != nil {
		// logrus.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
		logs.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
	}
	//info是发票信息的字符串形式，各个参数之间用逗号分割
	// fields := strings.Split(info, ",")
	temp := packedModifyInvoiceMessage{}
	//参数11是开票日期，参数8是发票类型，参数14是发票号码
	temp.uuid = UUID
	temp.sign = string(signed)
	// temp.params = fields[11] + "," + fields[8] + "," + fields[14] + "," + string(signed) + "," + ""
	temp.cipher = cipher
	temp.encryptionKey = encryptionKey
	temp.header = header
	temp.financingID = finangcingID
	p.encryptedPool.InsertModifyInvoice(temp, method, poolType)
}
func (p *Promoter) packFinancingInfo(UUID, header, info, poolType, method string) {
	cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(info))
	if err != nil {
		// logrus.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
		logs.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
	}
	temp := packedFinancingMessage{}
	fields := strings.Split(info, ",")
	temp.financingid = fields[9]
	temp.cipher = cipher
	temp.encryptionKey = encryptionKey
	temp.signed = signed
	temp.header = header
	temp.uuid = UUID
	p.encryptedPool.InsertFinancing(temp, method, poolType)
}

// 针对历史交易信息的packInfo
func (p *Promoter) packHistoricalInfos(UUID, header string, infos []string, poolType, method string) {
	var wg sync.WaitGroup
	for _, info := range infos {
		tempinfo := info
		wg.Add(1)
		go func(tempinfo string) {
			fields := strings.Split(tempinfo, ",")
			tradeYearMonth := fields[7] //交易年月
			tradeYearMonth = strings.Replace(tradeYearMonth, "[", "", -1)
			financeId := fields[4]
			fmt.Println(tradeYearMonth)
			cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(tempinfo))
			if err != nil {
				logs.Fatalln("数据加密失败,此条数据信息为:", header, tempinfo, "失败信息为:", err)
			}
			temp := packedHistoricalMessage{}
			temp.params = tradeYearMonth + "," + financeId + "," + string(signed) + "," + ""
			temp.cipher = cipher
			temp.encryptionKey = encryptionKey
			temp.header = header
			temp.uuid = UUID
			p.encryptedPool.InsertHistoricalTrans(temp, method, poolType)
			wg.Done()
		}(tempinfo)
	}
	wg.Wait()
}

// 针对入池信息的packInfo
func (p *Promoter) packPoolInfos(UUID, header string, infos []string, poolType, method string) {
	var wg sync.WaitGroup
	for _, info := range infos {
		tempinfo := info
		wg.Add(1)
		go func(header, tempinfo string) {
			fields := strings.Split(tempinfo, ",")
			tradeYearMonth := fields[5] //交易年月
			tradeYearMonth = strings.Replace(tradeYearMonth, "[", "", -1)
			cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(tempinfo))
			if err != nil {
				// logrus.Fatalln("数据加密失败,此条数据信息为:", header, tempinfo, "失败信息为:", err)
				logs.Fatalln("数据加密失败,此条数据信息为:", header, tempinfo, "失败信息为:", err)
			}
			temp := packedPoolMessage{}
			temp.params = tradeYearMonth + "," + string(signed) + "," + ""
			temp.cipher = cipher
			temp.encryptionKey = encryptionKey
			temp.header = header
			temp.uuid = UUID
			p.encryptedPool.InsertPoolData(temp, method, poolType)
			wg.Done()
		}(header, tempinfo)

	}
	wg.Wait()
}
func WriteToFile(info string) {
	filePath := "./configs/errorInfo.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(info)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}
