package gacha

type PoolInitStruct struct {
	ExchangeCnt int32
	Pieces      int32
	name        string
	NeedJing    bool
	IsAllGet    bool
}

type basePool struct {
	exchangeCnt int32
	pieces      int32
	gacha       []*GachaEntiy
	testCnt     int
	name        string
	needJing    bool
	iGachaPool  IGachaPool
	isAllGet    bool
	allGet      bool
}

func (this *basePool) Init(req *PoolInitStruct) {
	this.exchangeCnt = req.ExchangeCnt
	this.pieces = req.Pieces
	this.testCnt = TEST_CNT
	this.name = req.name
	if req.NeedJing {
		this.name += "@Jing"
	}
	if req.IsAllGet {
		this.name += "@AllGet"
	}
	this.needJing = req.NeedJing
	this.isAllGet = req.IsAllGet
	this.iGachaPool = this
}

func (this *basePool) AddGacha(entities []*GachaEntiy) {
	this.gacha = append(this.gacha, entities...)
}

func (this *basePool) NeedStop() bool {
	if this.needJing {
		return false
	}
	return this.allGet
}

func (this *basePool) BonusCallBack(totalGachaResult *GachaResult) {

}

func (this *basePool) JingCallBack(totalGachaResult *GachaResult, stop bool) {
	this.allGet = true
	totalGachaResult._charPieces += float64(this.pieces)
	if stop {
		totalGachaResult._pigPieces += float64(PIECES_3STARS)
	}
}

func (this *basePool) GetResult() *GachaResult {
	totalGachaResult := new(GachaResult)
	totalGachaResult._name = this.name

	for m := 0; m < this.testCnt; m++ {
		stop := false
		this.allGet = false
		//复刻抽两个最多2井
		for loop := 0; loop < 2; loop++ {
			//十连
			i := 0
			stopLoop := false
			for i = 0; i < int(this.exchangeCnt); i += 10 {
				//前九次
				for j := 0; j < 9; j++ {
					itemName := this.gacha[0].choicer.Pick()
					// fmt.Println(itemName)
					this.iGachaPool.CheckStatics(itemName, totalGachaResult)
					this.iGachaPool.BonusCallBack(totalGachaResult)
				}
				//第十次池子
				itemName := this.gacha[1].choicer.Pick()
				// fmt.Println(itemName)
				this.iGachaPool.CheckStatics(itemName, totalGachaResult)
				this.iGachaPool.BonusCallBack(totalGachaResult)
				//抽到了
				if this.iGachaPool.NeedStop() {
					totalGachaResult._cnt += float64(i + 10)
					totalGachaResult._allGet++
					stopLoop = true
					break
				}
			}

			//丼了
			if (stop && i == int(this.exchangeCnt-10)) || !stopLoop {
				this.iGachaPool.JingCallBack(totalGachaResult, stop)
				totalGachaResult._jing++
			}

			if this.allGet && !stopLoop {
				totalGachaResult._cnt += float64(int32(loop+1) * this.exchangeCnt)
				totalGachaResult._allGet++
				break
			}

			if stopLoop {
				break
			}
		}

	}
	totalGachaResult._charPieces /= float64(this.testCnt)
	totalGachaResult._pigPieces /= float64(this.testCnt)
	totalGachaResult._cnt /= float64(this.testCnt)
	totalGachaResult._jing /= float64(this.testCnt)
	totalGachaResult._heartBreak /= float64(this.testCnt)
	totalGachaResult._allGet /= float64(this.testCnt)
	return totalGachaResult
}

