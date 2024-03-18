docker-compose up -d council.$BASE_URL orderer.$BASE_URL soft.$BASE_URL web.$BASE_URL
sleep 10

# 创建 council 组织的 ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/council.$BASE_URL/ca/admin
# 使用 enroll 登录引导账户
fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@council.$BASE_URL
# 使用 register 注册用户
fabric-ca-client register -d --id.name orderer1 --id.secret orderer1 --id.type orderer -u https://council.$BASE_URL
fabric-ca-client register -d --id.name peer1soft --id.secret peer1soft --id.type peer -u https://council.$BASE_URL
fabric-ca-client register -d --id.name peer1web --id.secret peer1web --id.type peer -u https://council.$BASE_URL

# 创建 orderer 组织的 ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/orderer.$BASE_URL/ca/crypto/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/orderer.$BASE_URL/ca/admin
fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@orderer.$BASE_URL
fabric-ca-client register -d --id.name orderer1 --id.secret orderer1 --id.type orderer -u https://orderer.$BASE_URL
fabric-ca-client register -d --id.name admin1 --id.secret admin1 --id.type admin -u https://orderer.$BASE_URL

# 创建 soft 组织的 ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/soft.$BASE_URL/ca/crypto/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/soft.$BASE_URL/ca/admin
fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@soft.$BASE_URL
fabric-ca-client register -d --id.name peer1 --id.secret peer1 --id.type peer -u https://soft.$BASE_URL
fabric-ca-client register -d --id.name admin1 --id.secret admin1 --id.type admin -u https://soft.$BASE_URL
fabric-ca-client register -d --id.name user1 --id.secret user1 --id.type client -u https://soft.$BASE_URL

# 创建 web 组织的 ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/web.$BASE_URL/ca/crypto/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/web.$BASE_URL/ca/admin
fabric-ca-client enroll -d -u https://ca-admin:ca-adminpw@web.$BASE_URL
fabric-ca-client register -d --id.name peer1 --id.secret peer1 --id.type peer -u https://web.$BASE_URL
fabric-ca-client register -d --id.name admin1 --id.secret admin1 --id.type admin -u https://web.$BASE_URL
fabric-ca-client register -d --id.name user1 --id.secret user1 --id.type client -u https://web.$BASE_URL

# 创建 assets 目录，用于储存本组织根证书和用于组间通信的 LTS-CA 根证书
mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/assets
cp $LOCAL_CA_PATH/orderer.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/assets/ca-cert.pem
cp $LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/assets/tls-ca-cert.pem

mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/assets
cp $LOCAL_CA_PATH/soft.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/assets/ca-cert.pem
cp $LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/assets/tls-ca-cert.pem

mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/assets 
cp $LOCAL_CA_PATH/web.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/web.$BASE_URL/assets/ca-cert.pem
cp $LOCAL_CA_PATH/council.$BASE_URL/ca/crypto/ca-cert.pem $LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem

# 构造 orderer 组织成员证书
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/orderer.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin1:admin1@orderer.$BASE_URL

mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1/msp/admincerts
cp $LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1/msp/admincerts/cert.pem

export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/orderer.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://orderer1:orderer1@orderer.$BASE_URL
mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/msp/admincerts
cp $LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/msp/admincerts/cert.pem

