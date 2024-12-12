import asyncio
import requests
from pprint import pprint
from datetime import datetime
import urllib3

# Disable SSL verification warnings
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

PRODUCER_URL = "http://localhost:9083"
REAR_API_URL = "http://localhost:3002"

def print_separator():
    print("\n" + "="*80 + "\n")

async def test_producer():
    print("Testing Producer Node API Endpoints")
    print_separator()
    
    print("1. Testing DID Generation")
    did_request = {
        "name": "test-producer",
        "nattrs": 5
    }
    print("Request payload:")
    pprint(did_request)
    
    did_response = requests.post(
        f"{PRODUCER_URL}/fluidos/idm/generateDID",
        json=did_request,
        verify=False
    )
    print(f"\nResponse status: {did_response.status_code}")
    print("Response data:")
    pprint(did_response.json())
    print_separator()

    # Test REAR API Integration
    print("3. Testing REAR API Integration")
    
    # Create a flavor
    print("3.1 Creating a test flavor")
    flavor = {
        "flavorId": "test-flavor-001",
        "timestamp": datetime.utcnow().isoformat(),
        "location": {
            "latitude": "40.4168",
            "longitude": "-3.7038",
            "country": "Spain",
            "city": "Madrid",
            "additionalNotes": "Test datacenter"
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
            "nodeId": "test-node-001",
            "ip": "192.168.1.100",
            "additionalInformation": {}
        },
        "availability": True
    }
    
    created_flavor = requests.post(
        f"{PRODUCER_URL}/producer/flavors",
        json=flavor,
        verify=False
    )
    print(f"Response status: {created_flavor.status_code}")
    print("Created flavor:")
    pprint(created_flavor.json())
    print_separator()

    # List flavors
    print("3.2 Listing available flavors")
    flavors = requests.get(
        f"{PRODUCER_URL}/producer/flavors",
        verify=False
    )
    print(f"Response status: {flavors.status_code}")
    print("Available flavors:")
    pprint(flavors.json())
    print_separator()

    # List reservations
    print("3.3 Listing current reservations")
    reservations = requests.get(
        f"{PRODUCER_URL}/producer/reservations",
        verify=False
    )
    print(f"Response status: {reservations.status_code}")
    print("Current reservations:")
    pprint(reservations.json())
    print_separator()

    # Test purchase if there are any reservations
    if reservations.json().get("reservations"):
        print("3.4 Testing purchase for first reservation")
        first_reservation = reservations.json()["reservations"][0]
        purchase = requests.post(
            f"{PRODUCER_URL}/producer/reservations/{first_reservation['id']}/purchase",
            verify=False
        )
        print(f"Response status: {purchase.status_code}")
        print("Purchase result:")
        pprint(purchase.json())
        print_separator()

if __name__ == "__main__":
    asyncio.run(test_producer())