export LOCAL_ROOT_PATH=$(shell pwd)
export LOCAL_CA_PATH=$(LOCAL_ROOT_PATH)/orgs
export DOCKER_COMPOSE_PATH=$(LOCAL_ROOT_PATH)/compose
export LOCAL_TEMPLATE_PATH=$(LOCAL_ROOT_PATH)/template
export FABRIC_CFG_PATH=$(LOCAL_ROOT_PATH)/config

export FABRIC_BASE_VERSION=2.5
export FABRIC_CA_VERSION=1.5
export COUCHDB_VERSION=3.3

export COMPOSE_PROJECT_NAME=chain-a
export DOCKER_NETWORKS=chain-a
export CHANNEL_NAME=mychannel
export CHAINCODE_NAME=basic
# export CHAINCODE_PATH=$(LOCAL_ROOT_PATH)/asset-transfer-basic/chaincode-go
export CHAINCODE_PATH=$(LOCAL_ROOT_PATH)/chaincode

export BASE_URL=a.gzttc.top
export BASE_URL_SUBST:=$(subst .,-,$(BASE_URL))
# export SOFT_PASSWORD := $(openssl rand -base64 12)
# export WEB_PASSWORD := $(openssl rand -base64 12)
export SOFT_PASSWORD=$$X!KzqVcGt7FXwpC
export WEB_PASSWORD=UGwmTy$$5fJyN%9%E
export COUCHDB_PASSWORD=B5Vr6ecYta2W7bD9

check-root:
	@if [ `id -u` -ne 0 ]; then echo "You must be root to run this"; exit 1; fi

check-container:
	@if docker ps | grep -q 'council.$(BASE_URL)'; then echo "Container is running. Please down first!"; exit 1; fi

init:
	@envsubst < ${LOCAL_TEMPLATE_PATH}/configtx.yaml > ${FABRIC_CFG_PATH}/configtx.yaml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1soft > ${LOCAL_ROOT_PATH}/envpeer1soft
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1web > ${LOCAL_ROOT_PATH}/envpeer1web
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1hard > ${LOCAL_ROOT_PATH}/envpeer1hard
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/config.json > ${DOCKER_COMPOSE_PATH}/explorer/config.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/connection-profile/soft-network.json > ${DOCKER_COMPOSE_PATH}/explorer/connection-profile/soft-network.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/connection-profile/web-network.json > ${DOCKER_COMPOSE_PATH}/explorer/connection-profile/web-network.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/base.yml > ${DOCKER_COMPOSE_PATH}/base.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/ca.yml > ${DOCKER_COMPOSE_PATH}/ca.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/peer.yml > ${DOCKER_COMPOSE_PATH}/peer.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/db.yml > ${DOCKER_COMPOSE_PATH}/db.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer.yml > ${DOCKER_COMPOSE_PATH}/explorer.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/docker-compose.yml > ${LOCAL_ROOT_PATH}/docker-compose.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/caliper/caliper.yml > ${LOCAL_ROOT_PATH}/caliper/caliper.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/caliper/ccp.yml > ${LOCAL_ROOT_PATH}/caliper/networks/ccp.yaml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/caliper/networkConfig.yaml > ${LOCAL_ROOT_PATH}/caliper/networks/networkConfig.yaml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/caliper/report.html > ${LOCAL_ROOT_PATH}/caliper/report.html

up: check-root check-container init
	@scripts/up.sh

clean: check-root
	@if [ -e "data" ]; then rm -r "data"; fi
	@if [ -e "orgs" ]; then rm -r "orgs"; fi
	@if [ -e "basic.tar.gz" ]; then rm "basic.tar.gz"; fi

down: check-root clean
	@docker-compose down -v

code: check-root
	@scripts/code.sh

update: check-root
	@scripts/update.sh

test:
	@docker-compose up -d caliper

all: check-root down init up code
