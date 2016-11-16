package routers

import (
	"github.com/kataras/iris"
	"fmt"
)

func init () {

	fmt.Println("routers router2 setting !!")

	iris.Get("/fuck2", func(ctx *iris.Context) {
		ctx.Write("fuck2 ! %s", "KroFarm Broker v.2")
	})

}
