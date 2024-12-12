# FLUIDOS Privacy-Security Demo

This demo showcases the privacy and security features of the FLUIDOS framework, including DID management, Verifiable Presentations, and secure contract signing.

## Prerequisites

- Python 3.8+
- Running instances of:
  - Consumer Node (port 8083)
  - Producer Node (port 9083)
  - REAR API (port 3002)
  - Issuer Service (port 9082)

## Setup

1. Clone the repository
2. Install dependencies:

    ```bash
    pip install requests asyncio urllib3
    ```
3. Verify the following services are running and accessible:

CONSUMER_URL = "http://localhost:8083"
PRODUCER_URL = "http://localhost:9083"
REAR_API_URL = "http://localhost:3002"

If not, you have to run the services manually:

```bash
    python3 -m uvicorn nodes.consumer.consumer_node:app --host 0.0.0.0 --port 8083 --reload
    python3 -m uvicorn nodes.producer.producer_node:app --host 0.0.0.0 --port 9083 --reload
    python3 -m uvicorn rear_api.mock_api:app --host localhost --port 3002 --reload
```

## Running the Demo

1. Navigate to the demo directory:

    ```bash
    cd modules/privacy-security-demo/examples
    ```

2. Run the demo workflow:

    ```bash
    python3 demo_workflow.py
    ```

## Demo Steps

The demo follows a sequential workflow. Available steps:

### Important Notes:
1. Steps must be executed in order
2. The VP token is required for steps 5-7
3. Contract signing (steps 7-8) requires both parties

### Expected Flow:
1. Generate Consumer DID
2. Generate Producer DID
3. Request Consumer Credential from Issuer
4. Generate Verifiable Presentation
5. List Flavors (requires VP)
6. Create Reservation (requires VP)
7. Perform Purchase (Producer signs)
8. Consumer Signs Contract
9. Verify Contract Signatures