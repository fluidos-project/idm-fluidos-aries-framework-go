/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
)

var logger = log.New("aries-framework/vdr/fabric")

const (
	SC_METHOD_RESOLVE_DID = "get"
	SC_METHOD_CREATE_DID  = "set"
	SC_NAME               = "sacc"
	CHANNEL_NAME          = "mychannel"
	WALLET_ID             = "appUser"
	// StoreNamespace store name space for DID Store.
	StoreNamespace = "peer"
	// DefaultServiceType default service type.
	DefaultServiceType = "defaultServiceType"
	// DefaultServiceEndpoint default service endpoint.
	DefaultServiceEndpoint = "defaultServiceEndpoint"
	didMethod              = "fabric"
)

// VDR via HTTP(s) endpoint.
type VDR struct {
	endpointURL      string
	client           *http.Client
	accept           Accept
	resolveAuthToken string
	network          *gateway.Network
	contract         *gateway.Contract
	config           []byte
	wallet           *gateway.Wallet
	gw               *gateway.Gateway
}

// Accept is method to accept did method.
type Accept func(method string) bool

// New creates new DID Resolver.
func New(configURL string, opts ...Option) (*VDR, error) {
	v := &VDR{client: &http.Client{}, accept: func(method string) bool { return true }} // TODO change?

	req, err := http.NewRequest(http.MethodGet, configURL, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP create get request failed: %w", err)
	}

	resp, err := v.client.Do(req) // config.json
	if err != nil {
		return nil, fmt.Errorf("httpClient do: %w", err)
	}

	defer func() {
		e := resp.Body.Close()
		if e != nil {
			logger.Errorf("Failed to close response body: %s", e.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	var gotBody []byte

	gotBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to readAll bytes from resp.body: %w", err)
	}
	v.config = gotBody
	v.accept = func(method string) bool {
		return method == didMethod
	}

	return v, nil
}

// Accept did method - attempt to resolve any method.
func (v *VDR) Accept(method string, opts ...vdrapi.DIDMethodOption) bool {
	return method == DIDMethod
}

// Create did doc.
//func (v *VDR) Create(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
//return fmt.Errorf("not implemented in fabric VDR")
//}

// Close frees resources being maintained by vdr.
func (v *VDR) Close() error {
	v.gw.Close()
	return nil
}

// Update did doc.
func (v *VDR) Update(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) error {
	return fmt.Errorf("not implemented in fabric VDR")
}

// Deactivate did doc.
func (v *VDR) Deactivate(didID string, opts ...vdrapi.DIDMethodOption) error {
	return fmt.Errorf("not implemented in fabric VDR")
}

// Option configures the peer vdr.
type Option func(opts *VDR)

// WithTimeout option is for definition of HTTP(s) timeout value of DID Resolver.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *VDR) {
		opts.client.Timeout = timeout
	}
}

// WithHTTPClient option is for custom http client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(opts *VDR) {
		opts.client = httpClient
	}
}

// WithAccept option is for accept did method.
func WithAccept(accept Accept) Option {
	return func(opts *VDR) {
		opts.accept = accept
	}
}

// WithResolveAuthToken add auth token for resolve.
func WithResolveAuthToken(authToken string) Option {
	return func(opts *VDR) {
		opts.resolveAuthToken = "Bearer " + authToken
	}
}

func (v *VDR) connectGateway() (bool, error) {
	var dat map[string]interface{} // I converted my 2 elemts json onto map of 2 strings <-> 2 keys
	if err := json.Unmarshal(v.config, &dat); err != nil {
		panic(err)
	}
	// dat has map with keys: 1 - "identity" = 2 string elements ["pub", "priv"] , 2 - connection-profile = string

	priv := dat["identity"].(map[string]interface{})["priv"].(string)
	pub := dat["identity"].(map[string]interface{})["pub"].(string)

	//fmt.Println("priv: " + priv)
	//fmt.Println("pub: " + pub)

	bytes, err := json.Marshal(dat["connection-profile"])

	//fmt.Println("config: " + string(bytes))

	if err != nil {
		return false, fmt.Errorf("failed to convert connection-profile to string(Marshal) %w", err)
	}
	v.wallet = gateway.NewInMemoryWallet()
	if v.wallet == nil {
		return false, fmt.Errorf("failed to create wallet in memory")
	}
	v.wallet.Put("appUser",
		gateway.NewX509Identity(
			"Org1MSP",
			pub,
			priv))

	v.gw, err = gateway.Connect(gateway.WithConfig(config.FromRaw(bytes,
		"json")),
		gateway.WithIdentity(v.wallet, WALLET_ID))

	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		return false, err
		//os.Exit(1)
	}
	//defer gw.Close()
	v.network, err = v.gw.GetNetwork(CHANNEL_NAME)
	if err != nil {
		fmt.Printf("Failed to get network("+CHANNEL_NAME+"): %s\n", err)
		//os.Exit(1)
		return false, err
	}

	v.contract = v.network.GetContract(SC_NAME)
	if v.contract == nil {
		fmt.Printf("Failed to get smartContract: %s\n", SC_NAME)
	}
	return true, nil
}
func closeResponseBody(respBody io.Closer) {
	e := respBody.Close()
	if e != nil {
		logger.Errorf("Failed to close response body: %v", e)
	}
}
