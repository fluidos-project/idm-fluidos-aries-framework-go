/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package poc

import (
	"fmt"
	"net/http"

	poc"github.com/hyperledger/aries-framework-go/pkg/controller/command/poc"
	vcwalletc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vcwallet"
	vdrc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/controller/internal/cmdutil"
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
)

// constants for the VDR operations.
const (
	FluidosOperationID            = "/fluidos/idm"
	NewDIDPath                = FluidosOperationID + "/generateDID"
	DoDeviceEnrolmentPath     = FluidosOperationID + "/doEnrolment"
	GenerateVpPath            = FluidosOperationID + "/generateVP"
	AcceptDeviceEnrolmentPath = FluidosOperationID + "/acceptEnrolment"
	VerifyCredentialPath      = FluidosOperationID + "/verifyCredential"
	TestingCallPath     	  = FluidosOperationID + "/testingCall"
	GetTrustedIssuerListPath  = FluidosOperationID + "/trustedIssuers"
	GetVCredentialPath        = FluidosOperationID + "/getVCredential"
	SignJWTContentPath        = FluidosOperationID + "/signJWTContent"
	VerifyJWTContentPath      = FluidosOperationID + "/verifyJWTContent"
)



// Operation contains basic common operations provided by controller REST API.
type Operation struct {
	handlers []rest.Handler
	command  *poc.Command
}

// New returns new common operations rest client instance.
func New(vdrcommand *vdrc.Command, vcwalletcommand *vcwalletc.Command) (*Operation, error) {
	cmd, err := poc.New(vdrcommand, vcwalletcommand)
	if err != nil {
		return nil, fmt.Errorf("new vdr : %w", err)
	}

	o := &Operation{command: cmd}
	o.registerHandler()

	return o, nil
}



// GetRESTHandlers get all controller API handler available for this service.
func (o *Operation) GetRESTHandlers() []rest.Handler {
	return o.handlers
}

// registerHandler register handlers to be exposed from this protocol service as REST API endpoints.
func (o *Operation) registerHandler() {
	// Add more protocol endpoints here to expose them as controller API endpoints
	o.handlers = []rest.Handler{
		cmdutil.NewHTTPHandler(NewDIDPath, http.MethodPost, o.NewDID),
		cmdutil.NewHTTPHandler(DoDeviceEnrolmentPath, http.MethodPost, o.DoDeviceEnrolment),
		cmdutil.NewHTTPHandler(GenerateVpPath, http.MethodPost, o.GenerateVp),
		cmdutil.NewHTTPHandler(AcceptDeviceEnrolmentPath, http.MethodPost, o.AcceptDeviceEnrolment),
		cmdutil.NewHTTPHandler(VerifyCredentialPath, http.MethodPost, o.VerifyCredential),
		cmdutil.NewHTTPHandler(TestingCallPath, http.MethodPost, o.TestingCall),
		cmdutil.NewHTTPHandler(GetTrustedIssuerListPath, http.MethodGet, o.GetTrustedIssuerList),
		cmdutil.NewHTTPHandler(GetVCredentialPath, http.MethodPost, o.GetVCredential),
		cmdutil.NewHTTPHandler(SignJWTContentPath, http.MethodPost, o.SignJWTContent),
		cmdutil.NewHTTPHandler(VerifyJWTContentPath, http.MethodPost, o.VerifyJWTContent),
	}
}

// NewDID swagger:route POST /fluidos/idm/generateDID poc newDIDReq
//
// Create DID with keys/purposes as specified in request
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) NewDID(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.NewDID, rw, req.Body)
}

// DoEnrolment swagger:route POST /fluidos/idm/doEnrolment poc DoDeviceEnrolmentReq
//
// Do an enrolment process against the issuer, obtaining a new credential
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) DoDeviceEnrolment(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.DoDeviceEnrolment, rw, req.Body)
}

// GenerateVp swagger:route POST /fluidos/idm/generateVp poc GenerateVpReq
//
// Generate a VPresentation (for now VCredential?) for an authorization process
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) GenerateVp(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.GenerateVP, rw, req.Body)
}

// AcceptDeviceEnrolment swagger:route POST /fluidos/idm/acceptEnrolment poc AcceptDeviceEnrolmentReq
//
// Accepts enrolment requests, and if successful generates a Verifiable Credential for the enrolled device
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) AcceptDeviceEnrolment(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.AcceptEnrolment, rw, req.Body)
}

// VerifyCredential swagger:route POST /fluidos/idm/VerifyCredential poc VerifyCredentialReq
//
// Verify a Verifiable Credential, returns boolean of the verification result
//
// Responses:
//   default: genericError
// 		 200: documentRes
func (o *Operation) VerifyCredential(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.VerifyCredential, rw, req.Body)
}




func (o *Operation) TestingCall(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.TestingCall, rw, req.Body)
}

// GetTrustedIssuerList swagger:route GET /fluidos/idm/trustedIssuers poc GetTrustedIssuerListReq
//
// Get the list of trusted issuers
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) GetTrustedIssuerList(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.GetTrustedIssuerList, rw, req.Body)
}


// GetVCredential swagger:route POST /fluidos/idm/getVCredential poc GetVCredentialReq
//
// Get a Verifiable Credential
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) GetVCredential(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.GetVCredential, rw, req.Body)
}

// SignJWTContent swagger:route POST /fluidos/idm/signJWTContent poc SignJWTContentReq
//
// Sign the content of a JWT
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) SignJWTContent(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.SignJWTContent, rw, req.Body)
}

// VerifyJWTContent swagger:route POST /fluidos/idm/verifyJWTContent poc VerifyJWTContentReq
//
// Verify the content of a JWT
//
// Responses:
//    default: genericError
//        200: documentRes
func (o *Operation) VerifyJWTContent(rw http.ResponseWriter, req *http.Request) {
	rest.Execute(o.command.VerifyJWTContent, rw, req.Body)
}


