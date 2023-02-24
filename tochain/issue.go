package uptoChain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"ethereum/go-ethereum/common"

	"github.com/rebear077/changan/abi"
	chaininfos "github.com/rebear077/changan/chaininfos"
	"github.com/rebear077/changan/client"
	"github.com/rebear077/changan/conf"
	smartcontract "github.com/rebear077/changan/contract"
	logloader "github.com/rebear077/changan/logs"
	sql "github.com/rebear077/changan/sqlController"
	"github.com/sirupsen/logrus"
)

var Logs = logloader.NewLog()

type Controller struct {
	conn      *client.Client
	session   *smartcontract.HostFactoryControllerSession
	log       *sql.SqlCtr
	chaininfo *logrus.Logger
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

	chainld := chaininfos.NewChainlog()
	return &Controller{
		conn:      client,
		session:   hostFactoryControllerSession,
		log:       sql.NewSqlCtr(),
		chaininfo: chainld,
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
func (c *Controller) IssueSupplierFinancingApplication(id string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueSupplierFinancingApplication(invokeIssueSupplierFinancingApplicationHandler, id, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 上传发票信息
func (c *Controller) IssueInvoiceInformation(id string, data string, key string, hash string) error {
	// fmt.Println("key: ", []byte(key))
	// fmt.Println("hash: ", hash)
	_, err := c.session.AsyncIssueInvoiceInformationStorage(invokeIssueInvoiceInformationStorageHandler, id, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之入库信息
func (c *Controller) IssueHistoricalUsedInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueHistoricalUsedInformation(invokeIssueHistoricalUsedInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之结算信息
func (c *Controller) IssueHistoricalSettleInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueHistoricalSettleInformation(invokeIssueHistoricalSettleInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之订单信息
func (c *Controller) IssueHistoricalOrderInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueHistoricalOrderInformation(invokeIssueHistoricalOrderInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 历史交易信息之应收账款信息
func (c *Controller) IssueHistoricalReceivableInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssueHistoricalReceivableInformation(invokeIssueHistoricalReceivableInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 回款信息
func (c *Controller) IssuePushPaymentAccounts(id string, data string, key string, hash string) error {

	_, err := c.session.AsyncIssuePushPaymentAccounts(invokeIssuePushPaymentAccountsHandler, id, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 入池数据之供应商生产计划信息
func (c *Controller) IssuePoolPlanInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssuePoolPlanInformation(invokeIssuePoolPlanInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}

// 入池数据之供应商生产入库信息
func (c *Controller) IssuePoolUsedInformation(id string, time string, data string, key string, hash string) error {
	_, err := c.session.AsyncIssuePoolUsedInformation(invokeIssuePoolUsedInformationHandler, id, time, data, key, hash)
	if err != nil {
		return err
	}
	return nil
}
func (c *Controller) VerifyChainStatus() bool {
	if string(c.pendingTX) != "0x0" {
		return true
	} else {
		return false
	}
}
func (c *Controller) Start() {
	for {
		ticker1 := time.NewTicker(2 * time.Second)
		ctx, cancle := context.WithCancel(context.Background())
		go c.getInfor(ctx)
		select {
		case <-ticker1.C:
			cancle()
		}
	}

}
func (c *Controller) getInfor(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		chainID, err := c.conn.GetChainID(context.Background())
		if err != nil {
			// logrus.Errorln("监控器获取链ID失败:", err)
			c.chaininfo.Error("监控器获取链ID失败:", err)
		}
		// logrus.Infoln("区块链ID:", chainID)
		c.chaininfo.Infoln("区块链ID:", chainID)
		txCount, err := c.conn.GetTotalTransactionCount(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
		txNum, err := strconv.ParseInt(txCount.TxSum[2:], 16, 64)
		if err != nil {
			// logrus.Errorln("监控器获取区块链高度失败:", err)
			c.chaininfo.Errorln("监控器获取区块链高度失败:", err)
		}
		// logrus.Infoln("交易数量:", txNum)
		c.chaininfo.Infoln("交易数量:", txNum)
		pendingSize, err := c.conn.GetPendingTxSize(context.Background())
		if err != nil {
			// logrus.Errorln("监控器获取未上链交易数量失败:", err)
			c.chaininfo.Errorln("监控器获取未上链交易数量失败:", err)
		}
		pending, err := strconv.ParseInt(string(pendingSize)[3:len(pendingSize)-1], 16, 64)
		if err != nil {
			// logrus.Errorln("监控器获取未上链交易数量失败:", err)
			c.chaininfo.Errorln("监控器获取未上链交易数量失败:", err)
		}
		// logrus.Infof("交易池中未上链交易数量:%d\n", pending)
		c.chaininfo.Infof("交易池中未上链交易数量:%d", pending)
		c.pendingTX = pendingSize

	}
}
