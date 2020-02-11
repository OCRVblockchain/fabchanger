#!/bin/sh

go build -o fabchanger cmd/main.go
./fabchanger --mode getblock -o configOrderer.json
./fabchanger --mode configtxtojson -o extendOrderer.json --join orderer
./fabchanger --mode merge --f configOrderer.json --merge extendOrderer.json --o merged.json --join orderer
./fabchanger --mode jsontoproto --f configOrderer.json --o configOrderer.pb --join orderer
./fabchanger --mode jsontoproto --f merged.json --o merged.pb --join orderer
./fabchanger --mode delta -f configOrderer.pb -comparewith merged.pb -o delta.pb
./fabchanger --mode wrap -f delta.pb -o wrappedDelta.pb

