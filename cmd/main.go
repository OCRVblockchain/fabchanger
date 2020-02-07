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

	case "getblock":
		block, err := changer.FetchBlock()
		if err != nil {
			log.Fatal(err)
		}
		err = changer.BlockToJSON(block, "config.json")
		if err != nil {
			log.Fatal(err)
		}

	case "configtxtojson":
		var topLevelConfig *genesisconfig.TopLevel
		topLevelConfig = genesisconfig.LoadTopLevel(changer.Config.General.ConfigTxPath)
		if changer.Config.General.OrgToJoinMSP != "" {
			if err := changer.ConfigTxToJSON("extend.json", topLevelConfig); err != nil {
				log.Fatalf("Error on printOrg: %s", err)
			}
		}

	case "merge":
		err = changer.Merge("org", "config.json", "extend.json", "new.json")
		if err != nil {
			log.Fatal(err)
		}

	case "jsontoproto":
		err = changer.JSONToProtoConfig("config.json", "old.pb")
		if err != nil {
			log.Fatal(err)
		}

		err = changer.JSONToProtoConfig("new.json", "new.pb")
		if err != nil {
			log.Fatal(err)
		}

	case "delta":
		err = changer.ComputeDelta("old.pb", "new.pb", "delta.pb")
		if err != nil {
			log.Fatal(err)
		}

	case "sign":
		err = changer.Sign("delta.pb")
		if err != nil {
			log.Fatal(err)
		}
	}
}
