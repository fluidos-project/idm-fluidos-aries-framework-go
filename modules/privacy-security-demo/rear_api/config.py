from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    API_HOST: str = "0.0.0.0"
    API_PORT: int = 3002
    API_VERSION: str = "v2"
    
    class Config:
        env_file = ".env"

settings = Settings()