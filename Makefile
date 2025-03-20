# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0

GO_CMD ?= go
ARIES_AGENT_REST_PATH=cmd/aries-agent-rest
ARIES_AGENT_MOBILE_PATH=cmd/aries-agent-mobile
SIDETREE_CLI_PATH=test/bdd/cmd/sidetree
OPENAPI_DOCKER_IMG=quay.io/goswagger/swagger
OPENAPI_SPEC_PATH=build/rest/openapi/spec
OPENAPI_DOCKER_IMG_VERSION=v0.23.0

FABRIC_PATH=./modules/fabric-samples
FABRIC_VERSION=2.5.1
CONNECTION_PROFILE_PATH=${PWD}/test/bdd/fixtures/agent-rest/data/connection-profile.json

PABC_PATH=./modules/p-abc

# Namespace for the agent images
DOCKER_OUTPUT_NS   ?= aries-framework-go
AGENT_REST_IMAGE_NAME   ?= agent-rest
WEBHOOK_IMAGE_NAME ?= sample-webhook

# Tool commands (overridable)
DOCKER_CMD ?= docker
GO_CMD     ?= go
ALPINE_VER ?= 3.16
GO_TAGS    ?=
GO_VER ?= 1.19.2
PROJECT_ROOT = github.com/hyperledger/aries-framework-go
GOBIN_PATH=$(abspath .)/build/bin
MOCKGEN=$(GOBIN_PATH)/mockgen
GOMOCKS=pkg/internal/gomocks

ifdef PROFILE_DEV
 ENVIRON_SUFIX='.dev'
else
 ENVIRON_SUFIX=
endif

.PHONY: all
all: clean checks unit-test unit-test-wasm unit-test-mobile bdd-test

.PHONY: checks
checks: license lint generate-openapi-spec

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh

.PHONY: unit-test
unit-test: mocks
	@scripts/check_unit.sh

.PHONY: unit-test-ursa
unit-test-ursa: mocks
	@scripts/check_unit_ursa.sh

.PHONY: benchmark
benchmark:
	@scripts/check_bench.sh

.PHONY: unit-test-wasm
unit-test-wasm: export GOBIN=$(GOBIN_PATH)
unit-test-wasm: depend
	@scripts/check_unit_wasm.sh

.PHONY: unit-test-mobile
unit-test-mobile:
	@echo "Running unit tests for mobile"
	@cd ${ARIES_AGENT_MOBILE_PATH} && $(MAKE) unit-test

.PHONY: bdd-test
bdd-test: clean generate-test-keys agent-rest-docker sample-webhook-docker sidetree-cli bdd-test-js bdd-test-go

.PHONY: bdd-test-go
bdd-test-go:
	@scripts/check_go_integration.sh

.PHONY: bdd-test-js
bdd-test-js:
	@scripts/check_js_integration.sh

bdd-test-js-local:
	@scripts/check_js_integration.sh test-local

bdd-test-js-debug:
	@scripts/check_js_integration.sh test-debug


.PHONY: vc-test-suite
vc-test-suite: clean
	@scripts/run_vc_test_suite.sh

.PHONY: bbs-interop-test
bbs-interop-test:
	@scripts/check_bbs_interop.sh

.PHONY: generate-test-keys # TODO needed?
generate-test-keys: clean
	@mkdir -p -p test/bdd/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/go/src/$(PROJECT_ROOT) \
		--entrypoint "/opt/go/src/$(PROJECT_ROOT)/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: generate-test-keys-no-build
generate-test-keys-no-build: clean-fixtures-only-and-deploy
	@mkdir -p -p test/bdd/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/go/src/$(PROJECT_ROOT) \
		--entrypoint "/opt/go/src/$(PROJECT_ROOT)/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: generate-dpabc-clib
generate-dpabc-clib: clean
	@docker build -f ./images/agent-rest/Dockerfile_base_image_with_compilation_tools \
				--build-arg ALPINE_VER=$(ALPINE_VER) \
				-t basecompilationcontainer .
	@docker run -v $(abspath .):/opt/go/src/$(PROJECT_ROOT) \
    		-v $(abspath $(PABC_PATH)):/opt/c/src/p-abc \
    		--entrypoint "/opt/go/src/$(PROJECT_ROOT)/scripts/generate_compiled_dpabc_clib.sh" \
    		-i --rm -t basecompilationcontainer --no-cache
	@docker image rm basecompilationcontainer

