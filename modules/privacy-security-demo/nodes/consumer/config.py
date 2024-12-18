from pydantic_settings import BaseSettings

class ConsumerSettings(BaseSettings):
    HOST: str = "0.0.0.0"
    PORT: int = 8083
    CONSUMER_NAME: str = "consumer-node-1"
    CONSUMER_DOMAIN: str = "consumer.fluidos.eu"
    CONSUMER_NODE_ID: str = "consumer-001"
    CONSUMER_IP: str = "localhost:8083"
    
    # Privacy Security Manager API
<<<<<<< HEAD
    SECURITY_MANAGER_URL: str = "http://10.208.99.115:8082"
    
    # REAR API connection
    REAR_API_URL: str = "http://localhost:3002"
=======
    SECURITY_MANAGER_URL: str = "http://localhost:8082"
    
    # REAR API connection
    REAR_API_URL: str = "http://localhost:3003"
>>>>>>> origin/opencall-XADATU
    
    class Config:
        env_file = ".env"

settings = ConsumerSettings()