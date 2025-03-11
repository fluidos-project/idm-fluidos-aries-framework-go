from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # API Configuration
    NODE_HOST: str = "155.54.210.45"
    NODE_PORT: int = 8082  # Default to holder port, can be changed in .env
    
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