.PHONY: generate-openapi-spec
generate-openapi-spec: clean
	@echo "Generating and validating controller API specifications using Open API"
	@mkdir -p build/rest/openapi/spec
	@SPEC_META=$(ARIES_AGENT_REST_PATH) SPEC_LOC=${OPENAPI_SPEC_PATH}  \
	DOCKER_IMAGE=$(OPENAPI_DOCKER_IMG) DOCKER_IMAGE_VERSION=$(OPENAPI_DOCKER_IMG_VERSION)  \
	scripts/generate-openapi-spec.sh

.PHONY: generate-openapi-demo-specs
generate-openapi-demo-specs: clean generate-openapi-spec agent-rest-docker sample-webhook-docker
	@echo "Generate demo agent rest controller API specifications using Open API"
	@SPEC_PATH=${OPENAPI_SPEC_PATH} OPENAPI_DEMO_PATH=test/bdd/fixtures/demo/openapi \
    	DOCKER_IMAGE=$(OPENAPI_DOCKER_IMG) DOCKER_IMAGE_VERSION=$(OPENAPI_DOCKER_IMG_VERSION)  \
    	scripts/generate-openapi-demo-specs.sh

.PHONY: run-openapi-demo
run-openapi-demo: stop-openapi-demo generate-test-keys generate-dpabc-clib run-fabric generate-openapi-demo-specs
	@echo "Starting demo agent rest containers ..."
	@DEMO_COMPOSE_PATH=test/bdd/fixtures/demo/openapi SIDETREE_COMPOSE_PATH=test/bdd/fixtures/sidetree-mock AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        scripts/run-openapi-demo.sh

.PHONY: run-openapi-demo-no-build
run-openapi-demo-no-build: generate-test-keys-no-build
	@echo "Starting demo agent rest containers ..."
	@DEMO_COMPOSE_PATH=test/bdd/fixtures/demo/openapi SIDETREE_COMPOSE_PATH=test/bdd/fixtures/sidetree-mock AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        scripts/run-openapi-demo.sh

.PHONY: run-openapi-demo-build-no-clean
run-openapi-demo-build-no-clean: generate-test-keys-no-build agent-rest-docker
	@echo "Starting demo agent rest containers ..."
	@DEMO_COMPOSE_PATH=test/bdd/fixtures/demo/openapi SIDETREE_COMPOSE_PATH=test/bdd/fixtures/sidetree-mock AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        scripts/run-openapi-demo.sh

.PHONY: run-poc-demo # TODO needed?
run-poc-demo: generate-test-keys generate-dpabc-clib agent-rest-docker sample-webhook-docker
	@echo "Starting demo agent rest containers ..."
	@AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        scripts/run-poc-demo.sh 

.PHONY: run-poc-demo-no-build  # TODO needed?
run-poc-demo-no-build: generate-test-keys-no-build
	@echo "Starting demo agent rest containers ..."
	@AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        scripts/run-poc-demo.sh

.PHONY: stop-poc-demo # TODO needed?
stop-poc-demo:
	@echo "Stopping demo rest containers ..."
	@ AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        DEMO_COMPOSE_OP=down scripts/run-poc-demo.sh


.PHONY: stop-openapi-demo
stop-openapi-demo:
	@echo "Stopping demo agent rest containers ..."
	@DEMO_COMPOSE_PATH=test/bdd/fixtures/demo/openapi SIDETREE_COMPOSE_PATH=test/bdd/fixtures/sidetree-mock AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest  \
        DEMO_COMPOSE_OP=down scripts/run-openapi-demo.sh

.PHONY: agent-rest
agent-rest:
	@echo "Building aries-agent-rest"
	@echo "GO_FLAGS: ${GO_FLAGS}"
	@mkdir -p ./build/bin
	@cd ${ARIES_AGENT_REST_PATH} && go mod tidy &&  go build $(GO_FLAGS) -o ../../build/bin/aries-agent-rest main.go