func (this *basePool) CheckStatics(itemName string, totalGachaResult *GachaResult) bool {
	switch itemName {
	case "Up":
		defer func() {
			this.allGet = true
		}()
		if this.allGet {
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

type replicaPool struct {
	basePool
	bonusGacha *GachaEntiy
	firstUp    bool
	secondUp   bool
}

func (this *replicaPool) Init(req *PoolInitStruct) {
	this.basePool.Init(req)
	this.basePool.iGachaPool = this

	this.intiBouns()
}

func (this *replicaPool) intiBouns() {
	bonusGacha := make([]*GachaPro, 6)
	pro := []int32{50, 100, 500, 1000, 2350, 6000}
	bonusGacha[0] = newGachaPro("1", pro[0])
	bonusGacha[1] = newGachaPro("2", pro[1])
	bonusGacha[2] = newGachaPro("3", pro[2])
	bonusGacha[3] = newGachaPro("4", pro[3])
	bonusGacha[4] = newGachaPro("5", pro[4])
	bonusGacha[5] = newGachaPro("6", pro[5])
	this.bonusGacha = newGachaEntity(bonusGacha)
}

func (this *replicaPool) BonusCallBack(totalGachaResult *GachaResult) {
	itemName := this.bonusGacha.choicer.Pick()
	this.checkBonusStatics(itemName, totalGachaResult)
}

func (this *replicaPool) checkBonusStatics(itemName string, totalGachaResult *GachaResult) bool {
	switch itemName {
	case "1":
		totalGachaResult._charPieces += BONUS_1_CHARPIECES
		totalGachaResult._pigPieces += BONUS_1_PIGPIECES
		totalGachaResult._heartBreak += BONUS_1_HEARTBREAK
	case "2":
		totalGachaResult._charPieces += BONUS_2_CHARPIECES
		totalGachaResult._pigPieces += BONUS_2_PIGPIECES
		totalGachaResult._heartBreak += BONUS_2_HEARTBREAK
	case "3":
		totalGachaResult._charPieces += BONUS_3_CHARPIECES
		totalGachaResult._pigPieces += BONUS_3_PIGPIECES
		totalGachaResult._heartBreak += BONUS_3_HEARTBREAK
	case "4":
		totalGachaResult._charPieces += BONUS_4_CHARPIECES
		totalGachaResult._pigPieces += BONUS_4_PIGPIECES
		totalGachaResult._heartBreak += BONUS_4_HEARTBREAK
	case "5":
		totalGachaResult._charPieces += BONUS_5_CHARPIECES
		totalGachaResult._pigPieces += BONUS_5_PIGPIECES
		totalGachaResult._heartBreak += BONUS_5_HEARTBREAK
	case "6":
		totalGachaResult._charPieces += BONUS_6_CHARPIECES
		totalGachaResult._pigPieces += BONUS_6_PIGPIECES
		totalGachaResult._heartBreak += BONUS_6_HEARTBREAK
	}
	return false
}

func (this *replicaPool) CheckStatics(itemName string, totalGachaResult *GachaResult) bool {
	switch itemName {
	case "Up":
		if this.firstUp {
			totalGachaResult._pigPieces += PIECES_3STARS
		} else {
			this.firstUp = true
			//要求全获得并且第二个获得 或者 不要求全获得
			if this.isAllGet && this.secondUp {
				this.allGet = true
				return true
			}
			if !this.isAllGet {
				this.allGet = true
				return true
			}
		}
	case "SecondUp":
		if this.secondUp {
			totalGachaResult._pigPieces += PIECES_3STARS
		} else {
			this.secondUp = true
			//要求全获得并且第一个获得 或者 不要求全获得
			if this.isAllGet && this.firstUp {
				this.allGet = true
				return true
			}
		}
	case "3Stars", "2Stars", "1Stars":
		this.basePool.CheckStatics(itemName, totalGachaResult)
	}
	return false
}

func (this *replicaPool) JingCallBack(totalGachaResult *GachaResult, stop bool) {
	defer func() {
		if !this.isAllGet && (this.firstUp || this.secondUp) {
			this.allGet = true
		}
		if this.allGet {
			this.firstUp = false
			this.secondUp = false
		}
	}()
	if this.allGet {
		totalGachaResult._pigPieces += PIECES_3STARS
		return
	}
	if this.secondUp || this.firstUp {
		this.allGet = true
		return
	}

	this.firstUp = true

}

func (this *replicaPool) NeedStop() bool {
	if this.needJing || !this.allGet {
		return false
	}
	return this.allGet
}
