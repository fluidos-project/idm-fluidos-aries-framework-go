# Model Gateway

## Overview
This project implements a REST API server in Go that facilitates federated learning operations through Hyperledger Fabric blockchain. It manages two key aspects of federated learning:

1. **Model Aggregation**: Combines model updates from multiple nodes into a consolidated model, ensuring secure and verifiable aggregation through blockchain-based validation.

2. **Base Model Updates**: Handles individual model updates from participating nodes, tracking the evolution of the model across the federated learning network.

The API uses Application Gateways to establish secure SSL channels with Hyperledger Fabric peers, ensuring reliable and authenticated communication for all federated learning operations.

## Key Features
- Secure model update submission and aggregation
- Blockchain-based verification of model updates
- Traceable model evolution history
- Query capabilities for model updates and aggregations
- Support for signed proofs from participating nodes

## Prerequisites

### Required Software
* **Git**: install the latest version of Git if it is not already installed.

  ```bash
  $ sudo apt-get install git
  ```

* **Docker**: install the latest version of Docker if it is not already installed.

    ```bash
    $ sudo apt-get -y install docker-compose
    ```

    Once installed, you can confirm that the latest versions of both Docker and Docker compose executables were installed by running:

    ```bash
    $ docker --version
    Docker version 24.0.6, build ed223bc
    $ docker-compose --version
    docker-compose version 1.29.2, build unknown
    ```

    Make sure the Docker daemon is running

    ```bash
    $ sudo systemctl start docker
    ```

    To start the Docker daemon when the system starts, use the following:
    
    ```bash
    sudo systemctl enable docker
    ```

    Add your user to the Docker group.
    ```bash
    sudo usermod -a -G docker <username>
    ```

    After this, you must restart your session for changes to be applied. You can restart your PC or just close your session.

## **Hyperledger Fabric Certificates**

Update Hyperledger Fabric certificates to connect your Fabric deployment. There are two different procedures depending on `Hyperledger Fabric` and `REST API` are going to be deployed on the same machine or not. 

In this sense:

### Same machine

**IMPORTANT:** The script referended in this subseccion uses the information of `.env`, please define correctly the environment variables before execute it.

If you have deployed Hyperledger Fabric on the same machine on which you are going to deploy the REST-API, execute the following command:

```bash
$ ./updateCerts.sh
```

This will automatically copy the certificates from Hyperledger Fabric into the REST-API.

### Different machine

If your Hyperledger Fabric deployment is on a different machine, you will have to manually copy the certificates:

1. **CA CERT**: Copy the CA certificate from ``../modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt``to ``certs/ca.crt``.

2. **USER CERT**: Copy the user cert from ``../modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem`` to ``certs/cert.pem``

3. **USER KEY**: Copy the user private key from ``../modules/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/<keyName>`` to ``certs/key/`` where ``<keyName>`` is the name assigned to the key for that Fabric deployment.

**NOTE**: The directory ``FabricREST-Go/certs/key/`` should be empty before copying the key, so if it is not empty please delete all keys. If you are executing the script ``./updateCerts.sh`` you do not have to worry about this, as the script will delete all keys in the ``certs/key``directory.

# Run

Before running this service be sure that the privacy and security manager testbed is running. To do this, please refer to the [README](../../README.md) file in the privacy and security manager testbed.

As shortcut, you can go to the root path of this repository and run the following command:

```bash
$ make run-openapi-demo
```

or

```bash
$ make run-fabric
```

To run the REST-API server, you can use the provided Makefile with the following commands:

```bash
# To start the model treatment service (including building and running Docker containers)
$ make run-model-treatment

# To stop the model treatment service
$ make stop-model-treatment

# To clean up the environment (stop containers and remove volumes)
$ make clean
```

The `run-model-treatment` command will:
1. Run the `updateCerts.sh` script to update the certificates
2. Build the Docker containers
3. Start the Docker containers

```bash
$ docker-compose logs -f api
```

# Endpoints
Folder [postman](./postman/) contains a JSON file that can be exported to POSTMAN to test all endpoints described below:

### API Documentation

#### Model Aggregation
1. **Create Aggregated Model**
   - **Endpoint**: `/model/aggregated`
   - **Method**: POST
   - **Headers**: Content-Type: application/json
   - **Description**: Submit a model for aggregation

2. **Get Aggregated Model**
   - **Endpoint**: `/model/aggregated/{id}`
   - **Method**: GET
   - **Description**: Retrieve a specific aggregated model

3. **List All Aggregated Models**
   - **Endpoint**: `/models/aggregated`
   - **Method**: GET
   - **Description**: Get all aggregated models

4. **Query Aggregated Models by Date**
   - **Endpoint**: `/models/aggregated/query`
   - **Method**: GET
   - **Description**: Query models within a date range

#### Base Model Updates
1. **Create Base Model Update**
   - **Endpoint**: `/basemodel/update`
   - **Method**: POST
   - **Headers**: Content-Type: application/json
   - **Description**: Submit a base model update

2. **Get Base Model Update**
   - **Endpoint**: `/basemodel/{id}`
   - **Method**: GET
   - **Description**: Retrieve a specific base model update

3. **List All Base Model Updates**
   - **Endpoint**: `/basemodels`
   - **Method**: GET
   - **Description**: Get all base model updates

4. **Query Base Model Updates by Date**
   - **Endpoint**: `/basemodels/query`
   - **Method**: GET
   - **Description**: Query updates within a date range