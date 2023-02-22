package main

//把公钥从数据库中读取出来，并形成文件
import (
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func GetContractAddr(path string) (string, error) {
	file, err := os.Open(path)
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

// func main() {
// 	ctr := uptoChain.NewController()
// 	addr, err := GetContractAddr("configs/contractAddress.txt")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	pubkey, err := ctr.QueryPublic(addr, "银行")
// 	if err != nil {
// 		logrus.Fatalln(err)
// 	}
// 	fmt.Printf("pubkey: %s\n", pubkey)
// 	//创建文件
// 	file, _ := os.Create("public2.pem")

// 	_, err1 := io.WriteString(file, pubkey)
// 	if err1 != nil {
// 		fmt.Println("write error", err1)
// 		return
// 	}

// }
