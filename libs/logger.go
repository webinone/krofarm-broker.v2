package libs

import (
	_ "errors"
	"fmt"
	seelog "github.com/cihub/seelog"
	_ "io"
)

var Logger seelog.LoggerInterface

func LoadLogConfig() {

	fmt.Println(">>>>>>>>>>>>>> ConfigRoot : ", ConfigRoot)

	appConfig := `
	<seelog>
	    <outputs>
	    	<console formatid="common" />
	       <rollingfile formatid="common" type="date" filename="`+ ConfigRoot +`/logs/broker.log" datepattern="02.01.2006" maxrolls="3" />
	    </outputs>
	    <formats>
		<format id="common" format="%Date %Time [%LEV] %RelFile %Func %Msg%n" />
	    </formats>
	</seelog>
	`
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	UseLogger(logger)
}

func InitLogConfig() {
	DisableLog()
	LoadLogConfig()
}

// DisableLog disables all library log output
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}