package libs

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"flag"
)

type Configuration struct {
	MQTT 	mqttConfig
	INFO 	infoConfig
}

type mqttConfig struct {
	LocalUrl string 	`json:"localUrl"`
	KfarmUrl string 	`json:"kfarmUrl"`
}

type infoConfig struct {
	PrjctNo string 		`json:"prjctNo"`
	EndpntId string 	`json:"endpntId"`
}

var ConfigRoot string
var Config  = &Configuration{}

func LoadPathConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}

func LoadAutoConfig () {

	configRoot := flag.String("configRoot", "foo", "app.json path")
	flag.Parse()

	var path string
	if *configRoot == "foo" {
		path = "D:/Project/GO/src/krofarm-broker.v2/app.json"
		ConfigRoot = "D:/Project/GO/src/krofarm-broker.v2"
	} else {
		ConfigRoot = *configRoot
		path = *configRoot + "/app.json"
	}

	fmt.Println("configRoot : ", ConfigRoot)
	fmt.Println("path : ", path)

	LoadPathConfig(path)
}

