function ORDERER_NUMBER() {
    echo 1
}

if [ -z "$VERSION" ]; then
    VERSION=1.0
fi
if [ -z "$SEQUENCE" ]; then
    SEQUENCE=1
fi

echo "VERSION: $VERSION"
echo "SEQUENCE: $SEQUENCE"

peer lifecycle chaincode package basic.tar.gz --path $CHAINCODE_PATH --label basic

function install_code() {
    local orgs=("$@")

    for org in "${orgs[@]}"; do
        source $LOCAL_ROOT_PATH/envpeer1$org
        peer lifecycle chaincode install basic.tar.gz
        peer lifecycle chaincode queryinstalled
    done
}
install_code "soft" "web" "hard"

export CHAINCODE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | awk -F ', ' '{print $1}' | awk -F ': ' '{print $2}')

function approve_code() {
    local orgs=("$@")

    for org in "${orgs[@]}"; do
        source $LOCAL_ROOT_PATH/envpeer1$org
        peer lifecycle chaincode approveformyorg -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA  --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $VERSION --sequence $SEQUENCE --waitForEvent --init-required --package-id $CHAINCODE_ID
        peer lifecycle chaincode queryapproved -C $CHANNEL_NAME -n $CHAINCODE_NAME --sequence $SEQUENCE
    done
}
approve_code "soft" "web" "hard"

peer lifecycle chaincode checkcommitreadiness -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $VERSION --sequence $SEQUENCE --init-required

source $LOCAL_ROOT_PATH/envpeer1soft
peer lifecycle chaincode commit -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --init-required --version $VERSION --sequence $SEQUENCE --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.web.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name $CHAINCODE_NAME -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE
peer chaincode invoke --isInit -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.web.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c "{\"Args\":[\"InitLedger\", \"$CHAIN_ID\"]}"

# sleep 5

# peer chaincode invoke -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.web.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["GetAllAssets"]}'
