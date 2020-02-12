#!/bin/sh
if [ $2 = "commit" ]; then
  docker exec $1 peer channel update -f wrapped.pb -c $CHANNEL_NAME -o $3 --tls --cafile $ORDERER_CA
else
  docker cp wrappedDelta.pb $1:/opt/gopath/src/github.com/hyperledger/fabric/peer/wrapped.pb
  docker exec $1 /bin/sh -c "CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH CORE_PEER_TLS_ROOTCERT_FILE=$CORE_PEER_TLS_ROOTCERT_FILE  peer channel signconfigtx -f wrapped.pb"
  docker cp $1:/opt/gopath/src/github.com/hyperledger/fabric/peer/wrapped.pb wrappedDelta.pb
fi