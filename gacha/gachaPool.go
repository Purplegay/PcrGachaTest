package gacha

import (
	"sync"
)

const (
	MAXCNT = int32(10000)
)

//普通up池子
var gachaEntity1 *GachaEntiy
var gachaEntity10th *GachaEntiy

//3星翻倍卡池
var gachaEntity3StarsUp *GachaEntiy
var gachaEntity3StarsUp10th *GachaEntiy

//up翻倍+3星池子
var gachaEntityUpAnd3StarsUp *GachaEntiy
var gachaEntityUpAnd3StarsUp10th *GachaEntiy

var once sync.Once
var gachaPoolMgr *GachaPoolMgr

type GachaPoolMgr struct {
}

func GetInstance() *GachaPoolMgr {
	once.Do(func() {
		gachaPoolMgr = new(GachaPoolMgr)
		gachaPoolMgr.init()
	})
	return gachaPoolMgr
}

//exchangeCnt 兑换次数 pieces 角色碎片数
//MODE 模式 参考gacha_define.go
func (this *GachaPoolMgr) NewGachaPool(exchangeCnt, pieces int32, MODE int32) IGachaPool {
	pool := newGachaPool(exchangeCnt, pieces, MODE)
	return pool
}

func (this *GachaPoolMgr) init() {

	//普通up池子
	pool1Charas := make([]*GachaPro, 4)
	pro := []int32{70, 300, 1800, 7900}
	pool1Charas[0] = newGachaPro("Up", pro[0])
	pool1Charas[1] = newGachaPro("3Stars", pro[1]-pro[0])
	pool1Charas[2] = newGachaPro("2Stars", pro[2])
	pool1Charas[3] = newGachaPro("1Stars", pro[3])
	gachaEntity1 = newGachaEntity(pool1Charas)

	//普通池子第10次抽取
	poolThe10thCharas := make([]*GachaPro, 3)
	poolThe10thCharas[0] = newGachaPro("Up", pro[0])
	poolThe10thCharas[1] = newGachaPro("3Stars", pro[1]-pro[0])
	poolThe10thCharas[2] = newGachaPro("2Stars", MAXCNT-pro[1])
	gachaEntity10th = newGachaEntity(poolThe10thCharas)

	//3星翻倍卡池
	poolThreeStarsUpCharas := make([]*GachaPro, 4)
	poolThreeStarsUpCharas[0] = newGachaPro("Up", pro[0])
	poolThreeStarsUpCharas[1] = newGachaPro("3Stars", pro[1]*2-pro[0])
	poolThreeStarsUpCharas[2] = newGachaPro("2Stars", pro[2])
	poolThreeStarsUpCharas[3] = newGachaPro("1Stars", MAXCNT-pro[2]-pro[1]*2)
	gachaEntity3StarsUp = newGachaEntity(poolThreeStarsUpCharas)

	//3星翻倍卡池第10次抽取
	poolThreeStarsUpThe10thCharas := make([]*GachaPro, 3)
	poolThreeStarsUpThe10thCharas[0] = newGachaPro("Up", pro[0])
	poolThreeStarsUpThe10thCharas[1] = newGachaPro("3Stars", pro[1]*2-pro[0])
	poolThreeStarsUpThe10thCharas[2] = newGachaPro("2Stars", MAXCNT-pro[1]*2)
	gachaEntity3StarsUp10th = newGachaEntity(poolThreeStarsUpThe10thCharas)

	//up翻倍+3星池子
	poolUpAndThreeStarsUpCharas := make([]*GachaPro, 4)
	poolUpAndThreeStarsUpCharas[0] = newGachaPro("Up", pro[0]*2)
	poolUpAndThreeStarsUpCharas[1] = newGachaPro("3Stars", pro[1]*2-pro[0]*2)
	poolUpAndThreeStarsUpCharas[2] = newGachaPro("2Stars", pro[2])
	poolUpAndThreeStarsUpCharas[3] = newGachaPro("1Stars", MAXCNT-pro[2]-pro[1]*2)
	gachaEntityUpAnd3StarsUp = newGachaEntity(poolUpAndThreeStarsUpCharas)

	//up翻倍+3星池子第10次抽取
	poolUpAndThreeStarsUpThe10thCharas := make([]*GachaPro, 3)
	poolUpAndThreeStarsUpThe10thCharas[0] = newGachaPro("Up", pro[0]*2)
	poolUpAndThreeStarsUpThe10thCharas[1] = newGachaPro("3Stars", pro[1]*2-pro[0]*2)
	poolUpAndThreeStarsUpThe10thCharas[2] = newGachaPro("2Stars", MAXCNT-pro[1]*2)
	gachaEntityUpAnd3StarsUp10th = newGachaEntity(poolUpAndThreeStarsUpThe10thCharas)

}

