import os
import http.client 
import ssl
import json

# Get environment variables
SSI_CLIENT_HOST = os.environ.get('SSI_CLIENT_HOST')
SSI_CLIENT_PORT = os.environ.get('SSI_CLIENT_PORT')
NUMBER_ATTRIBUTES_PSMS = os.environ.get('NUMBER_ATTRIBUTES_PSMS')
if(SSI_CLIENT_HOST is None or SSI_CLIENT_PORT is None or NUMBER_ATTRIBUTES_PSMS is None):
    print("Environment variables not set, run with bash run_issuer_setup.sh and check that SSI_CLIENT_HOST, SSI_CLIENT_PORT and NUMBER_ATTRIBUTES_PSMS are set in issuer_config.env")
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
            },
                {
                "keyType":{
                    "type":"Bls12381G1Key2022",
                    "nattr": int(NUMBER_ATTRIBUTES_PSMS)
                },
                "purpose":"AssertionMethod"
            },
        ],
        "name":"didDevice"
    }
headers={'Accept': 'application/json','Content-Type':'application/json'}
try:
    conn = http.client.HTTPSConnection(SSI_CLIENT_HOST, SSI_CLIENT_PORT,context=ssl._create_unverified_context())
    conn.request(method, uri, json.dumps(body), headers)
    response = conn.getresponse()
    decoded=response.read().decode()
    jsonresponse = json.loads(decoded) # TODO UMU Check for error response...
    print(jsonresponse)
except Exception as e:
    print("Error while requesting newDID method from SSI client")
    print(e)
    exit(1)

