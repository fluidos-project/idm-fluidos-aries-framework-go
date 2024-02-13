#!/bin/bash
echo "FABRIC PATH -- ${FABRIC_PATH}"
# download docker samples and binaries 2.5.1
 cd $(dirname "$(realpath ${FABRIC_PATH})")
pwd
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples
git checkout v2.4.9
cd ..

./install-fabric.sh --fabric-version $FABRIC_VERSION d b 

rm ./install-fabric.sh

echo "Descargado samples, imagenes de docker de $FABRIC_VERSION y binarios de HYPERLEDGER FABRIC"
