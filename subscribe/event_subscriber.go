package subscriber

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"ethereum/go-ethereum/common"
	"ethereum/go-ethereum/crypto"

	"github.com/rebear077/changan/abi"
	"github.com/rebear077/changan/client"
	"github.com/rebear077/changan/conf"
	"github.com/rebear077/changan/core/types"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	Events           map[string][]TransactionInfor
	eventToSolMethod map[string]string
	stateMutex       sync.RWMutex
	subscribeClient  *client.Client
}

// 接口类型转string类型
func Strval(value interface{}) string {

	var key string
	if value == nil {
		return key
	}

	switch v := value.(type) {
	case float64:
		ft := v
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := v
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := v
		key = strconv.Itoa(it)
	case uint:
		it := v
		key = strconv.Itoa(int(it))
	case int8:
		it := v
		key = strconv.Itoa(int(it))
	case uint8:
		it := v
		key = strconv.Itoa(int(it))
	case int16:
		it := v
		key = strconv.Itoa(int(it))
	case uint16:
		it := v
		key = strconv.Itoa(int(it))
	case int32:
		it := v
		key = strconv.Itoa(int(it))
	case uint32:
		it := v
		key = strconv.Itoa(int(it))
	case int64:
		it := v
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := v
		key = strconv.FormatUint(it, 10)
	case string:
		key = v
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
func NewSubscriber() *Subscriber {
	mapEventToSol := map[string]string{
		"IssueAPChannelInfo(string,string)":      "IssueAPChannelInfo",
		"UpdateAPChannelInfo(string,string)":     "UpdateAPChannelInfo",
		"IssueBidingPriceInfo(string,string)":    "IssueBidingPriceInfo",
		"IssueChannelDealInfo(string,string)":    "IssueChannelDealInfo",
		"IssueChannelSwitchInfo(string,string)":  "IssueChannelSwitchInfo",
		"AutoUpdateAPChannelInfo(string,string)": "AutoUpdateAPChannelInfo",
	}
	return &Subscriber{
		Events:           make(map[string][]TransactionInfor),
		eventToSolMethod: mapEventToSol,
	}
}

// init the eventlogs params alonely
func InitEventLogParams(nodeURL string, groupID string, start string, end string, contractAddress string) types.EventLogParams {
	var eventLogParams types.EventLogParams
	eventLogParams.FromBlock = start
	eventLogParams.ToBlock = end
	eventLogParams.GroupID = groupID
	var addresses = make([]string, 1)
	addresses[0] = contractAddress
	eventLogParams.Addresses = addresses
	return eventLogParams
}
func (s *Subscriber) GetTxByMethod(eventName string) []TransactionInfor {
	return s.Events[eventName]
}

// subscribe a list of events by the name of method and return the struct of the parsed transactions or error
// the input params are nodeURL:for example 127.0.0.1:20200 start and end are the scope of searching events.contractAddress is the contraction adresss
func (s *Subscriber) GetEventsByMethodName(nodeURL string, groupID string, start string, end string, contractAddress string, needEvents string, abiFile string, signalChan chan bool) error {

	var txFinalList []TransactionInfor
	eventLogParams := InitEventLogParams(nodeURL, groupID, start, end, contractAddress)
	// endpoint := nodeURL
	// privateKey, _ := hex.DecodeString("145e247e170ba3afd6ae97e88f00dbc976c2345d511b0f6713355d19d8b80b58")
	// config := &conf.Config{IsHTTP: false, ChainID: 1, CAFile: "../../ca.crt", Key: "../../sdk.key", Cert: "../../sdk.crt",
	// 	IsSMCrypto: false, GroupID: 1, PrivateKey: privateKey, NodeURL: endpoint}
	const (
		indent = "  "
	)
	var topics = make([]string, 1)
	topics[0] = common.BytesToHash(crypto.Keccak256([]byte(needEvents))).Hex()
	eventLogParams.Topics = topics
	var addresses = make([]string, 1)
	addresses[0] = contractAddress
	eventLogParams.Addresses = addresses

	txChan := make(chan []TransactionInfor)
	done := make(chan bool)

	callBack := func(status int, logs []types.Log) {
		var txInfor TransactionInfor
		var txList []TransactionInfor
		for _, v := range logs {
			var tempABI abi.ABI

			tempABI, err := abi.JSON(strings.NewReader(abiFile))
			if err != nil {
				logrus.Println(err.Error())
			}
			switch needEvents {
			case "IssueAPChannelInfo(string,string)":
				var temp IssueAPChannelInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			case "UpdateAPChannelInfo(string,string)":
				var temp UpdateAPChannelInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			case "IssueBidingPriceInfo(string,string)":
				var temp IssueBidingPriceInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			case "IssueChannelDealInfo(string,string)":
				var temp IssueChannelDealInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			case "IssueChannelSwitchInfo(string,string)":
				var temp IssueChannelSwitchInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			case "AutoUpdateAPChannelInfo(string,string)":
				var temp AutoUpdateAPChannelInfo
				err = tempABI.Unpack(&temp, s.eventToSolMethod[needEvents], v.Data)
				if err != nil {
					logrus.Println(err.Error())
				}
				txInfor.Data = temp
			}
			logRes, err := json.MarshalIndent(v, "", indent)
			if err != nil {
				fmt.Printf("logs marshalIndent error: %v", err)
			}
			var temp1 types.Log
			if err := json.Unmarshal(logRes, &temp1); err != nil {
				fmt.Println("UnMarshal is err：", err)
			}
			txInfor.BlockNumber = temp1.BlockNumber
			txInfor.TransactionHash = temp1.TxHash
			txList = append(txList, txInfor)
		}
		txChan <- txList
		done <- true
	}

	err := s.subscribeClient.SubscribeEventLogs(eventLogParams, callBack)
	if err != nil {
		logrus.Printf("subscribe event failed, err: %v\n", err)
		return err
	}

	// go func() {
	// 	eventGetTicker := time.NewTicker(2 * time.Second)
	// 	for range eventGetTicker.C {
	// 		err := s.subscribeClient.SubscribeEventLogs(eventLogParams, callBack)
	// 		if err != nil {
	// 			logrus.Printf("subscribe event failed, err: %v\n", err)
	// 		}
	// 	}
	// }()

	for {
		select {
		case log := <-txChan:
			txFinalList = append(txFinalList, log...)
		case <-done:
			s.stateMutex.Lock()
			s.Events[needEvents] = append(s.Events[needEvents], txFinalList...)
			s.stateMutex.Unlock()
			txlength := len(txFinalList)
			txFinalList = make([]TransactionInfor, 0)
			if txlength != 0 {
				signalChan <- true
			}
		}
	}
}

func (sb *Subscriber) MarkEventRetrieved(eventType string) {
	sb.stateMutex.Lock()
	defer sb.stateMutex.Unlock()
	sb.Events[eventType] = make([]TransactionInfor, 0)
}

func (sb *Subscriber) InitDialClient(nodeURL string, caDirPath string) error {
	endpoint := nodeURL
	privateKey, _ := hex.DecodeString("145e247e170ba3afd6ae97e88f00dbc976c2345d511b0f6713355d19d8b80b58")

	if len(caDirPath) == 0 {
		caDirPath = "../.."
	}
	caCrtFile := caDirPath + "/" + "ca.crt"
	keyFile := caDirPath + "/" + "sdk.key"
	crtFile := caDirPath + "/" + "sdk.crt"
	config := &conf.Config{IsHTTP: false, ChainID: 1, CAFile: caCrtFile, Key: keyFile, Cert: crtFile,
		IsSMCrypto: false, GroupID: 1, PrivateKey: privateKey, NodeURL: endpoint}
	var c *client.Client
	var err error
	for i := 0; i < 3; i++ {
		c, err = client.Dial(config)
		if err != nil {
			logrus.Printf("init subscriber failed, err: %v, retrying\n", err)
			continue
		}
		break
	}
	if err != nil {
		logrus.Fatalf("init subscriber client failed, err: %v\n", err)
		return err
	}
	sb.subscribeClient = c
	return nil
}
