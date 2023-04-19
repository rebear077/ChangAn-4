package uptoChain

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"ethereum/go-ethereum/common"

	"github.com/rebear077/changan/abi"
	"github.com/rebear077/changan/client"
	"github.com/rebear077/changan/conf"
	smartcontract "github.com/rebear077/changan/contract"
	logloader "github.com/rebear077/changan/logs"
	sql "github.com/rebear077/changan/sqlController"
)

var Logs = logloader.NewLog()

type Controller struct {
	conn      *client.Client
	session   *smartcontract.HostFactoryControllerSession
	log       *sql.SqlCtr
	pendingTX []byte
}

func getContractAddr() (string, error) {
	file, err := os.Open("./configs/contractAddress.txt")
	if err != nil {
		return "", err
	}
	stat, _ := file.Stat()
	addr := make([]byte, stat.Size())
	_, err = file.Read(addr)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}
	return string(addr), nil
}

// 初始化
func NewController() *Controller {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}

	config := &configs[0]
	client, err := client.Dial(config)
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}
	contractAddr, err := getContractAddr()
	if err != nil {
		// logrus.Fatalln(contractAddr)
		Logs.Fatalln(contractAddr)
	}
	contractAddress := common.HexToAddress(contractAddr)
	instance, err := smartcontract.NewHostFactoryController(contractAddress, client)
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}
	hostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *client.GetCallOpts(), TransactOpts: *client.GetTransactOpts()}

	return &Controller{
		conn:    client,
		session: hostFactoryControllerSession,
		log:     sql.NewSqlCtr(),
	}
}

// 部署合约，写入configs/contractAddress.txt文件中
func (c *Controller) DeployContract() string {
	address, _, instance, err := smartcontract.DeployHostFactoryController(c.conn.GetTransactOpts(), c.conn) // deploy contract
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}
	_ = instance
	str := "./configs/contractAddress.txt"
	file, err := os.Create(str)
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}
	defer file.Close()
	_, err = file.WriteString(address.Hex())
	if err != nil {
		// logrus.Fatalln(err)
		Logs.Fatalln(err)
	}
	temp := fmt.Sprintf("合约部署成功，合约地址为%s，合约地址已写入./configs/contractAddress.txt文件中", address.Hex())
	c.log.InsertLogs(time.Now().String()[0:19], "info", temp)
	// logrus.Infof("合约部署成功，合约地址为%s，合约地址已写入./configs/contractAddress.txt文件中", address.Hex())
	logs.Infof("合约部署成功，合约地址为%s，合约地址已写入./configs/contractAddress.txt文件中", address.Hex())
	return address.Hex()
}

// 公钥上链
func (c *Controller) IssuePublicKeyStorage(id string, role string, key string) (bool, error) {
	_, receipt, err := c.session.IssuePublicKeyStorage(id, role, key)
	if err != nil {
		return false, err
	}
	if receipt.GetErrorMessage() != "" {
		return false, errors.New(receipt.GetErrorMessage())
	}
	parse, err := abi.JSON(strings.NewReader(smartcontract.HostFactoryControllerABI))
	if err != nil {
		return false, err
	}
	ret := big.NewInt(0)
	err = parse.Unpack(&ret, "issuePublicKeyStorage", common.FromHex(receipt.Output))
	if err != nil {
		return false, err
	}
	if ret.Cmp(big.NewInt(0)) == 1 {
		temp := fmt.Sprintf("func IssuePublicKeyStorage():,public key %s uploads to the block chain success", key)
		go c.log.InsertLogs(time.Now().String()[0:19], "debug", temp)
		return true, nil
	} else {
		return false, errors.New("smart contract error")
	}
}

// 上传融资意向请求
// 入口参数：id：供应商编号；financingid:融资意向申请id；data：加密后的数据；key：加密后的key值；hash：哈希值
func (c *Controller) IssueSupplierFinancingApplication(id string, financingid string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueSupplierFinancingApplication(invokeIssueSupplierFinancingApplicationHandler, id, financingid, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 上传发票信息
// id:供应商编号:发票的年月日日期(2022-09-08);params:复合参数包括time:发票时间；type:发票类型；num:发票后六位校验码；hash：哈希值；owner：空的
// data：加密后的数据；key：加密后的key值
func (c *Controller) IssueInvoiceInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssueInvoiceInformationStorage(invokeIssueInvoiceInformationStorageHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 验证并更新发票信息的owner字段
// 入口参数：id：供应商编号:发票的年月日;hash：发票哈希值；owner：融资意向申请编号
func (c *Controller) VerifyAndUpdateInvoiceInformation(id string, hash string, owner string) error {
	_, err := c.session.AsyncUpdateInvoiceInformationStorage(invokeVerifyAndUpdateInvoiceInformationStorageHandler, id, hash, owner)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之入库信息
// 入口参数：id：供应商id；params:复合参数包括tradeyearmonth:交易年月；financingid:融资意向申请编号；hash:哈希值；owner：空的
func (c *Controller) IssueHistoricalUsedInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssueHistoricalUsedInformation(invokeIssueHistoricalUsedInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之结算信息
func (c *Controller) IssueHistoricalSettleInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssueHistoricalSettleInformation(invokeIssueHistoricalSettleInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之订单信息
func (c *Controller) IssueHistoricalOrderInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssueHistoricalOrderInformation(invokeIssueHistoricalOrderInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之应收账款信息
func (c *Controller) IssueHistoricalReceivableInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssueHistoricalReceivableInformation(invokeIssueHistoricalReceivableInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 回款信息
// 长安业务服务器只负责修改回款账户信息
func (c *Controller) UpdatePushPaymentAccounts(id string, data string, key string, hash string) error {

	_, err := c.session.AsyncUpdatePushPaymentAccounts(invokeIssuePushPaymentAccountsHandler, id, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 入池数据之供应商生产计划信息
func (c *Controller) IssuePoolPlanInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssuePoolPlanInformation(invokeIssuePoolPlanInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}

// 入池数据之供应商生产入库信息
func (c *Controller) IssuePoolUsedInformation(id string, params string, data string, key string) error {
	_, err := c.session.AsyncIssuePoolUsedInformation(invokeIssuePoolUsedInformationHandler, id, params, data, key)
	if err != nil {
		return err
	}
	return nil
}
