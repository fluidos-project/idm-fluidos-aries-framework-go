import os
import http.client 
import ssl
import json

# Get environment variables
SSI_CLIENT_HOST = os.environ.get('SSI_CLIENT_HOST')
SSI_CLIENT_PORT = os.environ.get('SSI_CLIENT_PORT')
ISSUER_DID = os.environ.get('ISSUER_DID')
ISSUER_URL = os.environ.get('ISSUER_URL')
if(SSI_CLIENT_HOST is None or SSI_CLIENT_PORT is None or ISSUER_DID is None or ISSUER_URL is None):
    print("Environment variables not set, run with bash run_issuer_setup.sh and check that SSI_CLIENT_HOST, SSI_CLIENT_PORT, ISSUER_DID and ISSUER_URL are set in issuer_config.env")
    exit(1)

# Set up constants
ssi_client_url = "https://"+SSI_CLIENT_HOST+":"+SSI_CLIENT_PORT

# Generate DID
method = "POST"
uri = ssi_client_url+"/poc/newDID"
body = {
        "keys":[
            {
                "keyType":{
                    "type":"Ed25519VerificationKey2018"
                },
                "purpose":"Authentication"
            }
        ],
        "name":"didDevice"
    }
headers={'Accept': 'application/json','Content-Type':'application/json'}
didDocument=None
try:
    conn = http.client.HTTPSConnection(SSI_CLIENT_HOST, SSI_CLIENT_PORT,context=ssl._create_unverified_context())
    conn.request(method, uri, json.dumps(body), headers)
    response = conn.getresponse()
    decoded=response.read().decode()
    jsonresponse = json.loads(decoded) # TODO UMU More checks for errors?
    if "didDoc" not in jsonresponse:
        print("NewDID error")
        print(jsonresponse)
        exit(1)
    didDocument=jsonresponse["didDoc"]
except Exception as e:
    print("Error while requesting newDID method from SSI client")
    print(e)
    exit(1)

# Device enrolment IdM  
method = "POST"
uri = ssi_client_url+"/poc/doDeviceEnrolment"
body = {
        "idProofs":[
            {
                "attrName":"deviceName",
                "attrValue":"Vehicle"
            },
            {
                "attrName":"deviceRole",
                "attrValue":"EmergencyService"
            }
            
        ],
        "url":ISSUER_URL,
        "theirDID":ISSUER_DID
    }
headers={'Accept': 'application/json','Content-Type':'application/json'}

credential=None
credStorageId=None
try:
    conn = http.client.HTTPSConnection(SSI_CLIENT_HOST, SSI_CLIENT_PORT,context=ssl._create_unverified_context())
    conn.request(method, uri, json.dumps(body), headers)
    response = conn.getresponse()
    decoded=response.read().decode()
    jsonresponse = json.loads(decoded) 
    if "credential" not in jsonresponse or "credStorageId" not in jsonresponse:
        print("DoEnrolment error")
        print(jsonresponse)
        exit(1)
    credential=jsonresponse["credential"]
    credential_storage_id=jsonresponse["credStorageId"]
except Exception as e:
    print("Error while requesting doDeviceEnrolment method from SSI client")
    print(e)
    exit(1)

# New Device Trust Agent
#TODO The device's DID/DIDDocument and the Verifiable Credential from enrolment are available at this step, need to start the "newDevice" step through Trust Agent

# ---- Enrolment phase is finished, now we can go into operation phase -----

# Get VP from SSI client (for this PoC, simply a credential)
method = "POST"
uri = ssi_client_url+"/poc/generateVP"
body = {"credId":credential_storage_id}
headers={'Accept': 'application/json','Content-Type':'application/json'}
vp=None
try:
    conn = http.client.HTTPSConnection(SSI_CLIENT_HOST, SSI_CLIENT_PORT,context=ssl._create_unverified_context())
    conn.request(method, uri, json.dumps(body), headers)
    response = conn.getresponse()
    decoded=response.read().decode()
    jsonresponse = json.loads(decoded) # TODO UMU More checks for errors?
    if "credential" not in jsonresponse:
        print("Generate VP error")
        print(jsonresponse)
        exit(1)
    vp=jsonresponse["credential"]
except Exception as e:
    print("Error while requesting generateVP method from SSI client")
    print(e)
    exit(1)

# Make use of service, include VP from last step
#TODO This depends on how it is finally done, but as for last status it could be something simple like publishing an MQTT message that includes the contents of vp as "authentication" and a simple message like "request for green light"
print(vp)