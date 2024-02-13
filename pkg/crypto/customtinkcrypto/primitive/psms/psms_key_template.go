/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
//XXX Revised

package psms

import (
	"encoding/binary"
	"fmt"
	"strconv"

	tinkpb "github.com/google/tink/go/proto/tink_go_proto"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

// BLS12381G1KeyTemplate creates a Tink key template for PSMS on BLS12-381 curve with G1 group.
func BLS12381G1KeyTemplate(opts ...kms.KeyOpts) (*tinkpb.KeyTemplate, error) {
	keyOpts := kms.NewKeyOpt()
	for _, opt := range opts {
		opt(keyOpts)
	}
	att := 3
	if keyOpts.Attrs() != nil && len(keyOpts.Attrs()) == 1 {
		transAtt, err := strconv.Atoi(keyOpts.Attrs()[0])
		if err != nil {
			return nil, fmt.Errorf("bls12381g1 key template error: opt with number of attributes is mandatory: %w", err)
		}
		att = transAtt
	} else {
		//return nil, fmt.Errorf("bls12381g1 key template error: opt with number of attributes is mandatory") TODO UMU (crypto) Maybe require opts as mandatory, "issue" with Rotate and similar.
	}
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(att))
	return &tinkpb.KeyTemplate{
		TypeUrl:          PsmsSignerKeyTypeURL,
		Value:            bytes, //XXX UMU (serial) Proto serialization like in CL instead. Also
		OutputPrefixType: tinkpb.OutputPrefixType_RAW,
	}, nil
	//return createKeyTemplate(psmspb.PSMSCurveType_BLS12_381, psmspb.GroupField_G2, commonpb.HashType_SHA256)
}

/*
// createKeyTemplate for PSMS+ keys.
func createKeyTemplate(curve psmspb.PSMSCurveType, group psmspb.GroupField, hash commonpb.HashType) *tinkpb.KeyTemplate {
	format := &psmspb.PSMSKeyFormat{
		Params: &psmspb.PSMSParams{
			HashType: hash,
			Curve:    curve,
			Group:    group,
		},
	}

	serializedFormat, err := proto.Marshal(format)
	if err != nil {
		panic("failed to marshal PSMSKeyFormat proto")
	}

	return &tinkpb.KeyTemplate{
		TypeUrl:          psmsSignerKeyTypeURL,
		Value:            serializedFormat,
		OutputPrefixType: tinkpb.OutputPrefixType_RAW,
	}
}*/
