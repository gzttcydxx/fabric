export LOCAL_ROOT_PATH=$(shell pwd)
export LOCAL_CA_PATH=$(LOCAL_ROOT_PATH)/orgs
export DOCKER_COMPOSE_PATH=$(LOCAL_ROOT_PATH)/compose
export LOCAL_TEMPLATE_PATH=$(LOCAL_ROOT_PATH)/template
export FABRIC_CFG_PATH=$(LOCAL_ROOT_PATH)/config
export CHAINCODE_PATH=$(LOCAL_ROOT_PATH)/chaincode

include .env
export

export BASE_URL_SUBST:=$(subst .,-,$(BASE_URL))
export CHAIN_ID=$(shell uuidgen)

.PHONY: api
.DEFAULT_GOAL := all

check-root:
	@if [ `id -u` -ne 0 ]; then echo "You must be root to run this"; exit 1; fi

check-container:
	@if docker ps | grep -q 'council.$(BASE_URL)'; then echo "Container is running. Please down first!"; exit 1; fi

init:
	@mkdir -p ${DOCKER_COMPOSE_PATH}/explorer/connection-profile
	@envsubst < ${LOCAL_TEMPLATE_PATH}/configtx.yaml > ${FABRIC_CFG_PATH}/configtx.yaml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1soft > ${LOCAL_ROOT_PATH}/envpeer1soft
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1web > ${LOCAL_ROOT_PATH}/envpeer1web
	@envsubst < ${LOCAL_TEMPLATE_PATH}/envpeer1hard > ${LOCAL_ROOT_PATH}/envpeer1hard
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/config.json > ${DOCKER_COMPOSE_PATH}/explorer/config.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/connection-profile/soft-network.json > ${DOCKER_COMPOSE_PATH}/explorer/connection-profile/soft-network.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/explorer/connection-profile/web-network.json > ${DOCKER_COMPOSE_PATH}/explorer/connection-profile/web-network.json
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose/base.yml > ${DOCKER_COMPOSE_PATH}/base.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose/ca.yml > ${DOCKER_COMPOSE_PATH}/ca.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose/peer.yml > ${DOCKER_COMPOSE_PATH}/peer.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose/db.yml > ${DOCKER_COMPOSE_PATH}/db.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose/explorer.yml > ${DOCKER_COMPOSE_PATH}/explorer.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/compose.yml > ${LOCAL_ROOT_PATH}/compose.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/api/api.yml > ${LOCAL_ROOT_PATH}/api/api.yml
	@envsubst < ${LOCAL_TEMPLATE_PATH}/api/gateway/connection.gotmp > ${LOCAL_ROOT_PATH}/api/gateway/connection.go

up: check-root check-container init
	@scripts/up.sh

clean: check-root
	@if [ -e "data" ]; then rm -r "data"; fi
	@if [ -e "${LOCAL_CA_PATH}" ] && [ "${REMOVE_ORGS}" = "true" ]; then rm -r "${LOCAL_CA_PATH}"; fi
	@if [ -e "${DOCKER_COMPOSE_PATH}" ]; then rm -r "${DOCKER_COMPOSE_PATH}"; fi
	@if [ -e "compose.yml" ]; then rm "compose.yml"; fi
	@if [ -e "basic.tar.gz" ]; then rm "basic.tar.gz"; fi
	@if [ -e "envpeer1soft" ]; then rm "envpeer1soft"; fi
	@if [ -e "envpeer1web" ]; then rm "envpeer1web"; fi
	@if [ -e "envpeer1hard" ]; then rm "envpeer1hard"; fi
	@if [ "${DELETE_CHAINCODE}" = "true" ]; then docker images | awk '($1 ~ /dev-peer.*/) {print $3}' | xargs docker rmi; fi

down: check-root clean
	@docker compose down -v

code: check-root
	@scripts/code.sh

update: check-root
	@scripts/update.sh

explorer:
	@docker compose up -d explorerdb.${BASE_URL} explorer.${BASE_URL}

api:
	@docker compose down api.${BASE_URL}
	@docker compose up -d api.${BASE_URL}

api-log:
	@docker compose logs -f api.${BASE_URL}

build:
	@docker compose build

all: check-root down init up code all
