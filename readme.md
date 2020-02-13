 
## Fabchanger
 
 Tool for _easy_ Hyperledger Fabric network reconfiguration by committing a new configuration transaction

<br/>

#### Connect new organization
  1. Generate new crypto materials
  
      Add path to cryptogen bin to PATH env: 
      
         export PATH=$PATH:~/bin
         
      Edit `config/crypto-config.yaml` for your new configuration
     
      Generate crypto materials:
        
         ./generate.sh
      
      Copy generated crypto materials to your main `crypto-config` directory
  
  2. Create `configtx.yaml` in ./config dir using `configtx.yaml.org.sample` as sample and set config values in `/config.yaml`
  3. Run command:
      
         ./connectOrg.sh
  
  4. Sign configuration transaction
      
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
   1. Generate new crypto materials
     
         Add path to cryptogen bin to PATH env: 
         
            export PATH=$PATH:~/bin
            
         Edit `crypto-config.yaml` for your new configuration
        
         Generate crypto materials:
           
            ./generate.sh
         
         Copy generated crypto materials to your main `crypto-config` directory
  2. Create `configtx.yaml` in ./config dir using `configtx.yaml.orderer.sample` as sample and set config values in `/config.yaml`
  3. Run command:
      
         ./connectOrderer.sh
  
  4. Sign configuration transaction
      
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
         ./sign.sh cli commit orderer.org.ru:7050
         
You can sign not only from the cli container, just replace first script arg `./sign.sh <container> <action>`