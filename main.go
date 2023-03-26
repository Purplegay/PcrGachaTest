package main

import (
	"fmt"
	"gacha/gacha"
)

func main() {
	pool := gacha.GetInstance().NewGachaPool(200, 100, gacha.GACHA_MODE_NORMAL, false)
	result := pool.GetResult()
	fmt.Println(result.String())
}
