package libs

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	UUID "github.com/satori/go.uuid"
	"fmt"
	"strings"
)

// 온실 내부 Local Client 만들기
func NewMqttLocalClient() MQTT.Client {

	mqttUrl := Config.MQTT.LocalUrl
	mqttUrls := strings.Split(mqttUrl, ",")

	fmt.Println(">>>>>>>>>>>>>>>>>>>> mqttUrls : ", mqttUrls)

	opts := MQTT.NewClientOptions()
	for _, uri := range mqttUrls {
		opts.AddBroker(uri)
	}
	opts.SetClientID(UUID.NewV4().String())

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//LOG.Logger.Info("############### client.IsConnected() : " , client.IsConnected())

	return client
}

// 온실 FMS  Client 만들기
func NewMqttKFarmClient() MQTT.Client {

	mqttUrl := Config.MQTT.KfarmUrl

	//fmt.Println("MQTTURL", mqttUrl)

	opts := MQTT.NewClientOptions().AddBroker(mqttUrl)
	opts.SetClientID(UUID.NewV4().String())
	//opts.SetDefaultPublishHandler(messageArrived)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">>>>>>>>>>>>>>>>> 시발 여기가 에러냐????")
		panic(token.Error())
	}

	//fmt.Println("############### client.IsConnected() : " , client.IsConnected())

	return client
}