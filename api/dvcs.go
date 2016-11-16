package api

import (
	"github.com/kataras/iris"
	"fmt"
	"io/ioutil"
	"reflect"
	"unsafe"
	. "krofarm-broker.v2/libs"
)

type DvcsAPI struct {
	*iris.Context
}

func (this DvcsAPI) ListDvcs (ctx *iris.Context) {

	Logger.Debug(">>>>>>>>>> Dvcs List Handler !!!")

	b, err := ioutil.ReadFile(ConfigRoot + "/dvcList.json") // articles.json 파일의 내용을 읽어서 바이트 슬라이스에 저장
	if err != nil {
		fmt.Println(err)
		return
	}

	CleanJsonBody (&b)

	resultJson := BytesToString(b)

	Logger.Debug(">>>>>>>>>> resultJson : ", resultJson)

	ctx.SetHeader("Content-Type", "Application/json")
	ctx.Write(resultJson)
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
