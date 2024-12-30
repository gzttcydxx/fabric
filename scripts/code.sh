#!/bin/bash

function ORDERER_NUMBER() {
    echo $(((RANDOM % 3) + 1))
}

# 初始化版本和序列号
if [ -z "$VERSION" ]; then
    VERSION="1.0"
fi
if [ -z "$SEQUENCE" ]; then
    SEQUENCE=1
fi

echo "Initial VERSION: $VERSION"
echo "Initial SEQUENCE: $SEQUENCE"

# 检查链码是否已安装
function check_chaincode_installed() {
    local org=$1
    source $LOCAL_ROOT_PATH/envpeer1$org
    local installed=$(peer lifecycle chaincode queryinstalled | grep "$CHAINCODE_NAME" || true)
    if [ ! -z "$installed" ]; then
        echo "Chaincode already installed on peer1.$org"
        return 0
    fi
    return 1
}

# 获取已安装链码的最新版本和序列号
function get_chaincode_info() {
    local org=$1
    source $LOCAL_ROOT_PATH/envpeer1$org
    local info=$(peer lifecycle chaincode querycommitted -C $CHANNEL_NAME -n $CHAINCODE_NAME 2>/dev/null || true)
    echo "Query committed info: $info"
    if [ ! -z "$info" ]; then
        # 获取当前序列号并增加
        SEQUENCE=$(echo "$info" | tr ',' '\n' | grep "Sequence:" | awk '{print $2}')
        # 设置新的版本号
        VERSION="1.$SEQUENCE"
        # 增加序列号
        SEQUENCE=$((SEQUENCE + 1))

        rm -rf $LOCAL_ROOT_PATH/basic.tar.gz

        echo "Updating to Version: $VERSION, Sequence: $SEQUENCE"
        return 0
    fi
    return 1
}

get_chaincode_info "soft"
peer lifecycle chaincode package basic.tar.gz --path $CHAINCODE_PATH --label basic_${VERSION}

# 安装新打包的链码
function install_code() {
    local orgs=("$@")
    for org in "${orgs[@]}"; do
        (
            source $LOCAL_ROOT_PATH/envpeer1$org
            echo "Installing chaincode version ${VERSION} on peer1.$org..."
            peer lifecycle chaincode install basic.tar.gz
        ) &
    done
    wait
}
install_code "soft" "web" "hard"

# 获取 Package ID
export CHAINCODE_ID=$(peer lifecycle chaincode queryinstalled | grep "Package ID:" | grep "$CHAINCODE_NAME" | awk -F ', ' '{print $1}' | awk -F ': ' '{print $2}')

function approve_code() {
    local orgs=("$@")

    for org in "${orgs[@]}"; do
        source $LOCAL_ROOT_PATH/envpeer1$org
        peer lifecycle chaincode approveformyorg -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 \
            --tls --cafile $ORDERER_CA \
            --channelID $CHANNEL_NAME \
            --name $CHAINCODE_NAME \
            --version $VERSION \
            --sequence $SEQUENCE \
            --waitForEvent \
            --init-required \
            --package-id $CHAINCODE_ID
        peer lifecycle chaincode queryapproved -C $CHANNEL_NAME -n $CHAINCODE_NAME --sequence $SEQUENCE
    done
}
approve_code "soft" "web" "hard"

peer lifecycle chaincode checkcommitreadiness \
    -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 \
    --tls --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --version $VERSION \
    --sequence $SEQUENCE \
    --init-required

source $LOCAL_ROOT_PATH/envpeer1soft
peer lifecycle chaincode commit \
    -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 \
    --tls --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --init-required \
    --version $VERSION \
    --sequence $SEQUENCE \
    --peerAddresses peer1.soft.$BASE_URL:443 \
    --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE \
    --peerAddresses peer1.web.$BASE_URL:443 \
    --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE

peer lifecycle chaincode querycommitted \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 \
    --tls --cafile $ORDERER_CA \
    --peerAddresses peer1.soft.$BASE_URL:443 \
    --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE

# 初始化链码
peer chaincode invoke --isInit \
    -o orderer$(ORDERER_NUMBER).council.$BASE_URL:443 \
    --tls --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME \
    --name $CHAINCODE_NAME \
    --peerAddresses peer1.soft.$BASE_URL:443 \
    --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE \
    --peerAddresses peer1.web.$BASE_URL:443 \
    --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE \
    -c "{\"Args\":[\"InitLedger\", \"$CHAIN_ID\"]}"
