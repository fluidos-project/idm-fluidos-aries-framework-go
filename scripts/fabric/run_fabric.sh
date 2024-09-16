#!/bin/bash
cd $FABRIC_PATH
cd test-network

bash network.sh up -ca -s couchdb

bash network.sh up createChannel -c mychannel

bash network.sh deployCC -ccn sacc -ccp ../chaincode/sacc/ -ccv 1 -ccl go

echo "Desplegado e instalado Hyperledger Fabric, desplegado canal \"mychannel\" y desplegado smartcontract \"sacc\""

bash network.sh deployCC -ccn fluidosAccessHist -ccp ../chaincode/fluidosAccess-historical -ccv 1 -ccl go

echo "Desplegado e instalado Hyperledger Fabric, desplegado smartcontract \"fluidosAccessHist\" en el canal \"mychannel\""

bash network.sh deployCC -ccn xacml -ccp ../chaincode/xacml -ccv 1 -ccl go

echo "Desplegado e instalado Hyperledger Fabric, desplegado smartcontract \"xacml\" en el canal \"mychannel\""

privpath=./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore
certpath=./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts

certificate=$( sed  -z -e 's|\n|\\n|g' $certpath/$(ls $certpath | tr -d '\n'))
priv=$(sed  -z -e 's|\n|\\n|g' $privpath/$(ls $privpath | tr -d '\n'))

echo '{"identity":{"pub":"'"$certificate"'","priv":"'"$priv"'"},"connection-profile":' > $CONNECTION_PROFILE_PATH && cat ./organizations/peerOrganizations/org1.example.com/connection-org1.json >> $CONNECTION_PROFILE_PATH && echo '}' >> $CONNECTION_PROFILE_PATH

sed -i 's#grpcs://localhost:7051#grpcs://peer0.org1.example.com:7051#g' $CONNECTION_PROFILE_PATH
sed -i 's#https://localhost:7054#https://ca.org1.example.com:7054#g' $CONNECTION_PROFILE_PATH
