#!/bin/bash

export PROJECT_ROOT=github.com/hyperledger/aries-framework-go

# Stop demo agent rest containers
echo "Stopping demo agent rest containers ..."
DEMO_COMPOSE_PATH=test/bdd/fixtures/demo/openapi SIDETREE_COMPOSE_PATH=test/bdd/fixtures/sidetree-mock AGENT_REST_COMPOSE_PATH=test/bdd/fixtures/agent-rest

# Generate test keys
mkdir -p -p test/bdd/fixtures/keys/tls
docker run -i --rm \
  -v $(pwd):/opt/go/src/$(PROJECT_ROOT) \
  --entrypoint "/opt/go/src/$(PROJECT_ROOT)/scripts/generate_test_keys.sh" \
  frapsoft/openssl

# Generate dpabc clib
docker build -f ./images/agent-rest/Dockerfile_base_image_with_compilation_tools \
  --build-arg ALPINE_VER=$(ALPINE_VER) \
  -t basecompilationcontainer .
docker run -v $(pwd):/opt/go/src/$(PROJECT_ROOT)
docker image rm basecompilationcontainer

# Run Fabric
echo "Launching Fabric deployment."
FABRIC_PATH="$(FABRIC_PATH)" \
FABRIC_VERSION="$(FABRIC_VERSION)" \
CONNECTION_PROFILE_PATH="$(CONNECTION_PROFILE_PATH)" \
./scripts/fabric/run_fabric.sh

# Generate OpenAPI demo specs
echo "Generate demo agent rest controller API specifications using Open API"
SPEC_PATH=${OPENAPI_SPEC_PATH} OPENAPI_DEMO_PATH=test/bdd/fixtures/demo/openapi

# Run OpenAPI demo script
./scripts/run-openapi-demo.sh