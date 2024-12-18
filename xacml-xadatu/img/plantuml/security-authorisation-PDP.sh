#!/bin/sh

cat <<EOF > security-authorisation-PDP.pu
@startuml

skinparam monochrome false
skinparam shadowing true
skinparam roundcorner 10


participant "Capability\nManager"
participant "XACML\n(PDP)"

group Authorisation DCapBAC - Access Control - XACML
    "Capability\nManager" -> "XACML\n(PDP)" : XACML Authorisation Request
    "XACML\n(PDP)" <- "XACML\n(PDP)" : Validate AuthZ Request
    "Capability\nManager" <- "XACML\n(PDP)" : XACML Veredict
end

@enduml
EOF
plantuml security-authorisation-PDP.pu
