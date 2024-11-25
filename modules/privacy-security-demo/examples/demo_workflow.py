import asyncio
import requests
from pprint import pprint
from datetime import datetime
import json
import urllib3
import os
import string
import random

# Disable SSL verification warnings
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

CONSUMER_URL = "http://localhost:8083"
PRODUCER_URL = "http://localhost:9083"
REAR_API_URL = "http://localhost:3002"

# Global variables to store workflow state
consumer_did = None
producer_did = None
verifiable_presentation = None
selected_flavor = None
reservation_data = None
contract_data = None
cred_storage_id = None

def print_separator():
    print("\n" + "="*80 + "\n")

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def print_menu():
    clear_screen()
    print("FLUIDOS Demo Workflow")
    print_separator()
    print("Available steps:")
    print("1. Generate Consumer DID")
    print("2. Generate Producer DID")
    print("3. Request Consumer Credential")
    print("4. Generate Verifiable Presentation")
    print("5. List Flavors (with VP auth)")
    print("6. Create Reservation (with VP auth)")
    print("7. Perform Purchase and Producer Signs")
    print("8. Consumer Signs Contract")
    print("9. Verify Contract Signatures")
    print("0. Exit")
    print_separator()

def generate_random_string(length=8):
    """Generate a random string of fixed length"""
    letters = string.ascii_lowercase + string.digits
    return ''.join(random.choice(letters) for _ in range(length))

async def generate_consumer_did():
    global consumer_did
    print("Generating Consumer DID...")
    
    random_suffix = generate_random_string()
    did_request = {
        "name": f"consumer-{random_suffix}",
        "nattrs": 5
    }
    print("Request payload:")
    pprint(did_request)
    
    response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/generateDID",
        json=did_request,
        verify=False
    )
    consumer_did = response.json()
    print("Consumer DID generated:")
    pprint(consumer_did)
    return True

async def generate_producer_did():
    global producer_did
    print("Generating Producer DID...")
    
    random_suffix = generate_random_string()
    did_request = {
        "name": f"producer-{random_suffix}",
        "nattrs": 5
    }
    print("Request payload:")
    pprint(did_request)
    
    response = requests.post(
        f"{PRODUCER_URL}/fluidos/idm/generateDID",
        json=did_request,
        verify=False
    )
    producer_did = response.json()
    print("Producer DID generated:")
    pprint(producer_did)
    return True

async def request_consumer_credential():
    global cred_storage_id
    if not consumer_did:
        print("Error: Consumer DID not generated yet!")
        return False
        
    print("Requesting Consumer Credential...")
    
    enrolment_request = {
        "url": "https://issuer:9082",
        "idProofs": [
            {"attrName": "holderName", "attrValue": "FluidosNode"},
            {"attrName": "fluidosRole", "attrValue": "Customer"},
            {"attrName": "deviceType", "attrValue": "Server"},
            {"attrName": "orgIdentifier", "attrValue": "FLUIDOS_DEMO"},
            {"attrName": "physicalAddress", "attrValue": "00:11:22:33:44:55"}
        ]
    }
    
    response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/doEnrolment",
        json=enrolment_request,
        verify=False
    )
    enrollment_data = response.json()
    print("Enrollment response:")
    pprint(enrollment_data)
    
    cred_storage_id = enrollment_data.get("credStorageId")
    if not cred_storage_id:
        print("Error: No credStorageId in enrollment response")
        return False
    
    print(f"\nSaved credStorageId: {cred_storage_id}")
    return True

async def generate_verifiable_presentation():
    global verifiable_presentation
    if not cred_storage_id:
        print("Error: Complete enrollment first to get credStorageId!")
        return False
        
    print("Generating Verifiable Presentation...")
    
    frame = {
        "@context": [
            "https://www.w3.org/2018/credentials/v1",
            "https://ssiproject.inf.um.es/security/psms/v1"
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
        "credId": cred_storage_id,
        "querybyframe": {"frame": frame}
    }
    
    response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/generateVPresentation",
        json=vpresentation_request,
        verify=False
    )
    verifiable_presentation = response.json()
    print("Verifiable Presentation generated:")
    pprint(verifiable_presentation)
    return True

