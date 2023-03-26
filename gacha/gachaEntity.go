package gacha

import (
	"fmt"

	wrc "github.com/mroth/weightedrand/v2"
)

type GachaPro struct {
	_name string
	_pro  int32 //万分率
}

func (this *GachaPro) Init(name string, pro int32) {
	this._name = name
	this._pro = pro
}

func newGachaPro(name string, pro int32) *GachaPro {
	gacha := new(GachaPro)
	gacha.Init(name, pro)
	return gacha
}

type GachaEntiy struct {
	choicer *wrc.Chooser[string, int32]
}

func (this *GachaEntiy) init(items []*GachaPro) {
	choices := []wrc.Choice[string, int32]{}
	for _, item := range items {
		choice := wrc.NewChoice(item._name, item._pro)
		choices = append(choices, choice)
	}
	wrc, _ := wrc.NewChooser(choices...)
	this.choicer = wrc
}

func newGachaEntity(items []*GachaPro) *GachaEntiy {
	gacha := new(GachaEntiy)
	gacha.init(items)
	return gacha
}

type GachaResult struct {
	_name       string
	_charPieces float64
	_pigPieces  float64
	_cnt        float64
	_jing       float64
}

func (this *GachaResult) String() string {
	str := fmt.Sprintf("name:%s, charPieces:%f, pigPieces:%f, cnt:%f, jing:%f", this._name, this._charPieces, this._pigPieces, this._cnt, this._jing)
	return str
}
