import requests
import json
from datetime import datetime
from pprint import pprint

REAR_API_URL = "http://localhost:3002"

def print_separator():
    print("\n" + "="*80 + "\n")

def create_test_flavors():
    print("Creating test flavors...")
    print_separator()
    
    flavors = [
        {
            "flavorId": "flavor-001",
            "providerId": "did:provider:1",
            "timestamp": datetime.utcnow().isoformat(),
            "location": {
                "latitude": "40.4168",
                "longitude": "-3.7038",
                "country": "Spain",
                "city": "Madrid",
                "additionalNotes": "Main datacenter"
            },
            "networkPropertyType": "5G",
            "type": {
                "name": "k8slice",
                "data": {
                    "cpu": "4",
                    "memory": "8Gi",
                    "pods": "10"
                }
            },
            "price": {
                "amount": "100",
                "currency": "EUR",
                "period": "hourly"
            },
            "owner": {
                "domain": "fluidos.eu",
                "nodeId": "node-001",
                "ip": "192.168.1.100",
                "additionalInformation": {}
            },
            "availability": True
        },
        {
            "flavorId": "flavor-002",
            "providerId": "did:provider:2",
            "timestamp": datetime.utcnow().isoformat(),
            "location": {
                "latitude": "41.3851",
                "longitude": "2.1734",
                "country": "Spain",
                "city": "Barcelona",
                "additionalNotes": "Edge datacenter"
            },
            "networkPropertyType": "4G",
            "type": {
                "name": "k8slice",
                "data": {
                    "cpu": "2",
                    "memory": "4Gi",
                    "pods": "5"
                }
            },
            "price": {
                "amount": "50",
                "currency": "EUR",
                "period": "hourly"
            },
            "owner": {
                "domain": "fluidos.eu",
                "nodeId": "node-002",
                "ip": "192.168.1.101",
                "additionalInformation": {}
            },
            "availability": True
        }
    ]

    for flavor in flavors:
        print(f"Creating flavor {flavor['flavorId']}:")
        print("Request payload:")
        pprint(flavor)
        
        response = requests.post(f"{REAR_API_URL}/api/v2/flavors", json=flavor)
        print(f"\nResponse status: {response.status_code}")
        print("Response data:")
        pprint(response.json())
        print_separator()
        
def test_api_endpoints():
    print("Testing REAR API endpoints...")
    print_separator()
    
    # Test GET /api/v2/flavors
    print("1. Testing GET /api/v2/flavors")
    response = requests.get(f"{REAR_API_URL}/api/v2/flavors")
    print(f"Response status: {response.status_code}")
    print("Available flavors:")
    pprint(response.json())
    print_separator()
    
    # Test POST /api/v2/reservations
    print("2. Testing POST /api/v2/reservations")
    reservation_payload = {
        "FlavorID": "flavor-001",
        "Buyer": {
            "Domain": "test.fluidos.eu",
            "NodeID": "test-001",
            "IP": "192.168.1.200",
            "AdditionalInformation": {
                "test": True,
                "purpose": "Testing reservation"
            }
        },
        "Configuration": {
            "type": "k8slice",
            "data": {
                "cpu": "2",
                "memory": "4Gi",
                "storage": "100Gi"
            }
        }
    }
    
    print("Request payload:")
    pprint(reservation_payload)
    
    response = requests.post(f"{REAR_API_URL}/api/v2/reservations", json=reservation_payload)
    print(f"\nResponse status: {response.status_code}")
    print("Response data:")
    pprint(response.json())
    reservation_id = response.json().get("id")
    print_separator()
    
    # Test POST /api/v2/transactions/{id}/purchase
    print("3. Testing POST /api/v2/transactions/{id}/purchase")
    
    response = requests.post(
        f"{REAR_API_URL}/api/v2/transactions/{reservation_id}/purchase"
    )
    print(f"\nResponse status: {response.status_code}")
    print("Response data:")
    pprint(response.json())
    print_separator()

if __name__ == "__main__":
    test_api_endpoints() 