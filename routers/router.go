package routers

import (
	"github.com/kataras/iris"
	"krofarm-broker.v2/api"
	"fmt"
)

func init () {

	fmt.Println("routers router1 setting !!")

	iris.Get("/fuck1", func(ctx *iris.Context) {
		ctx.Write("fuck1 ! %s", "KroFarm Broker v.2")
	})

	platform := iris.Party("/csb/iotservice/v2")
	{
		// 디바이스 목록 가져오기
		platform.Get("/iotdvc/endpnts/:endpntId/dvcs/list", api.DvcsAPI{}.ListDvcs)
		// Udf 전송
		platform.Post("/iottransfer/infocmmnd/endpnts/:endpntId/post/udf", api.UdfAPI{}.SendUdf)
		// Exec 전송 -> 원격제어 전송
		platform.Post("/iottransfer/cntrlcmmnd/endpnts/:endpntId/dvcs/list/exec/actuators", api.ExecAPI{}.SendExec)
		// Conf 전송
		platform.Post("/iottransfer/infocmmnd/endpnts/:endpntId/dvcs/list/post/conf", api.ConfAPI{}.SendConf)
		// Reset 전송
		platform.Post("/iottransfer/cntrlcmmnd/endpnts/:endpntId/dvcs/list/reset/actuators", api.ResetAPI{}.SendReset)

	}
}
