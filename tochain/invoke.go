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

var (
	supplierCounter             = 0
	invoiceCounter              = 0
	historicalUsedCounter       = 0
	historicalSettleCounter     = 0
	historicalOrderCounter      = 0
	historicalReceivableCounter = 0
	paymentAccountsCounter      = 0
	poolUsedCounter             = 0
	poolPlanCounter             = 0

	supplierCounterMutex             sync.Mutex
	invoiceCounterMutex              sync.Mutex
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
		var params string
		ret, err := parsed.UnpackInput("issueSupplierFinancingApplication", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			// logrus.Fatalln("解析失败")
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueSupplierFinancingApplicationPool(receipt.TransactionHash, params)
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
		var params string
		ret, err := parsed.UnpackInput("issueInvoiceInformationStorage", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			// logrus.Fatalln("解析失败")
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueInvoiceInformationStoragePool(receipt.TransactionHash, params)
	} else {
		invoiceCounterMutex.Lock()
		invoiceCounter += 1
		invoiceCounterMutex.Unlock()
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
	// fmt.Println(setedLines)
	if setedLines.Int64() != 1 {
		var params string
		ret, err := parsed.UnpackInput("updateInvoiceInformationStorage", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			// logrus.Fatalln("解析失败")
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueInvoiceInformationStoragePool(receipt.TransactionHash, params)
	} else {
		invoiceCounterMutex.Lock()
		invoiceCounter += 1
		invoiceCounterMutex.Unlock()
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
		var params string
		ret, err := parsed.UnpackInput("issueHistoricalUsedInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueHistoricalUsedInformationPool(receipt.TransactionHash, params)
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
		var params string
		ret, err := parsed.UnpackInput("issueHistoricalSettleInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueHistoricalSettleInformationPool(receipt.TransactionHash, params)
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
		var params string
		ret, err := parsed.UnpackInput("issueHistoricalOrderInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueHistoricalOrderInformationPool(receipt.TransactionHash, params)
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
		var params string
		ret, err := parsed.UnpackInput("issueHistoricalReceivableInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssueHistoricalReceivableInformationPool(receipt.TransactionHash, params)
	} else {
		historicalReceivableCounterMutex.Lock()
		historicalReceivableCounter += 1
		historicalReceivableCounterMutex.Unlock()
	}
}

// 回款信息
func invokeIssuePushPaymentAccountsHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issuePushPaymentAccounts", receipt)
	if err != nil {
		log.Fatalf("error when transfer string to int: %v\n", err)
	}
	if setedLines.Int64() != 1 {
		var params string
		ret, err := parsed.UnpackInput("issuePushPaymentAccounts", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssuePushPaymentAccountsPool(receipt.TransactionHash, params)
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
		var params string
		ret, err := parsed.UnpackInput("issuePoolPlanInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssuePoolPlanInformationPool(receipt.TransactionHash, params)
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
		var params string

		ret, err := parsed.UnpackInput("issuePoolUsedInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			logs.Fatalln("解析失败")
		} else {
			params = parseRet[0].(string) + "," + parseRet[1].(string) + "," + parseRet[2].(string) + "," + parseRet[3].(string) + "," + parseRet[4].(string)
		}
		errorhandle.ERRDealer.InsertErrorIssuePoolUsedInformationPool(receipt.TransactionHash, params)
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
func QueryInvoiceSuccessCounter() int {
	invoiceCounterMutex.Lock()
	temp := invoiceCounter
	invoiceCounterMutex.Unlock()
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
func ResetInvoiceSuccessCounter() {
	invoiceCounterMutex.Lock()
	invoiceCounter = 0
	invoiceCounterMutex.Unlock()
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
