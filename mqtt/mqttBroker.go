package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	//UUID "github.com/satori/go.uuid"
	_ "os"
	"krofarm-broker/base/protobuf"
	"encoding/json"
	. "krofarm-broker.v2/libs"
	"log"
	"strconv"
	_ "strings"
)

var messageArrived MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {

	Logger.Debug(">>>>>>>>>> TOPIC: ", msg.Topic())

	//log.Println("MSG: %s\n", msg.Payload())
	prjctNo 	:= Config.INFO.PrjctNo
	endpntId 	:= Config.INFO.EndpntId

	var pubTopic string

	 //SENSORS나 ACTUATOR 데이터 인 경우...
	if msg.Topic() == "/ksn/v2/S/EVT/DATA/SENSORS" ||
	   msg.Topic() == "/ksn/v2/C/EVT/DATA/ACTUATORS" ||
	   msg.Topic() == "/ksn/v2/S/REPLY/GET/SENSORS" ||
	   msg.Topic() == "/ksn/v2/C/REPLY/GET/ACTUATORS" {

		dataPayload := &protobuf.DataPayload{}

		err := proto.Unmarshal(msg.Payload(), dataPayload)
		if err != nil {
			Logger.Critical("dataPayload unmarshaling error: ", err)
		}

		var dvcId int64
		var attrbVal string
		var attrbCd int32
		var createdAt int64
		var attrbStatCd int32
		var attrbStatMssage string
		var reqId int64
		var reqTy int32
		var curExecStep int32

		createdAt = dataPayload.GetCreatedAt()
		reqId = dataPayload.GetReqId()
		reqTy = dataPayload.GetReqTy()

		result := make([]map[string]interface{}, len(dataPayload.GetAttribute()), 100)

		for i , attribute := range dataPayload.GetAttribute() { // i에는 인덱스, value에는 배열 요소의 값이 들어감

			dvcId   	= attribute.GetDvcId()
			attrbCd 	= attribute.GetAttrbCd()
			attrbStatCd 	= attribute.GetAttrbStatCd()
			attrbStatMssage = attribute.GetAttrbStatMssage()
			curExecStep	= attribute.GetCurExecStep()

			switch attribute.GetAttrbVal().GetType() {
				case protobuf.ValueType_DOUBLE :
					attrbVal = strconv.FormatFloat(attribute.GetAttrbVal().GetDoubleValue(), 'f', -1, 32)
				case protobuf.ValueType_INT32 :
					attrbVal = strconv.FormatInt(int64(attribute.GetAttrbVal().GetIntValue()), 10)
				case protobuf.ValueType_INT64 :
					attrbVal = strconv.FormatInt(attribute.GetAttrbVal().GetLongValue(), 10)
			}

			Logger.Debug("dvcId : " , dvcId)
			Logger.Debug("attrbVal : " , attrbVal)
			Logger.Debug("attribute.GetAttrbVal().GetType() : " , attribute.GetAttrbVal().GetType())
			Logger.Debug("attrbCd : " , attrbCd)
			Logger.Debug("createdAt : " , createdAt)
			Logger.Debug("attrbStatCd : " , attrbStatCd)
			Logger.Debug("attrbStatMssage : " , attrbStatMssage)

			//Logger.Debug("reqId : " , reqId)
			//Logger.Debug("reqTy : " , reqTy)
			//Logger.Debug("curExecStep : " , curExecStep)

			jsonData := make(map[string]interface{})

			jsonData["prjctNo"] = prjctNo
			jsonData["endpntId"] = endpntId
			jsonData["dvcId"] = dvcId
			jsonData["attrbVal"] = attrbVal
			jsonData["attrbCd"] = attrbCd
			jsonData["createdAt"] = createdAt
			jsonData["attrbStatCd"] = attrbStatCd
			jsonData["attrbStatMssage"] = attrbStatMssage
			jsonData["reqId"] = reqId
			jsonData["reqTy"] = reqTy
			jsonData["curExecStep"] = curExecStep

			result[i] = jsonData

		}

		jsonString, _ := json.Marshal(result)

		Logger.Info(string(jsonString))
		//Logger.Debug(string(jsonString))

		client := NewMqttKFarmClient()

		if msg.Topic() == "/ksn/v2/S/EVT/DATA/SENSORS" || msg.Topic() == "/ksn/v2/S/REPLY/GET/SENSORS" {
			pubTopic = "/kcsb/v2/" + prjctNo + "/" + endpntId + "/EVT/DATA/SENSORS/json"
		} else if msg.Topic() == "/ksn/v2/C/REPLY/GET/ACTUATORS" ||
		msg.Topic() == "/ksn/v2/C/EVT/DATA/ACTUATORS"{
			pubTopic = "/kcsb/v2/" + prjctNo + "/" + endpntId + "/EVT/DATA/ACTUATORS/json"
		}

		//Logger.Debug("pubTopic : " + pubTopic)

		token := client.Publish(pubTopic, 0, false, jsonString)
		token.Wait()

		client.Disconnect(0)
	}

	// STATUS
	if msg.Topic() == "/ksn/v2/S/EVT/DATA/STATUS" ||
	   msg.Topic() == "/ksn/v2/C/EVT/DATA/STATUS" ||
	   msg.Topic() == "/ksn/v2/L/EVT/DATA/STATUS" ||
	   msg.Topic() == "/ksn/v2/S/REPLY/GET/STATUS" ||
	   msg.Topic() == "/ksn/v2/C/REPLY/GET/STATUS" ||
	   msg.Topic() == "/ksn/v2/L/REPLY/GET/STATUS" {

		statusPayload := &protobuf.StatusPayload{}

		err := proto.Unmarshal(msg.Payload(), statusPayload)
		if err != nil {
			Logger.Critical("statusPayload unmarshaling error: ", err)
		}

		pubTopic = "/kcsb/v2/" + prjctNo + "/" + endpntId + "/EVT/DATA/STATUS/json"

		var createdAt int64
		var reqId int64
		var reqTy int32
		var subDvcStatChg int32
		var dvcId int64
		var commStatCd int32
		var commStatMssage string
		var fnctngStatCd int32
		var fnctngStatMssage string
		var cntrlStat int32
		var dlgatStat int32

		createdAt = statusPayload.GetCreatedAt()
		reqId = statusPayload.GetReqId()
		reqTy = statusPayload.GetReqTy()
		subDvcStatChg = statusPayload.GetSubDvcStatChg()

		//Logger.Debug("createdAt", createdAt)
		Logger.Debug("reqId", reqId)
		Logger.Debug("reqTy", reqTy)
		Logger.Debug("subDvcStatChg", subDvcStatChg)

		result := make([]map[string]interface{}, len(statusPayload.GetStatDevice()), 100)

		//Logger.Debug("TOPIC: ", msg.Topic())

		//Logger.Debug(">>>>>>>>>>>>>>>>>>>>>> STATUS !!!! dvcId : ", dvcId)

		for i , attribute := range statusPayload.GetStatDevice() {
			// i에는 인덱스, value에는 배열 요소의 값이 들어감

			dvcId = attribute.GetDvcId()
			commStatCd = attribute.GetCommStatCd()
			commStatMssage = attribute.GetCommStatMssage()
			fnctngStatCd = attribute.GetFnctngStatCd()
			fnctngStatMssage = attribute.GetFnctngStatMssage()
			cntrlStat = attribute.GetCntrlStat()
			dlgatStat = attribute.GetDlgatStat()

			ioPoints := make([]map[string]interface{}, len(attribute.GetIoPoint()), 100)

			for x , ioPoint := range attribute.GetIoPoint() {

				ioPointJsonData := make(map[string]interface{})

				ioPointJsonData["ioPointId"] = ioPoint.IoPointId
				ioPointJsonData["fnctngStatCd"] = ioPoint.FnctngStatCd

				ioPoints[x] = ioPointJsonData
			}

			jsonData := make(map[string]interface{})

			jsonData["prjctNo"] = prjctNo
			jsonData["endpntId"] = endpntId
			jsonData["dvcId"] = dvcId
			jsonData["commStatCd"] = commStatCd
			jsonData["commStatMssage"] = commStatMssage
			jsonData["commStatCreatedAt"] = createdAt
			jsonData["fnctngStatCd"] = fnctngStatCd
			jsonData["fnctngStatMssage"] = fnctngStatMssage
			jsonData["fnctngStatCreatedAt"] = createdAt
			jsonData["cntrlStat"] = cntrlStat
			jsonData["dlgatStat"] = dlgatStat
			jsonData["cntrlIoPointStats"] = ioPoints

			result[i] = jsonData
		}

		//Logger.Debug("########################## STATUS PAYLOAD START ")

		jsonString, _ := json.Marshal(result)

		log.Println(string(jsonString))

		Logger.Info(string(jsonString))

		//Logger.Debug("########################## STATUS PAYLOAD END ")

		client := NewMqttKFarmClient()

		//Logger.Debug("pubTopic : " + pubTopic)

		token := client.Publish(pubTopic, 0, false, jsonString)
		token.Wait()

		client.Disconnect(0)
	}

	// UDF
	if msg.Topic() == "/ksn/v2/L/CMD/POST/UDF" {

		udfPayload := &protobuf.UDFPayload{}

		err := proto.Unmarshal(msg.Payload(), udfPayload)
		if err != nil {
			Logger.Critical("udfPayload unmarshaling error: ", err)
		}

		fnctTyCd := udfPayload.GetFnctTyCd()
		udfNm := udfPayload.GetUdfNm()
		createdAt := udfPayload.GetCreatedAt()

		Logger.Debug("fnctTyCd", fnctTyCd)
		Logger.Debug("udfNm", udfNm)
		Logger.Debug("createdAt", createdAt)
		//Logger.Debug("udfBody", createdAt)
	}

	// UDM
	if msg.Topic() == "/ksn/v2/L/EVT/MSG/UDM" {

		//Logger.Debug(">>>>>>>>>>>>>>> UDM !!!")

		udmPayload := &protobuf.UDMPayload{}

		err := proto.Unmarshal(msg.Payload(), udmPayload)
		if err != nil {
			Logger.Critical("udmPayload unmarshaling error: ", err)
		}

		var dvcId 	int64
		var createdAt 	int64
		var mssageTyCd 	int32
		var udmNm 	string
		var udmBody 	string
		var udmTy 	string

		dvcId 		= udmPayload.GetDvcId()
		createdAt 	= udmPayload.GetCreatedAt()
		mssageTyCd 	= udmPayload.GetMssageTyCd()
		udmNm		= udmPayload.GetUdmNm()
		udmBody		= udmPayload.GetUdmBody()
		udmTy		= udmPayload.GetUdmTy()

		jsonData := make(map[string]interface{})

		jsonData["prjctNo"] = prjctNo
		jsonData["endpntId"] = endpntId
		jsonData["dvcId"] = dvcId
		jsonData["createdAt"] = createdAt
		jsonData["mssageTyCd"] = mssageTyCd
		jsonData["udmNm"] = udmNm
		jsonData["udmBody"] = udmBody
		jsonData["udmTy"] = udmTy

		jsonString, _ := json.Marshal(jsonData)

		Logger.Debug(string(jsonString))

		client := NewMqttKFarmClient()

		if udmNm == "heatBoiler" {
			//Logger.Debug("########################### udmBody : ", udmBody)
			//Logger.Debug("########################### heatboiler udm : ", string(jsonString))
		}

		////Logger.Debug(">>>>>>>>>>>>>>> UDM 222")

		pubTopic = "/kcsb/v2/" + prjctNo + "/" + endpntId + "/EVT/MSG/UDM/json"

		////Logger.Debug("pubTopic : " + pubTopic)

		token := client.Publish(pubTopic, 0, false, jsonString)
		token.Wait()

		client.Disconnect(0)
	}

	if msg.Topic() == "/ksn/v2/L/REPLY/POST/CONF" ||
	   msg.Topic() == "/ksn/v2/C/REPLY/POST/CONF" ||
	   msg.Topic() == "/ksn/v2/S/REPLY/POST/CONF" ||
	   msg.Topic() == "/ksn/v2/G/REPLY/POST/CONF" ||
	   msg.Topic() == "/kgw/v2/L/REPLY/POST/CONF" ||
	   msg.Topic() == "/kgw/v2/C/REPLY/POST/CONF" ||
	   msg.Topic() == "/kgw/v2/S/REPLY/POST/CONF" {

		//Logger.Debug(">>>>>>>>>>>>>>> CONF REPLY !!!")
		Logger.Debug(">>>>>>>>>>>>>>> CONF REPLY TOPIC: ", msg.Topic())

		confPayload := &protobuf.ConfigPayload{}

		err := proto.Unmarshal(msg.Payload(), confPayload)
		if err != nil {
			Logger.Critical("confPayload unmarshaling error: ", err)
		}

		var createdAt int64
		var reqId int64
		var reqTy int32
		//var subDvcStatChg int32
		var dvcId int64
		var fnctngModeCd int32
		var totalOpenExecTime int32
		var totalCloseExecTime int32
		var execOffsetTime int32

		reqId = confPayload.GetReqId()
		reqTy = confPayload.GetReqTy()
		createdAt = confPayload.GetCreatedAt()

		//Logger.Debug("createdAt", createdAt)
		//Logger.Debug("reqId", reqId)
		//Logger.Debug("reqTy", reqTy)

		result := make([]map[string]interface{}, len(confPayload.GetConfDevice()), 100)

		for i , confDevice := range confPayload.GetConfDevice() {

			dvcId = confDevice.GetDvcId()
			fnctngModeCd = confDevice.GetFnctngModeCd()
			totalOpenExecTime = confDevice.GetTotalCloseExecTime()
			totalCloseExecTime = confDevice.GetTotalCloseExecTime()
			execOffsetTime = confDevice.GetExecOffsetTime()

			jsonData := make(map[string]interface{})

			jsonData["reqId"] = reqId
			jsonData["prjctNo"] = prjctNo
			jsonData["endpntId"] = endpntId
			jsonData["dvcId"] = dvcId
			jsonData["resCreatedAt"] = createdAt
			jsonData["reqTy"] = reqTy
			jsonData["fnctngModeCd"] = fnctngModeCd
			jsonData["totalOpenExecTime"] = totalOpenExecTime
			jsonData["totalCloseExecTime"] = totalCloseExecTime
			jsonData["execOffsetTime"] = execOffsetTime

			result[i] = jsonData
		}

		Logger.Debug("########################## CONF REPLY PAYLOAD START ")

		jsonString, _ := json.Marshal(result)

		Logger.Debug(string(jsonString))

		//log.Println(string(jsonString))

		//Logger.Info(string(jsonString))

		Logger.Debug("########################## CONF REPLY PAYLOAD END ")

		client := NewMqttKFarmClient()

		pubTopic = "/kcsb/v2/" + prjctNo + "/" + endpntId + "/REPLY/POST/CONF/json"

		//Logger.Debug("pubTopic : " + pubTopic)

		token := client.Publish(pubTopic, 0, false, jsonString)
		token.Wait()

		client.Disconnect(0)
	}
}

func init() {

}

func StartMqttReceiver() {

	//Logger.Debug("MQTT RECEIVER START !!!!")
	Logger.Info("MQTT RECEIVER START !!!!")

	client := NewMqttLocalClient()

	if token := client.Subscribe("/#", 0, messageArrived); token.Wait() && token.Error() != nil {
		//Logger.Debug(token.Error())
		StartMqttReceiver()
	}

	//Logger.Debug("MQTT RECEIVER END !!!!")
}