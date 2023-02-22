package errorhandle

import "sync"

var ERRDealer *ErrorHandler

func init() {
	ERRDealer = NewErrorHandler()
}

type ErrorHandler struct {
	errorIssueSupplierFinancingApplicationPool      map[string]string
	errorIssueSupplierFinancingApplicationPoolMutex sync.RWMutex

	errorIssueInvoiceInformationStoragePool  map[string]string
	errorIssueInvoiceInformationStorageMutex sync.RWMutex

	errorIssueHistoricalUsedInformationPool  map[string]string
	errorIssueHistoricalUsedInformationMutex sync.RWMutex

	errorIssueHistoricalSettleInformationPool  map[string]string
	errorIssueHistoricalSettleInformationMutex sync.RWMutex

	errorIssueHistoricalOrderInformationPool  map[string]string
	errorIssueHistoricalOrderInformationMutex sync.RWMutex

	errorIssueHistoricalReceivableInformationPool  map[string]string
	errorIssueHistoricalReceivableInformationMutex sync.RWMutex

	errorIssuePushPaymentAccountsPool  map[string]string
	errorIssuePushPaymentAccountsMutex sync.RWMutex

	errorIssuePoolPlanInformationPool  map[string]string
	errorIssuePoolPlanInformationMutex sync.RWMutex

	errorIssuePoolUsedInformationPool  map[string]string
	errorIssuePoolUsedInformationMutex sync.RWMutex
}

func NewErrorHandler() *ErrorHandler {
	errorIssueSupplierFinancingApplicationPool := make(map[string]string)
	errorIssueInvoiceInformationStoragePool := make(map[string]string)
	errorIssueHistoricalUsedInformationPool := make(map[string]string)
	errorIssueHistoricalSettleInformationPool := make(map[string]string)
	errorIssueHistoricalOrderInformationPool := make(map[string]string)
	errorIssueHistoricalReceivableInformationPool := make(map[string]string)
	errorIssuePushPaymentAccountsPool := make(map[string]string)
	errorIssuePoolPlanInformationPool := make(map[string]string)
	errorIssuePoolUsedInformationPool := make(map[string]string)
	return &ErrorHandler{
		errorIssueSupplierFinancingApplicationPool:    errorIssueSupplierFinancingApplicationPool,
		errorIssueInvoiceInformationStoragePool:       errorIssueInvoiceInformationStoragePool,
		errorIssueHistoricalUsedInformationPool:       errorIssueHistoricalUsedInformationPool,
		errorIssueHistoricalSettleInformationPool:     errorIssueHistoricalSettleInformationPool,
		errorIssueHistoricalOrderInformationPool:      errorIssueHistoricalOrderInformationPool,
		errorIssueHistoricalReceivableInformationPool: errorIssueHistoricalReceivableInformationPool,
		errorIssuePushPaymentAccountsPool:             errorIssuePushPaymentAccountsPool,
		errorIssuePoolPlanInformationPool:             errorIssuePoolPlanInformationPool,
		errorIssuePoolUsedInformationPool:             errorIssuePoolUsedInformationPool,
	}
}

