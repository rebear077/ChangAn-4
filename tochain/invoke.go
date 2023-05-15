package uptoChain

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"ethereum/go-ethereum/common"

	"github.com/rebear077/changan/abi"
	smartcontract "github.com/rebear077/changan/contract"
	"github.com/rebear077/changan/core/types"
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

// 融资意向
func invokeIssueSupplierFinancingApplicationHandler(receipt *types.Receipt, err error) {

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issueSupplierFinancingApplication", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {

		ret, err := parsed.UnpackInput("issueSupplierFinancingApplication", common.FromHex(receipt.Input)[4:])
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
		FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			FinancingApplicationIssueMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			FinancingApplicationIssueMapLock.Unlock()
			FinancingApplicationIssueMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {

		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			FinancingApplicationIssueMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			FinancingApplicationIssueMapLock.Unlock()
			FinancingApplicationIssueMap.LoadOrStore(uuid, mapping)
			return true
		})
	}
}

// 更新融资意向
func invokeUpdateSupplierFinancingApplicationHandler(receipt *types.Receipt, err error) {

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "updateSupplierFinancingApplication", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {

		ret, err := parsed.UnpackInput("updateSupplierFinancingApplication", common.FromHex(receipt.Input)[4:])
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
		ModifyFinancingMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			ModifyFinancingMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			ModifyFinancingMapLock.Unlock()
			ModifyFinancingMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {

		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		ModifyFinancingMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			ModifyFinancingMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			ModifyFinancingMapLock.Unlock()
			ModifyFinancingMap.LoadOrStore(uuid, mapping)
			return true
		})
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
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
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
		InvoiceMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			InvoiceMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			InvoiceMapLock.Unlock()
			InvoiceMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		// fmt.Println(receipt.BlockHash)
		InvoiceMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			InvoiceMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				packedMessage := mapping[receipt.TransactionHash]
				packedMessage.ok = true
				packedMessage.message = message
				mapping[receipt.TransactionHash] = packedMessage
			}
			InvoiceMapLock.Unlock()
			InvoiceMap.LoadOrStore(uuid, mapping)
			return true
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
		fmt.Printf("error when transfer string to int: %v\n", err)
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("updateInvoiceInformationStorage", common.FromHex(receipt.Input)[4:])
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
		ModifyInvoiceMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			ModifyInvoiceMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			ModifyInvoiceMapLock.Unlock()
			ModifyInvoiceMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		ModifyInvoiceMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			ModifyInvoiceMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			ModifyInvoiceMapLock.Unlock()
			ModifyInvoiceMap.LoadOrStore(uuid, mapping)
			return true
		})
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
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalUsedInformation", common.FromHex(receipt.Input)[4:])
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
		HistoricalUsedMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalUsedMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalUsedMapLock.Unlock()
			HistoricalUsedMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		HistoricalUsedMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalUsedMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalUsedMapLock.Unlock()
			HistoricalUsedMap.LoadOrStore(uuid, mapping)
			return true
		})
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
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalSettleInformation", common.FromHex(receipt.Input)[4:])
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
		HistoricalSettleMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalSettleMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalSettleMapLock.Unlock()
			HistoricalSettleMap.LoadOrStore(uuid, mapping)

			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		HistoricalSettleMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalSettleMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalSettleMapLock.Unlock()
			HistoricalSettleMap.LoadOrStore(uuid, mapping)

			return true
		})
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
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalOrderInformation", common.FromHex(receipt.Input)[4:])
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
		HistoricalOrderMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalOrderMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalOrderMapLock.Unlock()
			HistoricalOrderMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		HistoricalOrderMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalOrderMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalOrderMapLock.Unlock()
			HistoricalOrderMap.LoadOrStore(uuid, mapping)
			return true
		})
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
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issueHistoricalReceivableInformation", common.FromHex(receipt.Input)[4:])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("xxxxxxxxx")
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
		HistoricalReceivableMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalReceivableMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalReceivableMapLock.Unlock()
			HistoricalReceivableMap.LoadOrStore(uuid, mapping)

			return true
		})
	} else {
		fmt.Println("---------")
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		HistoricalReceivableMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			HistoricalReceivableMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			HistoricalReceivableMapLock.Unlock()
			HistoricalReceivableMap.LoadOrStore(uuid, mapping)
			return true
		})
	}
}

// 回款信息
func invokeUpdateAndLockPushPaymentAccountsHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "updateAndLockAccounts", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("updateAndLockAccounts", common.FromHex(receipt.Input)[4:])
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
		UpdateAndLockAccountMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			UpdateAndLockAccountMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			UpdateAndLockAccountMapLock.Unlock()
			UpdateAndLockAccountMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		UpdateAndLockAccountMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			UpdateAndLockAccountMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			UpdateAndLockAccountMapLock.Unlock()
			UpdateAndLockAccountMap.LoadOrStore(uuid, mapping)
			return true
		})
	}
}
func invokeLockPaymentAccountsHandler(receipt *types.Receipt, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "updateAndLockAccounts", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("updateAndLockAccounts", common.FromHex(receipt.Input)[4:])
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
		LockAccountsMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			LockAccountsMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			LockAccountsMapLock.Unlock()
			LockAccountsMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		LockAccountsMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			LockAccountsMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			LockAccountsMapLock.Unlock()
			LockAccountsMap.LoadOrStore(uuid, mapping)
			return true
		})
	}
}

// 入池数据之供应商生产计划信息
func invokeIssuePoolPlanInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		logrus.Errorf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issuePoolPlanInformation", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}
	if setedLines == nil || setedLines.Int64() != 1 {

		ret, err := parsed.UnpackInput("issuePoolPlanInformation", common.FromHex(receipt.Input)[4:])
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
		fmt.Println(message)
		fmt.Println(receipt.BlockHash)
		packedMessage := new(ResponseMessage)
		packedMessage.ok = false
		packedMessage.message = message
		PoolPlanMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			PoolPlanMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			PoolPlanMapLock.Unlock()
			PoolPlanMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		PoolPlanMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			PoolPlanMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			PoolPlanMapLock.Unlock()
			PoolPlanMap.LoadOrStore(uuid, mapping)
			return true
		})
	}
}

// 入池数据之供应商生产入库信息
func invokeIssuePoolUsedInformationHandler(receipt *types.Receipt, err error) {
	if err != nil {
		logrus.Errorf("%v\n", err)
		return
	}
	parsed, _ := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	setedLines, err := parseOutput(smartcontract.HostFactoryControllerABI, "issuePoolUsedInformation", receipt)
	if err != nil {
		log.Printf("error when transfer string to int: %v\n", err)
	}

	if setedLines == nil || setedLines.Int64() != 1 {
		ret, err := parsed.UnpackInput("issuePoolUsedInformation", common.FromHex(receipt.Input)[4:])
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
		PoolUsedMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			PoolUsedMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			PoolUsedMapLock.Unlock()
			PoolUsedMap.LoadOrStore(uuid, mapping)
			return true
		})
	} else {
		message := "success"
		packedMessage := new(ResponseMessage)
		packedMessage.ok = true
		packedMessage.message = message
		PoolUsedMap.Range(func(key, value interface{}) bool {
			uuid := key.(string)
			PoolUsedMapLock.Lock()
			mapping := value.(map[string]*ResponseMessage)
			if _, ok := mapping[receipt.TransactionHash]; ok {
				mapping[receipt.TransactionHash] = packedMessage
			}
			PoolUsedMapLock.Unlock()
			PoolUsedMap.LoadOrStore(uuid, mapping)
			return true
		})
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
