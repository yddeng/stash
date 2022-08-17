package config

import (
	"testing"

	"github.com/BurntSushi/toml"

	"github.com/smartystreets/goconvey/convey"

	"log"
)

type StringWriter struct {
	data []byte
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *StringWriter) ToStr() string {
	return string(w.data)
}

func TestDecodeConfig(t *testing.T) {
	var configData = `
			[Common]
			centerAddr     	= ["localhost:8010"]
			dbConfig      	= ["redis@dir@localhost:6379","flyfish@login@localhost:10012","flyfish@game@localhost:10012","flyfish@room@localhost:10012","flyfish@battle@localhost:10012"]
			excelPath       = "/Users/dev/svn-initialthree/RestoryGame03/GameDesign/Excel"
			missionXmlPath  = "/Users/dev/svn-initialthree/RestoryGame03/GameDesign/MissionXml"
			battleLevelXmlPath = "/Users/dev/svn-initialthree/RestoryGame03/GameDesign/BattleLevelXml"
			areaID         	= 5

			[[Dir]]
			filePath       	= "node_dir/config/center_config.json"
			logicAddr		= "1.1.1"
			clusterAddr    	= "localhost:8011"
			externalAddr   	= "10.128.2.244:9012"

			[[Login]]
			logicAddr		= "1.2.1"
			clusterAddr    	= "localhost:8012"
			externalAddr   	= "10.128.2.244:9010"

			[[Game]]
			logicAddr		= "1.3.1"
			clusterAddr    	= "localhost:8013"

			[[Map]]
			logicAddr		= "1.4.1"
			clusterAddr    	= "localhost:8014"

			[[Chat]]
			logicAddr		= "1.5.1"
			clusterAddr    	= "localhost:8015"

			[[Team]]
			logicAddr		= "1.6.1"
			clusterAddr    	= "localhost:8016"

			[[Battle]]
			logicAddr		= "1.7.1"
			clusterAddr    	= "localhost:8017"
			externalAddr   	= "10.128.2.244:9013"

			[[Gate]]
			logicAddr		= "1.8.1"
			clusterAddr    	= "localhost:8018"
			externalAddr   	= "10.128.2.244:9014"
			`
	convey.Convey("Test Decode Config ...", t, func() {
		conf := &Config{}
		_, err := toml.Decode(configData, conf)
		convey.So(err, convey.ShouldBeNil)
		convey.So(conf.Common.CenterAddr[0], convey.ShouldEqual, "localhost:8010")
		convey.So(conf.Login[0].LogicAddr, convey.ShouldEqual, "1.2.1")
		convey.So(conf.Login[0].ClusterAddr, convey.ShouldEqual, "localhost:8012")
		convey.So(conf.Gate[0].ExternalAddr, convey.ShouldEqual, "10.128.2.244:9014")
	})

	var configData2 = `
[[Game]]
logicAddr		= "1.3.1"
clusterAddr    	= "localhost:8013"
`

	convey.Convey("Test Decode Config ...", t, func() {
		conf := &Config{}
		_, err := toml.Decode(configData2, conf)

		convey.So(err, convey.ShouldBeNil)

		conf.Game = append(conf.Game, &Game{
			Addr: &Addr{
				LogicAddr:   "1.3.2",
				ClusterAddr: "localhost:8014",
			},
		})

		s := &StringWriter{}
		encoder := toml.NewEncoder(s)
		err = encoder.Encode(conf)

		convey.So(err, convey.ShouldBeNil)

		log.Println(s.ToStr())

	})

}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../../../upload/template/config.toml.template")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
	t.Log(cfg.Dir[0])
}
