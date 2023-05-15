package promote

import (
	"sync"

	uptoChain "github.com/rebear077/changan/tochain"
)

func (p *Promoter) invoiceInfoWaiter(length int) {
	for {
		counter := 0
		uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
			uptoChain.InvoiceMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.InvoiceMapLock.Unlock()
			return true
		})
		if counter == length {
			p.DataApi.IssueInvoiceOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.InvoiceMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
}

func (p *Promoter) historicalOrderInfoWaiter(orderLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalOrderMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalOrderMapLock.Unlock()
			return true
		})
		if counter == orderLength {
			p.DataApi.IssueHistoricalOrderInfoOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.HistoricalOrderMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) historicalReveivableInfoWaiter(receivableLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalReceivableMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalReceivableMapLock.Unlock()
			return true
		})
		if counter == receivableLength {
			p.DataApi.IssueHistoricalReceivableInfoOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.HistoricalReceivableMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) historicalSettleInfoWaiter(settleLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalSettleMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalSettleMapLock.Unlock()
			return true
		})
		if counter == settleLength {
			p.DataApi.IssueHistoricalSettleInfoOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.HistoricalSettleMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) historicalUsedInfoWaiter(usedLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
			uptoChain.HistoricalUsedMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.HistoricalUsedMapLock.Unlock()
			return true
		})
		if counter == usedLength {
			p.DataApi.IssueHistoryUsedInfoOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.HistoricalUsedMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) enterPoolPlanInfoWaiter(planLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
			uptoChain.PoolPlanMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.PoolPlanMapLock.Unlock()
			return true
		})
		if counter == planLength {
			p.DataApi.IssueEnterPoolPlanOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.PoolPlanMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) enterPoolUsedInfoWaiter(usedLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
			uptoChain.PoolUsedMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.PoolUsedMapLock.Unlock()
			return true
		})
		if counter == usedLength {
			p.DataApi.IssueEnterPoolUsedOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.PoolUsedMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) accountsUpdateInfoWaiter(accountsLength int) {
	for {
		counter := 0
		uptoChain.UpdateAndLockAccountMap.Range(func(key, value interface{}) bool {
			uptoChain.UpdateAndLockAccountMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.UpdateAndLockAccountMapLock.Unlock()
			return true
		})
		if counter == accountsLength {
			p.DataApi.UpdateAndLockAccountOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.UpdateAndLockAccountMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
}
func (p *Promoter) accountsLockInfoWaiter(accountsLength int) {
	for {
		counter := 0
		uptoChain.LockAccountsMap.Range(func(key, value interface{}) bool {
			uptoChain.LockAccountsMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.LockAccountsMapLock.Unlock()
			return true
		})
		if counter == accountsLength {
			p.DataApi.LockAccountOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.LockAccountsMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
}
func (p *Promoter) financingApplicationInfoWaiter(applicationLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
			uptoChain.FinancingApplicationIssueMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.FinancingApplicationIssueMapLock.Unlock()
			return true
		})
		if counter == applicationLength {
			p.DataApi.FinancingIntentionIssueOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.FinancingApplicationIssueMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) modifyInvoiceInfoWaiter(modifyLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyInvoiceMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyInvoiceMapLock.Unlock()
			return true
		})
		if counter == modifyLength {
			p.DataApi.ModifyInvoiceOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyInvoiceMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) modifyFinancingInfoWaiter(applicationLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyFinancingMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyFinancingMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyFinancingMapLock.Unlock()
			return true
		})
		if counter == applicationLength {
			p.DataApi.ModifyFinancingOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyFinancingMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
func (p *Promoter) modifyInvoiceInfoWhenModifyApplicationWaiter(modifyLength int, wg *sync.WaitGroup) {
	for {
		counter := 0
		uptoChain.ModifyInvoiceWhenMFAMap.Range(func(key, value interface{}) bool {
			uptoChain.ModifyInvoiceWhenMFAMapLock.Lock()
			mapping := value.(map[string]*uptoChain.ResponseMessage)
			counter += len(mapping)
			for _, message := range mapping {
				if message.GetMessage() == "" {
					counter = 0
					break
				}
			}
			uptoChain.ModifyInvoiceWhenMFAMapLock.Unlock()
			return true
		})
		if counter == modifyLength {
			p.DataApi.ModifyInvoiceWhenFinancingOKChan <- struct{}{}
			for {
				flag := 0
				uptoChain.ModifyInvoiceWhenMFAMap.Range(func(key, value interface{}) bool {
					if key != nil {
						flag++
						return false
					}
					return true
				})
				if flag == 0 {
					break
				}
			}
			break
		}
	}
	wg.Done()
}
