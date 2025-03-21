# FLUIDOS Privacy-Security Demo

This demo showcases the privacy and security features of the FLUIDOS framework, including DID management, Verifiable Presentations, and secure contract signing.

## Prerequisites

- Python 3.8+
- Running instances of:
  - Consumer Node (port 8083)
  - Producer Node (port 9083)
  - REAR API (port 3003)
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
REAR_API_URL = "http://localhost:3003"

If not, you have to run the services manually:

```bash
    python3 -m uvicorn nodes.consumer.consumer_node:app --host 0.0.0.0 --port 8083 --reload
    python3 -m uvicorn nodes.producer.producer_node:app --host 0.0.0.0 --port 9083 --reload
    python3 -m uvicorn rear_api.api:app --host localhost --port 3003 --reload
```

## Running the Demo

1. Replace IPs in [command.go](../../pkg/controller/command/poc/command.go) file:

    - Replace `<PRODUCER_IP>` with your local IP.
    - Replace `<XACML_IP>` with the IP where the XACML is going to be deployed (most probaly your local IP)

2. Run the IDM scenario

    ```bash
    make run-openapi-demo 
    ```

3. Navigate to the API REST server direcetory and deploy the component:

    ```bash
    cd restapi-server
    ```

    ```bash
    ./updateCerts.sh
    ```

    - In the [application_gateway.go](../../restapi-server/application-gateway/application_gateway.go) file, replace `<PEER_IP>` with the IP where the Blockchain is deployed (probably your local IP).

    ```bash
    docker-compose up -d --build --remove-orphans
    ```

4. Navigate to the XADATU XACML directory and deploy yhe component.

    ```bash
    cd xacml-xadatu
    ```

    - In the [.env](./../../xacml-xadatu/.env) file, replace `<rest-api_ip>` with the IP where the REST API server is deployed (probably your local IP).

    - In the [PDP.py](./../../xacml-xadatu/XACML_PDP_PYTHON/PDP.py) file, replace `<REST-API_IP>` with the IP were the REST API server is deployed (probably your local IP).

    ```bash
    docker-compose up -d --build --remove-orphans
    ```

5. Navigate to the PEP-Proxy directory for the demo and deploy de component:

    ```bash
    cd PEP-Proxy-Demo
    ```

    - In the [.env](../../PEP-Proxy-Demo/.env) file, replace all occurrence of `<YOUR_IP>` with your local IP.
 
    ```bash
    docker-compose up -d --build --remove-orphans
    ```

6. Navigate to the demo directory:

    ```bash
    cd modules/privacy-security-demo/examples
    ```

7. Run the demo workflow:

    - In the [demo_workflow.py](./examples/demo_workflow.py) file, replace `<YOUR_IP>` with your local IP.

    - Workflow 1:

    ```bash
    python3 demo_workflow.py
    ```

    - In the [demo_workflow2.py](./examples/demo_workflow2.py) file, replace `<YOUR_IP>` with your local IP.

    - Workflow 2:
    
    ```bash
    python3 demo_workflow2.py
    ```

## Demo Steps

The demo follows a sequential workflow. Available steps:

### Important Notes:
1. Steps must be executed in order
2. The VP token and Access Token is required for steps 5-7 in `demo_workflow2.py`, and for steps 6-8 in `demo_workflow.py`.
3. Contract signing (steps 7-8 or steps 8-9) requires both parties.
4. To update the expiration time of the Access Token (acutal expiration time is 2 minutes, or 120 seconds), go to the [command.go](../../pkg/controller/command/poc/command.go) file and modify the time (in seconds) in the `expiresAt := issuedAt + 120` line.

### Expected Flow for demo1:
1. Generate Consumer DID
2. Generate Producer DID
3. Request Consumer Credential from Issuer
4. Generate Verifiable Presentation
5. Obtain Access Token
5. List Flavors (requires VP and Access Token)
6. Create Reservation (requires VP and Access Token)
7. Perform Purchase (requires VP and Access Token, Producer signs)
8. Consumer Signs Contract
9. Verify Contract Signatures

### Expected Flow for demo2:
1. Generate Consumer DID
2. Generate Producer DID
3. Request Consumer Credential from Issuer
4. Generate Verifiable Presentation
5. List Flavors (requires VP and Access Token)
6. Create Reservation (requires VP and Access Token)
7. Perform Purchase (requires VP and Access Token, Producer signs)
8. Consumer Signs Contract
9. Verify Contract Signatures
