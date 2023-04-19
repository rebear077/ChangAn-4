package errorhandle

import "sync"

var ERRDealer *ErrorHandler

func init() {
	ERRDealer = NewErrorHandler()
}

type ErrorHandler struct {
	errorPools map[string]map[string]interface{}
	mutex      sync.Mutex
}

func NewErrorHandler() *ErrorHandler {
	errorpool := make(map[string]map[string]interface{})
	return &ErrorHandler{
		errorPools: errorpool,
	}
}
func (err *ErrorHandler) InsertError(poolName string, hash string, encryptedData interface{}) {
	subPool := make(map[string]interface{})
	subPool[hash] = encryptedData
	err.mutex.Lock()
	err.errorPools[poolName] = subPool
	err.mutex.Unlock()
}
func (err *ErrorHandler) DeleteError(poolName string) {
	err.mutex.Lock()
	delete(err.errorPools, poolName)
	err.mutex.Unlock()
}

func (err *ErrorHandler) GetErrorLength(poolName string) int {
	return len(err.errorPools[poolName])
}
func (err *ErrorHandler) GetErrorInfo(poolName string) map[string]interface{} {
	return err.errorPools[poolName]
}
