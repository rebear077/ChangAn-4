package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	server "github.com/rebear077/changan/backend"
	promote "github.com/rebear077/changan/promoter"
)

func TestByte(t *testing.T) {
	a := "123,456"
	b := []byte(a)
	fmt.Println(b)
}

// func init() {
// 	path := "./logs/output.log"
// 	/* 日志轮转相关函数
// 	`WithLinkName` 为最新的日志建立软连接
// 	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
// 	WithMaxAge 和 WithRotationCount二者只能设置一个
// 	  `WithMaxAge` 设置文件清理前的最长保存时间
// 	  `WithRotationCount` 设置文件清理前最多保存的个数
// 	*/
// 	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
// 	writer, _ := rotatelogs.New(
// 		path+".%Y%m%d%H%M",
// 		rotatelogs.WithLinkName(path),
// 		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
// 		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
// 	)
// 	log.SetOutput(writer)
// 	//log.SetFormatter(&log.JSONFormatter{})
// }
// func TestAll(t *testing.T) {
// 	serve := server.NewServer()
// 	strs := server.EncodefromFinancingResult(mockData())
// 	for _, str := range strs {
// 		res, err := serve.Issue([]byte(str), "issueFinancingResultFeedback")
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(res, "..")
// 	}

// }
func TestDeploy(t *testing.T) {
	serve := server.NewServer()
	fmt.Println("...")
	serve.DeployContract()
}

func TestIssue(t *testing.T) {
	promoter := promote.NewPromoter()
	go func() {
		http.HandleFunc("/asl/universal/push-invoice-info", promoter.DataApi.HandleInvoiceInformation)
		http.HandleFunc("/asl/universal/push-history", promoter.DataApi.HandleTransactionHistory)
		http.HandleFunc("/asl/universal/caqc/push-inpool", promoter.DataApi.HandleEnterpoolData)
		http.HandleFunc("/asl/universal/commmit-intention", promoter.DataApi.HandleFinancingIntention)
		http.HandleFunc("/asl/universal/back-account-lock", promoter.DataApi.HandleCollectionAccount)
		err := http.ListenAndServeTLS(":8443", "connApi/confs/server.pem", "connApi/confs/server.key", nil)
		if err != nil {
			log.Fatalf("启动 HTTPS 服务器失败: %v", err)
		}
	}()
	promoter.Start()
}
