from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    API_HOST: str = "0.0.0.0"
<<<<<<< HEAD
    API_PORT: int = 3002
=======
    API_PORT: int = 3003
>>>>>>> origin/opencall-XADATU
    API_VERSION: str = "v2"
    
    class Config:
        env_file = ".env"

settings = Settings()