from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # API Configurations
    PRODUCER_AGENT_HOST: str = "10.208.99.115"
    PRODUCER_AGENT_PORT: int = 9082
    CONSUMER_AGENT_HOST: str = "10.208.99.115"
    CONSUMER_AGENT_PORT: int = 8082
    
    # REAR API Configuration
    REAR_API_HOST: str = "localhost"
    REAR_API_PORT: int = 3002

    # API Endpoints
    GENERATE_DID_ENDPOINT: str = "/fluidos/idm/generateDID"
    DO_ENROLMENT_ENDPOINT: str = "/fluidos/idm/doEnrolment"
    VERIFY_CREDENTIAL_ENDPOINT: str = "/fluidos/idm/verifyCredential"
    SIGN_CONTRACT_ENDPOINT: str = "/fluidos/idm/signContract"
    VERIFY_CONTRACT_ENDPOINT: str = "/fluidos/idm/verifyContract"
    GENERATE_VPRESENTATION_ENDPOINT: str = "/fluidos/idm/generateVP"

    class Config:
        env_file = ".env"

settings = Settings()