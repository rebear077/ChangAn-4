package receive

// 返回错误信息
func sucessCode() string {
	jsonData := []byte(`{
		"msg":"",
		"result": "{}",
		"code":"SUC000000"
	}`)
	return string(jsonData)
}
func failedCode() string {
	jsonData := []byte(`{
		"msg":"",
		"result": "{}",
		"code":"failed"
	}`)
	return string(jsonData)
}
func timeExceeded() string {
	jsonData := []byte(`{
		"msg":"时间戳超时",
		"result": "{}",
		"code":"ERR03"
	}`)
	return string(jsonData)
}

func verySignatureFailed() string {
	jsonData := []byte(`{
		"msg":"签名信息验证失败，请确认是否使用正确私钥签名！！",
		"result": "{}",
		"code":"ERR01"
	}`)
	return string(jsonData)
}

func wrongVerifyMethod() string {
	jsonData := []byte(`{
		"msg":"验证方法无效",
		"result": "{}",
		"code":"ERR02"
	}`)
	return string(jsonData)
}
func wrongVerifyInvoice() string {
	jsonData := []byte(`{
		"msg":"发票验证无效",
		"result": "{}",
		"code":"ERR02"
	}`)
	return string(jsonData)
}
func wrongJsonType() string {
	jsonData := []byte(`{
		"msg":"错误的Json格式",
		"result": "{}",
		"code":"ERR04"
	}`)
	return string(jsonData)
}
