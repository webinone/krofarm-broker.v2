package api

import (
	"github.com/kataras/iris"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"time"
	"krofarm-broker.v2/models/protobuf"
	"krofarm-broker.v2/libs"
	_ "regexp"
	_ "strings"
)

type UdfAPI struct {
	*iris.Context
}

type Udf struct {
	UdfTy string `json:"udfTy" bson:"udfTy"`
	UdfNm string `json:"udfNm" bson:"udfNm"`
	UdfBody string `json:"udfBody" bson:"udfBody"`
}

func (this UdfAPI) SendUdf (ctx *iris.Context) {

	libs.Logger.Debug(">>>>>>>>>>>>>> UDF Handler !!!")

	body  := ctx.PostBody()

	// TODO : SSIBAL 내가 포인터 아규먼트 함수를 만들다니 ㅠㅠ
	libs.CleanJsonBody(&body)

	libs.Logger.Debug("############# BODY : ", string(body))

	reqId := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))

	var jsonData Udf
	json.Unmarshal([]byte(body), &jsonData)

	// MQTT 데이터 전송
	//--------------------------------------------------
	udfPayload := &protobuf.UDFPayload{
		FnctTyCd: proto.Int32(10160001),
		UdfTy: proto.String(jsonData.UdfTy),
		UdfNm: proto.String(jsonData.UdfNm),
		UdfBody:proto.String(jsonData.UdfBody),
		CreatedAt:proto.Int64(time.Now().UnixNano()),
	}

	client := libs.NewMqttLocalClient()

	pubTopic := "/kgw/v2/L/CMD/POST/UDF"

	message, err := proto.Marshal(udfPayload)
	if err != nil {
		libs.Logger.Debug(err)
	}

	token := client.Publish(pubTopic, 0, false, message)
	token.Wait()

	libs.Logger.Debug("############# UDF PUBLISH !!! ")

	client.Disconnect(0)

	libs.ResultAPI(ctx, reqId)

}