async def list_flavors_with_vp():
    global selected_flavor
    if not verifiable_presentation:
        print("Error: Generate Verifiable Presentation first!")
        return False
        
    print("Listing flavors with VP authentication...")
    
    # Get the first VP from results array
    vp = verifiable_presentation["results"][0]
    headers = {
        "Authorization": f"Bearer {json.dumps(vp)}"
    }
    
    response = requests.get(
        f"{PRODUCER_URL}/producer/flavors",
        headers=headers,
        verify=False
    )
    flavors = response.json()
    print("Available flavors:")
    pprint(flavors)
    
    # Let user select a flavor
    if flavors.get("flavors"):
        print("\nSelect a flavor by number:")
        for i, flavor in enumerate(flavors["flavors"]):
            print(f"{i+1}. {flavor['flavorId']}")
        
        choice = int(input("\nYour choice (1-{0}): ".format(len(flavors["flavors"]))))
        selected_flavor = flavors["flavors"][choice-1]
        print(f"\nSelected flavor: {selected_flavor['flavorId']}")
    return True

async def create_reservation():
    global reservation_data
    if not selected_flavor or not verifiable_presentation:
        print("Error: Select a flavor and generate VP first!")
        return False
        
    print("Creating reservation...")
    
    vp = verifiable_presentation["results"][0]
    headers = {
        "Authorization": f"Bearer {json.dumps(vp)}"
    }
    
    response = requests.post(
        f"{PRODUCER_URL}/producer/reservations",
        params={"flavor_id": selected_flavor["flavorId"]},
        headers=headers,
        verify=False
    )
    reservation_data = response.json()
    print("Reservation created:")
    pprint(reservation_data)
    return True

async def perform_purchase():
    global contract_data
    if not reservation_data or not verifiable_presentation:
        print("Error: Create reservation first!")
        return False
        
    print("Performing purchase...")
    
    vp = verifiable_presentation["results"][0]
    headers = {
        "Authorization": f"Bearer {json.dumps(vp)}"
    }
    
    # Step 1: Generate contract through purchase
    response = requests.post(
        f"{PRODUCER_URL}/producer/reservations/{reservation_data['reservation']['id']}/purchase",
        headers=headers,
        verify=False
    )
    contract_data = response.json()
    print("Purchase contract generated:")
    pprint(contract_data)

    print("\nWaiting for producer to sign contract...")
    await asyncio.sleep(2)  # Add small delay for better UX

    # Step 2: Producer signs contract automatically
    response = requests.post(
        f"{PRODUCER_URL}/fluidos/idm/signContract",
        json={"contract": contract_data["purchase"]},
        verify=False
    )
    contract_data = response.json()
    print("Producer signed contract:")
    pprint(contract_data)
    return True

async def consumer_sign_contract():
    global contract_data
    if not contract_data:
        print("Error: Producer must sign first!")
        return False
        
    print("Consumer signing contract...")
    
    # Handle nested JWT contract format
    contract_request = {
        "contract": {
            "JWTContract": contract_data["signedContract"]
        }
    }
    
    response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/signContract",
        json=contract_request,
        verify=False
    )
    contract_data = response.json()
    print("Consumer signed contract:")
    pprint(contract_data)
    return True

async def verify_contract():
    if not contract_data:
        print("Error: Both parties must sign first!")
        return False
        
    print("Verifying contract signatures...")
    
    response = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/verifyContract",
        json={"contract": contract_data["signedContract"]},
        verify=False
    )
    verification = response.json()
    print("Contract verification result:")
    pprint(verification)
    return True

async def main():
    while True:
        print_menu()
        choice = input("Select step (0-9): ")
        print_separator()
        
        try:
            if choice == "0":
                break
            elif choice == "1":
                await generate_consumer_did()
            elif choice == "2":
                await generate_producer_did()
            elif choice == "3":
                await request_consumer_credential()
            elif choice == "4":
                await generate_verifiable_presentation()
            elif choice == "5":
                await list_flavors_with_vp()
            elif choice == "6":
                await create_reservation()
            elif choice == "7":
                await perform_purchase()
            elif choice == "8":
                await consumer_sign_contract()
            elif choice == "9":
                await verify_contract()
        except Exception as e:
            print(f"Error: {str(e)}")
        
        input("\nPress Enter to continue...")

if __name__ == "__main__":
    asyncio.run(main())

