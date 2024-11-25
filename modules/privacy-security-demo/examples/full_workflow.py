import asyncio
import requests
from datetime import datetime

# API endpoints
PRODUCER_URL = "http://localhost:9083"
CONSUMER_URL = "http://localhost:8083"
REAR_API_URL = "http://localhost:3002"

async def run_workflow():
    # 1. Generate DIDs for both producer and consumer
    producer_did = requests.post(f"{PRODUCER_URL}/fluidos/idm/generateDID", 
        json={"name": "producer-1"}).json()
    
    consumer_did = requests.post(f"{CONSUMER_URL}/fluidos/idm/generateDID", 
        json={"name": "consumer-1"}).json()
    
    print("DIDs generated:", 
          f"\nProducer: {producer_did['didDoc']['id']}", 
          f"\nConsumer: {consumer_did['didDoc']['id']}")

    # 2. Create a flavor in REAR API
    flavor = {
        "flavorId": "flavor-001",
        "providerId": producer_did['didDoc']['id'],
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
    }
    
    created_flavor = requests.post(f"{REAR_API_URL}/api/v2/flavors", json=flavor).json()
    print("\nFlavor created:", created_flavor)

    # 3. Consumer lists and reserves flavor
    flavors = requests.get(f"{CONSUMER_URL}/consumer/flavors").json()
    print("\nAvailable flavors:", flavors)

    reservation = requests.post(
        f"{CONSUMER_URL}/consumer/reservations",
        params={"flavor_id": "flavor-001"}
    ).json()
    print("\nReservation created:", reservation)

    # 4. Create and sign contract
    contract = {
        "flavorId": "flavor-001",
        "reservationId": reservation["reservation"]["id"],
        "timestamp": datetime.utcnow().isoformat(),
        "signatures": []
    }

    # Consumer signs first
    consumer_signed = requests.post(
        f"{CONSUMER_URL}/fluidos/idm/signContract",
        json={"contract": contract}
    ).json()
    
    # Producer verifies and signs
    producer_signed = requests.post(
        f"{PRODUCER_URL}/fluidos/idm/signContract",
        json={"contract": consumer_signed["contract"]}
    ).json()
    
    print("\nFinal signed contract:", producer_signed)

if __name__ == "__main__":
    asyncio.run(run_workflow())