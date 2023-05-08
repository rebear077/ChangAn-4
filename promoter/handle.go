package promote

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	server "github.com/rebear077/changan/backend"
	chainloader "github.com/rebear077/changan/chaininfos"
	receive "github.com/rebear077/changan/connApi"
	logloader "github.com/rebear077/changan/logs"
	uptoChain "github.com/rebear077/changan/tochain"
	"github.com/sirupsen/logrus"
)

// const (
//
//	PoolPlanInfos             = "poolPlan"
//	PoolUsedInfos             = "poolUsed"
//	HistoricalOrderInfos      = "hisOrder"
//	HistoricalReceivableInfos = "hisReceivable"
//	HistoricalSettleInfos     = "hisSettle"
//	HistoricalUsedInfos       = "hisUsed"
//
// )
const (
	waitCheck = "待审批"
	again     = "重新申请"
)

var logs = logloader.NewLog()

type Promoter struct {
	server  *server.Server
	DataApi *receive.FrontEnd
	monitor *server.Monitor
	// encryptedPool *Pools
	loader    *logloader.Loader
	chaininfo *chainloader.ChainInfo
}

func NewPromoter() *Promoter {
	ser := server.NewServer()
	api := receive.NewFrontEnd()
	monitor := server.NewMonitor()
	// pool := NewPools()
	lder := logloader.NewLoader()
	chainld := chainloader.NewChainInfo()
	return &Promoter{
		server:  ser,
		DataApi: api,
		monitor: monitor,
		// encryptedPool: pool,
		loader:    lder,
		chaininfo: chainld,
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
			p.HistoricalInfoHandler()
			p.PushPaymentAccountsInfoHandler()
			p.EnterPoolInfoHandler()
			p.FinancingApplicationInfoWithSelectedInfosHandler()
			p.ModifySupplierFinancingApplicationInfoWithSelectedInfosHandler()
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

// 处理推送的发票信息
func (p *Promoter) InvoiceInfoHandler() {
	if len(p.DataApi.InvoicePool) != 0 {
		logs.Infoln("开始处理贸易数据-发票信息")
		// var wg sync.WaitGroup
		invoices := make(map[string]*receive.InvoiceInformation, 0)
		p.DataApi.IssueInvoicemutex.Lock()
		for uuid := range p.DataApi.InvoicePool {
			invoices[uuid] = p.DataApi.InvoicePool[uuid]
			delete(p.DataApi.InvoicePool, uuid)
		}
		p.DataApi.IssueInvoicemutex.Unlock()
		packedInvoices := p.server.PackedTradeData_InvoiceInfo(invoices)
		for _, packedInvoice := range packedInvoices {
			err := p.server.IssueInvoiceInformation(packedInvoice.Uuid, packedInvoice.Header, packedInvoice.Params, packedInvoice.Cipher, packedInvoice.EncryptionKey)
			if err != nil {
				logs.Errorln("发票信息上链失败:", packedInvoice.Header, "失败信息为:", err)
			}
		}
		p.invoiceInfoWaiter(len(packedInvoices))
	}
}

// 处理历史交易信息
func (p *Promoter) HistoricalInfoHandler() {
	if len(p.DataApi.TransactionHistoryPool) != 0 {
		logs.Infoln("开始c处理贸易数据-历史交易信息")
		hisinfos := make(map[string]*receive.TransactionHistory, 0)
		p.DataApi.TransactionHistorymutex.Lock()
		for uuid := range p.DataApi.TransactionHistoryPool {
			hisinfos[uuid] = p.DataApi.TransactionHistoryPool[uuid]
			delete(p.DataApi.TransactionHistoryPool, uuid)
		}
		p.DataApi.TransactionHistorymutex.Unlock()
		used, settle, order, receivable := p.server.PackedTradeData_HistoricalInfo(hisinfos)
		for _, usedInfo := range used {
			err := p.server.IssueHistoricalUsedInformation(usedInfo.Uuid, usedInfo.Header, usedInfo.Params, usedInfo.Cipher, usedInfo.EncryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", usedInfo.Header, "失败信息为:", err)
			}
		}
		for _, settleInfo := range settle {
			err := p.server.IssueHistoricalSettleInformation(settleInfo.Uuid, settleInfo.Header, settleInfo.Params, settleInfo.Cipher, settleInfo.EncryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", settleInfo.Header, "失败信息为:", err)
			}
		}
		for _, orderInfo := range order {
			err := p.server.IssueHistoricalOrderInformation(orderInfo.Uuid, orderInfo.Header, orderInfo.Params, orderInfo.Cipher, orderInfo.EncryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", orderInfo.Header, "失败信息为:", err)
			}
		}
		for _, receivableInfo := range receivable {
			err := p.server.IssueHistoricalReceivableInformation(receivableInfo.Uuid, receivableInfo.Header, receivableInfo.Params, receivableInfo.Cipher, receivableInfo.EncryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", receivableInfo.Header, "失败信息为:", err)
			}
		}
		var wg sync.WaitGroup
		wg.Add(4)
		go p.historicalOrderInfoWaiter(len(order), &wg)
		go p.historicalReveivableInfoWaiter(len(receivable), &wg)
		go p.historicalSettleInfoWaiter(len(settle), &wg)
		go p.historicalUsedInfoWaiter(len(used), &wg)
		wg.Wait()
	}
}

// 处理入池数据信息
func (p *Promoter) EnterPoolInfoHandler() {
	if len(p.DataApi.EnterpoolDataPool) != 0 {
		logrus.Infoln("开始入池数据信息")
		logs.Infoln("开始入池数据信息")
		poolInfos := make(map[string]*receive.EnterpoolData, 0)
		p.DataApi.EnterpoolDatamutex.Lock()
		for uuid := range p.DataApi.EnterpoolDataPool {
			poolInfos[uuid] = p.DataApi.EnterpoolDataPool[uuid]
			delete(p.DataApi.EnterpoolDataPool, uuid)
		}
		p.DataApi.EnterpoolDatamutex.Unlock()
		poolPlan, poolUsed := p.server.PackedTradeData_EnterPoolInfo(poolInfos)
		for _, plan := range poolPlan {
			err := p.server.IssuePoolPlanInformation(plan.Uuid, plan.Header, plan.Params, plan.Cipher, plan.EncryptionKey)
			if err != nil {
				logs.Errorln("入池计划信息上链失败:", plan.Header, "失败信息为:", err)
			}
		}
		for _, used := range poolUsed {
			err := p.server.IssuePoolUsedInformation(used.Uuid, used.Header, used.Params, used.Cipher, used.EncryptionKey)
			if err != nil {
				logs.Errorln("入池用户信息上链失败:", used.Header, "失败信息为:", err)
			}

		}
		var wg sync.WaitGroup
		go p.enterPoolPlanInfoWaiter(len(poolPlan), &wg)
		go p.enterPoolUsedInfoWaiter(len(poolUsed), &wg)
		wg.Wait()
		logrus.Infoln("outout")
	}
}

// 处理回款账户信息
func (p *Promoter) PushPaymentAccountsInfoHandler() {
	if len(p.DataApi.CollectionAccountPool) != 0 {
		logs.Infoln("开始同步回款信息")
		logrus.Infoln("开始同步回款信息")
		payinfos := make(map[string]*receive.CollectionAccount, 0)
		p.DataApi.CollectionAccountmutex.Lock()
		for uuid := range p.DataApi.CollectionAccountPool {
			payinfos[uuid] = p.DataApi.CollectionAccountPool[uuid]
			delete(p.DataApi.CollectionAccountPool, uuid)
		}
		p.DataApi.CollectionAccountmutex.Unlock()
		accounts := p.server.PackedTradeData_AccountInfo(payinfos)
		for _, account := range accounts {
			err := p.server.UpdatePushPaymentAccount(account.Uuid, account.Header, account.Cipher, account.EncryptionKey, account.Signed)
			if err != nil {
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		p.accountsInfoWaiter(len(accounts))
	}
}

// 处理融资意向申请信息
func (p *Promoter) FinancingApplicationInfoWithSelectedInfosHandler() {
	if len(p.DataApi.FinancingIntentionWithSelectedInfosPool) != 0 {
		logs.Infoln("开始同步融资意向请求信息")
		finintensWithSelectedInfos := make(map[string]*receive.SelectedInfosAndFinancingApplication, 0)
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Lock()
		for uuid := range p.DataApi.FinancingIntentionWithSelectedInfosPool {
			finintensWithSelectedInfos[uuid] = p.DataApi.FinancingIntentionWithSelectedInfosPool[uuid]
			delete(p.DataApi.FinancingIntentionWithSelectedInfosPool, uuid)
		}
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Unlock()
		financingInfos, modifyInvoices := p.server.PackedApplicationAndModifyInvoiceInfos(finintensWithSelectedInfos, waitCheck)
		for _, application := range financingInfos {
			err := p.server.IssueSupplierFinancingApplication(application.Uuid, application.Header, application.State, application.Cipher, application.EncryptionKey, application.Signed)
			if err != nil {
				logs.Errorln("融资意向请求上链失败,", "失败信息为:", err)
			}
		}
		for _, modify := range modifyInvoices {
			err := p.server.VerifyAndUpdateInvoiceInformation(modify.Uuid, modify.Header, modify.Sign, modify.FinancingID)
			if err != nil {
				logs.Errorln("发票信息上链失败:", modify.Header, "失败信息为:", err)
			}
		}
		var wg sync.WaitGroup
		go p.financingApplicationInfoWaiter(len(financingInfos), &wg)
		go p.modifyInvoiceInfoWaiter(len(modifyInvoices), &wg)
		wg.Done()
	}
}

// 处理修改融资意向申请信息
func (p *Promoter) ModifySupplierFinancingApplicationInfoWithSelectedInfosHandler() {
	if len(p.DataApi.ModifyFinancingWithSelectedInfosPool) != 0 {
		logs.Infoln("开始修改融资意向请求信息")
		finintensWithSelectedInfos := make(map[string]*receive.SelectedInfosAndFinancingApplication, 0)
		p.DataApi.ModifyFinancingWithSelectedInfosPoolMutex.Lock()
		for uuid := range p.DataApi.ModifyFinancingWithSelectedInfosPool {
			finintensWithSelectedInfos[uuid] = p.DataApi.ModifyFinancingWithSelectedInfosPool[uuid]
			delete(p.DataApi.ModifyFinancingWithSelectedInfosPool, uuid)
		}
		p.DataApi.ModifyFinancingWithSelectedInfosPoolMutex.Unlock()
		financingInfos, modifyInvoices := p.server.PackedApplicationAndModifyInvoiceInfos(finintensWithSelectedInfos, again)
		for _, financingInfo := range financingInfos {
			err := p.server.UpdateSupplierFinancingApplication(financingInfo.Uuid, financingInfo.Header, financingInfo.State, financingInfo.Cipher, financingInfo.EncryptionKey, financingInfo.Signed)
			if err != nil {
				logs.Errorln("融资意向请求上链失败,", "失败信息为:", err)
			}
		}
		for _, modifyInvoice := range modifyInvoices {
			err := p.server.VerifyAndUpdateInvoiceInformation(modifyInvoice.Uuid, modifyInvoice.Header, modifyInvoice.Sign, modifyInvoice.FinancingID)
			if err != nil {
				logs.Errorln("发票信息上链失败:", modifyInvoice.Header, "失败信息为:", err)
			}
		}
		var wg sync.WaitGroup
		go p.modifyFinancingInfoWaiter(len(financingInfos), &wg)
		go p.modifyInvoiceInfoWhenModifyApplicationWaiter(len(modifyInvoices), &wg)
		wg.Wait()

	}
}
func (p *Promoter) invoiceInfoWaiter(length int) {
	for {
		counter := 0
		uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
			uptoChain.InvoiceMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.InvoiceMapLock.Unlock()
			return true
		})
		if counter == length {
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

func (p *Promoter) historicalOrderInfoWaiter(orderLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalOrderMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalOrderMapLock.Unlock()
			return true
		})
		if counter == orderLength {
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
}
func (p *Promoter) historicalReveivableInfoWaiter(receivableLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalReceivableMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalReceivableMapLock.Unlock()
			return true
		})
		if counter == receivableLength {
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
}
func (p *Promoter) historicalSettleInfoWaiter(settleLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalSettleMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalSettleMapLock.Unlock()
			return true
		})
		if counter == settleLength {
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
}
func (p *Promoter) historicalUsedInfoWaiter(usedLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalUsedMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalUsedMapLock.Unlock()
			return true
		})
		if counter == usedLength {
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
}
func (p *Promoter) enterPoolPlanInfoWaiter(planLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
			uptoChain.PoolPlanMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.PoolPlanMapLock.Unlock()
			return true
		})
		if counter == planLength {
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
}
func (p *Promoter) enterPoolUsedInfoWaiter(usedLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
			uptoChain.PoolUsedMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.PoolUsedMapLock.Unlock()
			return true
		})
		if counter == usedLength {
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
}
func (p *Promoter) accountsInfoWaiter(accountsLength int) {
	for {
		counter := 0
		uptoChain.CollectionAccountMap.Range(func(key, value interface{}) bool {
			uptoChain.CollectionAccountMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.CollectionAccountMapLock.Unlock()
			return true
		})
		if counter == accountsLength {
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
func (p *Promoter) financingApplicationInfoWaiter(applicationLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
			uptoChain.FinancingApplicationIssueMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.FinancingApplicationIssueMapLock.Unlock()
			return true
		})
		if counter == applicationLength {
			p.DataApi.FinancingIntentionIssueOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
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
}
func (p *Promoter) modifyInvoiceInfoWaiter(modifyLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyInvoiceMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyInvoiceMapLock.Unlock()
			return true
		})
		if counter == modifyLength {
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
	wg.Done()
}
func (p *Promoter) modifyFinancingInfoWaiter(applicationLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyFinancingMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyFinancingMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyFinancingMapLock.Unlock()
			return true
		})
		if counter == applicationLength {
			p.DataApi.ModifyFinancingOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyFinancingMap.Range(func(key, value interface{}) bool {
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
}
func (p *Promoter) modifyInvoiceInfoWhenModifyApplicationWaiter(modifyLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyInvoiceWhenMFAMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyInvoiceWhenMFAMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyInvoiceWhenMFAMapLock.Unlock()
			return true
		})
		if counter == modifyLength {
			p.DataApi.ModifyInvoiceWhenFinancingOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyInvoiceWhenMFAMap.Range(func(key, value interface{}) bool {
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
}

// 写入文件
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
