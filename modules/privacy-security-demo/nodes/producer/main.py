import uvicorn
from .producer_node import app
from .config import settings

if __name__ == "__main__":
    uvicorn.run(
        "nodes.producer.src.producer_node:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=True
    )