.PHONY: agent-rest-debug
agent-rest-debug: GO_FLAGS = -gcflags="all=-N -l"
agent-rest-debug: agent-rest


.PHONY: agent-rest-acapy-interop
agent-rest-acapy-interop:
	@echo "Building aries-agent-rest for aca-py interop"
	@mkdir -p ./build/bin
	@cd ${ARIES_AGENT_REST_PATH} && go build -o ../../build/bin/aries-agent-rest -tags ACAPyInterop main.go

.PHONY: agent-mobile
agent-mobile:
	@echo "Building aries-agent-mobile"
	@cd ${ARIES_AGENT_MOBILE_PATH} && $(MAKE) all

.PHONY: sidetree-cli
sidetree-cli:
	@echo "Building sidetree-cli"
	@mkdir -p ./build/bin
	@cd ${SIDETREE_CLI_PATH} && go build -o ../../../../build/bin/sidetree main.go

.PHONY: agent-rest-docker
agent-rest-docker:
	@echo "Building aries agent rest docker image"
	@docker build -f ./images/agent-rest/Dockerfile$(ENVIRON_SUFIX) --no-cache -t $(DOCKER_OUTPUT_NS)/$(AGENT_REST_IMAGE_NAME):latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GO_TAGS=$(GO_TAGS) \
	--build-arg GOPROXY=$(GOPROXY) .

.PHONY: sample-webhook
sample-webhook:
	@echo "Building sample webhook server"
	@mkdir -p ./build/bin
	@go build -o ./build/bin/webhook-server test/bdd/webhook/main.go

.PHONY: sample-webhook-docker
sample-webhook-docker:
	@echo "Building sample webhook server docker image"
	@docker build -f ./images/mocks/webhook/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(WEBHOOK_IMAGE_NAME):latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GO_TAGS=$(GO_TAGS) \
	--build-arg GOPROXY=$(GOPROXY) .

comma:= ,
semicolon:= ;
mocks_dir =

