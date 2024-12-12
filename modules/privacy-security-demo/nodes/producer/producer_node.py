from fastapi import FastAPI, HTTPException
from typing import Dict, Any, List
import requests
from datetime import datetime
import sys
import os
import random
import string
from fastapi.middleware.cors import CORSMiddleware

sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))

from privacy_security_manager.src.manager import PrivacySecurityManager
from .config import settings

app = FastAPI()
security_manager = PrivacySecurityManager(agent_type="producer")

def generate_random_string(length=8):
    """Generate a random string of fixed length"""
    letters = string.ascii_lowercase + string.digits
    return ''.join(random.choice(letters) for _ in range(length))

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class RearApiClient:
    def __init__(self, base_url: str):
        self.base_url = base_url

    async def get_flavors(self) -> List[Dict[str, Any]]:
        response = requests.get(f"{self.base_url}/api/v2/flavors")
        response.raise_for_status()
        return response.json()

    async def create_flavor(self, flavor_data: Dict[str, Any]) -> Dict[str, Any]:
        response = requests.post(f"{self.base_url}/api/v2/flavors", json=flavor_data)
        response.raise_for_status()
        return response.json()

    async def get_reservations(self) -> List[Dict[str, Any]]:
        response = requests.get(f"{self.base_url}/api/v2/reservations")
        response.raise_for_status()
        return response.json()

    async def create_purchase(self, reservation_id: str) -> Dict[str, Any]:
        response = requests.post(
            f"{self.base_url}/api/v2/transactions/{reservation_id}/purchase"
        )
        response.raise_for_status()
        return response.json()

    async def create_reservation(self, flavor_id: str, producer_info: Dict[str, Any]) -> Dict[str, Any]:
        payload = {
            "FlavorID": flavor_id,
            "Buyer": producer_info,
            "Configuration": {
                "type": "k8slice",
                "data": {"cpu": "2", "memory": "4Gi"}
            }
        }
        response = requests.post(f"{self.base_url}/api/v2/reservations", json=payload)
        response.raise_for_status()
        return response.json()

rear_client = RearApiClient(settings.REAR_API_URL)

@app.post("/fluidos/idm/generateDID")
async def generate_did(request: Dict[str, Any]) -> Dict[str, Any]:
    """Generate DID for producer node with random suffix"""
    try:
        random_suffix = generate_random_string()
        name = f"{request.get('name', 'producer')}-{random_suffix}"
        return await security_manager.generate_did(
            name=name,
            nattrs=int(request.get("nattrs", 5))
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/fluidos/idm/doEnrolment")
async def do_enrolment(request: Dict[str, Any]) -> Dict[str, Any]:
    """Request enrolment from issuer"""
    try:
        if not request.get("url"):
            raise HTTPException(status_code=400, detail="URL is required")
        if not request.get("idProofs"):
            raise HTTPException(status_code=400, detail="ID proofs are required")
            
        return await security_manager.do_enrolment(
            url=request.get("url"),
            id_proofs=request.get("idProofs")
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/fluidos/idm/verifyCredential")
async def verify_credential(request: Dict[str, Any]) -> Dict[str, Any]:
    """Verify a credential"""
    try:
        if not request.get("credential"):
            raise HTTPException(status_code=400, detail="Credential is required")
            
        return await security_manager.verify_credential(
            credential=request.get("credential")
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/fluidos/idm/signContract")
async def sign_contract(request: Dict[str, Any]) -> Dict[str, Any]:
    """Sign a contract"""
    try:
        if not request.get("contract"):
            raise HTTPException(status_code=400, detail="Contract is required")
            
        return await security_manager.sign_contract(request.get("contract"))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/fluidos/idm/verifyContract")
async def verify_contract(request: Dict[str, Any]) -> Dict[str, Any]:
    """Verify contract signatures"""
    try:
        if not request.get("contract"):
            raise HTTPException(status_code=400, detail="Contract is required")
            
        return await security_manager.verify_contract(request.get("contract"))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/producer/flavors")
async def list_flavors() -> Dict[str, Any]:
    """Get available flavors from REAR API"""
    try:
        flavors = await rear_client.get_flavors()
        return {"status": "success", "flavors": flavors}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/producer/flavors")
async def create_flavor(flavor: Dict[str, Any]) -> Dict[str, Any]:
    """Create a new flavor in REAR API"""
    try:
        if not security_manager.current_did:
            raise HTTPException(status_code=400, detail="Generate DID first")
        
        flavor["providerId"] = security_manager.current_did
        flavor["timestamp"] = datetime.utcnow().isoformat()
        
        created_flavor = await rear_client.create_flavor(flavor)
        return {"status": "success", "flavor": created_flavor}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/producer/reservations")
async def list_reservations() -> Dict[str, Any]:
    """Get all reservations from REAR API"""
    try:
        reservations = await rear_client.get_reservations()
        return {"status": "success", "reservations": reservations}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/producer/reservations/{reservation_id}/purchase")
async def purchase_reservation(reservation_id: str) -> Dict[str, Any]:
    """Create purchase contract for a reservation"""
    try:
        purchase = await rear_client.create_purchase(reservation_id)
        return {"status": "success", "purchase": purchase}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/producer/reservations")
async def create_reservation(flavor_id: str) -> Dict[str, Any]:
    """Create a reservation for a flavor"""
    try:
        if not security_manager.current_did:
            raise HTTPException(status_code=400, detail="Generate DID first")
            
        producer_info = {
            "Domain": settings.PRODUCER_DOMAIN,
            "NodeID": settings.PRODUCER_NODE_ID,
            "IP": settings.PRODUCER_IP,
            "AdditionalInformation": {
                "DID": security_manager.current_did
            }
        }
        
        reservation = await rear_client.create_reservation(flavor_id, producer_info)
        return {"status": "success", "reservation": reservation}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))