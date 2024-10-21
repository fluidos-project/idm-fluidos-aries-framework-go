package web

import (
	"encoding/json"
	application_gateway "fabricrest-go/application-gateway"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type XACMLEntity struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp TypeValue `json:"timestamp"`
	Version   TypeValue `json:"version"`
	Xacml     TypeValue `json:"xacml"`
}

type AttributeEntity struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Timestamp  TypeValue `json:"timestamp"`
	Version    TypeValue `json:"version"`
	Attributes TypeValue `json:"attributes"`
}

type UpdatePolicy struct {
	Version   TypeValue `json:"version"`
	Xacml     TypeValue `json:"xacml"`
	Timestamp TypeValue `json:"timestamp"`
}

type UpdateAttribute struct {
	Version    TypeValue `json:"version"`
	Attributes TypeValue `json:"attributes"`
	Timestamp  TypeValue `json:"timestamp"`
}

type VersionEntity struct {
	Id      string    `json:"id"`
	Type    string    `json:"type"`
	Version TypeValue `json:"version"`
}

type TypeXACML struct {
	Type string `json:"type"`
}

// XACML
func ManageEntities(w http.ResponseWriter, r *http.Request) {
	//Verify if HTTP method is valid
	switch r.Method {
	case http.MethodPost:
		RegisterXACMLEntities(w, r)
	case http.MethodGet:
		queryParamType := r.URL.Query().Get("type")
		queryParamId := r.URL.Query().Get("id")
		if queryParamType != "" {
			QueryAllXACMLEntities(w, r, queryParamType)
		} else if queryParamId != "" {
			QueryXACMLEntity(w, r, queryParamId)
		} else {
			http.Error(w, "Invalid Query parameter", http.StatusBadRequest)
			return
		}
	case http.MethodPatch:
		UpdateXACMLEntity(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func RegisterXACMLEntities(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - RegisterXACMLEntities")
	var policyPayload XACMLEntity
	var attributePayload AttributeEntity
	var id string
	var registered bool
	var typeXACML TypeXACML

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	entityString := string(body)

	err = json.Unmarshal(body, &typeXACML)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if typeXACML.Type == "xacml" {
		err = json.Unmarshal(body, &policyPayload)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if policyPayload.Id == "" || policyPayload.Type == "" || policyPayload.Version.Type == "" || policyPayload.Version.Value == "" ||
			policyPayload.Timestamp.Type == "" || policyPayload.Timestamp.Value == "" || policyPayload.Xacml.Type == "" ||
			policyPayload.Xacml.Value == "" {

			http.Error(w, "Invalid Policy JSON", http.StatusBadRequest)
			return
		}
		id = policyPayload.Id
	} else if typeXACML.Type == "attributes" {
		err := json.Unmarshal(body, &attributePayload)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if attributePayload.Id == "" || attributePayload.Type == "" || attributePayload.Version.Type == "" || attributePayload.Version.Value == "" ||
			attributePayload.Timestamp.Type == "" || attributePayload.Timestamp.Value == "" || attributePayload.Attributes.Type == "" ||
			attributePayload.Attributes.Value == "" {

			http.Error(w, "Invalid Attributes JSON", http.StatusBadRequest)
			return
		}
		id = attributePayload.Id
	} else {
		http.Error(w, "Invalid Type", http.StatusBadRequest)
		return
	}

	registered, _ = application_gateway.SetXACMLEntity(id, entityString)

	if registered {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "XACML Entity '%s' successfully registered in Hyperledger Fabric\n", id)
	} else {
		http.Error(w, fmt.Sprintf("There was a problem while registering the XACML Entity '%s' in Fabric", id), http.StatusInternalServerError)
		return
	}
}

func QueryAllXACMLEntities(w http.ResponseWriter, r *http.Request, queryParamType string) {
	fmt.Println("Received request - QueryAllEntities")

	if queryParamType != "xacml" && queryParamType != "attributes" {
		http.Error(w, fmt.Sprintf("Invalid type '%s' query parameter", queryParamType), http.StatusBadRequest)
		return
	}

	entities, err := application_gateway.IndexQueryXACML(fmt.Sprintf("{\"selector\":{\"type\":\"%s\"}, \"use_index\":[\"_design/indexTypeDoc\", \"indexType\"]}", queryParamType))
	if err != nil {
		http.Error(w, "No Policies/Attributes found", http.StatusNotFound)
		return
	}

	if entities != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(entities))
	} else {
		http.Error(w, "No Policies/Attributes found", http.StatusNotFound)
		return
	}
}

