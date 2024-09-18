# GoSimhash 相似度计算



## Usage

```
go get github.com/cnmade/gosimhash
```

### Usage of Package

```golang
import (
	"github.com/cnmade/gosimhash"
)

func getSimhash() {
    hasher := gosimhash.NewJiebaSimhasher(nil)
    defer hasher.Free()

    var sentence string = "今天的天气确实适合户外运动"
    var another string = "今年的气候确实很糟糕"
    var topN int = 5
    var limit int = 3

    // make simhash in uint64, like: 0xfa596a42bb35f945
    var first uint64 = hasher.MakeSimhash(&sentence, topN)
    var second uint64 = hasher.MakeSimhash(&another, topN)
    var dist1 int = gosimhash.CalculateDistanceBySimhash(first, second)
    var duplicated bool = gosimhash.IsSimhashDuplicated(first, second, limit)
    
    // make simhash in binary string, like: "10101110101111010101..."
    var firstStr string = hasher.MakeSimhashBinString(&sentence, topN)
    var secondStr string = hasher.MakeSimhashBinString(&another, topN)
    dist2, err := gosimhash.CalculateDistanceBySimhashBinString(firstStr, secondStr)
    if err != nil {
        fmt.Printf(err.Error())
    }
    duplicated, anotherErr := gosimhash.IsSimhashBinStringDuplicated(firstStr, secondStr, limit)
    if anotherErr != nil {
        fmt.Printf(anotherErr.Error())
    }
}
```

What's more, you can customize the hash algorithm(currently support siphash and jenkins) in simhash and dicts for jieba.

```golang
import (
	"github.com/cnmade/gosimhash"
	"github.com/cnmade/gosimhash/utils"
)
...
sip := utils.NewSipHasher([]byte(gosimhash.DEFAULT_HASH_KEY))
// jenkins := utils.NewJenkinsHasher()

hasher := gosimhash.NewSimhasher(sip, "./dict/jieba.dict.utf8", "./dict/hmm_model.utf8", "", "./dict/idf.utf8", "./dict/stop_words.utf8")
```

### Usage of Command

See example in [example/example.go](example/example.go)

```
cd example
go build

./example -help
```
