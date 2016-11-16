package api

import (
	"github.com/kataras/iris"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"time"
	"krofarm-broker.v2/models/protobuf"
	"strconv"
	"krofarm-broker.v2/libs"
)

type ExecAPI struct {
	*iris.Context
}

type Exec struct {
	DvcId		int64 	`json:"dvcId"`
	AttrbCd		int32 	`json:"attrbCd"`
	AttrbVal	string 	`json:"attrbVal"`
	AttrbExecs 	[]ExecAttrbExecs `json:"attrbExecs"`
}

type ExecAttrbExecs struct {
	StepSeq 	int32 	`json:"stepSeq"`
	StepDelay 	int32 	`json:"stepDelay"`
	StepFactor 	int32 	`json:"stepFactor"`
}

func (this ExecAPI) SendExec (ctx *iris.Context) {

	libs.Logger.Debug("#################### ExecHandler !!!!!!!")

	body  := ctx.PostBody()
	libs.CleanJsonBody(&body)

	libs.Logger.Debug(string(body))

	var jsonData []Exec

	json.Unmarshal([]byte(body), &jsonData)

	libs.Logger.Debug("#################### JSON DATA LENGTH : ", len(jsonData))

	reqId := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))

	dataPayload := &protobuf.DataPayload{}

	dataPayload.ReqId = proto.Int64(reqId)
	dataPayload.CreatedAt = proto.Int64(reqId)
	// 수동제어
	dataPayload.ReqTy = proto.Int32(1)

	var attributes []*protobuf.Attribute = make([]*protobuf.Attribute, len(jsonData))
	int32ValueType  := protobuf.ValueType_INT32

	var dvcId int64
	var attrbCd int32
	var attrbVal int32

	for i , exec := range jsonData {

		libs.Logger.Debug("i : ", i)

		dvcId = exec.DvcId
		attrbCd = exec.AttrbCd
		tmpAttrbVal, _ := strconv.ParseInt(exec.AttrbVal, 10, 32)
		attrbVal = int32(tmpAttrbVal)

		libs.Logger.Debug("DVC ID : ",   dvcId)
		libs.Logger.Debug("AttrbCd : ",  attrbCd)
		libs.Logger.Debug("AttrbVal : ", attrbVal)
		libs.Logger.Debug("len(execRoot.AttrbExecs) : ", len(exec.AttrbExecs))

		var execAttributes []*protobuf.ExecAttribute = make([]*protobuf.ExecAttribute, len(exec.AttrbExecs))

		for x, attrbExec := range exec.AttrbExecs {

			libs.Logger.Debug(" >>>>> x : " , x)

			libs.Logger.Debug("StepSeq : ",   attrbExec.StepSeq)
			libs.Logger.Debug("StepDelay : ",  attrbExec.StepDelay)
			libs.Logger.Debug("StepFactor : ", attrbExec.StepFactor)

			execAttribute := &protobuf.ExecAttribute{
				StepSeq : proto.Int32(attrbExec.StepSeq),
				StepDelay : proto.Int32(attrbExec.StepDelay),
				StepFactor :  &protobuf.AttributeValue {
					Type: &int32ValueType,
					IntValue:proto.Int32(attrbExec.StepFactor),
				},
			}
			execAttributes[x] = execAttribute
		}

		attribute := &protobuf.Attribute{
			DvcId:proto.Int64(dvcId),
			AttrbCd : proto.Int32(attrbCd),
			AttrbVal:&protobuf.AttributeValue{
				Type: &int32ValueType,
				IntValue:proto.Int32(attrbVal),
			},
			ExecAttribute: execAttributes,
		}
		//
		libs.Logger.Debug(attribute)

		attributes[i] = attribute
	}

	libs.Logger.Debug("attributes : ", attributes)

	dataPayload.Attribute = attributes

	libs.Logger.Debug("dataPayload : ", dataPayload)

	// TODO : 실행될때 열자...
	client := libs.NewMqttLocalClient()
	pubTopic := "/kgw/v2/C/CMD/EXEC/ACTUATORS"

	message, err := proto.Marshal(dataPayload)
	if err != nil {
		//log.Fatal("marshaling error: ", err)
		libs.Logger.Debug(err)
	}

	token := client.Publish(pubTopic, 0, false, message)
	token.Wait()

	libs.Logger.Debug("############# EXEC PUBLISH !!! ")
	client.Disconnect(0)

	libs.ResultAPI(ctx, reqId)
}