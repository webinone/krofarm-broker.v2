package libs

import (
	"github.com/kataras/iris"
)

type APIResult struct {
	Header headerResult	`json:"header" bson:"header"`
	Body   bodyResult	`json:"body" bson:"body"`
}

type headerResult struct {
	ResultCode 	int32 	`json:"resultCode" bson:"resultCode"`
	ResultMessage 	string 	`json:"resultMessage" bson:"resultMessage"`
}

type bodyResult struct {
	Item 	itemResult	`json:"item" bson:"item"`
}

type itemResult struct {
	ReqId		int64	`json:"reqId" bson:"reqId"`
	SndngCd		int32	`json:"sndngCd" bson:"sndngCd"`
}

func ResultAPI(ctx *iris.Context, reqId int64) error {

	apiResult := APIResult{
		Header : headerResult{
			ResultCode: 201,
			ResultMessage: "Created",
		},
		Body: bodyResult{
			Item: itemResult{
				ReqId: reqId,
				SndngCd: 10150002,
			},
		},
	}
	// CORS 설정
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
	return ctx.JSON(iris.StatusOK, &apiResult)

}
