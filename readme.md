 
## Fabchanger
 
 Tool for _easy_ reconfiguration of Hyperledger Fabric network by committing a new configuration transaction

<br/>

#### Connect new organization
  1. Create `configtx.yaml` in ./config dir using `configtx.yaml.org.sample` as sample and set config values in `/config.yaml`
  2. Run command:
      
         ./connectOrg.sh
  
  3. Sign configuration transaction
      
      Export environment variables:
  
         export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
         export CHANNEL_NAME=mychannel
         export CORE_PEER_LOCALMSPID="Org1MSP"
         export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
         export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
         
      Sign:
      
         ./sign.sh cli sign
         
      Repeat this step for all orgs, then commit update transaction:
        
         ./sign.sh cli commit
         
         
<br/><br/>         
#### Connect new RAFT orderer
  1. Create `configtx.yaml` in ./config dir using `configtx.yaml.orderer.sample` as sample and set config values in `/config.yaml`
  2. Run command:
      
         ./connectOrderer.sh
  
  3. Sign configuration transaction
      
      Export environment variables:
  
         export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
         export CHANNEL_NAME=mychannel
         export CORE_PEER_LOCALMSPID="Org1MSP"
         export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
         export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
         
      Sign:
      
         ./sign.sh cli sign
         
      Repeat this step for all orgs, then commit update transaction:
        
         export CORE_PEER_LOCALMSPID="OrdererMSP"
         export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/users/Admin@example.com/msp
         ./sign.sh cli commit
         
You can sign not only from the cli container, just replace first script arg `./sign.sh <container> <action>`