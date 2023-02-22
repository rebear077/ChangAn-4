package main

import (
	uptoChain "github.com/rebear077/changan/tochain"
)

type Uploader struct {
	ctr *uptoChain.Controller
}

//	func getContractAddr(path string) (string, error) {
//		file, err := os.Open(path)
//		if err != nil {
//			return "", err
//		}
//		stat, _ := file.Stat()
//		addr := make([]byte, stat.Size())
//		_, err = file.Read(addr)
//		if err != nil {
//			return "", err
//		}
//		err = file.Close()
//		if err != nil {
//			return "", err
//		}
//		return string(addr), nil
//	}
func NewUploader() *Uploader {

	return &Uploader{
		ctr: uptoChain.NewController(),
	}
}

// func getRSAPublicKey(path string) ([]byte, error) {
// 	pubKey, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pubKey, nil
// }
// func getContractAddr(path string) (string, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	stat, _ := file.Stat()
// 	addr := make([]byte, stat.Size())
// 	_, err = file.Read(addr)
// 	if err != nil {
// 		return "", err
// 	}
// 	err = file.Close()
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(addr), nil
// }

// func main() {
// 	up := NewUploader()
// 	addr, err := getContractAddr("configs/contractAddress.txt")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	pubKey, err := getRSAPublicKey("configs/public.pem")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	up.ctr.IssuePublicKeyStorage(addr, "银行", "长安", string(pubKey))

// }
