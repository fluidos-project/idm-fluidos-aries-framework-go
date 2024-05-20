# Moded Aries-Framework-Go for FLUIDOS Identity Management

Aries Agent with modification for:

- Use dp-abc crypto
- Use vdr fabric for hyperledger fabric as verfiable storage
- Simplified Demo
- Usage example
- iDM-poc for FLUIDOS

## Requisites

- Go(last check with version 1.19.8)

```
apt-get install go
```


- Docker-compose 1.28.5+

```
curl -L https://github.com/docker/compose/releases/download/1.28.5/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

- JQ command

```
apt install jq
```

## Quick start

To obtain the repo with submodules run clone like this:

```
git clone git@github.com:fluidos-project/p-abc.git
cd idm-fluidos-aries-framework-go
git submodule  update --init --recursive
```

### Hyperledger Fabric

There are new Make Rules related to the deployment of Hyperledger Fabric:

- build-fabric
- run-fabric
- stop-fabric
- restart-fabric # TODO include and test

To run this mod aries-framework-go which uses a basic vdr to connect with Hyperledger Fabric, first of all and only once, you need to run: `make build-fabric`. It will download docker and resources needed to ./modules/fabric-samples.

### Deploy agent

The new Fabric VDR is automatically added by default. Besides, the rules related to launching de Aries Demo `run-openapi-demo` were modified in order to run hyperledger fabric automatically with the necessary rules.

```
run-openapi-demo: stop-openapi-demo generate-test-keys generate-dpabc-clib run-fabric generate-openapi-demo-specs
```

So, doing `make run-openapi-demo` will now, stop it, generate dp-abc C libraries and run Hyperledger Fabric deployment, already with an SmartContract installed. See scripts `./scripts/fabric/*`

`NOTE:` if you have problems deploying everything and errors appears, may is a problem related to volumes or containers that didn't get correctly stopped, so try:

```
make clean
make stop-openapi-demo
```

Open a browser and go to http://localhost:8089/openapi or http://localhost:9089/openapi to see API Swagger if you deploy in your local machine

### FLUIDOS iDM POC Calls

```
/fluidos/idm/VerifyCredential
/fluidos/idm/acceptDeviceEnrolment
/fluidos/idm/doDeviceEnrolment
/fluidos/idm/generateVp
/fluidos/idm/newDID
```


### Test and Development environment

To test the changes made in the implementation restarting only the aries services do:

`make run-openapi-demo-build-no-clean`

#### Remote DEBUG :arrow_down::arrow_down::arrow_down:

Before running `make run-openapi-demo`

If you want a remote debug mode with remote delve run `export PROFILE_DEV='dev'` before `make run-openapi-demo`, it will open 2 ports in the API containers. To setup your VSCode remote debugger add this to your VSCode run configuration launch.json :


```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [

        {
            "name": "holder",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "",
            "port": 4000,
            "host": "10.208.99.115",
            "showLog": true,
            "trace": "log",
            "logOutput": "rpc",
            "dlvLoadConfig": {
                "followPointers": true,
                "maxVariableRecurse": 1,
                "maxStringLen": 3000,
                "maxArrayValues": 64,
                "maxStructFields": -1
            }
        },
                {
            "name": "issuer",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "",
            "port": 5000,
            "host": "10.208.99.115",
            "showLog": true,
            "trace": "log",
            "logOutput": "rpc",
            "dlvLoadConfig": {
                "followPointers": true,
                "maxVariableRecurse": 1,
                "maxStringLen": 3000,
                "maxArrayValues": 64,
                "maxStructFields": -1
            }
            
        }
    ]
}
```