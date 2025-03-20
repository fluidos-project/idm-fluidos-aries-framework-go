#!/bin/sh

cat <<EOF > security-authorisation-PDPDLT.pu
@startuml

skinparam monochrome false
skinparam shadowing true
skinparam roundcorner 10

participant "Capability\nManager"
participant "XACML\n(PDP)"
participant "DLT"

group Authorisation DCapBAC - Access Control
    
    "Capability\nManager" -> "XACML\n(PDP)" : XACML Authorisation Request

    group HASH Validation\n(DLT Integration)
        "XACML\n(PDP)" -> "DLT" : GET_HASH(Domain)
        "XACML\n(PDP)" -> "DLT" : Hash
    end

    "XACML\n(PDP)" <- "XACML\n(PDP)" : Validate(XACML_POL_FILE,HASH)
    "XACML\n(PDP)" <- "XACML\n(PDP)" : Validate AuthZ Request

    "Capability\nManager" <- "XACML\n(PDP)" : XACML Veredict
end

@enduml
EOF
plantuml security-authorisation-PDPDLT.pu
