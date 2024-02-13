/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms

/*
import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	commonpb "github.com/google/tink/go/proto/common_go_proto"
	tinkpb "github.com/google/tink/go/proto/tink_go_proto"
	"github.com/stretchr/testify/require"

	psmspb "github.com/hyperledger/aries-framework-go/pkg/crypto/tinkcrypto/primitive/proto/psms_go_proto"
)

func TestPSMSignerKeyManager_Primitive(t *testing.T) {
	km := newPSMSSignerKeyManager()

	t.Run("Test signer key manager Primitive() with empty serialized key", func(t *testing.T) {
		p, err := km.Primitive([]byte(""))
		require.EqualError(t, err, errInvalidPSMSSignerKey.Error(),
			"psmsSignerKeyManager primitive from empty serialized key must fail")
		require.Empty(t, p)
	})

	t.Run("Test signer key manager Primitive() with bad serialize key", func(t *testing.T) {
		p, err := km.Primitive([]byte("bad.data"))
		require.Contains(t, err.Error(), errInvalidPSMSSignerKey.Error())
		require.Contains(t, err.Error(), "invalid proto: proto:")
		require.Contains(t, err.Error(), "cannot parse invalid wire-format data")
		require.Empty(t, p)
	})

	flagTests := []struct {
		tcName     string
		version    uint32
		hashType   commonpb.HashType
		curveType  psmspb.PSMSCurveType
		groupField psmspb.GroupField
	}{
		{
			tcName:     "signer key manager Primitive() success",
			version:    0,
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager Primitive() using key with bad version",
			version:    9999,
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager Primitive() using key with bad hash type",
			version:    0,
			hashType:   commonpb.HashType_UNKNOWN_HASH,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager Primitive() using key with bad curve",
			version:    0,
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_UNKNOWN_PSMS_CURVE_TYPE,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager Primitive() using key with bad group",
			version:    0,
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_UNKNOWN_GROUP_FIELD,
		},
	}

	for _, tc := range flagTests {
		tt := tc
		t.Run("Test "+tt.tcName, func(t *testing.T) {
			v := tt.version

			privKeyProto := &psmspb.PSMSPrivateKey{
				Version: v,
				PublicKey: &psmspb.PSMSPublicKey{
					Version: v,
					Params: &psmspb.PSMSParams{
						HashType: tt.hashType,
						Curve:    tt.curveType,
						Group:    tt.groupField,
					},
				},
			}

			sPrivKey, err := proto.Marshal(privKeyProto)
			require.NoError(t, err)

			p, err := km.Primitive(sPrivKey)
			if strings.Contains(tt.tcName, "success") {
				require.NoError(t, err)
				require.NotEmpty(t, p)
				return
			}

			require.Errorf(t, err, tt.tcName)
			require.Empty(t, p)
		})
	}
}

func TestPSMSSignerKeyManager_DoesSupport(t *testing.T) {
	km := newPSMSSignerKeyManager()
	require.False(t, km.DoesSupport("bad/url"))
	require.True(t, km.DoesSupport(psmsSignerKeyTypeURL))
}

func TestPSMSSignerKeyManager_NewKey(t *testing.T) {
	km := newPSMSSignerKeyManager()

	t.Run("Test signer key manager NewKey() with nil key", func(t *testing.T) {
		k, err := km.NewKey(nil)
		require.EqualError(t, err, errInvalidPSMSSignerKeyFormat.Error())
		require.Empty(t, k)
	})

	t.Run("Test signer key manager NewKey() with bad serialize key", func(t *testing.T) {
		p, err := km.NewKey([]byte("bad.data"))
		require.Contains(t, err.Error(), errInvalidPSMSSignerKey.Error())
		require.Contains(t, err.Error(), "invalid proto: proto:")
		require.Contains(t, err.Error(), "cannot parse invalid wire-format data")
		require.Empty(t, p)
	})

	flagTests := []struct {
		tcName     string
		hashType   commonpb.HashType
		curveType  psmspb.PSMSCurveType
		groupField psmspb.GroupField
	}{
		{
			tcName:     "success signer key manager NewKey() and NewKeyData()",
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager NewKey() and NewKeyData() using key with bad hash",
			hashType:   commonpb.HashType_UNKNOWN_HASH,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager NewKey() and NewKeyData() using key with bad curve",
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_UNKNOWN_PSMS_CURVE_TYPE,
			groupField: psmspb.GroupField_G2,
		},
		{
			tcName:     "signer key manager NewKey() and NewKeyData() using key with bad group",
			hashType:   commonpb.HashType_SHA256,
			curveType:  psmspb.PSMSCurveType_BLS12_381,
			groupField: psmspb.GroupField_UNKNOWN_GROUP_FIELD,
		},
	}

	for _, tc := range flagTests {
		tt := tc
		t.Run("Test "+tt.tcName, func(t *testing.T) {
			privKeyProto := &psmspb.PSMSKeyFormat{
				Params: &psmspb.PSMSParams{
					HashType: tt.hashType,
					Curve:    tt.curveType,
					Group:    tt.groupField,
				},
			}

			sPrivKey, err := proto.Marshal(privKeyProto)
			require.NoError(t, err)

			p, err := km.NewKey(sPrivKey)
			if strings.Contains(tt.tcName, "success") {
				require.NoError(t, err)
				require.NotEmpty(t, p)

				sp, e := proto.Marshal(p)
				require.NoError(t, e)
				require.NotEmpty(t, sp)

				// try PublicKeyData() with bad serialized private key
				pubK, e := km.PublicKeyData([]byte("bad serialized private key"))
				require.Error(t, e)
				require.Empty(t, pubK)

				// try PublicKeyData() with valid serialized private key
				pubK, e = km.PublicKeyData(sp)
				require.NoError(t, e)
				require.NotEmpty(t, pubK)
			}

			kd, err := km.NewKeyData(sPrivKey)
			if strings.Contains(tt.tcName, "success") {
				require.NoError(t, err)
				require.NotEmpty(t, kd)
				require.Equal(t, kd.TypeUrl, psmsSignerKeyTypeURL)
				require.Equal(t, kd.KeyMaterialType, tinkpb.KeyData_ASYMMETRIC_PRIVATE)
				return
			}

			require.Errorf(t, err, tt.tcName)
			require.Empty(t, p)
		})
	}
}


*/
