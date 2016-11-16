package api

import (
	"github.com/kataras/iris"
	"encoding/json"
	"krofarm-broker.v2/models/protobuf"
	"github.com/golang/protobuf/proto"
	"time"
	"krofarm-broker.v2/libs"
)

type ConfAPI struct {
	*iris.Context
}

type Conf struct {
	DvcId 			int64 `json:"dvcId"`
	FnctngModeCd 		*int32 `json:"fnctngModeCd"`
	TotalOpenExecTime 	*int32 `json:"totalOpenExecTime"`
	TotalCloseExecTime 	*int32 `json:"totalCloseExecTime"`
	ExecOffsetTime 		*int32 `json:"execOffsetTime"`
}

func (this ConfAPI) SendConf (ctx *iris.Context) {

	libs.Logger.Debug(">>>>>>>>>> Conf Handler !!!")

	body  := ctx.PostBody()
	libs.CleanJsonBody(&body)

	libs.Logger.Debug(string(body))

	var jsonData []Conf
	json.Unmarshal([]byte(body), &jsonData)

	libs.Logger.Debug(">>>>>>>>>> len(jsonData) : " ,len(jsonData))

	reqId := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))

	configPayload := &protobuf.ConfigPayload{}
	configPayload.ReqId = proto.Int64(reqId)
	configPayload.CreatedAt = proto.Int64(reqId)
	// 수동제어
	configPayload.ReqTy = proto.Int32(1)

	var confDevices []*protobuf.ConfDevice = make([]*protobuf.ConfDevice, len(jsonData))

	for i , conf := range jsonData {

		libs.Logger.Debug(">>>>>> i : ", i)
		libs.Logger.Debug(">>>>>> CONF DVCID : ", conf.DvcId)

		if conf.FnctngModeCd == nil {

			// 시간 변경 관련 conf
			libs.Logger.Debug(">>>>>> CONF TotalOpenExecTime : ", *conf.TotalOpenExecTime)
			libs.Logger.Debug(">>>>>> CONF TotalCloseExecTime : ", *conf.TotalCloseExecTime)
			libs.Logger.Debug(">>>>>> CONF ExecOffsetTime : ", *conf.ExecOffsetTime)

			confDevice := &protobuf.ConfDevice{
				DvcId:proto.Int64(conf.DvcId),
				TotalOpenExecTime:proto.Int32(*conf.TotalOpenExecTime),
				TotalCloseExecTime:proto.Int32(*conf.TotalCloseExecTime),
				ExecOffsetTime:proto.Int32(*conf.ExecOffsetTime),
			}
			//
			libs.Logger.Debug(confDevice)
			////
			confDevices[i] = confDevice
		}

		if conf.TotalOpenExecTime == nil {

			// 모드 변경 관련
			libs.Logger.Debug(">>>>>> CONF FnctngModeCd : ", *conf.FnctngModeCd)

			confDevice := &protobuf.ConfDevice{
				DvcId:proto.Int64(conf.DvcId),
				FnctngModeCd:proto.Int32(*conf.FnctngModeCd),
			}
			//
			libs.Logger.Debug(confDevice)
			////
			confDevices[i] = confDevice
		}

	}

	configPayload.ConfDevice = confDevices

	libs.Logger.Debug("#################### confPayload : ", configPayload)

	client := libs.NewMqttLocalClient()

	message, err := proto.Marshal(configPayload)
	if err != nil {
		libs.Logger.Critical(err)
	}

	token := client.Publish("/kgw/v2/C/CMD/POST/CONF", 0, false, message)
	token.Wait()

	libs.Logger.Debug("############# CONF PUBLISH !!! ")

	client.Disconnect(0)

	libs.ResultAPI(ctx, reqId)

}
