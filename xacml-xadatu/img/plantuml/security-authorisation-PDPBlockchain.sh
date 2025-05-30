#!/bin/sh

cat <<EOF > security-authorisation-PDPBlockchain.pu
@startuml

skinparam monochrome false
skinparam shadowing true
skinparam roundcorner 10

participant "Capability\nManager"
participant "XACML\n(PDP)"
participant "Blockchain"

group Authorisation DCapBAC - Access Control
    
    "Capability\nManager" -> "XACML\n(PDP)" : XACML Authorisation Request

    group Blockchain Integration
        "XACML\n(PDP)" -> "Blockchain" : GET_HASH(Domain)
        "XACML\n(PDP)" -> "Blockchain" : Hash
    end

    "XACML\n(PDP)" <- "XACML\n(PDP)" : Validate(XACML_POL_FILE,HASH)
    "XACML\n(PDP)" <- "XACML\n(PDP)" : Validate AuthZ Request

    "Capability\nManager" <- "XACML\n(PDP)" : XACML Veredict
end

@enduml
EOF
plantuml security-authorisation-PDPBlockchain.pu
