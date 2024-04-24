if [ ! -d "$LOCAL_CA_PATH" ]; then
    docker-compose up -d council.$BASE_URL soft.$BASE_URL web.$BASE_URL hard.$BASE_URL
    sleep 5

    # 创建 council 组织的 ca
    export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem
    export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/council.$BASE_URL/ca/admin
    # 使用 enroll 登录引导账户
    fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@council.$BASE_URL
    # 使用 register 注册用户
    fabric-ca-client register -d --id.name admin1 --id.secret admin1 --id.type admin -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer1 --id.secret orderer1 --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer2 --id.secret orderer2 --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer3 --id.secret orderer3 --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer1-admin --id.secret orderer1-admin --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer2-admin --id.secret orderer2-admin --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name orderer3-admin --id.secret orderer3-admin --id.type orderer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name peer1soft --id.secret peer1soft --id.type peer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name peer1web --id.secret peer1web --id.type peer -u https://council.$BASE_URL
    fabric-ca-client register -d --id.name peer1hard --id.secret peer1hard --id.type peer -u https://council.$BASE_URL

    # 创建组织的 ca
    function register_orgs_ca() {
        local orgs=("$@")

        for org in "${orgs[@]}"; do
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/$org.$BASE_URL/ca/crypto/ca-cert.pem
            export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/$org.$BASE_URL/ca/admin
            fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@$org.$BASE_URL
            fabric-ca-client register -d --id.name peer1 --id.secret peer1 --id.type peer -u https://$org.$BASE_URL
            fabric-ca-client register -d --id.name admin1 --id.secret admin1 --id.type admin -u https://$org.$BASE_URL
            fabric-ca-client register -d --id.name user1 --id.secret user1 --id.type client -u https://$org.$BASE_URL
        done
    }
    register_orgs_ca "soft" "web" "hard"

    # 创建 assets 目录，用于储存本组织根证书和用于组间通信的 LTS-CA 根证书
    function create_assets() {
        local orgs=("$@")

        for org in "${orgs[@]}"; do
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/assets
            cp $LOCAL_CA_PATH/$org.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/assets/ca-cert.pem
            cp $LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/assets/tls-ca-cert.pem
        done
    }
    create_assets "council" "soft" "web" "hard"

    # 构造 council 组织成员证书
    # admin
    export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/council.$BASE_URL/registers/admin1
    export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/council.$BASE_URL/assets/ca-cert.pem
    export FABRIC_CA_CLIENT_MSPDIR=msp
    fabric-ca-client enroll -d -u https://admin1:admin1@council.$BASE_URL
    mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/registers/admin1/msp/admincerts
    cp $LOCAL_CA_PATH/council.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/council.$BASE_URL/registers/admin1/msp/admincerts/cert.pem

    # orderer
    function enroll_and_setup_orderer() {
        local orgs=("$@")

        for org in "${orgs[@]}"; do
            export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/council.$BASE_URL/registers/$org
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/council.$BASE_URL/assets/ca-cert.pem
            export FABRIC_CA_CLIENT_MSPDIR=msp
            fabric-ca-client enroll -d -u https://$org:$org@council.$BASE_URL
            mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/registers/$org/msp/admincerts
            cp $LOCAL_CA_PATH/council.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/council.$BASE_URL/registers/$org/msp/admincerts/cert.pem

            export FABRIC_CA_CLIENT_MSPDIR=tls-msp
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/council.$BASE_URL/assets/tls-ca-cert.pem
            fabric-ca-client enroll -d -u https://$org:$org@council.$BASE_URL --enrollment.profile tls --csr.hosts $org.council.$BASE_URL
            cp $LOCAL_CA_PATH/council.$BASE_URL/registers/$org/tls-msp/keystore/*_sk $LOCAL_CA_PATH/council.$BASE_URL/registers/$org/tls-msp/keystore/key.pem
        done
    }
    enroll_and_setup_orderer "orderer1" "orderer2" "orderer3" "orderer1-admin" "orderer2-admin" "orderer3-admin"

    mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/msp/admincerts
    mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/msp/cacerts
    mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/msp/tlscacerts
    mkdir -p $LOCAL_CA_PATH/council.$BASE_URL/msp/users
    cp $LOCAL_CA_PATH/council.$BASE_URL/assets/ca-cert.pem $LOCAL_CA_PATH/council.$BASE_URL/msp/cacerts/
    cp $LOCAL_CA_PATH/council.$BASE_URL/assets/tls-ca-cert.pem $LOCAL_CA_PATH/council.$BASE_URL/msp/tlscacerts/
    cp $LOCAL_CA_PATH/council.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/council.$BASE_URL/msp/admincerts/cert.pem
    cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/council.$BASE_URL/msp/config.yaml

    # 构造组织成员证书
    function enroll_and_setup_peer() {
        local orgs=("$@")

        for org in "${orgs[@]}"; do
            export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/$org.$BASE_URL/registers/user1
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/$org.$BASE_URL/assets/ca-cert.pem
            export FABRIC_CA_CLIENT_MSPDIR=msp
            fabric-ca-client enroll -d -u https://user1:user1@$org.$BASE_URL

            export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/$org.$BASE_URL/assets/ca-cert.pem
            export FABRIC_CA_CLIENT_MSPDIR=msp
            fabric-ca-client enroll -d -u https://admin1:admin1@$org.$BASE_URL
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/admincerts
            cp $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/admincerts/cert.pem

            export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/$org.$BASE_URL/assets/ca-cert.pem
            export FABRIC_CA_CLIENT_MSPDIR=msp
            fabric-ca-client enroll -d -u https://peer1:peer1@$org.$BASE_URL

            export FABRIC_CA_CLIENT_MSPDIR=tls-msp
            export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/$org.$BASE_URL/assets/tls-ca-cert.pem
            fabric-ca-client enroll -d -u https://peer1$org:peer1$org@council.$BASE_URL --enrollment.profile tls --csr.hosts peer1.$org.$BASE_URL
            cp $LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1/tls-msp/keystore/*_sk $LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1/tls-msp/keystore/key.pem
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1/msp/admincerts
            cp $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1/msp/admincerts/cert.pem
            # ?是否安全
            mv $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/keystore/*_sk $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/keystore/key.pem
            mv $LOCAL_CA_PATH/$org.$BASE_URL/registers/user1/msp/keystore/*_sk $LOCAL_CA_PATH/$org.$BASE_URL/registers/user1/msp/keystore/key.pem

            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/msp/admincerts
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/msp/cacerts
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/msp/tlscacerts
            mkdir -p $LOCAL_CA_PATH/$org.$BASE_URL/msp/users
            cp $LOCAL_CA_PATH/$org.$BASE_URL/assets/ca-cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/msp/cacerts/
            cp $LOCAL_CA_PATH/$org.$BASE_URL/assets/tls-ca-cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/msp/tlscacerts/
            cp $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/$org.$BASE_URL/msp/admincerts/cert.pem

            cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/$org.$BASE_URL/msp/config.yaml
            cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/$org.$BASE_URL/registers/user1/msp/config.yaml
            cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/$org.$BASE_URL/registers/admin1/msp/config.yaml
            cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/$org.$BASE_URL/registers/peer1/msp/config.yaml
        done
    }
    enroll_and_setup_peer "soft" "web" "hard"

    # 替换cacerts文件名
    find $LOCAL_CA_PATH -type f -regex ".+cacerts.+.pem" -not -regex ".+tlscacerts.+" -exec bash -c 'if [[ "$(basename "$1")" != "ca-cert.pem" ]]; then mv "$1" "$(dirname "$1")/ca-cert.pem"; fi' _ {} \;
else
    docker-compose up -d council.$BASE_URL soft.$BASE_URL web.$BASE_URL hard.$BASE_URL
    sleep 5
fi

docker-compose up -d peer1.soft.$BASE_URL peer1.web.$BASE_URL peer1.hard.$BASE_URL
docker-compose up -d orderer1.council.$BASE_URL orderer2.council.$BASE_URL orderer3.council.$BASE_URL
sleep 10

mkdir $LOCAL_ROOT_PATH/data

configtxgen -profile OrgsChannel -outputCreateChannelTx $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
configtxgen -profile OrgsChannel -outputBlock $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block -channelID $CHANNEL_NAME

cp $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block $LOCAL_CA_PATH/soft.$BASE_URL/assets/
cp $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block $LOCAL_CA_PATH/web.$BASE_URL/assets/
cp $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block $LOCAL_CA_PATH/hard.$BASE_URL/assets/

source $LOCAL_ROOT_PATH/envpeer1soft
function orderer_join_channel() {
    local orderers=("$@")

    for orderer in "${orderers[@]}"; do
        export ORDERER_ADMIN_TLS_SIGN_CERT=$LOCAL_CA_PATH/council.$BASE_URL/registers/$orderer-admin/tls-msp/signcerts/cert.pem
        export ORDERER_ADMIN_TLS_PRIVATE_KEY=$LOCAL_CA_PATH/council.$BASE_URL/registers/$orderer-admin/tls-msp/keystore/key.pem
        osnadmin channel join -o $orderer-admin.council.$BASE_URL --channelID $CHANNEL_NAME --config-block $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
        osnadmin channel list -o $orderer-admin.council.$BASE_URL --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
    done
}
orderer_join_channel "orderer1" "orderer2" "orderer3"

function peer_join_channel() {
    local peers=("$@")

    for peer in "${peers[@]}"; do
        source $LOCAL_ROOT_PATH/envpeer1$peer
        peer channel join -b $LOCAL_CA_PATH/$peer.$BASE_URL/assets/$CHANNEL_NAME.block
        peer channel list
    done
}
peer_join_channel "soft" "web" "hard"

# 添加当前用户访问权限，不能用于生产环境
chown -R 1000:1000 $LOCAL_CA_PATH
