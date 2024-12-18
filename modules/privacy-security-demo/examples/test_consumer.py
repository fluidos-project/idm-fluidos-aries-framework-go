import asyncio
import requests
from pprint import pprint
from datetime import datetime
import uuid
import urllib3
import random
import string
import json

# Disable SSL verification warnings
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

CONSUMER_URL = "http://localhost:8083"
PRODUCER_URL = "http://localhost:9083"

def generate_random_string(length=8):
    """Generate a random string of fixed length"""
    letters = string.ascii_lowercase + string.digits
    return ''.join(random.choice(letters) for _ in range(length))

def print_separator():
    print("\n" + "="*80 + "\n")

async def test_consumer():
    print("Testing Consumer Node API Endpoints")
    print_separator()
    
    # Step 1: Generate DID
    print("1. Testing DID Generation")
    random_suffix = generate_random_string()
    did_request = {
        "name": f"test-consumer-{random_suffix}",
        "nattrs": 5
    }
    print("Request payload:")
    pprint(did_request)
    
    did_response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/generateDID",
        json=did_request,
        verify=False
    )
    print(f"\nResponse status: {did_response.status_code}")
    print("Response data:")
    pprint(did_response.json())
    print_separator()
    
    # Step 2: Enrollment
    print("2. Testing Enrollment")
    enrolment_request = {
        "url": "https://issuer:9082",
        "idProofs": [
            {
                "attrName": "holderName",
                "attrValue": "FluidosNode"
            },
            {
                "attrName": "fluidosRole",
                "attrValue": "Customer"
            },
            {
                "attrName": "deviceType",
                "attrValue": "Server"
            },
            {
                "attrName": "orgIdentifier",
                "attrValue": "FLUIDOS_id_23241231412"
            },
            {
                "attrName": "physicalAddress",
                "attrValue": "50:80:61:82:ab:c9"
            }
        ]
    }
    
    print("Request payload:")
    pprint(enrolment_request)
    
    enrolment = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/doEnrolment",
        json=enrolment_request,
        verify=False
    )
    print(f"\nResponse status: {enrolment.status_code}")
    print("Response data:")
    pprint(enrolment.json())
    
    # Save credStorageId from enrollment response
    cred_storage_id = enrolment.json().get("credStorageId")
    if not cred_storage_id:
        raise Exception("No credStorageId in enrollment response")
    print(f"\nSaved credStorageId: {cred_storage_id}")
    print_separator()
    
    # Step 3: Generate VP using the saved credStorageId
    print("3. Testing Verifiable Presentation Generation")
    frame = {
        "@context": [
            "https://www.w3.org/2018/credentials/v1",
            "https://www.w3.org/2018/credentials/examples/v1",
            "https://ssiproject.inf.um.es/security/psms/v1",
            "https://ssiproject.inf.um.es/poc/context/v1"
        ],
        "type": ["VerifiableCredential", "FluidosCredential"],
        "@explicit": True,
        "identifier": {},
        "issuer": {},
        "issuanceDate": {},
        "credentialSubject": {
            "@explicit": True,
            "fluidosRole": {},
            "holderName": {}
        }
    }
    
    vpresentation_request = {
        "credId": cred_storage_id,  # Use the saved credStorageId
        "querybyframe": {"frame": frame}
    }
    
    print("Request payload:")
    pprint(vpresentation_request)
    
    vpresentation = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/generateVPresentation",
        json=vpresentation_request,
        verify=False
    )
    print(f"\nResponse status: {vpresentation.status_code}")
    print("Response data:")
    pprint(vpresentation.json())
    
    # Step 4: Send VP to producer for verification
    print("\n4. Testing Producer Verification of VP")
    
    # Get the first VP from results array and escape quotes
    vp_response = vpresentation.json()
    if not vp_response.get("results"):
        raise Exception("No results in VP response")
    
    first_vp = vp_response["results"][0]
    # Convert to string and escape quotes
    first_vp_str = json.dumps(first_vp).replace('"', '\\"')
    
    verify_request = {
        "credential": first_vp_str
    }
    print("\nSending VP to producer for verification")
    print("Request payload:")
    pprint(verify_request)
    
    producer_verification = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/verifyCredential",
        json=verify_request,
        verify=False
    )
    print(f"\nProducer verification status: {producer_verification.status_code}")
    print("Producer verification response:")
    pprint(producer_verification.json())
    print_separator()

if __name__ == "__main__":
    asyncio.run(test_consumer())