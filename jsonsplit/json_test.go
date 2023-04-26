package jsonsplit

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	str := "数据时点,RMB,00890,120105199910221810,1000,[,,],[2022-09-19,1,1]"
	RES := SplitEnterpoolDataPool(str)
	fmt.Println(RES)
}
