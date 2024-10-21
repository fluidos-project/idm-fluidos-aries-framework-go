#!/bin/bash
BASE_DIR=$(realpath "$(dirname "$0")/..")
echo "$BASE_DIR"

# Certs files route - Origin
CERT_DIR="$BASE_DIR/modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem"
CA_CERT_DIR="$BASE_DIR/modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
KEYSTORE_DIR="$BASE_DIR/modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore"

# Find Key File
KEY_FILE=$(find "$KEYSTORE_DIR" -type f -name "*_sk")

# Certs files route - Destiny
CERT_DIR_REST="$BASE_DIR/restapi-server/certs/cert.pem"
CA_CERT_REST="$BASE_DIR/restapi-server/certs/ca.crt"
KEY_REST="$BASE_DIR/restapi-server/certs/key/"

# Verify if files exists before copy them

# CERT.PEM
if [ -f "$CERT_DIR" ]; then
	cp "$CERT_DIR" "$CERT_DIR_REST"
	echo "*** Certified cert.pem copied to $CERT_DIR_REST"
	cat "$CERT_DIR"
	echo
else
	echo "*** cert.pem file not found in $CERT_DIR"
fi

# CA.CRT
if [ -f "$CA_CERT_DIR" ]; then
	cp "$CA_CERT_DIR" "$CA_CERT_REST"
	echo "*** CA Certified ca.crt copied in $CA_CERT_REST"
	cat "$CA_CERT_DIR"
	echo
else
	echo "*** ca.crt file not found in $CA_CERT_DIR"
fi

# KEY FILE
if [ -f "$KEY_FILE" ]; then
    if [ -d "$KEY_REST" ]; then
        # Delete old keys
        rm -rf "$KEY_REST"/*
        echo "*** Directory $KEY_REST cleaned"
    fi
    
    cp "$KEY_FILE" "$KEY_REST"
    echo "*** Key copied in $KEY_REST"
    cat "$KEY_FILE"
else
    echo "*** Key not found in $KEYSTORE_DIR"
fi