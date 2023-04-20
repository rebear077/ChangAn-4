package uptoChain

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"ethereum/go-ethereum/common"

	"github.com/rebear077/changan/abi"
	smartcontract "github.com/rebear077/changan/contract"
	"github.com/rebear077/changan/core/types"
	"github.com/rebear077/changan/errorhandle"
	logloader "github.com/rebear077/changan/logs"
	"github.com/sirupsen/logrus"
)

var logs = logloader.NewLog()

const (
	SupplierFinancingApplicationInfo = "SupplierFinancingApplication"
	IssueInvoiceInformation          = "IssueInvoiceInformation"
	UpdateInvoiceInformation         = "UpdateInvoiceInformation"
	HistoricalUsedInformation        = "HistoricalUsedInformation"
	HistoricalSettleInformation      = "HistoricalSettleInformation"
	HistoricalReceivableInformation  = "HistoricalReceivableInformation"
	HistoricalOrderInformation       = "HistoricalOrderInformation"
	UpdatePushPaymentAccounts        = "UpdatePushPaymentAccounts"
	PoolPlanInfo                     = "PoolPlanInfo"
	PoolUsedInfo                     = "PoolUsedInfo"
)

var (
	supplierCounter             = 0
	issueinvoiceCounter         = 0
	updateinvoiceCounter        = 0
	historicalUsedCounter       = 0
	historicalSettleCounter     = 0
	historicalOrderCounter      = 0
	historicalReceivableCounter = 0
	paymentAccountsCounter      = 0
	poolUsedCounter             = 0
	poolPlanCounter             = 0

	supplierCounterMutex             sync.Mutex
	issueinvoiceCounterMutex         sync.Mutex
	updateinvoiceCounterMutex        sync.Mutex
	historicalUsedCounterMutex       sync.Mutex
	historicalSettleCounterMutex     sync.Mutex
	historicalOrderCounterMutex      sync.Mutex
	historicalReceivableCounterMutex sync.Mutex
	paymentAccountsCounterMutex      sync.Mutex
	poolUsedCounterMutex             sync.Mutex
	poolPlanCounterMutex             sync.Mutex
)

// 融资意向
func invokeIssueSupplierFinancingApplicationHandler(receipt *types.Receipt, err error) {

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueSupplierFinancingApplication", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueSupplierFinancingApplication", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(SupplierFinancingApplicationInfo, receipt.TransactionHash, parseRet)
	} else {
		supplierCounterMutex.Lock()
		supplierCounter += 1
		supplierCounterMutex.Unlock()
	}
}

// 发布发票信息回调函数
func invokeIssueInvoiceInformationStorageHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueInvoiceInformationStorage", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	// fmt.Println(setedLines)
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueInvoiceInformationStorage", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		var message string
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			message = parseRet[0].(string) + "," + parseRet[1].(string)
		}
		packedMessage := new(ResponseMessage)
		packedMessage.ok = false
		packedMessage.message = message
		M.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			_, ok := M.LoadOrStore(uuid, mapping)
			if !ok {
				logs.Fatalln("sync.map error")
			}
			return ok
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		M.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			_, ok := M.LoadOrStore(uuid, mapping)
			if !ok {
				logs.Fatalln("sync.map error")
			}
			return ok
		})
	}
}

// 验证并更新发票信息回调函数
func invokeVerifyAndUpdateInvoiceInformationStorageHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "updateInvoiceInformationStorage", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("updateInvoiceInformationStorage", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(UpdateInvoiceInformation, receipt.TransactionHash, parseRet)
	} else {
		updateinvoiceCounterMutex.Lock()
		updateinvoiceCounter += 1
		updateinvoiceCounterMutex.Unlock()
	}
}

// 历史交易信息之入库信息
func invokeIssueHistoricalUsedInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueHistoricalUsedInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalUsedInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(HistoricalUsedInformation, receipt.TransactionHash, parseRet)
	} else {
		historicalUsedCounterMutex.Lock()
		historicalUsedCounter += 1
		historicalUsedCounterMutex.Unlock()
	}
}

// 历史交易信息之结算信息
func invokeIssueHistoricalSettleInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueHistoricalSettleInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalSettleInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(HistoricalSettleInformation, receipt.TransactionHash, parseRet)
	} else {
		historicalSettleCounterMutex.Lock()
		historicalSettleCounter += 1
		historicalSettleCounterMutex.Unlock()
	}
}

// 历史交易信息之订单信息
func invokeIssueHistoricalOrderInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueHistoricalOrderInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalOrderInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(HistoricalOrderInformation, receipt.TransactionHash, parseRet)
	} else {
		historicalOrderCounterMutex.Lock()
		historicalOrderCounter += 1
		historicalOrderCounterMutex.Unlock()
	}
}

// 历史交易信息之应收账款信息
func invokeIssueHistoricalReceivableInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueHistoricalReceivableInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalReceivableInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(HistoricalReceivableInformation, receipt.TransactionHash, parseRet)
	} else {
		historicalReceivableCounterMutex.Lock()
		historicalReceivableCounter += 1
		historicalReceivableCounterMutex.Unlock()
	}
}

