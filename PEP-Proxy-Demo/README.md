# License

PEP-Proxy Project source code files are made avaialable under the Apache License, Version 2.0 (Apache-2.0), located into the LICENSE file.

# What is a PEP-Proxy

PEP-Proxy is the component responsible for receiving the queries from a client entity (applications, services, users, or devices) aimed to access a resource,  accompanied by the corresponding authorisation token (JWT Access Token)  and forwarding requests to the corresponding endpoint of a target component and the responses back to the requester.

## Project details

This project contains:

- Main source files & folders.

    - [PEP-Proxy.py](./PEP-Proxy.py): Main file of the component, contains the functionality and allows the execution of the HTTP/HTTPS server.

- Folder to ssl (offer HTTPS).

    - [certs](./certs/): Contains the certificates needed to ssl.

- Files to deploy the component.

    - [requirements.txt](./requirements.txt): Contain the auxiliar python modules needed by the application. Dockerfile uses it. 
    - [Dockerfile](./Dockerfile): Contain the actions and commands to create the Docker image needed to run the component.
    - [docker-compose.yml](./docker-compose.yml): To deploy the component.

- File to monitor.

    - [monitor_PEP-Proxy.sh](./monitor_PEP-Proxy.sh): To test if the component is running or not.

- File to clean docker logs.

    - [deleteLogs.sh](./deleteLogs.sh): To clean logs of the component (registered and used by docker).
    **NOTE:** This script remove the logs of all dockerized services present in the machine not only the logs of Capability Manager.

**NOTE**: Please note that the certificates are expired. If you want to use HTTPS, you may create new certificates or not verify the expired certificates (this is recommended if you are only testing the component).

# PEP-Proxy functionality

This component receives queries aimed at accessing a resource. These queries contain a header with an Access Token. The PEP-Proxy validates this token, and if the evaluation is positive, it forwards requests to the specific endpoint’s API. When a request to access a resource is received by  the PEP-Proxy:

- Recovers the x-auth-token header (Access Token).
- Validate Access Token (signature, resource, and expiration date).
- Forward the request to a specific endpoint's API.
- Finally, PEP-Proxy forwards the message and sends responses back to the requester.

# How to deploy/test

## Prerequisites to deploy PEP-Proxy

To run this component it is neccessary to install the docker-compose tool.

https://docs.docker.com/compose/install/

Launch the following components before deploying the PEP-Proxy:

- FLUIDOS IdM must be deployed.
- [REST-API server](../restapi-server/) must be deployed.
- [XACML](../xacml-xadatu/) must be deployed.

## Configuration .env/docker-compose.yml files

Review the content of [.env](./.env) and [docker-compose.yml](./docker-compose.yml) to configure the instance.

**NOTE**: Please change all `<YOUR_IP>` occurrences in the `.env` file with your local IP.

# Installation / Execution.

After the review of [.env](./.env) and [docker-compose.yml](docker-compose.yml) files, the next step is to obtain the Docker image. To do this, you have to build a local one:

```sh
docker-compose build
```

Finally, AFTER REVIEW docker-compose.yml file and especially its environment variables, to launch the connector image we use the next command:

```sh
docker-compose up -d
```

# Testing

To test the PEP-Proxy component, it has been deployed in the 'Producer', so all the requests made to the ThirdPNode/Producer now are sent to the PEP-Proxy instead. Now, requests sent to the `/fluidos/idm/signContract` and `/fluidos/idm/verifyContract` endpoints of the `Producer` include the `x-auth-token` header with the Access Token previously requested. The PEP-Proxy will verify the Access Token, and if the token is successfully verified, the PEP-Proxy will forward the request to the Producer.

The Holder/Consumer request an Access Token to access the resource `https://<pepproxy_ip>:1027/fluidos/idm/.*` using POST. The `.*` means that the Access Token will be valid for all endpoints under /fluidos/idm. For example, the token will be valid for `https://<pepproxy_ip>:1027/fluidos/idm/signContract`.

There is an updated POSTMAN collection that tests all the integrated XADATU security components in [POSTMAN COLLECTION](./postman/FLUIDOS-IDM-REAR.postman_collection.json).

# How to monitoring.

- To test if the container is running:

```sh
docker ps -as
```

The system must return that the status of the PEP-Proxy container is up.

- To show the PEP-Proxy container logs.

```sh
docker logs -f pepproxy
```

# Troubleshooting

**If the certificates were renovated and "pep_protocol"="https", you need to refresh them in pepproxy service. Once certificate files are ubicated in the corresponding folder following prerequisites indications, access the project directory and run:**

```bash  
docker-compose restart pepproxy
```

To solve it automatically, you must add something like this in the crontab, to restart the component every 4 hours:
```
0 0,4,8,12,16,20 * * * cd {{projectFolder}}; docker-compose restart pepproxy
```

**RECOMMENDED, monitor each minute (using the GET endpoint of PEP-Proxy), through crontab, if PEP-Proxyis running:**

```
#To monitor PEP-Proxy.
* * * * * cd {{projectFolder}}; ./monitor_PEP-Proxy.sh >> nohup_PEP-Proxy.out 2>&1
```
*NOTE:* Before configuring it, be sure to define:

- {{Protocol}}://{{IP}}
- {{projectPath}}

**RECOMMENDED, to prevent PEP-Proxy logs from filling up the hard drive include in crontab to clean the logs periodically. --> Logs stored in the container of PEP-Proxy.**

```
#To remove PEP-Proxy logs.
0 0 * * 6 cd {{projectFolder}}; cat /dev/null > pepproxy.log
```

**RECOMMENDED, to prevent PEP-Proxy logs from filling up the hard drive include in crontab to clean the logs periodically. --> Logs stored by docker technology related to services.**

**NOTE:** This script removes the logs of all dockerized services present in the machine not only the logs of PEP-Proxy.

```
# Assign permissions
sudo chmod 755 deleteLogs.sh

# Copy to /usr/bin
sudo su
cd /usr/bin/
cp {{projectFolder}}/deleteLogs.sh ./

# Configure schedule in /etc/crontab
cd /etc
vi crontab
  # Add this line
  0 0     * * 6   root    /usr/bin/deleteLogs.sh
```

# PEP-Proxy REST API

The PEP-Proxy component supports multiple REST APIs. This section defines all REST APIs supported by PEP-Proxy. In this sense:

- Accessing to resource
- Test PEP-Proxy is running

## Accessing to resource

The PEP-Proxy component does not have a specific API, it uses the same API that will forward requests. The unique difference is you must add a new header `x-auth-token` where the Access Token must be included as a string.