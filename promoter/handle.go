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
	"github.com/rebear077/changan/errorhandle"
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
		invoices := make(map[string][]*receive.InvoiceInformation, 0)
		p.DataApi.IssueInvoicemutex.Lock()
		invoices = p.DataApi.InvoicePool
		for uuid := range p.DataApi.InvoicePool {
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
			uptoChain.M.Range(func(key, value interface{}) bool {
				mapping := value.(map[string]*uptoChain.ResponseMessage)
				counter += len(mapping)
				return true
			})
			if counter == len(messages) {
				p.DataApi.IssueInvoiceOKChan <- true
				for {
					flag := 0
					uptoChain.M.Range(func(key, value interface{}) bool {
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
		hisinfos := make(map[string][]*receive.TransactionHistory, 0)
		p.DataApi.TransactionHistorymutex.Lock()
		hisinfos = p.DataApi.TransactionHistoryPool
		for uuid := range p.DataApi.TransactionHistoryPool {
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
			err := p.server.IssueHistoricalUsedInformation(tempUsed.header, tempUsed.params, tempUsed.cipher, tempUsed.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}
		}
		for _, message := range hisSettleMessage {
			tempSettle, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalSettleInformation(tempSettle.header, tempSettle.params, tempSettle.cipher, tempSettle.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempSettle.header, "失败信息为:", err)
			}
		}
		for _, message := range hisOrderMessage {
			tempOrder, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalOrderInformation(tempOrder.header, tempOrder.params, tempOrder.cipher, tempOrder.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempOrder.header, "失败信息为:", err)
			}
		}
		for _, message := range hisReceivableMessage {
			tempReceivable, _ := message.(packedHistoricalMessage)
			err := p.server.IssueHistoricalReceivableInformation(tempReceivable.header, tempReceivable.params, tempReceivable.cipher, tempReceivable.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempReceivable.header, "失败信息为:", err)
			}
		}
		wg.Add(4)
		var hisUsedTotal int
		var hisUsedSuccess int
		var hisUsedError int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.HistoricalUsedInformation)
				success := uptoChain.QueryHistoricalUsedCounter()
				if errNum+success == len(hisUsedMessage) {
					// if errNum != 0 {
					// 	errorHisUsed := errorhandle.ERRDealer.GetErrorInfo(uptoChain.HistoricalUsedInformation)
					// 	for transactionHash, data := range errorHisUsed {
					// 		parseRet, ok := data.([]interface{})
					// 		if !ok {
					// 			logs.Fatalln("解析失败")
					// 		}
					// 		info := transactionHash +
					// 			WriteToFile(value+"\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssueHistoricalUsedInformationPool()
					// }
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisUsedMessage), success, errNum)
					uptoChain.ResetHistoricalUsedCounter()
					hisUsedTotal = len(hisUsedMessage)
					hisUsedSuccess = success
					hisUsedError = errNum
					break
				}
			}
			wg.Done()

		}()

		var hisSettleTotal int
		var hisSettleError int
		var hisSettleSuccess int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.HistoricalSettleInformation)
				success := uptoChain.QueryHistoricalSettleCounter()
				if errNum+success == len(hisSettleMessage) {
					// if errNum != 0 {
					// 	mapping := errorhandle.ERRDealer.QueryHistoricalSettleInfoPool()
					// 	for _, value := range mapping {
					// 		WriteToFile(value + "\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssueHistoricalSettleInformationPool()
					// }
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisSettleMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisSettleMessage), success, errNum)
					uptoChain.ResetHistoricalSettleCounter()
					hisSettleTotal += len(hisSettleMessage)
					hisSettleError += errNum
					hisSettleSuccess += success
					break
				}
			}
			wg.Done()
		}()
		var hisOrderTotal int
		var hisOrderSuccess int
		var hisOrderError int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.HistoricalOrderInformation)
				success := uptoChain.QueryHistoricalOrderCounter()
				if errNum+success == len(hisOrderMessage) {
					// if errNum != 0 {
					// 	mapping := errorhandle.ERRDealer.QueryHistoricalOrderInfoPool()
					// 	for _, value := range mapping {
					// 		WriteToFile(value + "\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssueHistoricalOrderInformationPool()
					// }
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisOrderMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisOrderMessage), success, errNum)
					uptoChain.ResetHistoricalOrderCounter()
					hisOrderTotal = len(hisOrderMessage)
					hisOrderSuccess = success
					hisOrderError = errNum
					break
				}
			}
			wg.Done()
		}()
		var hisReceivableTotal int
		var hisReceivableError int
		var hisReceivableSuccess int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.HistoricalReceivableInformation)
				success := uptoChain.QueryHistoricalReceivableCounter()
				if errNum+success == len(hisReceivableMessage) {
					// if errNum != 0 {
					// 	mapping := errorhandle.ERRDealer.QueryHistoricalReceivableInfoPool()
					// 	for _, value := range mapping {
					// 		WriteToFile(value + "\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssueHistoricalReceivableInformationPool()
					// }
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisReceivableMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisReceivableMessage), success, errNum)
					uptoChain.ResetHistoricalReceivableCounter()
					hisReceivableTotal = len(hisReceivableMessage)
					hisReceivableSuccess = success
					hisReceivableError = errNum
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		hisOrder := [3]int{hisOrderTotal, hisOrderSuccess, hisOrderError}
		hisReceivable := [3]int{hisReceivableTotal, hisReceivableSuccess, hisReceivableError}
		hisSettle := [3]int{hisSettleTotal, hisSettleSuccess, hisSettleError}
		hisUsed := [3]int{hisUsedTotal, hisUsedSuccess, hisUsedError}
		historyInfos := make(map[string][3]int)
		historyInfos[HistoricalOrderInfos] = hisOrder
		historyInfos[HistoricalReceivableInfos] = hisReceivable
		historyInfos[HistoricalSettleInfos] = hisSettle
		historyInfos[HistoricalUsedInfos] = hisUsed
		p.DataApi.HistoryInfoChan <- historyInfos
		logs.Println("退出")
	}
}
func (p *Promoter) PoolInfoHandler() {
	if len(p.DataApi.EnterpoolDataPool) != 0 {
		// logrus.Infoln("开始入池数据信息")
		logs.Infoln("开始入池数据信息")
		var wg sync.WaitGroup
		poolinfos := make([]*receive.EnterpoolData, 0)
		p.DataApi.EnterpoolDatamutex.Lock()
		poolinfos = append(poolinfos, p.DataApi.EnterpoolDataPool...)
		p.DataApi.EnterpoolDataPool = nil
		p.DataApi.EnterpoolDatamutex.Unlock()
		mapping := server.EncodeEnterpoolData(poolinfos)
		for index := range mapping {
			for header, info := range mapping[index] {
				tempheader := header
				tempinfo := info
				wg.Add(1)
				go func(tempheader string, tempinfo string) {
					var wwg sync.WaitGroup
					planvalue, providerusedvalue := server.PoolInformationSlice(tempheader, tempinfo, 1)
					wwg.Add(2)
					go func(tempheader string, planvalue []string) {
						p.packPoolInfos(tempheader, planvalue, "fast", "poolPlan")
						wwg.Done()
					}(tempheader, planvalue)
					go func(tempheader string, providerusedvalue []string) {
						p.packPoolInfos(tempheader, providerusedvalue, "fast", "poolUsed")
						wwg.Done()
					}(tempheader, providerusedvalue)
					wwg.Wait()
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		planMessages := p.encryptedPool.QueryMessages("poolPlan", "fast")
		usedMessages := p.encryptedPool.QueryMessages("poolUsed", "fast")
		for _, message := range planMessages {
			tempPlan, _ := message.(packedPoolMessage)
			err := p.server.IssuePoolPlanInformation(tempPlan.header, tempPlan.params, tempPlan.cipher, tempPlan.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempPlan.header, "失败信息为:", err)
			}
		}
		for _, message := range usedMessages {
			tempUsed, _ := message.(packedPoolMessage)
			err := p.server.IssuePoolUsedInformation(tempUsed.header, tempUsed.params, tempUsed.cipher, tempUsed.encryptionKey)
			if err != nil {
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}

		}
		wg.Add(2)
		var poolPlanTotal int
		var poolPlanSuccess int
		var poolPlanError int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.PoolPlanInfo)
				success := uptoChain.QueryPoolPlanCounter()
				if errNum+success == len(planMessages) {
					// if errNum != 0 {
					// 	mapping := errorhandle.ERRDealer.QueryPoolPlanInfoPool()
					// 	for _, value := range mapping {
					// 		WriteToFile(value + "\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssuePoolPlanInformationPool()
					// }
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(planMessages), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(planMessages), success, errNum)
					uptoChain.ResetPoolPlanCounter()
					poolPlanTotal = len(planMessages)
					poolPlanSuccess = success
					poolPlanError = errNum
					break
				}
			}
			wg.Done()
		}()
		var poolUsedTotal int
		var poolUsedSuccess int
		var poolUsedError int
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.PoolUsedInfo)
				success := uptoChain.QueryPoolUsedCounter()
				if errNum+success == len(usedMessages) {
					// if errNum != 0 {
					// 	mapping := errorhandle.ERRDealer.QueryHistoricalUsedInfoPool()
					// 	for _, value := range mapping {
					// 		WriteToFile(value + "\n")
					// 	}
					// 	errorhandle.ERRDealer.DeleteErrorIssuePoolUsedInformationPool()
					// }
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(usedMessages), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(usedMessages), success, errNum)
					uptoChain.ResetPoolUsedCounter()
					poolUsedTotal = len(usedMessages)
					poolUsedSuccess = success
					poolUsedError = errNum
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		poolPlan := [3]int{poolPlanTotal, poolPlanSuccess, poolPlanError}
		poolUsed := [3]int{poolUsedTotal, poolUsedSuccess, poolUsedError}
		poolInfos := make(map[string][3]int)
		poolInfos[PoolPlanInfos] = poolPlan
		poolInfos[PoolUsedInfos] = poolUsed
		p.DataApi.PoolInfoChan <- poolInfos
		logs.Println("退出")
	}
}

