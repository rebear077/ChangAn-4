package subscriber

import "ethereum/go-ethereum/common"

type IssueAPChannelInfo struct {
	ApchannelInfo string
	Addr          string
}
type UpdateAPChannelInfo struct {
	ApchannelInfo string
	Addr          string
}
type IssueBidingPriceInfo struct {
	BidingPrice string
	Time        string
}
type IssueChannelDealInfo struct {
	Channeldeal string
	Time        string
}
type IssueChannelSwitchInfo struct {
	Channelswitch string
	Time          string
}
type AutoUpdateAPChannelInfo struct {
	Buyeraddr  string
	Selleraddr string
}

type TransactionInfor struct {
	Data            interface{}
	BlockNumber     uint64
	TransactionHash common.Hash
}
