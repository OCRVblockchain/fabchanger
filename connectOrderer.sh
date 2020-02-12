#!/bin/sh

go build -o fabchanger cmd/main.go
mkdir artifacts
./fabchanger --mode getblock -o ./artifacts/configOrderer.json
./fabchanger --mode configtxtojson -o ./artifacts/extendOrderer.json --join orderer
./fabchanger --mode merge --f ./artifacts/configOrderer.json --merge ./artifacts/extendOrderer.json --o ./artifacts/merged.json --join orderer
./fabchanger --mode jsontoproto --f ./artifacts/configOrderer.json --o ./artifacts/configOrderer.pb --join orderer
./fabchanger --mode jsontoproto --f ./artifacts/merged.json --o ./artifacts/merged.pb --join orderer
./fabchanger --mode delta -f ./artifacts/configOrderer.pb -comparewith ./artifacts/merged.pb -o ./artifacts/delta.pb
./fabchanger --mode wrap -f ./artifacts/delta.pb -o wrappedDelta.pb
rm -rf ./artifacts/