export FABRIC_CA_CLIENT_MSPDIR=tls-msp
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/orderer.$BASE_URL/assets/tls-ca-cert.pem
fabric-ca-client enroll -d -u https://orderer1:orderer1@council.$BASE_URL --enrollment.profile tls --csr.hosts orderer1.orderer.$BASE_URL
cp $LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/keystore/*_sk $LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/keystore/key.pem

mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/msp/admincerts
mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/msp/cacerts
mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/msp/tlscacerts
mkdir -p $LOCAL_CA_PATH/orderer.$BASE_URL/msp/users
cp $LOCAL_CA_PATH/orderer.$BASE_URL/assets/ca-cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/msp/cacerts/
cp $LOCAL_CA_PATH/orderer.$BASE_URL/assets/tls-ca-cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/msp/tlscacerts/
cp $LOCAL_CA_PATH/orderer.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/orderer.$BASE_URL/msp/admincerts/cert.pem
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/orderer.$BASE_URL/msp/config.yaml

# 构造 soft 组织成员证书
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/soft.$BASE_URL/registers/user1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/soft.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://user1:user1@soft.$BASE_URL

export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/soft.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin1:admin1@soft.$BASE_URL
mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/admincerts
cp $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/admincerts/cert.pem

export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/soft.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://peer1:peer1@soft.$BASE_URL

export FABRIC_CA_CLIENT_MSPDIR=tls-msp
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/soft.$BASE_URL/assets/tls-ca-cert.pem
fabric-ca-client enroll -d -u https://peer1soft:peer1soft@council.$BASE_URL --enrollment.profile tls --csr.hosts peer1.soft.$BASE_URL
cp $LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1/tls-msp/keystore/*_sk $LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1/tls-msp/keystore/key.pem
mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1/msp/admincerts
cp $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1/msp/admincerts/cert.pem
# ?是否安全
cp $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/keystore/*_sk $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/keystore/key.pem

mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/msp/admincerts
mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/msp/cacerts
mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/msp/tlscacerts
mkdir -p $LOCAL_CA_PATH/soft.$BASE_URL/msp/users
cp $LOCAL_CA_PATH/soft.$BASE_URL/assets/ca-cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/msp/cacerts/
cp $LOCAL_CA_PATH/soft.$BASE_URL/assets/tls-ca-cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/msp/tlscacerts/
cp $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/soft.$BASE_URL/msp/admincerts/cert.pem
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/soft.$BASE_URL/msp/config.yaml

cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/soft.$BASE_URL/registers/user1/msp/config.yaml
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp/config.yaml
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/soft.$BASE_URL/registers/peer1/msp/config.yaml

# 构造 web 组织成员证书
export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/web.$BASE_URL/registers/admin1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/web.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://admin1:admin1@web.$BASE_URL
mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/admincerts
cp $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/admincerts/cert.pem

export FABRIC_CA_CLIENT_HOME=$LOCAL_CA_PATH/web.$BASE_URL/registers/peer1
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/web.$BASE_URL/assets/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=msp
fabric-ca-client enroll -d -u https://peer1:peer1@web.$BASE_URL

export FABRIC_CA_CLIENT_MSPDIR=tls-msp
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem
fabric-ca-client enroll -d -u https://peer1web:peer1web@council.$BASE_URL --enrollment.profile tls --csr.hosts peer1.web.$BASE_URL
cp $LOCAL_CA_PATH/web.$BASE_URL/registers/peer1/tls-msp/keystore/*_sk $LOCAL_CA_PATH/web.$BASE_URL/registers/peer1/tls-msp/keystore/key.pem
mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/registers/peer1/msp/admincerts
cp $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/web.$BASE_URL/registers/peer1/msp/admincerts/cert.pem
# ?是否安全
cp $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/keystore/*_sk $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/keystore/key.pem

mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/msp/admincerts
mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/msp/cacerts
mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/msp/tlscacerts
mkdir -p $LOCAL_CA_PATH/web.$BASE_URL/msp/users
cp $LOCAL_CA_PATH/web.$BASE_URL/assets/ca-cert.pem $LOCAL_CA_PATH/web.$BASE_URL/msp/cacerts/
cp $LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem $LOCAL_CA_PATH/web.$BASE_URL/msp/tlscacerts/
cp $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/signcerts/cert.pem $LOCAL_CA_PATH/web.$BASE_URL/msp/admincerts/cert.pem
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/web.$BASE_URL/msp/config.yaml

cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/web.$BASE_URL/registers/user1/msp/config.yaml
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp/config.yaml
cp $LOCAL_ROOT_PATH/config/config-msp.yaml $LOCAL_CA_PATH/web.$BASE_URL/registers/peer1/msp/config.yaml

find $LOCAL_CA_PATH/ -regex ".+cacerts.+.pem" -not -regex ".+tlscacerts.+" | rename 's/cacerts\/.+\.pem/cacerts\/ca-cert\.pem/'

export LOCAL_ROOT_PATH=$PWD
export LOCAL_CA_PATH=$LOCAL_ROOT_PATH/orgs
export FABRIC_CA_CLIENT_MSPDIR=tls-msp
export FABRIC_CA_CLIENT_TLS_CERTFILES=$LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem

mkdir $LOCAL_ROOT_PATH/data

configtxgen -configPath $LOCAL_ROOT_PATH/config -profile OrgsOrdererGenesis -outputBlock $LOCAL_ROOT_PATH/data/genesis.block -channelID syschannel
configtxgen -configPath $LOCAL_ROOT_PATH/config -profile OrgsChannel -outputCreateChannelTx $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME

docker-compose up -d orderer1.orderer.$BASE_URL peer1.soft.$BASE_URL peer1.web.$BASE_URL
sleep 10

export FABRIC_CFG_PATH=$LOCAL_ROOT_PATH/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="softMSP"
export CORE_PEER_ADDRESS=peer1.soft.$BASE_URL:443
export CORE_PEER_TLS_ROOTCERT_FILE=$LOCAL_CA_PATH/soft.$BASE_URL/assets/tls-ca-cert.pem
export CORE_PEER_MSPCONFIGPATH=$LOCAL_CA_PATH/soft.$BASE_URL/registers/admin1/msp
export ORDERER_CA=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/tlscacerts/tls-council-$BASE_URL_SUBST.pem
peer channel create -c $CHANNEL_NAME -f $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.tx -o orderer1.orderer.$BASE_URL:443 --tls --cafile $ORDERER_CA --outputBlock $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block

cp $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block $LOCAL_CA_PATH/soft.$BASE_URL/assets/
cp $LOCAL_ROOT_PATH/data/$CHANNEL_NAME.block $LOCAL_CA_PATH/web.$BASE_URL/assets/

peer channel join -b $LOCAL_CA_PATH/soft.$BASE_URL/assets/$CHANNEL_NAME.block

export CORE_PEER_LOCALMSPID="webMSP"
export CORE_PEER_ADDRESS=peer1.web.$BASE_URL:443
export CORE_PEER_TLS_ROOTCERT_FILE=$LOCAL_CA_PATH/web.$BASE_URL/assets/tls-ca-cert.pem
export CORE_PEER_MSPCONFIGPATH=$LOCAL_CA_PATH/web.$BASE_URL/registers/admin1/msp
export ORDERER_CA=$LOCAL_CA_PATH/orderer.$BASE_URL/registers/orderer1/tls-msp/tlscacerts/tls-council-$BASE_URL_SUBST.pem

peer channel join -b $LOCAL_CA_PATH/web.$BASE_URL/assets/$CHANNEL_NAME.block

docker-compose up -d explorerdb.$BASE_URL explorer.$BASE_URL
