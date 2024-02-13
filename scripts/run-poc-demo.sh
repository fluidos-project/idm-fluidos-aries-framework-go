#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

DEMO_COMPOSE_OP="${DEMO_COMPOSE_OP:-up --force-recreate -d}"
AGENT_PATH="${AGENT_REST_COMPOSE_PATH}"
AGENT_COMPOSE_FILE="$PWD/$AGENT_PATH"

set -o allexport
[[ -f $AGENT_PATH/.env ]] && source $AGENT_PATH/.env
set +o allexport

cd $AGENT_COMPOSE_FILE
docker-compose -f docker-compose.yml  ${DEMO_COMPOSE_OP}
