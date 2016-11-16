package api

import (
	"github.com/kataras/iris"
	"encoding/json"
	"time"
	"krofarm-broker.v2/models/protobuf"
	"github.com/golang/protobuf/proto"
	"krofarm-broker.v2/libs"
)

type ResetAPI struct {
	*iris.Context
}

type Reset struct {
	DvcId		int64 	`json:"dvcId"`
}

func (this ResetAPI) SendReset (ctx *iris.Context) {

	libs.Logger.Debug(">>>>>>>>>>>>>>> Reset Handler !!!")

	body  := ctx.PostBody()
	libs.CleanJsonBody(&body)

	libs.Logger.Debug(string(body))

	var jsonData []Reset
	json.Unmarshal([]byte(body), &jsonData)

	reqId := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))

	dataPayload := &protobuf.DataPayload{}
	dataPayload.ReqId = proto.Int64(reqId)
	dataPayload.CreatedAt = proto.Int64(reqId)
	// 수동제어
	dataPayload.ReqTy = proto.Int32(1)

	var dvcId int64

	var attributes []*protobuf.Attribute = make([]*protobuf.Attribute, len(jsonData))

	for i , resetDvcId := range jsonData {

		dvcId = resetDvcId.DvcId

		libs.Logger.Debug("dvcId : ", dvcId)

		attribute := &protobuf.Attribute{
			DvcId:proto.Int64(dvcId),
		}

		attributes[i] = attribute
	}

	libs.Logger.Debug("attributes : ", attributes)
	dataPayload.Attribute = attributes

	libs.Logger.Debug("dataPayload : ", dataPayload)

	client := libs.NewMqttLocalClient()
	pubTopic := "/kgw/v2/C/CMD/EXEC/ACTUATORS"

	message, err := proto.Marshal(dataPayload)
	if err != nil {
		//log.Fatal("marshaling error: ", err)
		libs.Logger.Debug(err)
	}

	token := client.Publish(pubTopic, 0, false, message)
	token.Wait()

	libs.Logger.Debug("############# RESET PUBLISH !!! ")
	client.Disconnect(0)

	libs.ResultAPI(ctx, reqId)
}