func UpdateXACMLEntity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - UpdateXACMLEntity")
	var entityType string
	var updateEntity UpdatePolicy
	var updateAttribute UpdateAttribute
	var policyEntity XACMLEntity
	var attributeEntity AttributeEntity

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	entityId := r.URL.Query().Get("id")
	parts := strings.Split(entityId, ":")
	if len(parts) >= 3 {
		entityType = parts[2]
		if entityType == "xacml" {
			err := json.Unmarshal(body, &updateEntity)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			if updateEntity.Version.Type == "" || updateEntity.Version.Value == "" || updateEntity.Timestamp.Type == "" ||
				updateEntity.Timestamp.Value == "" || updateEntity.Xacml.Type == "" || updateEntity.Xacml.Value == "" {

				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
		} else if entityType == "attributes" {
			err := json.Unmarshal(body, &updateAttribute)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			if updateAttribute.Version.Type == "" || updateAttribute.Version.Value == "" || updateAttribute.Timestamp.Type == "" ||
				updateAttribute.Timestamp.Value == "" || updateAttribute.Attributes.Type == "" || updateAttribute.Attributes.Value == "" {

				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Invalid entity", http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, "Invalid Entity ID", http.StatusBadRequest)
		return
	}

	entity, err := application_gateway.GetXACMLEntity(entityId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Entity '%s' not found", entityId), http.StatusNotFound)
		return
	}

	if entityType == "xacml" {
		err := json.Unmarshal([]byte(entity), &policyEntity)
		if err != nil {
			http.Error(w, "There was an error while reading entity JSON", http.StatusInternalServerError)
			return
		}
		policyEntity.Version = updateEntity.Version
		policyEntity.Xacml = updateEntity.Xacml
		policyEntity.Timestamp = updateEntity.Timestamp

		policyEntityJson, err := json.Marshal(policyEntity)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		updated, err := application_gateway.UpdateEntity(policyEntity.Id, string(policyEntityJson))

		if updated {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "XACML Entity %s successfully updated in Hyperledger Fabric\n", policyEntity.Id)
		} else {
			fmt.Fprintf(w, "Failed to update XACML Entity %s in Fabric: %s", policyEntity.Id, err)
			http.Error(w, fmt.Sprintf("There was a problem while updating the XACML Entity %s in Fabric", policyEntity.Id), http.StatusInternalServerError)
			return
		}
	} else {
		err := json.Unmarshal([]byte(entity), &attributeEntity)
		if err != nil {
			http.Error(w, "There was an error while reading entity JSON", http.StatusInternalServerError)
			return
		}
		attributeEntity.Version = updateAttribute.Version
		attributeEntity.Attributes = updateAttribute.Attributes
		attributeEntity.Timestamp = updateAttribute.Timestamp

		attributesEntityJson, err := json.Marshal(attributeEntity)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		updated, err := application_gateway.UpdateEntity(attributeEntity.Id, string(attributesEntityJson))

		if updated {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "XACML Entity %s successfully updated in Hyperledger Fabric\n", attributeEntity.Id)
		} else {
			fmt.Fprintf(w, "Failed to update XACML Entity %s in Fabric: %s", attributeEntity.Id, err)
			http.Error(w, fmt.Sprintf("There was a problem while updating the XACML Entity %s in Fabric", attributeEntity.Id), http.StatusInternalServerError)
			return
		}
	}
}

func QueryXACMLEntity(w http.ResponseWriter, r *http.Request, queryParamId string) {
	fmt.Println("Received request - QueryAllEntities")
	var entity string
	var versionEntity VersionEntity
	var policyEntity XACMLEntity
	var attributesEntity AttributeEntity

	attrs := r.URL.Query().Get("attrs")
	entity, err := application_gateway.GetXACMLEntity(queryParamId)
	if err != nil {
		http.Error(w, fmt.Sprintf("XACML Entity '%s' not found", queryParamId), http.StatusNotFound)
		return
	}
	if attrs == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(entity))
		return
	}

	if attrs == "version" {
		err = json.Unmarshal([]byte(entity), &policyEntity)
		if err != nil {
			err2 := json.Unmarshal([]byte(entity), &attributesEntity)
			if err2 != nil {
				http.Error(w, "There was an error while reading the entity JSON", http.StatusInternalServerError)
				return
			}
			versionEntity.Id = attributesEntity.Id
			versionEntity.Type = attributesEntity.Type
			versionEntity.Version = attributesEntity.Version
		} else {
			versionEntity.Id = policyEntity.Id
			versionEntity.Type = policyEntity.Type
			versionEntity.Version = policyEntity.Version
		}

		versionEntityJson, errJson := json.Marshal(versionEntity)
		if errJson != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(versionEntityJson))
	} else {
		http.Error(w, fmt.Sprintf("Query attrs parameter '%s' not supported", attrs), http.StatusBadRequest)
		return
	}
}
