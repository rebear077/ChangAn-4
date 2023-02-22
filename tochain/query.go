package uptoChain

import (
	"errors"
	"strings"

	"ethereum/go-ethereum/common"
	"github.com/rebear077/changan/abi"
	smartcontract "github.com/rebear077/changan/contract"
)

// 查询公钥
func (c *Controller) QueryPublic(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryPublicKey(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryPublicKey", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 历史交易信息之入库信息
func (c *Controller) QueryHIstoricalUsedInformation(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryHIstoricalUsedInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryHIstoricalUsedInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 历史交易信息之结算信息
func (c *Controller) QueryHIstoricalSettleInformation(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryHIstoricalSettleInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryHIstoricalSettleInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 历史交易信息之订单信息
func (c *Controller) QueryHIstoricalOrderInformation(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryHIstoricalOrderInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryHIstoricalOrderInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 历史交易信息之应收账款信息
func (c *Controller) QueryHIstoricalReceivableInformation(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryHIstoricalReceivableInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryHIstoricalReceivableInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 发票信息
func (c *Controller) QueryInvoiceInformation(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryInvoiceInformationInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryInvoiceInformationInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}

}

// 回款信息
func (c *Controller) QueryPushPaymentAccounts(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryPushPaymentAccountsInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryPushPaymentAccountsInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}

}

// 入池数据之供应商生产计划信息
func (c *Controller) QueryPoolPlanInfo(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryPoolPlanInfoInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryPoolPlanInfoInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 入池数据之供应商生产入库信息
func (c *Controller) QueryPoolUsedInfo(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QueryPoolUsedInfoInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "queryPoolUsedInfoInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

// 融资意向
func (c *Controller) QuerySupplierFinancingApplication(contractAddr string, id string) (string, error) {
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
	if err != nil {
		return "", err
	}
	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
	_, receipt, err := HostFactoryControllerSession.QuerySupplierFinancingApplicationInJson(id)
	if err != nil {
		return "", err
	}
	if receipt.GetErrorMessage() != "" {
		return "", errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return "", err
	}
	var ret string
	err = parse.Unpack(&ret, "querySupplierFinancingApplicationInJson", common.FromHex(receipt.Output))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}

}
