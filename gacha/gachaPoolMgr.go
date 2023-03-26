package gacha

import (
	"sync"
)

const (
	MAXCNT = int32(10000)
)

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
func (this *GachaPoolMgr) NewGachaPool(req *PoolInitStruct, MODE int32) IGachaPool {
	pool := newGachaPool(req, MODE)
	return pool
}

type IGachaPool interface {
	GetResult() *GachaResult
	NeedStop() bool
	Init(req *PoolInitStruct)
	AddGacha([]*GachaEntiy)
	BonusCallBack(totalGachaResult *GachaResult)
	JingCallBack(totalGachaResult *GachaResult, stop bool)
	CheckStatics(itemName string, totalGachaResult *GachaResult, stop bool) bool
}

func newNormalPool(req *PoolInitStruct) IGachaPool {
	pool := new(basePool)
	req.name = "Normal"
	pool.Init(req)
	pool.AddGacha([]*GachaEntiy{gachaEntity1, gachaEntity10th})
	return pool
}

func newThreeStarsUpPool(req *PoolInitStruct) IGachaPool {
	pool := new(basePool)
	req.name = "ThreeStarsUp"
	pool.Init(req)
	pool.AddGacha([]*GachaEntiy{gachaEntity3StarsUp, gachaEntity3StarsUp10th})
	return pool
}

func newUpAndThreeStarsUpPool(req *PoolInitStruct) IGachaPool {
	pool := new(basePool)
	req.name = "UpAndThreeStarsUp"
	pool.Init(req)
	pool.AddGacha([]*GachaEntiy{gachaEntityUpAnd3StarsUp, gachaEntityUpAnd3StarsUp10th})
	return pool
}

func newReplicaPool(req *PoolInitStruct) IGachaPool {
	pool := new(replicaPool)
	req.name = "Replica"
	pool.Init(req)
	pool.AddGacha([]*GachaEntiy{gachaReplica, gachaReplica10th})
	return pool
}

func newGachaPool(req *PoolInitStruct, MODE int32) IGachaPool {
	switch MODE {
	case GACHA_MODE_NORMAL:
		return newNormalPool(req)
	case GACHA_MODE_THREESTARSUP:
		return newThreeStarsUpPool(req)
	case GACHA_MODE_UPANDTHREES:
		return newUpAndThreeStarsUpPool(req)
	case GACHA_MODE_REPLICA:
		return newReplicaPool(req)
	}

	return newNormalPool(req)
}
