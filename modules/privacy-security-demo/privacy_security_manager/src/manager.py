import requests
from typing import Dict, Any, Optional, List
from ..config import settings

class PrivacySecurityManager:
    def __init__(self, agent_type: str = "consumer"):
        self.agent_type = agent_type
        self.base_url = self._get_base_url()
        self.current_did = None

    def _get_base_url(self) -> str:
        if self.agent_type == "consumer":
            return f"https://{settings.CONSUMER_AGENT_HOST}:{settings.CONSUMER_AGENT_PORT}"
        return f"https://{settings.PRODUCER_AGENT_HOST}:{settings.PRODUCER_AGENT_PORT}"

    async def generate_did(self, name: str, nattrs: int = 5) -> Dict[str, Any]:
        """Generate a new DID for the entity"""
        if self.agent_type == "producer":
            payload = {
                "keys": [
                    {
                        "keyType": {
                            "keytype": "Ed25519VerificationKey2018"
                        },
                        "purpose": "Authentication"
                    },
                    {
                        "keyType": {
                            "keytype": "Bls12381G1Key2022",
                            "attrs": [str(nattrs)]
                        },
                        "purpose": "AssertionMethod"
                    }
                ],
                "name": name
            }
        else:
            payload = {
                "keys": [{
                    "keyType": {
                        "keytype": "Ed25519VerificationKey2018"
                    },
                    "purpose": "Authentication"
                }],
                "name": name
            }
        
        response = requests.post(
            f"{self.base_url}{settings.GENERATE_DID_ENDPOINT}",
            json=payload,
            verify=False
        )
        response.raise_for_status()
        result = response.json()
        self.current_did = result["didDoc"]["id"]
        return result

    async def do_enrolment(self, url: str, id_proofs: List[Dict[str, str]]) -> Dict[str, Any]:
        """Request enrolment from issuer"""
        if not self.current_did:
            raise ValueError("No DID available. Generate one first.")
        
        payload = {
            "url": url,
            "idProofs": id_proofs
        }
        
        response = requests.post(
            f"{self.base_url}{settings.DO_ENROLMENT_ENDPOINT}",
            json=payload,
            verify=False
        )
        response.raise_for_status()
        return response.json()

    async def generate_verifiable_presentation(self, cred_id: str, frame: Dict[str, Any]) -> Dict[str, Any]:
        """Generate a verifiable presentation from a credential using a frame"""
        if not self.current_did:
            raise ValueError("No DID available. Generate one first.")
        
        payload = {
            "credId": cred_id,
            "querybyframe": {
                "frame": frame
            }
        }
        
        response = requests.post(
            f"{self.base_url}{settings.GENERATE_VPRESENTATION_ENDPOINT}",
            json=payload,
            verify=False
        )
        response.raise_for_status()
        return response.json()

    async def verify_credential(self, credential: Dict[str, Any], endpoint, method) -> Dict[str, Any]:
        """Verify a credential or verifiable presentation"""
        if not self.current_did:
            raise ValueError("No DID available. Generate one first.")
        
        payload = {
            "credential": credential,
            "endpoint": endpoint,
            "method": method
        }
        
        response = requests.post(
            f"{self.base_url}{settings.VERIFY_CREDENTIAL_ENDPOINT}",
            json=payload,
            verify=False
        )
        response.raise_for_status()
        return response.json()

    async def sign_contract(self, contract: Dict[str, Any]) -> Dict[str, Any]:
        """Sign a contract"""
        if not self.current_did:
            raise ValueError("No DID available. Generate one first.")
        
        response = requests.post(
            f"{self.base_url}{settings.SIGN_CONTRACT_ENDPOINT}",
            json={"contract": contract},
            verify=False
        )
        response.raise_for_status()
        return response.json()

    async def verify_contract(self, contract: Dict[str, Any]) -> Dict[str, Any]:
        response = requests.post(
            f"{self.base_url}{settings.VERIFY_CONTRACT_ENDPOINT}",
            json={"contract": contract},
            verify=False
        )
        response.raise_for_status()
        return response.json()