func (e *ErrorHandler) InsertErrorIssueSupplierFinancingApplicationPool(key string, value string) {
	e.errorIssueSupplierFinancingApplicationPoolMutex.Lock()
	e.errorIssueSupplierFinancingApplicationPool[key] = value
	e.errorIssueSupplierFinancingApplicationPoolMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssueInvoiceInformationStoragePool(key string, value string) {
	e.errorIssueInvoiceInformationStorageMutex.Lock()
	e.errorIssueInvoiceInformationStoragePool[key] = value
	e.errorIssueInvoiceInformationStorageMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssueHistoricalUsedInformationPool(key string, value string) {
	e.errorIssueHistoricalUsedInformationMutex.Lock()
	e.errorIssueHistoricalUsedInformationPool[key] = value
	e.errorIssueHistoricalUsedInformationMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssueHistoricalSettleInformationPool(key string, value string) {
	e.errorIssueHistoricalSettleInformationMutex.Lock()
	e.errorIssueHistoricalSettleInformationPool[key] = value
	e.errorIssueHistoricalSettleInformationMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssueHistoricalOrderInformationPool(key string, value string) {
	e.errorIssueHistoricalOrderInformationMutex.Lock()
	e.errorIssueHistoricalOrderInformationPool[key] = value
	e.errorIssueHistoricalOrderInformationMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssueHistoricalReceivableInformationPool(key string, value string) {
	e.errorIssueHistoricalReceivableInformationMutex.Lock()
	e.errorIssueHistoricalReceivableInformationPool[key] = value
	e.errorIssueHistoricalReceivableInformationMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssuePushPaymentAccountsPool(key string, value string) {
	e.errorIssuePushPaymentAccountsMutex.Lock()
	e.errorIssuePushPaymentAccountsPool[key] = value
	e.errorIssuePushPaymentAccountsMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssuePoolPlanInformationPool(key string, value string) {
	e.errorIssuePoolPlanInformationMutex.Lock()
	e.errorIssuePoolPlanInformationPool[key] = value
	e.errorIssuePoolPlanInformationMutex.Unlock()
}
func (e *ErrorHandler) InsertErrorIssuePoolUsedInformationPool(key string, value string) {
	e.errorIssuePoolUsedInformationMutex.Lock()
	e.errorIssuePoolUsedInformationPool[key] = value
	e.errorIssuePoolUsedInformationMutex.Unlock()
}

func (e *ErrorHandler) DeleteErrorIssueSupplierFinancingApplicationPool() {
	e.errorIssueSupplierFinancingApplicationPoolMutex.Lock()
	for k := range e.errorIssueSupplierFinancingApplicationPool {
		delete(e.errorIssueSupplierFinancingApplicationPool, k)
	}
	e.errorIssueSupplierFinancingApplicationPoolMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssueInvoiceInformationStoragePool() {
	e.errorIssueInvoiceInformationStorageMutex.Lock()
	for k := range e.errorIssueInvoiceInformationStoragePool {
		delete(e.errorIssueInvoiceInformationStoragePool, k)
	}
	e.errorIssueInvoiceInformationStorageMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssueHistoricalUsedInformationPool() {
	e.errorIssueHistoricalUsedInformationMutex.Lock()
	for k := range e.errorIssueHistoricalUsedInformationPool {
		delete(e.errorIssueHistoricalUsedInformationPool, k)
	}
	e.errorIssueHistoricalUsedInformationMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssueHistoricalSettleInformationPool() {
	e.errorIssueHistoricalSettleInformationMutex.Lock()
	for k := range e.errorIssueHistoricalSettleInformationPool {
		delete(e.errorIssueHistoricalSettleInformationPool, k)
	}
	e.errorIssueHistoricalSettleInformationMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssueHistoricalOrderInformationPool() {
	e.errorIssueHistoricalOrderInformationMutex.Lock()
	for k := range e.errorIssueHistoricalOrderInformationPool {
		delete(e.errorIssueHistoricalOrderInformationPool, k)
	}
	e.errorIssueHistoricalOrderInformationMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssueHistoricalReceivableInformationPool() {
	e.errorIssueHistoricalReceivableInformationMutex.Lock()
	for k := range e.errorIssueHistoricalReceivableInformationPool {
		delete(e.errorIssueHistoricalReceivableInformationPool, k)
	}
	e.errorIssueHistoricalReceivableInformationMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssuePushPaymentAccountsPool() {
	e.errorIssuePushPaymentAccountsMutex.Lock()
	for k := range e.errorIssuePushPaymentAccountsPool {
		delete(e.errorIssuePushPaymentAccountsPool, k)
	}
	e.errorIssuePushPaymentAccountsMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssuePoolPlanInformationPool() {
	e.errorIssuePoolPlanInformationMutex.Lock()
	for k := range e.errorIssuePoolPlanInformationPool {
		delete(e.errorIssuePoolPlanInformationPool, k)
	}
	e.errorIssuePoolPlanInformationMutex.Unlock()
}
func (e *ErrorHandler) DeleteErrorIssuePoolUsedInformationPool() {
	e.errorIssuePoolUsedInformationMutex.Lock()
	for k := range e.errorIssuePoolUsedInformationPool {
		delete(e.errorIssuePoolUsedInformationPool, k)
	}
	e.errorIssuePoolUsedInformationMutex.Unlock()
}
func (e *ErrorHandler) GetSupplierFinancingApplicationPoolLength() int {
	return len(e.errorIssueSupplierFinancingApplicationPool)
}
func (e *ErrorHandler) GetInvoiceInfoPoolLength() int {
	return len(e.errorIssueInvoiceInformationStoragePool)
}
func (e *ErrorHandler) GetHistoricalUsedInfoPoolLength() int {
	return len(e.errorIssueHistoricalUsedInformationPool)
}
func (e *ErrorHandler) GetHistoricalSettleInfoPoolLength() int {
	return len(e.errorIssueHistoricalSettleInformationPool)
}
func (e *ErrorHandler) GetHistoricalOrderInfoPoolLength() int {
	return len(e.errorIssueHistoricalOrderInformationPool)
}
func (e *ErrorHandler) GetHistoricalReceivableInfoPoolLength() int {
	return len(e.errorIssueHistoricalReceivableInformationPool)
}
func (e *ErrorHandler) GetPushPaymentAccountPoolLength() int {
	return len(e.errorIssuePushPaymentAccountsPool)
}

func (e *ErrorHandler) GetPoolPlanInfoPoolLength() int {
	return len(e.errorIssuePoolPlanInformationPool)
}
func (e *ErrorHandler) GetPoolUsedInfoPoolLength() int {
	return len(e.errorIssuePoolUsedInformationPool)
}

func (e *ErrorHandler) QuerySupplierFinancingApplicationPool() map[string]string {
	return e.errorIssueSupplierFinancingApplicationPool
}
func (e *ErrorHandler) QueryInvoiceInfoPool() map[string]string {
	return e.errorIssueInvoiceInformationStoragePool
}
func (e *ErrorHandler) QueryHistoricalUsedInfoPool() map[string]string {
	return e.errorIssueHistoricalUsedInformationPool
}
func (e *ErrorHandler) QueryHistoricalSettleInfoPool() map[string]string {
	return e.errorIssueHistoricalSettleInformationPool
}
func (e *ErrorHandler) QueryHistoricalOrderInfoPool() map[string]string {
	return e.errorIssueHistoricalOrderInformationPool
}
func (e *ErrorHandler) QueryHistoricalReceivableInfoPool() map[string]string {
	return e.errorIssueHistoricalReceivableInformationPool
}
func (e *ErrorHandler) QueryPushPaymentAccountPool() map[string]string {
	return e.errorIssuePushPaymentAccountsPool
}

func (e *ErrorHandler) QueryPoolPlanInfoPool() map[string]string {
	return e.errorIssuePoolPlanInformationPool
}
func (e *ErrorHandler) QueryPoolUsedInfoPool() map[string]string {
	return e.errorIssuePoolUsedInformationPool
}