type IGachaPool interface {
	GetResult() *GachaResult
}

type basePool struct {
	exchangeCnt int32
	pieces      int32
	gacha       []*GachaEntiy
	testCnt     int
	name        string
}

func (this *basePool) init(exchangeCnt, pieces int32, name string) {
	this.exchangeCnt = exchangeCnt
	this.pieces = pieces
	this.testCnt = TEST_CNT
	this.name = name
}

func (this *basePool) GetResult() *GachaResult {
	totalGachaResult := new(GachaResult)
	totalGachaResult._name = this.name

	for m := 0; m < this.testCnt; m++ {
		stop := false
		//十连
		i := 0
		for i = 0; i < int(this.exchangeCnt); i += 10 {
			//前九次
			for j := 0; j < 9; j++ {
				itemName := this.gacha[0].choicer.Pick()
				// fmt.Println(itemName)
				stop = this.checkStatics(itemName, totalGachaResult, stop) || stop
			}
			//第十次池子
			itemName := this.gacha[1].choicer.Pick()
			// fmt.Println(itemName)
			stop = this.checkStatics(itemName, totalGachaResult, stop) || stop
			//抽到了
			if stop {
				totalGachaResult._cnt += float64(i + 10)
				break
			}
		}
		//丼了
		if (stop && i == int(this.exchangeCnt-10)) || !stop {
			totalGachaResult._cnt += float64(this.exchangeCnt)
			totalGachaResult._charPieces += float64(this.pieces)
			totalGachaResult._jing++
		}
	}
	totalGachaResult._charPieces /= float64(this.testCnt)
	totalGachaResult._pigPieces /= float64(this.testCnt)
	totalGachaResult._cnt /= float64(this.testCnt)
	totalGachaResult._jing /= float64(this.testCnt)
	return totalGachaResult
}

func (this *basePool) checkStatics(itemName string, totalGachaResult *GachaResult, stop bool) bool {
	switch itemName {
	case "Up":
		if stop {
			totalGachaResult._charPieces += float64(this.pieces)
			totalGachaResult._pigPieces += PIECES_3STARS
			return false
		} else {
			totalGachaResult._charPieces += float64(this.pieces)
			return true
		}

	case "3Stars":
		totalGachaResult._pigPieces += PIECES_3STARS
	case "2Stars":
		totalGachaResult._pigPieces += PIECES_2STARS
	case "1Stars":
		totalGachaResult._pigPieces += PIECES_1STARS
	}
	return false
}
func newNormalPool(exchangeCnt, pieces int32) IGachaPool {
	pool := new(basePool)
	pool.init(exchangeCnt, pieces, "NormalPool")
	pool.gacha = []*GachaEntiy{gachaEntity1, gachaEntity10th}
	return pool
}

func newThreeStarsUpPool(exchangeCnt, pieces int32) IGachaPool {
	pool := new(basePool)
	pool.init(exchangeCnt, pieces, "ThreeStarsUpPool")
	pool.gacha = []*GachaEntiy{gachaEntity3StarsUp, gachaEntity3StarsUp10th}
	return pool
}

func newUpAndThreeStarsUpPool(exchangeCnt, pieces int32) IGachaPool {
	pool := new(basePool)
	pool.init(exchangeCnt, pieces, "UpAndThreeStarsUpPool")
	pool.gacha = []*GachaEntiy{gachaEntityUpAnd3StarsUp, gachaEntityUpAnd3StarsUp10th}
	return pool
}

func newGachaPool(exchangeCnt, pieces int32, MODE int32) IGachaPool {
	switch MODE {
	case GACHA_MODE_NORMAL:
		return newNormalPool(exchangeCnt, pieces)
	case GACHA_MODE_THREESTARSUP:
		return newThreeStarsUpPool(exchangeCnt, pieces)
	case GACHA_MODE_UPANDTHREES:
		return newUpAndThreeStarsUpPool(exchangeCnt, pieces)
	}

	return newNormalPool(exchangeCnt, pieces)
}
