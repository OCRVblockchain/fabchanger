package main

import (
	"github.com/OCRVblockchain/fabchanger"
	"github.com/OCRVblockchain/fabchanger/configtxgen/genesisconfig"
	log "github.com/sirupsen/logrus"
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
		err = changer.BlockToJSON(block, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "configtxtojson":
		var topLevelConfig *genesisconfig.TopLevel
		topLevelConfig = genesisconfig.LoadTopLevel(changer.Config.General.ConfigTxPath)
		if changer.Config.Connect.OrgToJoinMSP != "" {
			if err := changer.ConfigTxToJSON(changer.Config.Output, topLevelConfig); err != nil {
				log.Fatalf("Error on printOrg: %s", err)
			}
		}

	case "merge":
		err = changer.Merge(changer.Config.Input, changer.Config.Merge, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "jsontoproto":
		err = changer.JSONToProtoConfig(changer.Config.Input, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "delta":
		err = changer.ComputeDelta(changer.Config.Input, changer.Config.CompareWith, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "wrap":
		err = changer.Wrap(changer.Config.Input, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "sign":
		err = changer.Sign(changer.Config.Input, changer.Config.Output)
		if err != nil {
			log.Fatal(err)
		}

	case "generate":
		err = changer.GenerateConfigs()
		if err != nil {
			log.Fatal(err)
		}
	}

}
