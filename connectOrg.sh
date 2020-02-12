#!/bin/sh

go build -o fabchanger cmd/main.go
mkdir artifacts
./fabchanger --mode getblock -o ./artifacts/configOrg.json
./fabchanger --mode configtxtojson -o ./artifacts/extendOrg.json --join org
./fabchanger --mode merge --f ./artifacts/configOrg.json --merge ./artifacts/extendOrg.json --o ./artifacts/merged.json --join org
./fabchanger --mode jsontoproto --f ./artifacts/configOrg.json --o ./artifacts/configOrg.pb
./fabchanger --mode jsontoproto --f ./artifacts/merged.json --o ./artifacts/merged.pb
./fabchanger --mode delta -f ./artifacts/configOrg.pb -comparewith ./artifacts/merged.pb -o ./artifacts/delta.pb
./fabchanger --mode wrap -f ./artifacts/delta.pb -o wrappedDelta.pb
rm -rf ./artifacts/


