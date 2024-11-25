from pydantic_settings import BaseSettings

class ProducerSettings(BaseSettings):
    HOST: str = "0.0.0.0"
    PORT: int = 9083
    PRODUCER_DOMAIN: str = "producer.fluidos.eu"
    PRODUCER_NAME: str = "producer-node-1"
    PRODUCER_NODE_ID: str = "producer-001"
    PRODUCER_IP: str = "localhost:9083"
    
    # Privacy Security Manager API
    SECURITY_MANAGER_URL: str = "http://10.208.99.115:8082"
    
    # REAR API connection
    REAR_API_URL: str = "http://localhost:3002"
    
    class Config:
        env_file = ".env"

settings = ProducerSettings()