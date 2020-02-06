package fabchanger

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	fabricconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/common/tools/protolator/protoext/ordererext"
	"github.com/pkg/errors"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/config"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxgen/encoder"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxgen/genesisconfig"
	"io/ioutil"
	"os"
)

type FabChanger struct {
	Config *config.Config
}

func New() (*FabChanger, error) {
	configuration, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return &FabChanger{Config: configuration}, nil
}

func (f *FabChanger) ConfigTxToJSON(JSONFileName string, t *genesisconfig.TopLevel) error {
	for _, org := range t.Organizations {
		if org.Name == f.Config.General.OrgToJoinMSP {
			og, err := encoder.NewOrdererOrgGroup(org)
			if err != nil {
				return errors.Wrapf(err, "bad org definition for org %s", org.Name)
			}
			file, err := os.OpenFile(JSONFileName, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			if err := protolator.DeepMarshalJSON(file, &ordererext.DynamicOrdererOrgGroup{ConfigGroup: og}); err != nil {
				return errors.Wrapf(err, "malformed org definition for org: %s", org.Name)
			}
			return nil
		}
	}
	return errors.Errorf("organization %s not found", f.Config.General.OrgToJoinMSP)
}

func (f *FabChanger) FetchBlock(blockName string) (*common.Block, error) {
	fabConfig := fabricconfig.FromFile(f.Config.General.ConnectionProfile)
	sdk, err := fabsdk.New(fabConfig)
	if err != nil {
		return nil, err
	}
	defer sdk.Close()

	clientChannelContext := sdk.ChannelContext(f.Config.Channel, fabsdk.WithUser(f.Config.Identity), fabsdk.WithOrg(f.Config.MyOrg))
	ledgerClient, err := ledger.New(clientChannelContext)
	if err != nil {
		return nil, err
	}

	block, err := ledgerClient.QueryConfigBlock()
	if err != nil {
		return nil, err
	}

	//b, err := proto.Marshal(block)
	//if err != nil {
	//	return err
	//}
	//
	//if err = ioutil.WriteFile(blockName, b, 0644); err != nil {
	//	return err
	//}

	return block, nil
}

func (f *FabChanger) BlockToJSON(blockFileName string, b *common.Block) error {

	var buffer bytes.Buffer
	err := protolator.DeepMarshalJSON(&buffer, b)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(blockFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	_, err = buffer.WriteTo(file)

	if err := file.Close(); err != nil {
		return err
	}

	return err
}

func (f *FabChanger) Merge(oldConfig, extendConfig, newFile string) error {
	oldFileBytes, err := ioutil.ReadFile(oldConfig)
	if err != nil {
		return err
	}
	extendConfigBytes, err := ioutil.ReadFile(extendConfig)
	if err != nil {
		return err
	}

	var oldConfigJson = make(map[string]interface{})
	err = json.Unmarshal(oldFileBytes, &oldConfigJson)
	if err != nil {
		return err
	}
	var extendConfigJson = make(map[string]interface{})
	err = json.Unmarshal(extendConfigBytes, &extendConfigJson)
	if err != nil {
		return err
	}
	//[0]["payload"].(map[string]interface{})["data"].(map[string]interface{})["config"]
	//.(map[string]interface{}["payload"]
	newConfigJSON := oldConfigJson
	newConfigJSON["data"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["payload"].(map[string]interface{})["data"].(map[string]interface{})["config"].(map[string]interface{})["channel_group"].(map[string]interface{})["groups"].(map[string]interface{})["Application"].(map[string]interface{})["groups"].(map[string]interface{})[f.Config.OrgToJoinMSP] = extendConfigJson

	bytesJson, err := json.Marshal(newConfigJSON)
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(newFile, bytesJson, 0755)
	if err != nil {
		return err
	}

	return nil
}
