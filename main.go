package main

import (
	"github.com/kataras/iris"
	. "krofarm-broker.v2/libs"
	//"krofarm-broker.v2/api"
	"github.com/iris-contrib/middleware/recovery"
	"krofarm-broker.v2/mqtt"
	_ "krofarm-broker.v2/routers"
)

func init() {
	// Config Loading (app.json)
	LoadAutoConfig()
	InitLogConfig() // logging
}

func main() {

	// ProjectNo로 전환한다.
	Logger.Debug(">>>>>>>>>> KroFarm Broker v.2 Start !!!")

	// MQTT Broker goroutine call
	go mqtt.StartMqttReceiver()

	// 이건 kist때문에 추가...!!!!!!!! KIST !!!!

	// Recovery Middleware
	//----------------------------------------------------------------
	iris.Use(recovery.Handler)
	//----------------------------------------------------------------

	iris.Get("/hi", func(ctx *iris.Context) {
		ctx.Write("Hello! %s", "KroFarm Broker v.2")
	})

	iris.Listen(":8989")
}
