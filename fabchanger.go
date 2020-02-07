package fabchanger

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	fabricconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/common/tools/protolator/protoext/ordererext"
	"github.com/pkg/errors"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/config"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxgen/encoder"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxgen/genesisconfig"
	"gitlab.sch.ocrv.com.rzd/blockchain/fabchanger/configtxlator/update"
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

func (f *FabChanger) FetchBlock(blockName string) error {
	fabConfig := fabricconfig.FromFile(f.Config.General.ConnectionProfile)
	sdk, err := fabsdk.New(fabConfig)
	if err != nil {
		return err
	}
	defer sdk.Close()

	clientChannelContext := sdk.ChannelContext(f.Config.Channel, fabsdk.WithUser(f.Config.Identity), fabsdk.WithOrg(f.Config.MyOrg))
	ledgerClient, err := ledger.New(clientChannelContext)
	if err != nil {
		return err
	}

	block, err := ledgerClient.QueryConfigBlock()
	if err != nil {
		return err
	}

	b, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(blockName, b, 0644); err != nil {
		return err
	}

	return nil
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
			if err := file.Close(); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.Errorf("organization %s not found", f.Config.General.OrgToJoinMSP)
}

func (f *FabChanger) BlockToJSON(blockFileName, newFileName string) error {

	var block *common.Block

	blockFileBytes, err := ioutil.ReadFile(blockFileName)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(blockFileBytes, block)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = protolator.DeepMarshalJSON(&buffer, block)
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

func (f *FabChanger) Merge(mode, oldConfig, extendConfig, newFile string) error {
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

	if mode == "org" {
		newConfigJSON["data"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["payload"].(map[string]interface{})["data"].(map[string]interface{})["config"].(map[string]interface{})["channel_group"].(map[string]interface{})["groups"].(map[string]interface{})["Application"].(map[string]interface{})["groups"].(map[string]interface{})[f.Config.OrgToJoinMSP] = extendConfigJson
	} else if mode == "orderer" {
		newConfigJSON["data"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["payload"].(map[string]interface{})["data"].(map[string]interface{})["config"].(map[string]interface{})["channel_group"].(map[string]interface{})["groups"].(map[string]interface{})["Orderer"].(map[string]interface{})["groups"].(map[string]interface{})[f.Config.OrgToJoinMSP] = extendConfigJson
	} else {
		return errors.New("Mode not specified")
	}

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

func (f *FabChanger) JSONToProtoConfig(source, newName string) error {

	var msg *common.Config

	file, err := os.Open(source)
	if err != nil {
		return err
	}

	err = protolator.DeepUnmarshalJSON(file, msg)
	if err != nil {
		return err
	}

	blockBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	newFile, err := os.OpenFile(newName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(blockBytes)
	_, err = buf.WriteTo(newFile)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	if err := newFile.Close(); err != nil {
		return err
	}

	return err
}

func (f *FabChanger) ComputeDelta(original, updated, output string) error {
	originalFile, err := os.Open(original)
	if err != nil {
		return err
	}
	updatedFile, err := os.Open(updated)
	if err != nil {
		return err
	}
	outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	origIn, err := ioutil.ReadAll(originalFile)
	if err != nil {
		return errors.Wrapf(err, "error reading original config")
	}

	origConf := &cb.Config{}
	err = proto.Unmarshal(origIn, origConf)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling original config")
	}

	updtIn, err := ioutil.ReadAll(updatedFile)
	if err != nil {
		return errors.Wrapf(err, "error reading updated config")
	}

	updtConf := &cb.Config{}
	err = proto.Unmarshal(updtIn, updtConf)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling updated config")
	}

	cu, err := update.Compute(origConf, updtConf)
	if err != nil {
		return errors.Wrapf(err, "error computing config update")
	}

	cu.ChannelId = f.Config.Channel

	outBytes, err := proto.Marshal(cu)
	if err != nil {
		return errors.Wrapf(err, "error marshaling computed config update")
	}

	_, err = outputFile.Write(outBytes)
	if err != nil {
		return errors.Wrapf(err, "error writing config update to output")
	}

	if err := originalFile.Close(); err != nil {
		return err
	}
	if err := updatedFile.Close(); err != nil {
		return err
	}
	if err := outputFile.Close(); err != nil {
		return err
	}

	return nil
}
