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
	"github.com/sirupsen/logrus"
)

const (
	waitCheck = "待审批"
)

var logs = logloader.NewLog()

type Promoter struct {
	server    *server.Server
	DataApi   *receive.FrontEnd
	monitor   *server.Monitor
	loader    *logloader.Loader
	chaininfo *chainloader.ChainInfo
}

func NewPromoter() *Promoter {
	ser := server.NewServer()
	api := receive.NewFrontEnd()
	monitor := server.NewMonitor()
	lder := logloader.NewLoader()
	chainld := chainloader.NewChainInfo()
	return &Promoter{
		server:    ser,
		DataApi:   api,
		monitor:   monitor,
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
			p.UpdatePushPaymentAccountsInfoHandler()
			p.LockPushPaymentAccountsInfoHandler()
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
		wg.Add(2)
		go p.enterPoolPlanInfoWaiter(len(poolPlan), &wg)
		go p.enterPoolUsedInfoWaiter(len(poolUsed), &wg)
		wg.Wait()
		logrus.Infoln("outout")
	}
}

// 处理回款账户信息
func (p *Promoter) UpdatePushPaymentAccountsInfoHandler() {
	if len(p.DataApi.UpdateCollectionAccountPool) != 0 {
		logs.Infoln("开始同步回款信息")
		logrus.Infoln("开始同步回款信息")
		payinfos := make(map[string]*receive.UpdateCollectionAccount, 0)
		p.DataApi.UpdateCollectionAccountPoolMutex.Lock()
		for uuid := range p.DataApi.UpdateCollectionAccountPool {
			payinfos[uuid] = p.DataApi.UpdateCollectionAccountPool[uuid]
			delete(p.DataApi.UpdateCollectionAccountPool, uuid)
		}
		p.DataApi.UpdateCollectionAccountPoolMutex.Unlock()
		accounts := p.server.PackedTradeData_UpdateAccountInfo(payinfos)
		for _, account := range accounts {
			err := p.server.UpdateAndLockPushPaymentAccounts(account.Uuid, account.Header+","+account.FinanceID, account.Cipher, account.EncryptionKey, account.NewHash, account.OldHash)
			if err != nil {
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		p.accountsUpdateInfoWaiter(len(accounts))
	}
}
func (p *Promoter) LockPushPaymentAccountsInfoHandler() {
	if len(p.DataApi.LockAccountPool) != 0 {
		logs.Infoln("开始同步回款信息")
		logrus.Infoln("开始同步回款信息")
		payinfos := make(map[string]*receive.LockAccount, 0)
		p.DataApi.LockAccountPoolMutex.Lock()
		for uuid := range p.DataApi.LockAccountPool {
			payinfos[uuid] = p.DataApi.LockAccountPool[uuid]
			delete(p.DataApi.LockAccountPool, uuid)
		}
		p.DataApi.LockAccountPoolMutex.Unlock()
		accounts := p.server.PackedTradeData_LockAccountInfo(payinfos)
		for _, account := range accounts {
			err := p.server.LockPaymentAccounts(account.Uuid, account.Header, account.FinanceID, string(account.Signed))
			if err != nil {
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		p.accountsLockInfoWaiter(len(accounts))
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
			err := p.server.IssueSupplierFinancingApplication(application.Uuid, application.Header, application.CustomerID, application.Cipher, application.EncryptionKey, application.Signed)
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
		financingInfos, modifyInvoices := p.server.PackedApplicationAndModifyInvoiceInfos(finintensWithSelectedInfos, waitCheck)
		for _, financingInfo := range financingInfos {
			err := p.server.UpdateSupplierFinancingApplication(financingInfo.Uuid, financingInfo.Header, financingInfo.CustomerID, financingInfo.Cipher, financingInfo.EncryptionKey, financingInfo.Signed)
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
