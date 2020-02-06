package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxgen/genesisconfig"
)

func main() {
	changer, err := fabchanger.New()
	if err != nil {
		log.Fatal(err)
	}

	switch changer.Config.Mode {

	case "configtxtojson":
		var topLevelConfig *genesisconfig.TopLevel
		topLevelConfig = genesisconfig.LoadTopLevel(changer.Config.General.ConfigTxPath)
		if changer.Config.General.OrgToJoinMSP != "" {
			if err := changer.ConfigTxToJSON("extend.json", topLevelConfig); err != nil {
				log.Fatalf("Error on printOrg: %s", err)
			}
		}

	case "getblock":
		block, err := changer.FetchBlock("block")
		if err != nil {
			log.Fatal(err)
		}
		err = changer.BlockToJSON("config.json", block)
		if err != nil {
			log.Fatal(err)
		}

	case "merge":
		err = changer.Merge("config.json", "extend.json", "new.json")
		if err != nil {
			log.Fatal(err)
		}
	}
}