define create_mock
  $(eval mocks_dir := $(subst pkg,$(GOMOCKS),$(1)))
  @echo Creating $(mocks_dir)
  @mkdir -p $(mocks_dir) && rm -rf $(mocks_dir)/*
  @$(MOCKGEN) -destination $(mocks_dir)/mocks.gen.go -self_package mocks -package mocks $(PROJECT_ROOT)/$(1) $(subst $(semicolon),$(comma),$(2))
endef

define create_spi_provider_mocks
  $(eval mocks_dir := $(GOMOCKS)/spi/storage)
  @echo Creating $(mocks_dir)
  @mkdir -p $(mocks_dir) && rm -rf $(mocks_dir)/*
  @$(MOCKGEN) -destination $(mocks_dir)/mocks.gen.go -self_package mocks -package mocks $(PROJECT_ROOT)/$(1) $(subst $(semicolon),$(comma),$(2))
endef

depend:
	@mkdir -p ./build/bin
	GOBIN=$(GOBIN_PATH) go install github.com/golang/mock/mockgen@v1.5.0
	GOBIN=$(GOBIN_PATH) go install github.com/agnivade/wasmbrowsertest@v0.3.5

.PHONY: mocks
mocks: depend clean-mocks
	$(call create_mock,pkg/framework/aries/api/vdr,Registry)
	$(call create_mock,pkg/kms,Provider;KeyManager)
	$(call create_mock,pkg/didcomm/protocol/issuecredential,Provider)
	$(call create_mock,pkg/didcomm/protocol/middleware/issuecredential,Provider;Metadata)
	$(call create_mock,pkg/didcomm/protocol/middleware/presentproof,Provider;Metadata)
	$(call create_mock,pkg/client/outofband,Provider;OobService)
	$(call create_mock,pkg/client/outofbandv2,Provider;OobService)
	$(call create_mock,pkg/didcomm/protocol/presentproof,Provider)
	$(call create_mock,pkg/client/introduce,Provider;ProtocolService)
	$(call create_mock,pkg/client/issuecredential,Provider;ProtocolService)
	$(call create_mock,pkg/client/presentproof,Provider;ProtocolService)
	$(call create_mock,pkg/didcomm/protocol/introduce,Provider)
	$(call create_mock,pkg/didcomm/common/service,DIDComm;Event;Messenger;MessengerHandler)
	$(call create_mock,pkg/didcomm/dispatcher,Outbound)
	$(call create_mock,pkg/didcomm/messenger,Provider)
	$(call create_mock,pkg/store/verifiable,Store)
	$(call create_mock,pkg/store/did,ConnectionStore)
	$(call create_mock,pkg/controller/command/presentproof,Provider)
	$(call create_mock,pkg/controller/command/issuecredential,Provider)
	$(call create_mock,pkg/controller/webnotifier,Notifier)
	$(call create_spi_provider_mocks,spi/storage,Provider;Store)

.PHONY: clean-mocks
clean-mocks:
	@if [ -d $(GOMOCKS) ]; then rm -r $(GOMOCKS); echo "Folder $(GOMOCKS) was removed!"; fi

.PHONY: clean
clean: stop-fabric clean-fixtures clean-build clean-images

.PHONY: clean-images
clean-images: clean-fixtures
clean-images: IMAGES=$(shell docker image ls | grep aries-framework-go | awk '{print $$3}')
clean-images:
	@if [ ! -z "$(IMAGES)" ]; then \
		echo "Cleaning aries-framework-go docker images ..."; \
		docker rmi -f $(IMAGES); \
	fi;

.PHONY: clean-build
clean-build:
	@rm -f coverage.out
	@rm -Rf ./build
	@rm -Rf $(ARIES_AGENT_MOBILE_PATH)/build
	@rm -Rf ./test/bdd/db
	@rm -Rf ./test/bdd/*.log
	# Remove the line that deletes the fabric folder

.PHONY: clean-fixtures
clean-fixtures:
	@rm -Rf ./test/bdd/fixtures/keys/tls
	@rm -Rf ./test/bdd/fixtures/demo/openapi/specs
	@cd test/bdd/fixtures/sidetree-mock && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/demo/openapi && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/agent-rest && docker-compose down 2> /dev/null
	
	

.PHONY: clean-fixtures-no-build
clean-fixtures-no-build:
	@rm -Rf ./test/bdd/fixtures/keys/tls
	@cd test/bdd/fixtures/demo/openapi && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/sidetree-mock && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/agent-rest && docker-compose down 2> /dev/null

.PHONY: clean-fixtures-only-and-deploy
	@rm -Rf ./test/bdd/fixtures/keys/tls
	@cd test/bdd/fixtures/demo/openapi && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/agent-rest && docker-compose down 2> /dev/null
	@cd test/bdd/fixtures/sidetree-mock && docker-compose down 2> /dev/null

.PHONY: build-fabric
build-fabric:
	@echo "Downloading necessary resources."
	@chmod +x -R scripts/
	@FABRIC_PATH="$(FABRIC_PATH)" \
	FABRIC_VERSION="$(FABRIC_VERSION)" \
	CONNECTION_PROFILE_PATH="$(CONNECTION_PROFILE_PATH)" \
	./scripts/fabric/build_fabric.sh

.PHONY: stop-fabric
stop-fabric:
	@cd ${FABRIC_PATH}
	@echo "Stopping Fabric deployment"
	FABRIC_PATH="$(FABRIC_PATH)" \
	FABRIC_VERSION="$(FABRIC_VERSION)" \
	CONNECTION_PROFILE_PATH="$(CONNECTION_PROFILE_PATH)" \
	./scripts/fabric/stop_fabric.sh

.PHONY: run-fabric
run-fabric:
	@echo "Launching Fabric deployment."
	FABRIC_PATH="$(FABRIC_PATH)" \
	FABRIC_VERSION="$(FABRIC_VERSION)" \
	CONNECTION_PROFILE_PATH="$(CONNECTION_PROFILE_PATH)" \
	./scripts/fabric/run_fabric.sh

.PHONY: restart-fabric
restart-fabric:
	@cd $(FABRIC_PATH)
	@echo "Restarting Fabric deployment"
	@FABRIC_PATH="$(FABRIC_PATH)"; \
	FABRIC_VERSION="$(FABRIC_VERSION)"; \
	CONNECTION_PROFILE_PATH="$(CONNECTION_PROFILE_PATH)"; \
	./scripts/fabric/restart_fabric.sh
