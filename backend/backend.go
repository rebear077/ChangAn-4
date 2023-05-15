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
func (s *Server) IssueInvoiceInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueInvoiceInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (s *Server) VerifyAndUpdateInvoiceInformation(UUID, id, hash, owner string) error {
	err := s.ctr.VerifyAndUpdateInvoiceInformation(UUID, id, string(hash), owner)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之入库信息
func (s *Server) IssueHistoricalUsedInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalUsedInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之结算信息
func (s *Server) IssueHistoricalSettleInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalSettleInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之订单信息
func (s *Server) IssueHistoricalOrderInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalOrderInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 历史交易信息之应收账款信息
func (s *Server) IssueHistoricalReceivableInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssueHistoricalReceivableInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产计划信息
func (s *Server) IssuePoolPlanInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssuePoolPlanInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 入池数据之供应商生产入库信息
func (s *Server) IssuePoolUsedInformation(UUID, id, params string, cipher, encryptionKey []byte) error {
	err := s.ctr.IssuePoolUsedInformation(UUID, id, params, string(cipher), string(encryptionKey))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 上传融资意向请求
func (s *Server) IssueSupplierFinancingApplication(UUID, id, customerID string, cipher, encryptionKey, signed []byte) error {
	err := s.ctr.IssueSupplierFinancingApplication(UUID, id, customerID, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 更新融资意向请求
func (s *Server) UpdateSupplierFinancingApplication(UUID, id, customerID string, cipher, encryptionKey, signed []byte) error {
	err := s.ctr.UpdateSupplierFinancingApplication(UUID, id, customerID, string(cipher), string(encryptionKey), string(signed))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 回款信息
func (s *Server) UpdateAndLockPushPaymentAccounts(UUID, idAndFinanceID string, cipher, encryptionKey, newHash, oldHash []byte) error {
	err := s.ctr.UpdateAndLockPushPaymentAccounts(UUID, idAndFinanceID, string(cipher), string(encryptionKey), string(newHash), string(oldHash))
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (s *Server) LockPaymentAccounts(UUID, id, financeID, hash string) error {
	err := s.ctr.LockPaymentAccounts(UUID, id, financeID, hash)
	if err != nil {
		return err
	} else {
		return nil
	}
}
