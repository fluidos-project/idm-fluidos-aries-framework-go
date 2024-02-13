/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms12381g1pub_test

import (
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPsmsWrapperIndividualFlow(t *testing.T) {
	err1, err2, err3, err4, pbytes, err5 := psms12381g1pub.TestPsmsWrapperIndividualFlowCgoWrapper()
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	require.NoError(t, err4)
	require.NotEmpty(t, pbytes)
	require.NoError(t, err5)
}

func TestPsmsWrapperDistributedFlow(t *testing.T) {
	err1, err2, err3, err4, err5, err6, pbytes, err7 := psms12381g1pub.TestPsmsWrapperDistributedFlowCgoWrapper()
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	require.NoError(t, err4)
	require.NoError(t, err5)
	require.NoError(t, err6)
	require.NotEmpty(t, pbytes)
	require.NoError(t, err7)
}
