package server

// import (
// 	"context"
// 	"strconv"
// 	"time"

// 	chainloader "github.com/rebear077/changan/chaininfos"
// 	"github.com/rebear077/changan/client"
// 	"github.com/rebear077/changan/conf"
// 	sql "github.com/rebear077/changan/sqlController"
// 	"github.com/sirupsen/logrus"
// )

// var chaininfos = chainloader.NewChainlog()

// const indent = "  "

// type Monitor struct {
// 	conn      *client.Client
// 	sql       *sql.SqlCtr
// 	pendingTX []byte
// }

// func NewMonitor() *Monitor {
// 	configs, err := conf.ParseConfigFile("./configs/config.toml")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	config := &configs[0]
// 	client, err := client.Dial(config)
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	sql := sql.NewSqlCtr()
// 	return &Monitor{
// 		conn: client,
// 		sql:  sql,
// 	}
// }

// func (m *Monitor) Start() {
// 	for {
// 		ticker1 := time.NewTicker(2 * time.Second)
// 		ctx, cancle := context.WithCancel(context.Background())
// 		go m.getInfor(ctx)
// 		select {
// 		case <-ticker1.C:
// 			cancle()
// 		}
// 	}

// }
// func (m *Monitor) getInfor(ctx context.Context) {
// 	ticker := time.NewTicker(10 * time.Second)
// 	select {
// 	case <-ctx.Done():
// 		return
// 	case <-ticker.C:
// 		chainID, err := m.conn.GetChainID(context.Background())
// 		if err != nil {
// 			// logrus.Errorln("监控器获取链ID失败:", err)
// 			chaininfos.Errorln("监控器获取链ID失败:", err)
// 		}
// 		// logrus.Infoln("区块链ID:", chainID)
// 		chaininfos.Infoln("区块链ID:", chainID)
// 		txCount, err := m.conn.GetTotalTransactionCount(context.Background())
// 		if err != nil {
// 			logrus.Errorln(err)
// 		}
// 		txNum, err := strconv.ParseInt(txCount.TxSum[2:], 16, 64)
// 		if err != nil {
// 			// logrus.Errorln("监控器获取区块链高度失败:", err)
// 			chaininfos.Errorln("监控器获取区块链高度失败:", err)
// 		}
// 		// logrus.Infoln("交易数量:", txNum)
// 		chaininfos.Infoln("交易数量:", txNum)
// 		pendingSize, err := m.conn.GetPendingTxSize(context.Background())
// 		if err != nil {
// 			// logrus.Errorln("监控器获取未上链交易数量失败:", err)
// 			chaininfos.Errorln("监控器获取未上链交易数量失败:", err)
// 		}
// 		pending, err := strconv.ParseInt(string(pendingSize)[3:len(pendingSize)-1], 16, 64)
// 		if err != nil {
// 			// logrus.Errorln("监控器获取未上链交易数量失败:", err)
// 			chaininfos.Errorln("监控器获取未上链交易数量失败:", err)
// 		}
// 		// logrus.Infof("交易池中未上链交易数量:%d\n", pending)
// 		chaininfos.Infof("交易池中未上链交易数量:%d", pending)
// 		m.pendingTX = pendingSize

// 	}
// }
// func (m *Monitor) VerifyChainStatus() bool {
// 	if string(m.pendingTX) != "0x0" {
// 		return true
// 	} else {
// 		return false
// 	}
// }