// 回款信息
func invokeUpdatePushPaymentAccountsHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "updatePushPaymentAccounts", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("updatePushPaymentAccounts", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(UpdatePushPaymentAccounts, receipt.TransactionHash, parseRet)
	} else {
		paymentAccountsCounterMutex.Lock()
		paymentAccountsCounter += 1
		paymentAccountsCounterMutex.Unlock()
	}
}

// 入池数据之供应商生产计划信息
func invokeIssuePoolPlanInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issuePoolPlanInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issuePoolPlanInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(PoolPlanInfo, receipt.TransactionHash, parseRet)
	} else {
		poolPlanCounterMutex.Lock()
		poolPlanCounter += 1
		poolPlanCounterMutex.Unlock()
	}
}

// 入池数据之供应商生产入库信息
func invokeIssuePoolUsedInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issuePoolUsedInformation", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issuePoolUsedInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		}
		errorhandle.ERRDealer.InsertError(PoolUsedInfo, receipt.TransactionHash, parseRet)
	} else {
		poolUsedCounterMutex.Lock()
		poolUsedCounter += 1
		poolUsedCounterMutex.Unlock()
	}
}

func parseOutput(abiStr, name string, receipt *types.Receipt) (*big.Int, error) {
	var ret *big.Int
	if receipt.Output == "" {
		logrus.Errorln("empty output")
		logrus.Errorln(receipt.TransactionHash)
		ret = big.NewInt(0)
		return ret, nil
	}
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	b, err := hex.DecodeString(receipt.Output[2:])
	if err != nil {
		return nil, err
	}
	err = parsed.Unpack(&ret, name, b)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func QuerySupplierSuccessCounter() int {
	supplierCounterMutex.Lock()
	temp := supplierCounter
	supplierCounterMutex.Unlock()
	return temp
}
func QueryIssueInvoiceSuccessCounter() int {
	issueinvoiceCounterMutex.Lock()
	temp := issueinvoiceCounter
	issueinvoiceCounterMutex.Unlock()
	return temp
}
func QueryUpdateInvoiceSuccessCounter() int {
	updateinvoiceCounterMutex.Lock()
	temp := updateinvoiceCounter
	updateinvoiceCounterMutex.Unlock()
	return temp
}
func QueryHistoricalUsedCounter() int {
	historicalUsedCounterMutex.Lock()
	temp := historicalUsedCounter
	historicalUsedCounterMutex.Unlock()
	return temp
}
func QueryHistoricalOrderCounter() int {
	historicalOrderCounterMutex.Lock()
	temp := historicalOrderCounter
	historicalOrderCounterMutex.Unlock()
	return temp
}
func QueryHistoricalSettleCounter() int {
	historicalSettleCounterMutex.Lock()
	temp := historicalSettleCounter
	historicalSettleCounterMutex.Unlock()
	return temp
}
func QueryHistoricalReceivableCounter() int {
	historicalReceivableCounterMutex.Lock()
	temp := historicalReceivableCounter
	historicalReceivableCounterMutex.Unlock()
	return temp
}
func QueryPaymentAccountsCounter() int {
	paymentAccountsCounterMutex.Lock()
	temp := paymentAccountsCounter
	paymentAccountsCounterMutex.Unlock()
	return temp
}
func QueryPoolPlanCounter() int {
	poolPlanCounterMutex.Lock()
	temp := poolPlanCounter
	poolPlanCounterMutex.Unlock()
	return temp
}
func QueryPoolUsedCounter() int {
	poolUsedCounterMutex.Lock()
	temp := poolUsedCounter
	poolUsedCounterMutex.Unlock()
	return temp
}

func ResetSupplierSuccessCounter() {
	supplierCounterMutex.Lock()
	supplierCounter = 0
	supplierCounterMutex.Unlock()

}
func ResetIssueInvoiceSuccessCounter() {
	issueinvoiceCounterMutex.Lock()
	issueinvoiceCounter = 0
	issueinvoiceCounterMutex.Unlock()
}
func ResetUpdateInvoiceSuccessCounter() {
	updateinvoiceCounterMutex.Lock()
	updateinvoiceCounter = 0
	updateinvoiceCounterMutex.Unlock()
}
func ResetHistoricalUsedCounter() {
	historicalUsedCounterMutex.Lock()
	historicalUsedCounter = 0
	historicalUsedCounterMutex.Unlock()

}
func ResetHistoricalOrderCounter() {
	historicalOrderCounterMutex.Lock()
	historicalOrderCounter = 0
	historicalOrderCounterMutex.Unlock()

}
func ResetHistoricalSettleCounter() {
	historicalSettleCounterMutex.Lock()
	historicalSettleCounter = 0
	historicalSettleCounterMutex.Unlock()

}
func ResetHistoricalReceivableCounter() {
	historicalReceivableCounterMutex.Lock()
	historicalReceivableCounter = 0
	historicalReceivableCounterMutex.Unlock()

}
func ResetPaymentAccountsCounter() {
	paymentAccountsCounterMutex.Lock()
	paymentAccountsCounter = 0
	paymentAccountsCounterMutex.Unlock()

}
func ResetPoolPlanCounter() {
	poolPlanCounterMutex.Lock()
	poolPlanCounter = 0
	poolPlanCounterMutex.Unlock()
}
func ResetPoolUsedCounter() {
	poolUsedCounterMutex.Lock()
	poolUsedCounter = 0
	poolUsedCounterMutex.Unlock()
}