func (p *Promoter) SupplierFinancingApplicationInfoWithSelectedInfosHandler() {
	if len(p.DataApi.FinancingIntentionWithSelectedInfosPool) != 0 {
		logs.Infoln("开始同步融资意向请求信息")
		var wg sync.WaitGroup
		finintensWithSelectedInfos := make([]*receive.SelectedInfosAndFinancingApplication, 0)
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Lock()
		finintensWithSelectedInfos = append(finintensWithSelectedInfos, p.DataApi.FinancingIntentionWithSelectedInfosPool...)
		p.DataApi.FinancingIntentionWithSelectedInfosPool = nil
		p.DataApi.FinancingIntentionWithSelectedInfosMutex.Unlock()
		selectedInfos, financing := server.HandleFinancingIntentionAndSelectedInfos(finintensWithSelectedInfos)
		for index := range mapping {
			for header, info := range mapping[index] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(tempheader string, tempinfo string) {
					p.packFinancingInfo(tempheader, tempinfo, "fast", "application")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("application", "fast")
		for _, message := range messages {
			temp, _ := message.(packedFinancingMessage)
			err := p.server.IssueSupplierFinancingApplication(temp.header, temp.financingid, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				// logrus.Errorln("融资意向请求上链失败,", "失败信息为:", err)
				logs.Errorln("融资意向请求上链失败,", "失败信息为:", err)
			}
		}
		for {
			errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.SupplierFinancingApplicationInfo)
			success := uptoChain.QuerySupplierSuccessCounter()
			if errNum+success == len(messages) {
				// if errNum != 0 {
				// 	mapping := errorhandle.ERRDealer.QuerySupplierFinancingApplicationPool()
				// 	for _, value := range mapping {
				// 		WriteToFile(value + "\n")
				// 	}
				// 	errorhandle.ERRDealer.DeleteErrorIssueSupplierFinancingApplicationPool()
				// }
				// logrus.Infof("同步融资意向完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				logs.Infof("同步融资意向完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				uptoChain.ResetSupplierSuccessCounter()
				result := [3]int{len(messages), success, errNum}
				p.DataApi.FinancingIntentionChan <- result
				break
			}
		}
	}
}

func (p *Promoter) PushPaymentAccountsInfoHandler() {
	if len(p.DataApi.CollectionAccountPool) != 0 {
		logs.Infoln("开始同步回款信息")
		var wg sync.WaitGroup
		payinfos := make([]*receive.CollectionAccount, 0)
		p.DataApi.CollectionAccountmutex.Lock()
		payinfos = append(payinfos, p.DataApi.CollectionAccountPool...)
		p.DataApi.CollectionAccountPool = nil
		p.DataApi.CollectionAccountmutex.Unlock()
		mapping := server.EncodeCollectionAccount(payinfos)
		for index := range mapping {
			for header, info := range mapping[index] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(tempheader string, tempinfo string) {
					p.packInfo(tempheader, tempinfo, "fast", "payment")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("payment", "fast")
		for _, message := range messages {
			temp, ok := message.(packedMessage)
			if !ok {
				fmt.Println("errorerror")
			}
			err := p.server.UpdatePushPaymentAccount(temp.header, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		for {
			errNum := errorhandle.ERRDealer.GetErrorLength(uptoChain.UpdatePushPaymentAccounts)
			success := uptoChain.QueryPaymentAccountsCounter()
			if errNum+success == len(messages) {
				// if errNum != 0 {
				// 	mapping := errorhandle.ERRDealer.QueryPushPaymentAccountPool()
				// 	for _, value := range mapping {
				// 		WriteToFile(value + "\n")
				// 	}
				// 	errorhandle.ERRDealer.DeleteErrorIssuePushPaymentAccountsPool()
				// }
				logs.Infof("回款信息同步完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				uptoChain.ResetPaymentAccountsCounter()
				resluts := [3]int{len(messages), success, errNum}
				p.DataApi.PushPaymentAccountChan <- resluts
				break
			}
		}
	}
}

func (p *Promoter) packInfo(header string, info string, poolType string, method string) {
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

func (p *Promoter) packFinancingInfo(header string, info string, poolType string, method string) {
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
func (p *Promoter) packPoolInfos(header string, infos []string, poolType string, method string) {
	var wg sync.WaitGroup
	for _, info := range infos {
		tempinfo := info
		wg.Add(1)
		go func(header string, tempinfo string) {
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
