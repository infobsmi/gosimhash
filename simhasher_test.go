package gosimhash

import (
	"testing"

	"github.com/infobsmi/gosimhash/utils"
)

func TestNewJiebaSimhasher(t *testing.T) {
	hasher := NewJiebaSimhasher(nil)
	defer hasher.Free()

	var sentence = "【王者荣耀：限制未成年用户游戏时间 单机模式玩法暂时关闭】今日，王者荣耀在其官网发布公告称，在游戏时长限制方面"
	//var another string = "今年的气候确实很糟糕"
	var another = "【王者荣耀：限制未成年用户游戏时间 今日，王者荣耀在其官网发布公告称，限400元人民币。"
	var topN = 20

	// make simhash in uint64, like: 0xfa596a42bb35f945
	var first = hasher.MakeSimhash(&sentence, topN)
	var second = hasher.MakeSimhash(&another, topN)
	dist1 := CalculateDistanceBySimhash(first, second)

	t.Log(first, second, dist1)

	another = "【王者荣耀：限制未成年用户游戏时间 单机模式玩法暂时关闭】今日，王者荣耀在其官网发布公告称"
	f0 := hasher.MakeSimhash(&sentence, topN)
	s0 := hasher.MakeSimhash(&another, topN)
	dist1 = CalculateDistanceBySimhash(f0, s0)
	t.Log(first, second, dist1)

	another = "王者荣耀在其官网发布公告称，限400元人民币，限制未成年用户游戏时间"
	f0 = hasher.MakeSimhash(&sentence, topN)
	s0 = hasher.MakeSimhash(&another, topN)
	dist1 = CalculateDistanceBySimhash(f0, s0)
	t.Log(first, second, dist1)

	another = "我们之间是完全不可能有任何结果的"
	f0 = hasher.MakeSimhash(&sentence, topN)
	s0 = hasher.MakeSimhash(&another, topN)
	dist1 = CalculateDistanceBySimhash(f0, s0)
	t.Log(first, second, dist1)
}
func TestSimhashWithJenkins(t *testing.T) {
	hasher := NewSimpleSimhasher()
	defer hasher.Free()

	var sentence = "我来到北京清华大学"
	var topN = 5

	func() {
		var expected uint64
		var actual uint64

		expected = 0xfa596a42bb35f945
		actual = hasher.MakeSimhash(&sentence, topN)
		if expected != actual {
			t.Error(expected, "!=", actual)
		}
	}()

	func() {
		var expected string
		var actual string

		expected = "1111101001011001011010100100001010111011001101011111100101000101"
		actual = hasher.MakeSimhashBinString(&sentence, topN)
		if expected != actual {
			t.Error(expected, "!=", actual)
		}
	}()

	func() {
		var first uint64
		var second uint64

		first = hasher.MakeSimhash(&sentence, topN)
		second = hasher.MakeSimhash(&sentence, topN)
		if first != second {
			t.Error(first, "!=", second)
		}
	}()

	func() {
		distance := CalculateDistanceBySimhash(0x812e5cf1b47eb66, 0x812e5cf1b47eb61)
		t.Logf("distance: %+v", distance)
		if distance != 3 {
			t.Error(distance, "!= 3")
		}
	}()

	func() {
		distance, err := CalculateDistanceBySimhashBinString(
			"100000010010111001011100111100011011010001111110101101100110",
			"100000010010111001011100111100011011010001111110101101100001")
		if err != nil {
			t.Error(err.Error())
			return
		}
		if distance != 3 {
			t.Error(distance, "!= 3")
		}
	}()
}

func TestSimhashWithSipHash(t *testing.T) {
	sip := utils.NewSipHasher([]byte(DefaultHashKey))
	hasher := NewSimhasher(sip, "./dict/jieba.dict.utf8", "./dict/hmm_model.utf8", "",
		"./dict/idf.utf8", "./dict/stop_words.utf8")
	defer hasher.Free()

	var sentence = "我来到北京清华大学"
	var topN = 5

	func() {
		var expected uint64
		var actual uint64

		expected = 0x812e5cf1b47eb66
		actual = hasher.MakeSimhash(&sentence, topN)
		if expected != actual {
			t.Error(expected, "!=", actual)
		}
	}()

	func() {
		var expected string
		var actual string

		expected = "100000010010111001011100111100011011010001111110101101100110"
		actual = hasher.MakeSimhashBinString(&sentence, topN)
		if expected != actual {
			t.Error(expected, "!=", actual)
		}
	}()

	func() {
		var first uint64
		var second uint64

		first = hasher.MakeSimhash(&sentence, topN)
		second = hasher.MakeSimhash(&sentence, topN)
		if first != second {
			t.Error(first, "!=", second)
		}
	}()

	func() {
		distance := CalculateDistanceBySimhash(0x812e5cf1b47eb66, 0x812e5cf1b47eb61)
		if distance != 3 {
			t.Error(distance, "!= 3")
		}
	}()

	func() {
		distance, err := CalculateDistanceBySimhashBinString(
			"100000010010111001011100111100011011010001111110101101100110",
			"100000010010111001011100111100011011010001111110101101100001")
		if err != nil {
			t.Error(err.Error())
			return
		} else if distance != 3 {
			t.Error(distance, "!= 3")
		}

		distance, err = CalculateDistanceBySimhashBinString(
			"100000010010111001011100111100011011010001111110101101100113",
			"100000010010111001011100111100011011010001111110101101100001")
		if err == nil {
			t.Error("Should throw error at CalculateDistanceBySimhashBinString")
		}

		distance, err = CalculateDistanceBySimhashBinString(
			"100000010010111001011100111100011011010001111110101101100110",
			"100000010010111001011100111100011011010001111110101101100003")
		if err == nil {
			t.Error("Should throw error at CalculateDistanceBySimhashBinString")
		}
	}()

	func() {
		duplicated := IsSimhashDuplicated(0x812e5cf1b47eb66, 0x812e5cf1b47eb61, 3)
		if !duplicated {
			t.Error("Should be duplicated at IsSimhashDuplicated")
		}
		duplicated, err := IsSimhashBinStringDuplicated(
			"100000010010111001011100111100011011010001111110101101100110",
			"100000010010111001011100111100011011010001111110101101100001", 3)
		if err != nil {
			t.Error(err.Error())
		} else if !duplicated {
			t.Error("Should be duplicated at IsSimhashBinStringDuplicated")
		}
		duplicated, err = IsSimhashBinStringDuplicated(
			"100000010010111001011100111100011011010001111110101101100113",
			"100000010010111001011100111100011011010001111110101101100001", 3)
		if err == nil {
			t.Error("Should throw error at IsSimhashBinStringDuplicated")
		}
		duplicated, err = IsSimhashBinStringDuplicated(
			"100000010010111001011100111100011011010001111110101101100110",
			"100000010010111001011100111100011011010001111110101101100003", 3)
		if err == nil {
			t.Error("Should throw error at IsSimhashBinStringDuplicated")
		}
	}()
}
