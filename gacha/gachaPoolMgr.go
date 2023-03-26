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
func (this *GachaPoolMgr) NewGachaPool(exchangeCnt, pieces int32, MODE int32, needJing bool) IGachaPool {
	pool := newGachaPool(exchangeCnt, pieces, MODE, needJing)
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
	NeedStop(stepStop bool) bool
	Init(exchangeCnt, pieces int32, name string, needJing bool)
	AddGacha([]*GachaEntiy)
}

func newNormalPool(exchangeCnt, pieces int32, needJing bool) IGachaPool {
	pool := new(basePool)
	pool.Init(exchangeCnt, pieces, "NormalPool", needJing)
	pool.AddGacha([]*GachaEntiy{gachaEntity1, gachaEntity10th})
	return pool
}

func newThreeStarsUpPool(exchangeCnt, pieces int32, needJing bool) IGachaPool {
	pool := new(basePool)
	pool.Init(exchangeCnt, pieces, "ThreeStarsUpPool", needJing)
	pool.AddGacha([]*GachaEntiy{gachaEntity3StarsUp, gachaEntity3StarsUp10th})
	return pool
}

func newUpAndThreeStarsUpPool(exchangeCnt, pieces int32, needJing bool) IGachaPool {
	pool := new(basePool)
	pool.Init(exchangeCnt, pieces, "UpAndThreeStarsUpPool", needJing)
	pool.AddGacha([]*GachaEntiy{gachaEntityUpAnd3StarsUp, gachaEntityUpAnd3StarsUp10th})
	return pool
}

func newGachaPool(exchangeCnt, pieces int32, MODE int32, needJing bool) IGachaPool {
	switch MODE {
	case GACHA_MODE_NORMAL:
		return newNormalPool(exchangeCnt, pieces, needJing)
	case GACHA_MODE_THREESTARSUP:
		return newThreeStarsUpPool(exchangeCnt, pieces, needJing)
	case GACHA_MODE_UPANDTHREES:
		return newUpAndThreeStarsUpPool(exchangeCnt, pieces, needJing)
	}

	return newNormalPool(exchangeCnt, pieces, needJing)
}
