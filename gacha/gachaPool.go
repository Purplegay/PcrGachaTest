package gacha

type basePool struct {
	exchangeCnt int32
	pieces      int32
	gacha       []*GachaEntiy
	testCnt     int
	name        string
	needJing    bool
}

func (this *basePool) Init(exchangeCnt, pieces int32, name string, needJing bool) {
	this.exchangeCnt = exchangeCnt
	this.pieces = pieces
	this.testCnt = TEST_CNT
	this.name = name
	if needJing {
		this.name += "@Jing"
	}
	this.needJing = needJing
}

func (this *basePool) AddGacha(entities []*GachaEntiy) {
	this.gacha = append(this.gacha, entities...)
}

func (this *basePool) NeedStop(stepStop bool) bool {
	if this.needJing {
		return false
	}
	return stepStop
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
			if this.NeedStop(stop) {
				totalGachaResult._cnt += float64(i + 10)
				break
			}
		}
		//丼了
		if (stop && i == int(this.exchangeCnt-10)) || !stop || this.needJing {
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
