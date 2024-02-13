package poc

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/internal/logutil"
	"net/http"
	"strconv"
	"strings"
)



func (o *Command) ValidateProof(proof IdProof) bool {
	for idx,validator:= range o.idProofValidators{
		logutil.LogInfo(logger, CommandName, "PufIdProof", fmt.Sprintf("index: %d", idx))
		if validator.Accept(proof.AttrName) {
			return validator.Validate(proof)
		}
	}
	return false
}

//IdProof validators
type IdProofValidator interface {
	Accept(string) bool
	Validate(IdProof) bool
}

//Default validator, always true
type DefaultValidator struct {}

func (v *DefaultValidator) Accept(attrName string) bool {
	return true
}

func (v *DefaultValidator) Validate(proof IdProof) bool {
	return true
}


//Puf Validator
type PufProofValidator struct {
	ValidationUrl string
	UseService bool
}

type PufProof struct {
	Key string `json:"key,omitempty"`
}

func (v *PufProofValidator) Accept(attrName string) bool {
	return strings.EqualFold("pufId",attrName)
}

func (v *PufProofValidator) Validate(proof IdProof) bool {
	if !v.UseService{
		return true
	}
	//TODO UMU Parse proof
	var pp PufProof
	err := json.Unmarshal(proof.ProofData,&pp)
	if err != nil || pp.Key=="" {
		logutil.LogInfo(logger, CommandName, "PufIdProof", fmt.Sprintf("error parsing idproof: %s", err))
		return false
	}

	req, err := http.NewRequest("GET", v.ValidationUrl, nil)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "PufIdProof", fmt.Sprintf("error preparing puf validation request: %s", err))
		return false
	}
	attrStr, ok := proof.AttrValue.(string)
	if !ok {
		logutil.LogInfo(logger, CommandName, "PufIdProof", fmt.Sprintf("error attribute value type: %s", err))
		return false
	}
	q := req.URL.Query()
	q.Add("id", attrStr) 
	q.Add("tk", pp.Key)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "PufIdProof", fmt.Sprintf("error executing puf validation request: %s", err))
		return false
	}
	defer resp.Body.Close()
	logutil.LogInfo(logger, CommandName, "PufIdProof", "HTTP Response Status:"+strconv.Itoa(resp.StatusCode))
	return resp.StatusCode == 200
}

