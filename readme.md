 
## Fabchanger
 
 Tool for _easy_ Hyperledger Fabric network reconfiguration by committing a new configuration transaction

<br/>

#### Futher development

Fabchanger just works right now, but we must make the utility even simpler. For this purpose, it is necessary to simplify the creation of configuration artifacts: `configtx.yaml`, `crypto-config.yaml`, `connection.yaml`.  

<br/>

**Reference**
 
 - [Connect new organization](#org)

 - [Connect new RAFT orderer](#orderer)

<br/>

#### <a name=org>Connect new organization</a>
  1. Build Fabchanger:
  
         go build -o fabchanger cmd/main.go
  
  2. Generate new crypto materials
      
      Edit `config/config.yaml`
     
      Generate configs: 
      
          ./fabchanger --mode generate --join org
  
  3. Create connection profile and refer to it from `config/config.yaml`.
  
  4. Run command:
      
         ./connectOrg.sh
  
  5. Sign configuration transaction
     
      Sign as admin specified in config/config.yaml:
      
         ./fabchanger --mode sign -f ./wrappedDelta.pb -o ./wrappedDelta.pb 
         
      Repeat this step for all orgs (change credentials in config/config.yaml and run command specified above)
       
      Commit tx to orderer:
        
         ./fabchanger --mode update -f ./wrappedDelta.pb
         
         
<br/><br/>         
#### <a name=orderer>Connect new RAFT orderer</a>
   1. Build Fabchanger:
     
            go build -o fabchanger cmd/main.go
     
   2. Generate new crypto materials
          
       Edit `config/config.yaml`
        
       Generate configs: 
       
          ./fabchanger --mode generate --join orderer
  
  3. Create connection profile and refer to it from `config/config.yaml`.
       
  4. Run command:
      
         ./connectOrderer.sh
  
  5. Sign configuration transaction
       
        Sign as admin specified in config/config.yaml:
        
           ./fabchanger --mode sign -f ./wrappedDelta.pb -o ./wrappedDelta.pb 
           
        Repeat this step for all orgs (change credentials in config/config.yaml and run command specified above)
         
        Commit tx to orderer:
          
           ./fabchanger --mode update -f ./wrappedDelta.pb