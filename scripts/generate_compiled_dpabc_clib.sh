#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Generating dpabc C library"
cd /opt/c/src/p-abc

sh setupscript.sh -b -d /opt/go/src/github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub/clib

echo "done generating dpabc C library"
