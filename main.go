package main

import (
	"fmt"
	"gacha/gacha"
)

func main() {
	//没做复刻池井完继续抽的情况
	pool := gacha.GetInstance().NewGachaPool(&gacha.PoolInitStruct{
		ExchangeCnt: 200,
		Pieces:      100,
		NeedJing:    false,
		IsAllGet:    false,
	}, gacha.GACHA_MODE_REPLICA)
	result := pool.GetResult()
	fmt.Println(result.String())
}
