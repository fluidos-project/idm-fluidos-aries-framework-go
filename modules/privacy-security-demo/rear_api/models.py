from pydantic import BaseModel
from typing import Dict, Any, Optional
from datetime import datetime

class Location(BaseModel):
    latitude: str
    longitude: str
    country: str
    city: str
    additionalNotes: Optional[str] = None

class Owner(BaseModel):
    domain: str
    nodeId: str
    ip: str
    additionalInformation: Dict[str, Any] = {}

class Price(BaseModel):
    amount: str
    currency: str
    period: str

class FlavorType(BaseModel):
    name: str
    data: Dict[str, Any]

class Flavor(BaseModel):
    flavorId: str
    providerId: str
    timestamp: datetime
    location: Location
    networkPropertyType: str
    type: FlavorType
    price: Price
    owner: Owner
    availability: bool

class NodeIdentity(BaseModel):
    Domain: str
    NodeID: str
    IP: str
    AdditionalInformation: Dict[str, Any] = {}

class Configuration(BaseModel):
    type: str
    data: Dict[str, Any]

class Reservation(BaseModel):
    FlavorID: str
    Buyer: NodeIdentity
    Configuration: Configuration