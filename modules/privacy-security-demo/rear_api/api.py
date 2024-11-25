from fastapi import FastAPI, HTTPException
from typing import List, Dict, Any
import uuid
from datetime import datetime
from .models import Flavor, Reservation

app = FastAPI()

# Default flavors for testing
DEFAULT_FLAVORS = [
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

# In-memory storage
flavors_db: Dict[str, Flavor] = {
    flavor["flavorId"]: Flavor(**flavor) for flavor in DEFAULT_FLAVORS
}
reservations_db: Dict[str, Dict] = {}
contracts_db: Dict[str, Dict] = {}

@app.get("/api/v2/flavors")
async def get_flavors() -> List[Flavor]:
    return list(flavors_db.values())

@app.post("/api/v2/flavors")
async def create_flavor(flavor: Flavor) -> Flavor:
    flavors_db[flavor.flavorId] = flavor
    return flavor

@app.post("/api/v2/reservations")
async def create_reservation(reservation: Reservation) -> Dict[str, Any]:
    reservation_id = str(uuid.uuid4())
    reservation_data = {
        "id": reservation_id,
        "status": "pending",
        "timestamp": datetime.utcnow().isoformat(),
        **reservation.dict()
    }
    reservations_db[reservation_id] = reservation_data
    return reservation_data

@app.post("/api/v2/transactions/{transaction_id}/purchase")
async def purchase_transaction(transaction_id: str) -> Dict[str, Any]:
    if transaction_id not in contracts_db:
        # Get reservation and flavor details
        reservation = reservations_db.get(transaction_id, {})
        flavor_id = reservation.get("FlavorID")
        flavor = flavors_db.get(flavor_id)
        
        # Generate contract with proper structure
        generated_contract = {
            "apiVersion": "reservation.fluidos.eu/v1alpha1",
            "kind": "Contract",
            "metadata": {
                "creationTimestamp": datetime.utcnow().isoformat(),
                "generation": 1,
                "name": f"contract-fluidos.eu-k8s-fluidos-{str(uuid.uuid4())[:8]}",
                "namespace": "fluidos",
                "resourceVersion": str(int(datetime.utcnow().timestamp())),
                "uid": str(uuid.uuid4())
            },
            "spec": {
                "buyer": {
                    "domain": reservation.get("Buyer", {}).get("Domain", "fluidos.eu"),
                    "ip": reservation.get("Buyer", {}).get("IP", "172.18.0.4:30000"),
                    "nodeID": reservation.get("Buyer", {}).get("NodeID", str(uuid.uuid4())[:10])
                },
                "buyerClusterID": str(uuid.uuid4()),
                "expirationTime": (datetime.utcnow().replace(year=datetime.utcnow().year + 1)).isoformat(),
                "flavour": {
                    "metadata": {
                        "name": f"fluidos.eu-k8s-fluidos-{str(uuid.uuid4())[:8]}",
                        "namespace": "fluidos"
                    },
                    "spec": {
                        "characteristics": {
                            "architecture": "amd64",
                            "cpu": "7985105637n",
                            "ephemeral-storage": "0",
                            "gpu": "0",
                            "memory": "32386980Ki",
                            "persistent-storage": "0",
                            "pods": "110"
                        },
                        "optionalFields": {
                            "availability": True,
                            "workerID": "fluidos-provider-1-worker2"
                        },
                        "owner": {
                            "domain": "fluidos.eu",
                            "ip": "172.18.0.2:30001",
                            "nodeID": "jgmewzljr9"
                        },
                        "policy": {
                            "aggregatable": {
                                "maxCount": 0,
                                "minCount": 0
                            },
                            "partitionable": {
                                "cpuMin": "0",
                                "cpuStep": "1",
                                "memoryMin": "0",
                                "memoryStep": "100Mi",
                                "podsMin": "0",
                                "podsStep": "0"
                            }
                        },
                        "price": {
                            "amount": "",
                            "currency": "",
                            "period": ""
                        },
                        "providerID": "jgmewzljr9",
                        "type": "k8s-fluidos"
                    },
                    "status": {
                        "creationTime": datetime.utcnow().isoformat(),
                        "expirationTime": (datetime.utcnow().replace(year=datetime.utcnow().year + 1)).isoformat(),
                        "lastUpdateTime": datetime.utcnow().isoformat()
                    }
                },
                "partition": {
                    "architecture": "amd64",
                    "cpu": "1",
                    "ephemeral-storage": "0",
                    "gpu": "0",
                    "memory": "1Gi",
                    "pods": "50",
                    "storage": "0"
                },
                "seller": {
                    "domain": "fluidos.eu",
                    "ip": "172.18.0.2:30001",
                    "nodeID": "jgmewzljr9"
                },
                "sellerCredentials": {
                    "clusterID": str(uuid.uuid4()),
                    "clusterName": "fluidos-provider-1",
                    "endpoint": "https://172.18.0.2:32197",
                    "token": str(uuid.uuid4()) + str(uuid.uuid4())
                },
                "transactionID": f"{str(uuid.uuid4())}-{int(datetime.utcnow().timestamp())}"
            }
        }
        
        contracts_db[transaction_id] = {
            "status": "completed",
            "timestamp": datetime.utcnow().isoformat(),
            "contract": generated_contract
        }
    
    return contracts_db[transaction_id]