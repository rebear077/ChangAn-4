package server

import (
	"io/ioutil"
	"os"
	"time"

	encrypter "github.com/rebear077/changan/encryption"
	sql "github.com/rebear077/changan/sqlController"
	uptoChain "github.com/rebear077/changan/tochain"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctr      *uptoChain.Controller
	encrypte *encrypter.Encrypter
	sql      *sql.SqlCtr
	symKey   []byte
	pubKey   []byte
	priKey   []byte
}

func NewServer() *Server {
	ctr := uptoChain.NewController()
	en := encrypter.NewEncrypter()
	symkey, err := getSymKey("./configs/symPri.txt")
	if err != nil {
		logrus.Fatalln(err)
	}
	pubkey, err := getRSAPublicKey("./configs/public.pem")
	if err != nil {
		logrus.Fatalln(err)
	}
	prikey, err := getRSAPrivateKey("./configs/private.pem")
	if err != nil {
		logrus.Fatalln(err)
	}

	return &Server{
		ctr:      ctr,
		encrypte: en,
		sql:      sql.NewSqlCtr(),
		symKey:   symkey,
		pubKey:   pubkey,
		priKey:   prikey,
	}
}
func (s *Server) StartMonitor() {
	s.ctr.Start()
}
func (s *Server) VerifyChainstatus() bool {
	return s.ctr.VerifyChainStatus()
}
func (s *Server) DeployContract() string {
	res := s.ctr.DeployContract()
	return res
}
func (s *Server) ValidateHash(hash []byte, plain []byte) bool {
	resHash := s.encrypte.Signature(plain)
	if string(resHash) == string(hash) {
		return true
	} else {
		return false
	}
}
func getSymKey(path string) ([]byte, error) {
	filesymPrivate, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := filesymPrivate.Stat()
	if err != nil {
		return nil, err
	}
	symkey := make([]byte, stat.Size())
	filesymPrivate.Read(symkey)
	filesymPrivate.Close()
	return symkey, nil
}
func getRSAPublicKey(path string) ([]byte, error) {
	pubKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
func getRSAPrivateKey(path string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}
func (s *Server) IssuePubilcKey(role string) (bool, error) {
	res, err := s.ctr.IssuePublicKeyStorage(role, role, string(s.pubKey))
	if err != nil {
		return false, err
	}
	return res, nil

}

// 数据加密
func (s *Server) DataEncryption(data []byte) ([]byte, []byte, []byte, error) {
	cipher, err := s.encrypte.SymEncrypt(data, s.symKey)
	// fmt.Println("s.symKey:", string(s.symKey))
	if err != nil {
		logrus.Errorln("数据加密失败，退出")

		return nil, nil, nil, err
	}
	encryptionKey, err := s.encrypte.AsymEncrypt(s.symKey, s.pubKey)
	// fmt.Println("encryptionKey: ", encryptionKey)
	if err != nil {
		logrus.Infoln("数据加密失败，退出")
		return nil, nil, nil, err
	}
	signed := s.encrypte.Signature(data)
	return cipher, encryptionKey, signed, nil
}

// 发票信息
func (s *Server) IssueInvoiceInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueInvoiceInformation(id, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之入库信息
func (s *Server) IssueHistoricalUsedInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueHistoricalUsedInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之结算信息
func (s *Server) IssueHistoricalSettleInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueHistoricalSettleInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之订单信息
func (s *Server) IssueHistoricalOrderInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueHistoricalOrderInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之应收账款信息
func (s *Server) IssueHistoricalReceivableInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueHistoricalReceivableInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产计划信息
func (s *Server) IssuePoolPlanInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssuePoolPlanInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产入库信息
func (s *Server) IssuePoolUsedInformation(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssuePoolUsedInformation(id, time.Now().String()[0:19], string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 上传融资意向请求
func (s *Server) IssueSupplierFinancingApplication(id string, cipher []byte, encryptionKey []byte, signed []byte) error {
	err := s.ctr.IssueSupplierFinancingApplication(id, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 回款信息
func (s *Server) IssuePushPaymentAccount(id string, cipher []byte, encryptionKey []byte, signed []byte) error {

	err := s.ctr.IssuePushPaymentAccounts(id, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 插入日志
func (s *Server) InsertLog(level string, info string) {
	time := time.Now().String()[0:19]
	err := s.sql.InsertLogs(time, level, info)
	if err != nil {
		logrus.Errorln(err)
	}
}

// 插入日志
func (s *Server) InsertChainLog(level string, title string, info string) {
	time := time.Now().String()[0:19]
	err := s.sql.InserChainInfos(time, level, title, info)
	if err != nil {
		logrus.Errorln(err)
	}
}
