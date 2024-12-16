import uvicorn
from .api import app
from .config import settings

if __name__ == "__main__":
    uvicorn.run(
        "rear_api.api:app",
        host=settings.API_HOST,
        port=settings.API_PORT,
        reload=True
    )