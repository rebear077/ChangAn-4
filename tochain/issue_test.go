package uptoChain

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"ethereum/go-ethereum/common"

// 	"github.com/rebear077/changan/client"
// 	"github.com/rebear077/changan/conf"
// 	smartcontract "github.com/rebear077/changan/contract"
// 	"github.com/sirupsen/logrus"
// )

// type testController struct {
// 	conn *client.Client
// }

// func newController() *testController {
// 	configs, err := conf.ParseConfigFile("../configs/config.toml")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	config := &configs[0]
// 	client, err := client.Dial(config)
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	return &testController{
// 		conn: client,
// 	}
// }
// func (c *testController) deployContract() string {
// 	address, _, instance, err := smartcontract.DeployHostFactoryController(c.conn.GetTransactOpts(), c.conn) // deploy contract
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	_ = instance
// 	str := "../configs/contractAddress.txt"
// 	file, err := os.Create(str)
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	defer file.Close()
// 	_, err = file.WriteString(address.Hex())
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	fmt.Printf("合约部署成功，合约地址为%s，合约地址已写入./configs/contractAddress.txt文件中", address.Hex())

// 	return address.Hex()
// }
// func (c *testController) IssueInvoiceInformation(contractAddr string, id string, timeandtype string, data string, key string, hash string) (bool, error) {
// 	contractAddress := common.HexToAddress(contractAddr)
// 	instance, err := smartcontract.NewHostFactoryController(contractAddress, c.conn)
// 	if err != nil {
// 		return false, err
// 	}
// 	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *c.conn.GetCallOpts(), TransactOpts: *c.conn.GetTransactOpts()}
// 	_, err = HostFactoryControllerSession.AsyncIssueInvoiceInformationStorage(invokeIssueInvoiceInformationStorageHandler, id, timeandtype, data, key, hash)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
// func TestDeployContract(t *testing.T) {
// 	ctr := newController()
// 	res := ctr.deployContract()
// 	fmt.Println(res)
// }
// func TestIssueInvoiceInfo(t *testing.T) {
// 	ctr := newController()
// 	contractAddress := common.HexToAddress("0x8cf400be91bbf4e9dccd232d77ebbb5510681dcf")
// 	instance, err := smartcontract.NewHostFactoryController(contractAddress, ctr.conn)
// 	if err != nil {
// 		panic(err)
// 	}
// 	txCount, err := ctr.conn.GetTotalTransactionCount(context.Background())
// 	if err != nil {
// 		logrus.Errorln(err)
// 	}
// 	txNum, err := strconv.ParseInt(txCount.TxSum[2:], 16, 64)
// 	if err != nil {
// 		logrus.Errorln("监控器获取区块链高度失败:", err)
// 	}
// 	logrus.Infoln("交易数量:", txNum)
// 	HostFactoryControllerSession := &smartcontract.HostFactoryControllerSession{Contract: instance, CallOpts: *ctr.conn.GetCallOpts(), TransactOpts: *ctr.conn.GetTransactOpts()}
// 	start := time.Now()
// 	fmt.Println(start)
// 	for index := 0; index < 1; index++ {
// 		// ctr.IssueInvoiceInformation("0xCa6C5Bd02a0b4da0aD7cC62DcD10796930781a18", "xxx", "xxxxxxxxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxx")
// 		HostFactoryControllerSession.AsyncIssueHistoricalUsedInformation(invokeIssueHistoricalUsedInformationHandler, "xxx", "xxxxxxxxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxx", "xxxxxxx")
// 	}

// 	for {
// 		if QueryHistoricalUsedCounter() == 1 {
// 			end := time.Now()
// 			fmt.Println(end)
// 			dura := end.Sub(start).Seconds()
// 			fmt.Println(dura)
// 			txCountend, err := ctr.conn.GetTotalTransactionCount(context.Background())
// 			if err != nil {
// 				logrus.Errorln(err)
// 			}
// 			txNumend, err := strconv.ParseInt(txCountend.TxSum[2:], 16, 64)
// 			if err != nil {
// 				logrus.Errorln("监控器获取区块链高度失败:", err)
// 			}
// 			logrus.Infoln("交易数量:", txNumend)

// 			diff := txNumend - txNum
// 			logrus.Infof("共上链%d交易", diff)
// 			tps := dura / float64(diff) * 1000
// 			logrus.Infoln("单笔交易上链时延约为", tps, "毫秒")
// 			break
// 		}

// 	}

// }
