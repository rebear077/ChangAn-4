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
		http.HandleFunc("/asl/universal/commmit-intention", promoter.DataApi.HandleFinancingIntentionWithSelectedInfos)
		http.HandleFunc("/asl/universal/update-lock-back-account", promoter.DataApi.HandleUpdateCollectionAccount)
		http.HandleFunc("/asl/universal/lock-back-account", promoter.DataApi.HandleLockAccount)
		http.HandleFunc("/asl/universal/modify-financing-intension", promoter.DataApi.HandleModifyFinancingIntentionWithSelectedInfos)
		// http.HandleFunc("/asl/universal/selected-to-application", promoter.DataApi.HandleCollectionAccount)
		err := http.ListenAndServeTLS(":8443", "connApi/confs/server.pem", "connApi/confs/server.key", nil)
		if err != nil {
			log.Fatalf("启动 HTTPS 服务器失败: %v", err)
		}
	}()
	promoter.Start()
}
