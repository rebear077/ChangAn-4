package main

import (
	"log"
	"net/http"

	promote "github.com/rebear077/changan/promoter"
)

func main() {

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
