package server

import (
	"io/ioutil"
	"os"

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

func (s *Server) Signature(data []byte) []byte {
	signed := s.encrypte.Signature(data)
	return signed
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
func (s *Server) IssueInvoiceInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueInvoiceInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (s *Server) VerifyAndUpdateInvoiceInformation(uuid, id, hash, owner string) error {
	err := s.ctr.VerifyAndUpdateInvoiceInformation(uuid, id, string(hash), owner)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之入库信息
func (s *Server) IssueHistoricalUsedInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalUsedInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之结算信息
func (s *Server) IssueHistoricalSettleInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalSettleInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之订单信息
func (s *Server) IssueHistoricalOrderInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalOrderInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之应收账款信息
func (s *Server) IssueHistoricalReceivableInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalReceivableInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产计划信息
func (s *Server) IssuePoolPlanInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssuePoolPlanInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产入库信息
func (s *Server) IssuePoolUsedInformation(uuid, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssuePoolUsedInformation(uuid, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 上传融资意向请求
func (s *Server) IssueSupplierFinancingApplication(uuid, id, financingid string, cipher, encryptionKey, signed []byte) error {
	err := s.ctr.IssueSupplierFinancingApplication(uuid, id, financingid, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 回款信息
func (s *Server) UpdatePushPaymentAccount(uuid, id string, cipher, encryptionKey, signed []byte) error {
	err := s.ctr.UpdatePushPaymentAccounts(uuid, id, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// // 插入日志
// func (s *Server) InsertLog(level string, info string) {
// 	time := time.Now().String()[0:19]
// 	err := s.sql.InsertLogs(time, level, info)
// 	if err != nil {
// 		logrus.Errorln(err)
// 	}
// }

// // 插入日志
// func (s *Server) InsertChainLog(level string, title string, info string) {
// 	time := time.Now().String()[0:19]
// 	err := s.sql.InserChainInfos(time, level, title, info)
// 	if err != nil {
// 		logrus.Errorln(err)
// 	}
// }
