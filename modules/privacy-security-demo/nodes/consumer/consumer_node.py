from fastapi import FastAPI, HTTPException
from typing import Dict, Any, List
import requests
from datetime import datetime
import sys
import os

# Add the project root to Python path
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))

from privacy_security_manager.src.manager import PrivacySecurityManager
from .config import settings

app = FastAPI()
security_manager = PrivacySecurityManager()

# REAR API client for consumer operations
class RearApiClient:
    def __init__(self, base_url: str):
        self.base_url = base_url

    async def get_flavors(self) -> List[Dict[str, Any]]:
        response = requests.get(f"{self.base_url}/api/v2/flavors")
        response.raise_for_status()
        return response.json()

    async def create_reservation(self, flavor_id: str, buyer_info: Dict[str, Any]) -> Dict[str, Any]:
        payload = {
            "FlavorID": flavor_id,
            "Buyer": buyer_info,
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
    """Generate DID for consumer node"""
    try:
        return await security_manager.generate_did(request.get("name", "consumer"))
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

@app.get("/consumer/flavors")
async def list_flavors() -> Dict[str, Any]:
    """Get available flavors from REAR API"""
    try:
        flavors = await rear_client.get_flavors()
        return {"status": "success", "flavors": flavors}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/consumer/reservations")
async def create_reservation(flavor_id: str) -> Dict[str, Any]:
    """Create a reservation for a flavor"""
    try:
        if not security_manager.current_did:
            raise HTTPException(status_code=400, detail="Generate DID first")
            
        buyer_info = {
            "Domain": settings.CONSUMER_DOMAIN,
            "NodeID": settings.CONSUMER_NODE_ID,
            "IP": settings.CONSUMER_IP,
            "AdditionalInformation": {
                "DID": security_manager.current_did
            }
        }
        
        reservation = await rear_client.create_reservation(flavor_id, buyer_info)
        return {"status": "success", "reservation": reservation}
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

@app.post("/fluidos/idm/generateVPresentation")
async def generate_verifiable_presentation(request: Dict[str, Any]) -> Dict[str, Any]:
    """Generate a verifiable presentation from a credential using a frame"""
    try:
        if not request.get("credId"):
            raise HTTPException(status_code=400, detail="Credential ID is required")
        if not request.get("querybyframe", {}).get("frame"):
            raise HTTPException(status_code=400, detail="Frame is required")
            
        return await security_manager.generate_verifiable_presentation(
            cred_id=request.get("credId"),
            frame=request.get("querybyframe", {}).get("frame")
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/fluidos/idm/verifyCredential")
async def verify_credential(request: Dict[str, Any]) -> Dict[str, Any]:
    """Verify a credential or verifiable presentation"""
    try:
        if not request.get("credential"):
            raise HTTPException(status_code=400, detail="Credential is required")
            
        return await security_manager.verify_credential(
            credential=request.get("credential")
        )
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