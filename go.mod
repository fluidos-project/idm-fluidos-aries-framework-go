// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/hyperledger/aries-framework-go

// TODO (#2815): Remove circular dependency between the main module and component/storage/edv

require (
	github.com/PaesslerAG/gval v1.1.0
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/VictoriaMetrics/fastcache v1.5.7
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cenkalti/backoff/v4 v4.1.2
	github.com/go-jose/go-jose/v3 v3.0.1-0.20221117193127-916db76e8214
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/tink/go v1.7.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.3
	github.com/hyperledger/aries-framework-go/component/storage/edv v0.0.0-20221025204933-b807371b6f1e
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20220322085443-50e8f9bd208b
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20221025204933-b807371b6f1e
	github.com/hyperledger/fabric-sdk-go v1.0.1-0.20210603143513-14047c6d88f0
	github.com/hyperledger/ursa-wrapper-go v0.3.1
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kawamuray/jsonpath v0.0.0-20201211160320-7483bafabd7e
	github.com/kilic/bls12-381 v0.1.1-0.20210503002446-7b7597926c69
	github.com/mitchellh/mapstructure v1.5.0
	github.com/multiformats/go-multibase v0.1.1
	github.com/multiformats/go-multihash v0.0.13
	github.com/piprate/json-gold v0.4.2
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.9.0
	github.com/teserakt-io/golang-ed25519 v0.0.0-20210104091850-3888c087a4c8
	github.com/tidwall/gjson v1.6.7
	github.com/tidwall/sjson v1.1.4
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.1.0
	golang.org/x/sys v0.1.0
	google.golang.org/protobuf v1.28.1
	nhooyr.io/websocket v1.8.3
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-kivik/couchdb/v3 v3.2.6 // indirect
	github.com/go-kivik/kivik/v3 v3.2.3 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.8.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.0.6 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.7.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0 // indirect
	github.com/jackc/puddle v1.1.3 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/valyala/fastjson v1.6.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.8.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cloudflare/cfssl v1.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hyperledger/aries-framework-go-ext/component/storage/couchdb v0.0.0-20230523133839-1441aced7b6f
	github.com/hyperledger/aries-framework-go-ext/component/storage/mongodb v0.0.0-20230523133839-1441aced7b6f
	github.com/hyperledger/aries-framework-go-ext/component/storage/mysql v0.0.0-20230523133839-1441aced7b6f
	github.com/hyperledger/aries-framework-go-ext/component/storage/postgresql v0.0.0-20230523133839-1441aced7b6f
	github.com/hyperledger/aries-framework-go/component/storage/leveldb v0.0.0-20230605134811-1970e7e22384
	github.com/hyperledger/fabric-config v0.0.5 // indirect
	github.com/hyperledger/fabric-lib-go v1.0.0 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-varint v0.0.5 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.3.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.1.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tidwall/match v1.0.3 // indirect
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/weppos/publicsuffix-go v0.5.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/zmap/zcrypto v0.0.0-20190729165852-9051775e6a2e // indirect
	github.com/zmap/zlint v0.0.0-20190806154020-fd021b4cfbeb // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20220218161850-94dd64e39d7c // indirect
	google.golang.org/grpc v1.44.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.19

//replace github.com/square/go-jose/v3 => github.com/go-jose/go-jose/v3 v3.0.1-0.20221117193127-916db76e8214
//
//replace github.com/square/go-jose/v3/json => github.com/go-jose/go-jose/v3/json v1.0.1-0.20221117193127-916db76e8214
//
//replace github.com/square/go-jose/v3/jwt => github.com/go-jose/go-jose/v3/jwt v1.0.1-0.20221117193127-916db76e8214
//
//replace github.com/square/go-jose/v3/cipher => github.com/go-jose/go-jose/v3/cipher v1.0.1-0.20221117193127-916db76e8214
