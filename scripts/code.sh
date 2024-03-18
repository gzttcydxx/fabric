export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="softMSP"
export CORE_PEER_ADDRESS=peer1.soft.$BASE_URL:443
export CORE_PEER_TLS_ROOTCERT_FILE=$LOCAL_CA_PATH/soft.$BASE_URL/assets/tls-ca-cert.pem
export CORE_PEER_MSPCONFIGPATH=$LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp
export ORDERER_CA=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/tlscacerts/tls-council-$BASE_URL_SUBST.pem

peer lifecycle chaincode package basic.tar.gz --path $LOCAL_ROOT_PATH/asset-transfer-basic/chaincode-go --label basic
peer lifecycle chaincode install basic.tar.gz
export CHAINCODE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | awk -F ', ' '{print $1}' | awk -F ': ' '{print $2}')
peer lifecycle chaincode approveformyorg -o orderer1.orderer.$BASE_URL:443 --tls --cafile $ORDERER_CA  --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID

export CORE_PEER_LOCALMSPID="webMSP"
export CORE_PEER_ADDRESS=peer1.web.$BASE_URL:443
export CORE_PEER_TLS_ROOTCERT_FILE=$LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem
export CORE_PEER_MSPCONFIGPATH=$LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp
export ORDERER_CA=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/tlscacerts/tls-council-$BASE_URL_SUBST.pem

peer lifecycle chaincode install basic.tar.gz
export CHAINCODE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | awk -F ', ' '{print $1}' | awk -F ': ' '{print $2}')
peer lifecycle chaincode approveformyorg -o orderer1.orderer.$BASE_URL:443 --tls --cafile $ORDERER_CA  --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID

peer lifecycle chaincode commit -o orderer1.orderer.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --init-required --version 1.0 --sequence 1 --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.web.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE

peer chaincode invoke --isInit -o orderer1.orderer.$BASE_URL:443 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --peerAddresses peer1.soft.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.web.$BASE_URL:443 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["InitLedger"]}'
