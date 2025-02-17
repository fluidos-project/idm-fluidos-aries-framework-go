#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from http.server import HTTPServer, BaseHTTPRequestHandler
import ssl
import http.client
import logging
import sys
import json
import configparser
import UtilsPEP
from subprocess import Popen, PIPE
import html
import os
from socketserver import ThreadingMixIn
#import threading
from threading import Thread

import time
import jwt
import re

sys.path.insert(0, './pyCapabilityLib')

import pyCapabilityLib.pyCapabilityEvaluator

logPath="./"
fileName="out"

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(threadName)-12.12s] [%(levelname)-5.5s]  %(message)s",
    handlers=[
        logging.FileHandler("{0}/{1}.log".format(logPath, fileName)),
        logging.StreamHandler(sys.stdout)
    ])

#import numpy as np


#  "ProxyPrivKey": "certs/server-priv-rsa.pem",
#  "ProxyPubKey": "certs/server-pub-rsa.pem",
#  "ProxyCert": "certs/server-public-cert.crt",

#Obtain configuracion from config.cfg file.
cfg = configparser.ConfigParser()  
cfg.read(["./config.cfg"])  
pep_host = cfg.get("GENERAL", "pep_host")
pep_port = int(cfg.get("GENERAL", "pep_port"))

#APIVersion = cfg.get("GENERAL", "APIVersion")

#target_protocol = cfg.get("GENERAL", "target_protocol")
#target_host = cfg.get("GENERAL", "target_host")
#target_port = int(cfg.get("GENERAL", "target_port"))

#blockchain_usevalidation=int(cfg.get("GENERAL", "blockchain_usevalidation"))
#blockchain_protocol = cfg.get("GENERAL", "blockchain_protocol")
#blockchain_host = cfg.get("GENERAL", "blockchain_host")
#blockchain_port = int(cfg.get("GENERAL", "blockchain_port"))

try:
    pep_protocol = str(os.getenv('pep_protocol'))
except Exception as e:
    logging.error(e)
    pep_protocol = "https"

if (str(pep_protocol).upper() == "None".upper()) :
    pep_protocol = "https"

try:
    pep_authtoken_type = str(os.getenv('pep_authtoken_type'))
except Exception as e:
    logging.error(e)
    pep_authtoken_type = "jwt"

if (str(pep_authtoken_type).upper() == "None".upper()) :
    pep_authtoken_type = "jwt"

try:
    node_verifier_protocol = str(os.getenv('node_verifier_protocol'))
except Exception as e:
    logging.error(e)
    node_verifier_protocol = "https"

if (str(node_verifier_protocol).upper() == "None".upper()) :
    node_verifier_protocol = "https"

try:
    producer_node_host = str(os.getenv('producer_node_host'))
except Exception as e:
    logging.error(e)
    producer_node_host = "localhost"

if (str(producer_node_host).upper() == "None".upper()) :
    producer_node_host = "localhost"

producer_node_port = 9082
try:
    producer_node_port = int(os.getenv('producer_node_port'))
except Exception as e:
    logging.error(e)
    producer_node_port = 9082

try:
    node_verifier_post_verifycredential = str(os.getenv('node_verifier_post_verifycredential'))
except Exception as e:
    logging.error(e)
    node_verifier_post_verifycredential = "/fluidos/idm/verifyCredential"

if (str(node_verifier_post_verifycredential).upper() == "None".upper()) :
    node_verifier_post_verifycredential = "/fluidos/idm/verifyCredential"

try:
    node_verifier_post_verifyjwtcontent = str(os.getenv('node_verifier_post_verifyjwtcontent'))
except Exception as e:
    logging.error(e)
    node_verifier_post_verifyjwtcontent = "/fluidos/idm/verifyJWTContent"

if (str(node_verifier_post_verifyjwtcontent).upper() == "None".upper()) :
    node_verifier_post_verifyjwtcontent = "/fluidos/idm/verifyJWTContent"


node_jwt_validatesignature = 1
try:
    node_jwt_validatesignature = int(os.getenv('node_jwt_validatesignature'))
except Exception as e:
    logging.error(e)
    node_jwt_validatesignature = 1

try:
    node_jwt_algorithms = str(os.getenv('node_jwt_algorithms'))
except Exception as e:
    logging.error(e)
    node_jwt_algorithms = "ES256K"

if (str(node_jwt_algorithms).upper() == "None".upper()) :
    node_jwt_algorithms = "ES256K"

try:
    target_protocol = str(os.getenv('target_protocol'))
except Exception as e:
    logging.error(e)
    target_protocol = "http"

if (str(target_protocol).upper() == "None".upper()) :
    target_protocol = "http"

try:
    target_host = str(os.getenv('target_host'))
except Exception as e:
    logging.error(e)
    target_host = "localhost"

if (str(target_host).upper() == "None".upper()) :
    target_host = "localhost"

target_port = 9082
try:
    target_port = int(os.getenv('target_port'))
except Exception as e:
    logging.error(e)
    target_port = 9082

try:
    APIVersion = str(os.getenv('target_API'))
except Exception as e:
    logging.error(e)
    APIVersion = "GenericAPI"

if (str(APIVersion).upper() == "None".upper()) :
    APIVersion = "GenericAPI"

try:
    target2_protocol = str(os.getenv('target2_protocol'))
except Exception as e:
    logging.error(e)
    target2_protocol = "http"

if (str(target2_protocol).upper() == "None".upper()) :
    target2_protocol = "http"

try:
    target2_host = str(os.getenv('target2_host'))
except Exception as e:
    logging.error(e)
    target2_host = "localhost"

if (str(target2_host).upper() == "None".upper()) :
    target2_host = "localhost"

target2_port = 1026
try:
    target2_port = int(os.getenv('target2_port'))
except Exception as e:
    logging.error(e)
    target2_port = 1026

try:
    API2Version = str(os.getenv('target2_API'))
except Exception as e:
    logging.error(e)
    API2Version = "GenericAPI"

if (str(API2Version).upper() == "None".upper()) :
    API2Version = "GenericAPI"

try:
    target2_thingdescription = str(os.getenv('target2_thingdescription'))
except Exception as e:
    logging.error(e)
    target2_thingdescription = "/thingdescription"

if (str(target2_thingdescription).upper() == "None".upper()) :
    target2_thingdescription = "/thingdescription"

blockchain_usevalidation = 0
try:
    blockchain_usevalidation = int(os.getenv('blockchain_usevalidation'))
except Exception as e:
    logging.error(e)
    blockchain_usevalidation = 0

try:
    blockchain_api = str(os.getenv('blockchain_api'))
except Exception as e:
    logging.error(e)
    blockchain_api = "NativeAPI"

if (str(blockchain_api).upper() == "None".upper()) :
    blockchain_api = "NativeAPI"

try:
    blockchain_protocol = str(os.getenv('blockchain_protocol'))
except Exception as e:
    logging.error(e)
    blockchain_protocol = "https"

if (str(blockchain_protocol).upper() == "None".upper()) :
    blockchain_protocol = "https"

try:
    blockchain_host = str(os.getenv('blockchain_host'))
except Exception as e:
    logging.error(e)
    blockchain_host = "localhost"

if (str(blockchain_host).upper() == "None".upper()) :
    blockchain_host = "localhost"

blockchain_port = 8000
try:
    blockchain_port = int(os.getenv('blockchain_port'))
except Exception as e:
    logging.error(e)
    blockchain_port = 8000

try:
    blockchain_get_token = str(os.getenv('blockchain_get_token'))
except Exception as e:
    logging.error(e)
    blockchain_get_token = "/token"

if (str(blockchain_get_token).upper() == "None".upper()) :
    blockchain_get_token = "/token"

chunk_size=int(cfg.get("GENERAL", "chunk_size"))

allApiHeaders=json.loads(cfg.get("GENERAL", "allApiHeaders"))

#Obtain API headers
for m in range(len(allApiHeaders)):
    if(allApiHeaders[m][0].upper()==APIVersion.upper()):
        apiHeaders = allApiHeaders[m][1]
        break

allSeparatorPathAttributeEncriptation=json.loads(cfg.get("GENERAL", "allSeparatorPathAttributeEncriptation"))

#Obtain API separator
for m in range(len(allSeparatorPathAttributeEncriptation)):
    if(allSeparatorPathAttributeEncriptation[m][0].upper()==APIVersion.upper()):
        sPAE = allSeparatorPathAttributeEncriptation[m][1]
        break

rPAE=json.loads(cfg.get("GENERAL", "relativePathAttributeEncriptation"))
noEncryptedKeys = json.loads(cfg.get("GENERAL", "noEncryptedKeys"))

tabLoggingString = "\t\t\t\t\t\t\t"

pep_device = str(os.getenv('PEP_ENDPOINT'))
pep_device_list = pep_device.split(",")

PEPPROXY_CORS_ENABLED = 1
try:
    PEPPROXY_CORS_ENABLED = int(os.getenv('PEPPROXY_CORS_ENABLED'))
except Exception as e:
    logging.error(e)
    PEPPROXY_CORS_ENABLED = 1

logginKPI = cfg.get("GENERAL", "logginKPI")

def CBConnection(method, uri,headers,body = None):

    try:

        #logging.info("")
        #logging.info("")
        #logging.info("********* CBConnection *********")

        uri = UtilsPEP.obtainValidUri(APIVersion,uri)

        #logging.info("CBConnection: Sending the Rquest")

        # send some data
   
        if(target_protocol.upper() == "http".upper() or target_protocol.upper() == "https".upper()):

            state = True

            if (body != None and str(uri).upper().startswith("/v1/subscribeContext".upper()) == False 
                and str(uri).upper().startswith("/v2/subscriptions".upper()) == False 
                and str(uri).upper().startswith("/ngsi-ld/v1/subscriptions".upper()) == False):
            
                state = False

                #logging.info("Original body request: ")
                #logging.info(str(body))
               
                milli_secEP=0
                milli_secEP2=0

                milli_secEP=int(round(time.time() * 1000))

                body, state = UtilsPEP.encryptProcess(APIVersion,method,uri,body,sPAE,rPAE,noEncryptedKeys)
                #body = html.escape(body)
                
                #logging.info("Body request AFTER encryption process: ")
                #logging.info(str(body))

                milli_secEP2=int(round(time.time() * 1000))

                if(logginKPI.upper()=="Y".upper()):
                    #logging.info("")
                    #logging.info("")
                    logging.info("Total(ms) Encrypt process: " + str(milli_secEP2 - milli_secEP))

            if (state):

                if(target_protocol.upper() == "http".upper()):
                    conn = http.client.HTTPConnection(target_host, target_port)
                else:                    
                    #gcontext = ssl.SSLContext()
                    #conn = http.client.HTTPSConnection(target_host,target_port,
                    #                            context=gcontext)
                    # Secure Client SSL Context
                    #gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    #gcontext.check_hostname = True
                    #gcontext.verify_mode = ssl.CERT_REQUIRED
                    #gcontext.load_default_certs()

                    # Unsecure Client SSL Context
                    gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    gcontext.check_hostname = False
                    gcontext.verify_mode = ssl.CERT_NONE

                    conn = http.client.HTTPSConnection(target_host,target_port,
                                                context=gcontext)

                    # Issue detected in https scenarios caused by the presence of the Host header
                    if headers.get("Host"):
                       headers.pop("Host")



                #Deleting "x-auth-token" header, before NGSILD, REQUEST.

                if headers.get("x-auth-token"): 
                    headers.pop("x-auth-token")

                #Python will handle Content-Length header
                if headers.get("Content-Length"): 
                    headers.pop("Content-Length")

                #logging.info("BROKER REQUEST:\n" + 
                #tabLoggingString + "- Host: " + target_host + "\n" + 
                #tabLoggingString + "- Port: " + str(target_port) + "\n" + 
                #tabLoggingString + "- Method: " + method + "\n" + 
                #tabLoggingString + "- URI: " + uri + "\n" + 
                #tabLoggingString + "- Headers: " + str(headers) + "\n" + 
                #tabLoggingString + "- Body: " + str(body))

                conn.request(method, uri, body, headers)
                response = conn.getresponse()

                #logging.info("CBConnection - RESPONSE")
                logging.info("BROKER RESPONSE CODE: "      + str(response.code))
                #logging.info("Headers: ")
                #logging.info(response.headers)
            else:
                return -1
        
        #logging.info("********* CBConnection - END ********* ")

        return response

    except Exception as e:
        logging.error(e)
        return -1

def CB2Connection(method, uri,headers,body = None):

    try:

        #logging.info("")
        #logging.info("")
        #logging.info("********* CB2Connection *********")

        uri = UtilsPEP.obtainValidUri(APIVersion,uri)

        #logging.info("CB2Connection: Sending the Rquest")

        # send some data
   
        if(target2_protocol.upper() == "http".upper() or target2_protocol.upper() == "https".upper()):

            state = True

            if (body != None and str(uri).upper().startswith("/v1/subscribeContext".upper()) == False 
                and str(uri).upper().startswith("/v2/subscriptions".upper()) == False 
                and str(uri).upper().startswith("/ngsi-ld/v1/subscriptions".upper()) == False):
            
                state = False

                #logging.info("Original body request: ")
                #logging.info(str(body))
               
                milli_secEP=0
                milli_secEP2=0

                milli_secEP=int(round(time.time() * 1000))

                body, state = UtilsPEP.encryptProcess(APIVersion,method,uri,body,sPAE,rPAE,noEncryptedKeys)
                #body = html.escape(body)
                
                #logging.info("Body request AFTER encryption process: ")
                #logging.info(str(body))

                milli_secEP2=int(round(time.time() * 1000))

                if(logginKPI.upper()=="Y".upper()):
                    #logging.info("")
                    #logging.info("")
                    logging.info("Total(ms) Encrypt process: " + str(milli_secEP2 - milli_secEP))

            if (state):

                if(target2_protocol.upper() == "http".upper()):
                    conn = http.client.HTTPConnection(target2_host, target2_port)
                else:                    
                    #gcontext = ssl.SSLContext()
                    #conn = http.client.HTTPSConnection(target2_host,target2_port,
                    #                            context=gcontext)
                    # Secure Client SSL Context
                    #gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    #gcontext.check_hostname = True
                    #gcontext.verify_mode = ssl.CERT_REQUIRED
                    #gcontext.load_default_certs()

                    # Unsecure Client SSL Context
                    gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    gcontext.check_hostname = False
                    gcontext.verify_mode = ssl.CERT_NONE

                    conn = http.client.HTTPSConnection(target2_host,target2_port,
                                                context=gcontext)

                    # Issue detected in https scenarios caused by the presence of the Host header
                    if headers.get("Host"):
                       headers.pop("Host")


                #Deleting "x-auth-token" header, before NGSILD, REQUEST.

                if headers.get("x-auth-token"): 
                    headers.pop("x-auth-token")

                #Python will handle Content-Length header
                if headers.get("Content-Length"): 
                    headers.pop("Content-Length")

                #logging.info("BROKER REQUEST:\n" + 
                #tabLoggingString + "- Host: " + target2_host + "\n" + 
                #tabLoggingString + "- Port: " + str(target2_port) + "\n" + 
                #tabLoggingString + "- Method: " + method + "\n" + 
                #tabLoggingString + "- URI: " + uri + "\n" + 
                #tabLoggingString + "- Headers: " + str(headers) + "\n" + 
                #tabLoggingString + "- Body: " + str(body))

                conn.request(method, uri, body, headers)
                response = conn.getresponse()

                #logging.info("CB2Connection - RESPONSE")
                logging.info("BROKER RESPONSE CODE: "      + str(response.code))
                #logging.info("Headers: ")
                #logging.info(response.headers)
            else:
                return -1
        
        #logging.info("********* CB2Connection - END ********* ")

        return response

    except Exception as e:
        logging.error(e)
        return -1

def BlockChainConnection(method, uri, headers = {}, body = None):

    try:

        #logging.info("")
        #logging.info("")
        #logging.info("********* BlockChainConnection *********")

        # send some data
        
        if(blockchain_protocol.upper() == "http".upper() or blockchain_protocol.upper() == "https".upper()):

            if(blockchain_protocol.upper() == "http".upper()):
                conn = http.client.HTTPConnection(blockchain_host, blockchain_port)
            else:                    
                
                #gcontext = ssl.SSLContext()
                #conn = http.client.HTTPSConnection(blockchain_host,blockchain_port,
                #                                context=gcontext)
                # Secure Client SSL Context
                #gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                #gcontext.check_hostname = True
                #gcontext.verify_mode = ssl.CERT_REQUIRED
                #gcontext.load_default_certs()
                # Unsecure Client SSL Context
                gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                gcontext.check_hostname = False
                gcontext.verify_mode = ssl.CERT_NONE
                conn = http.client.HTTPSConnection(blockchain_host,blockchain_port,
                                            context=gcontext)
                # Issue detected in https scenarios caused by the presence of the Host header
                if headers.get("Host"):
                   headers.pop("Host")

                

            #logging.info("BLOCKCHAIN REQUEST:\n" +
            #             tabLoggingString + "- Host: " + blockchain_host + "\n" + 
            #             tabLoggingString + "- Port: " + str(blockchain_port) + "\n" + 
            #             tabLoggingString + "- Method: " + method + "\n" + 
            #             tabLoggingString + "- URI: " + uri + "\n" + 
            #             tabLoggingString + "- Headers: " + str(headers) + "\n" + 
            #             tabLoggingString + "- Body: " + str(body))

            conn.request(method, uri, body, headers)

            response = conn.getresponse()

            #logging.info("BlockChainConnection - RESPONSE")
            logging.info(" SUCCESS : BlockChain response - code: "      + str(response.code))
            #logging.info("Headers: ")
            #logging.info(response.headers)
        
        #logging.info("********* BlockChainConnection - END *********")

        return response

    except Exception as e:
        logging.error(e)
        return -1

def Node_Verifier_Connection(method, uri,headers,body = None):

    try:

        #logging.info("")
        #logging.info("")
        #logging.info("********* Node_Verifier_Connection *********")

        uri = UtilsPEP.obtainValidUri(APIVersion,uri)

        #logging.info("Node_Verifier_Connection: Sending the Rquest")

        # send some data

        if(target_protocol.upper() == "http".upper() or target_protocol.upper() == "https".upper()):

                if(node_verifier_protocol.upper() == "http".upper()):
                    conn = http.client.HTTPConnection(producer_node_host, producer_node_port)
                else:
                    #gcontext = ssl.SSLContext()
                    #conn = http.client.HTTPSConnection(producer_node_host,producer_node_port,
                    #                            context=gcontext)
                    # Secure Client SSL Context
                    #gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    #gcontext.check_hostname = True
                    #gcontext.verify_mode = ssl.CERT_REQUIRED
                    #gcontext.load_default_certs()

                    # Unsecure Client SSL Context
                    gcontext = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                    gcontext.check_hostname = False
                    gcontext.verify_mode = ssl.CERT_NONE

                    conn = http.client.HTTPSConnection(producer_node_host,producer_node_port,
                                                context=gcontext)

                    # Issue detected in https scenarios caused by the presence of the Host header
                    if headers.get("Host"):
                       headers.pop("Host")

                #Python will handle Content-Length header
                if headers.get("Content-Length"): 
                    headers.pop("Content-Length")

                #tabLoggingString + "- Host: " + target_host + "\n" +
                #tabLoggingString + "- Port: " + str(target_port) + "\n" +
                #tabLoggingString + "- Method: " + method + "\n" +
                #tabLoggingString + "- URI: " + uri + "\n" +
                #tabLoggingString + "- Headers: " + str(headers) + "\n" +
                #tabLoggingString + "- Body: " + str(body))

                conn.request(method, uri, body, headers)
                response = conn.getresponse()

                #logging.info("Node_Verifier_Connection - RESPONSE")
                logging.info("NODE VERIFIER RESPONSE CODE: "      + str(response.code))
                #logging.info("Headers: ")
                #logging.info(response.headers)

                #logging.info("********* Node_Verifier_Connection - END ********* ")

                return response

        return -1

    except Exception as e:
        logging.error(e)
        return -1

def getstatusoutput(command):
    process = Popen(command, stdout=PIPE,stderr=PIPE)
    out, err = process.communicate()

    return (process.returncode, out)

def obtainRequestHeaders(RequestHeaders):

    headers = dict()

    content_length = 0

    try:
        # We get the headers
        
        #logging.info ("********* HEADERS BEFORE obtainRequestHeaders *********")
        #logging.info (RequestHeaders)
        
        for key in RequestHeaders:
            #logging.info("Procesando: " + str(key) + ":" + str(RequestHeaders[key]))

            #value_index=-1
            #
            #try:
            #    #To find only admittable headers from request previously configured in config.cfg file.
            #    value_index = apiHeaders.index(key.lower())
            #
            #    #Supporting all headers
            #    if (value_index == -1 ):
            #        value_index = apiHeaders.index("*")
            #
            #except:
            #    value_index = -1
            #
            ##If the header key was found, it will be considered after.
            #if (value_index > -1 ):
            #
            #    #logging.info("Incluido: " + str(key) + ":" + str(RequestHeaders[key]))
            #
            #    headers[key] = RequestHeaders[key]

            headers[key] = RequestHeaders[key]

            if(key.upper()=="Content-Length".upper()):
                content_length = int(RequestHeaders[key])

    except Exception as e:
        logging.error(e)

        headers["Error"] = str(e)

    #logging.info ("********* HEADERS AFTER obtainRequestHeaders *********")
    #logging.info (headers)

    return headers, content_length

#Validates the correct format of the String before parsing it to JSON format
def validateJSONString(response):

    validated = response.replace('\'', '\"')
    validated = validated.replace('\n', '')
    #Deletes the first and last characters of the String if they are double quotes
    if (validated[0] == '"' and validated[len(validated)-1] == '"'):
        validated = validated[1:-1]

    return validated

def validationToken(headers,method,uri,body = None):

    if (pep_authtoken_type.upper() == "capability".upper()):

        validationCapabilityToken = False
        validationBlockChain = False
        validationResult = False

        isRevoked = False
        strRevoked = ""

        outTypeProcessed = ""

        #print("uri: " + uri)

        milli_secValCT=0
        milli_secValCT2=0

        milli_secBC=0
        milli_secBC2=0

        milli_secValCTT=0
        milli_secValCTT2=0

        try:

            milli_secValCTT=int(round(time.time() * 1000))

            for key in headers:

                if(key.upper()=="x-auth-token".upper()):

                    headersStr = json.dumps(headers)

                    # DEPRECATED NOT USEFULL AND FAILS WITH FORM-DATA BODIES.
                    #if (body == None):
                    #    bodyStr = "{}"
                    #else:
                    #    bodyStr = body.decode('utf8').replace("'", '"').replace("\t", "").replace("\n", "")
                    #    #bodyStr = body.decode('utf8').replace("'", '"')
                    bodyStr = "{}"
                    
                    #print(type(str(method)))
                    #print(type(str(uri)))
                    #print(type(headersStr))
                    #print(type(bodyStr))
                    #print(type(str(headers[key])))

                    #print(str(method))
                    #print(str(uri))
                    #print(headersStr)
                    #print(bodyStr)
                    #print(str(headers[key]))

                    
                    ##Validating token (v1)
                    ##Observation: str(uri).replace("&",";") --> for PDP error: "The reference to entity "***" must end with the ';' delimiter.""
                    #codeType, outType = getstatusoutput(["java","-jar","CapabilityEvaluator_old.jar",
                    ##str(pep_device),
                    #str(method),
                    #str(uri).replace("&",";"),
                    #headersStr, # "{}", #headers
                    #bodyStr,
                    #str(headers[key])])
                    #
                    #logging.info("codeType_v0: " + str(codeType))
                    #logging.info("outType_v0: " + str(outType))
            
                    milli_secValCT=int(round(time.time() * 1000))

                    try:

                        outTypeProcessed = pyCapabilityLib.pyCapabilityEvaluator.pyCapabilityEvaluator(str(pep_device), str(method), 
                            str(uri).replace("&",";"), str(headers[key]), "", headersStr, bodyStr)

                        outTypeProcessed = outTypeProcessed.replace("'", '"').replace("CODE: ","").replace("\n", "")

                    except Exception as e:
                        logging.error(e)

                    if (outTypeProcessed.upper()=="AUTHORIZED".upper()):

                        validationCapabilityToken = True

                        if (blockchain_usevalidation == 1):

                            milli_secBC=int(round(time.time() * 1000))

                            capabilityTokenId = json.loads(headers[key])["id"]

                            #Send requests to blockchain to obtain if Capability Token id exists and its 
                            resultGet = BlockChainConnection("GET", blockchain_get_token + "urn:ngsi-ld:blockchaincaptoken:" + capabilityTokenId, {}, None)
                            #resultGet = BlockChainConnection("GET", "/token/43oi76utr62o8fuad5v79duhia", {}, None)


                            errorBlockChainConnectionGET = False
                            try:
                                if(resultGet==-1):
                                    errorBlockChainConnectionGET = True
                            except:
                                errorBlockChainConnectionGET = False

                            if (errorBlockChainConnectionGET==False):

                                strDataGET = resultGet.read(chunk_size).decode('utf8')

                                if (resultGet.code!=200):
                                    
                                    #This request is sent by Capability Manager component.
                                    ##If Capability Token id don't exist, send requests to register it.
                                    #resultPost = BlockChainConnection("POST", "/token/register", {}, "{\"id\":\"" + capabilityTokenId + "\"}")
                                    #
                                    #errorBlockChainConnectionPOST = False
                                    #try:
                                    #    if(resultPost==-1):
                                    #        errorBlockChainConnectionPOST = True
                                    #except:
                                    #    errorBlockChainConnectionPOST = False
                                    #
                                    #if (errorBlockChainConnectionPOST==False):
                                    #
                                    #    strDataPOST = resultPost.read(chunk_size).decode('utf8')
                                    #
                                    #    if (resultPost.code==201):
                                    #        validationBlockChain = True
                                    #    else:
                                    #        headers["Error"] = str("Can't confirm validity state of the registered token.(2)")
                                    #        validationBlockChain = False
                                    #else:
                                    #    headers["Error"] = str("Can't confirm validity state of the registered token.(1)")
                                    #    validationBlockChain = False

                                    headers["Error"] = str("Can't confirm validity state of the registered token.")
                                    validationBlockChain = False
                                else:
                                    if (blockchain_api == "NativeAPI"):
                                        stateValue = json.loads(validateJSONString(strDataGET))["state"]

                                        if (str(stateValue)=="1"):
                                            validationBlockChain = True
                                        else:
                                            isRevoked = True
                                            validationBlockChain = False

                                    else:
                                        if (blockchain_api == "NGSIv2" or blockchain_api == "NGSI-LD"):
                                            stateValue = json.loads(validateJSONString(strDataGET))["state"]["value"]

                                            if (str(stateValue)=="1"):
                                                validationBlockChain = True
                                            else:
                                                isRevoked = True
                                                validationBlockChain = False

                            else:
                                headers["Error"] = str("Can't confirm validity state of the registered token.(0)")
                                validationBlockChain = False

                            milli_secBC2=int(round(time.time() * 1000))

                        else:
                            validationBlockChain = True

                    break

        except Exception as e:
            logging.error(e)

            headers["Error"] = str(e)

        if (validationCapabilityToken and validationBlockChain):
            validationResult = True

        if (isRevoked):
            strRevoked = " - Code: REVOKED"
        else:
            if (blockchain_usevalidation == 0):
                strRevoked = " - Validation no configured (always true)."

        if headers.get("Error"):
            logging.error("Error: " + str(headers["Error"]))

        #logging.info("CAPABILITY TOKEN'S VALIDATION:\n" +
        #            tabLoggingString + "1) Request... - Result: " + str(validationCapabilityToken) + " - Code: " + str(outTypeProcessed) + "\n" +
        #            tabLoggingString + "2) BlockChain - Result: " + str(validationBlockChain) + strRevoked + "\n" +
        #            tabLoggingString + "SUCCESS : Capability token's validation response - " + str(validationResult).upper()
        #)

        milli_secValCTT2=int(round(time.time() * 1000))

        if(logginKPI.upper()=="Y".upper()):
            #logging.info("")
            #logging.info("")
            #logging.info("Total(ms) Validacion CT: " + str(milli_secValCT2 - milli_secValCT))
            #logging.info("Total(ms) BlockChain: " + str(milli_secBC2 - milli_secBC))
            logging.info("Total(ms) Validacion process: " + str(milli_secValCTT2 - milli_secValCTT))

        return validationResult
    
    if (pep_authtoken_type.upper() == "jwt".upper()):

        milli_secValCTT=0
        milli_secValCTT2=0
        validationResult = False
        milli_secValCTT=int(round(time.time() * 1000))

        try:

            for key in headers:

                if(key.upper()=="x-auth-token".upper()):

                    jwtToken = str(headers[key])

                    isValidJWT = True
                    body = {
                            "jwt": jwtToken
                        }
                    json_body = json.dumps(body)
                    result = Node_Verifier_Connection("POST", "/fluidos/idm/verifyJWTContent", headers={}, body=json_body)
                    if result != -1:
                        if result.code == 200:
                            #log
                            response_body = result.read()
                            try:
                                response_data = json.loads(response_body)
                                if response_data.get("verified") == True:
                                    logging.info("Verified is True")
                                else:
                                    logging.error("Verified is False or not present")
                                    isValidJWT = False
                            except json.JSONDecodeError:
                                logging.error("error decoding JSON response")
                                isValidJWT = False
                        else:
                            logging.error(f"Failed with status: {result.code}")
                            isValidJWT = False
                    else:
                        logging.error("Error: No response received from Node Verifier.")
                        isValidJWT = False

                    if (isValidJWT):
                        decoded = {}
                        decoded = jwt.decode(jwtToken, options={"verify_signature": False})
                        if method.upper() in decoded["method"].upper().split('/'):

                            currentTimestamp = int(round(time.time()))

                            if (currentTimestamp >= decoded["exp"]):
                                logging.error("(Error: " + "TOKEN : Expiration date exceeded.")
                            else:

                                for i in range(len(pep_device_list)):

                                    pattern1 = re.compile("^" + decoded["resource"].split('?')[0] + "$")
                                    matcher1 = pattern1.match(pep_device_list[i]+uri.split('?')[0])

                                    if ( matcher1 == None):

                                        pattern2 = re.compile("^" +  decoded["resource"] + "$")
                                        matcher2 = pattern2.match("\\Q" + pep_device_list[i]+uri + "\\E" )

                                        if ( matcher2 != None): 
                                            validationResult = True
                                            break
                                    else:
                                        validationResult = True
                                        break

                                    #if (decoded["aud"].upper() == pep_device_list[i].upper()+uri.upper()):
                                    #    validationResult = True
                                    #    break
                                
                                if (validationResult == False):
                                    logging.error ("Error: " + "device/resource is different: \n"
                                        "\t\t - HTTP: ('" + str(pep_device) + "', '" + str(uri) + "')\n"
                                        "\t\t - TOKEN: ('" + str(decoded["resource"]) + "')"
                                    )                         
                        
                        else:
                            logging.error ("Error: " + "Method is different: \n"
                                "\t\t - HTTP: ('" + str(method) + "')\n"
                                "\t\t - TOKEN: ('" + str(decoded["method"]) + "')"
                            )

                        #index = -1
                        #try: 
                        #    index = pep_device.index(decoded["device"])
                        #except ValueError as e: 
                        #    index = -1
                        #    
                        ##if (decoded["device"].upper() == pep_device.upper() and 
                        #if ( index >= 0 and 
                        #    decoded["method"].upper() == method.upper() and 
                        #    decoded["aud"].upper() == uri.upper()):
                        #
                        #    currentTimestamp = int(round(time.time()))
                        #
                        #    #logging.info("currentTimestamp: " + str(currentTimestamp))
                        #    #logging.info(type(currentTimestamp))
                        #    #logging.info("decoded['exp']: " + str(decoded["exp"]))
                        #    #logging.info(type(decoded["exp"]))
                        #
                        #    if (currentTimestamp < decoded["exp"]):
                        #        validationResult = True
                        #    else:
                        #        logging.error("(Error: " + "TOKEN : Expiration date exceeded.")
                        #else:
                        #    logging.error ("Error: " + "Method/device/resource are differents: \n"
                        #        "\t\t - HTTP: ('" + str(method) + "', '" + str(pep_device) + "', '" + str(uri) + "')\n"
                        #        "\t\t - TOKEN: ('" + str(decoded["method"]) + "', '" + str(decoded["device"]) + "', '" 
                        #        + str(decoded["aud"]) + "')"
                        #    )


        except Exception as e:
            logging.error(e)

        milli_secValCTT2=int(round(time.time() * 1000))
        if(logginKPI.upper()=="Y".upper()):
            logging.info("Total(ms) Validacion process: " + str(milli_secValCTT2 - milli_secValCTT))

        #logging.info("validationResult: " + str(validationResult))

        return validationResult  
    
def obtainResponseHeaders(ResponseHeaders):

    #logging.info ("********* HEADERS BEFORE obtainResponseHeaders *********")
    #logging.info (ResponseHeaders)


    headers = dict()

    target_chunkedResponse = False

    try:
        for key in ResponseHeaders:
            #logging.info(str(key) + ":" + str(ResponseHeaders[key]))

            if(key.upper()=="Transfer-Encoding".upper() and ResponseHeaders[key].upper()=="chunked".upper()):
                target_chunkedResponse = True

            if(key.upper()!="Date".upper() and key.upper()!="Server".upper()):
                headers[key] = ResponseHeaders[key]
                
    except Exception as e:

        headers["Error"] = str(e)

    #logging.info ("********* HEADERS AFTER obtainResponseHeaders *********")
    #logging.info (headers)

    return  headers, target_chunkedResponse

def loggingPEPRequest(req):
    #logging.info("")
    #logging.info (" ********* PEP-REQUEST ********* ")
    #logging.info(req.address_string())
    #logging.info(req.date_time_string())
    #logging.info(req.path)
    #logging.info(req.protocol_version)
    #logging.info(req.raw_requestline)
    logging.info("******* PEP-REQUEST : " + req.address_string() + " - " + str(req.raw_requestline) + " *******")  

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):

    def do_OPTIONS(self):
        self.send_response(200)
        self.send_headers()
        self.end_headers()
        
    def send_headers (self):
        if (PEPPROXY_CORS_ENABLED == 1):
	        #self.send_header('Access-Control-Allow-Credentials', 'true')
            self.send_header('Access-Control-Allow-Origin', '*')
            #self.send_header('Access-Control-Allow-Methods', '*')
            self.send_header('Access-Control-Allow-Methods', 'GET,POST,PATCH,PUT,DELETE,OPTIONS')
            self.send_header('Access-Control-Allow-Headers', '*')
            self.send_header('Access-Control-Expose-Headers', '*')
            #self.send_header('Cache-Control','no-store,no-cache,must-revalidate')
            self.send_header('Cache-Control','no-cache,private,no-store,must-revalidate,max-stale=0,post-check=0,pre-check=0')

    def do_HandleError(self,method,code,title,details):
        messageBody = UtilsPEP.obtainErrorResponseBody(APIVersion,method,code,title,details)
        #code = UtilsPEP.obtainErrorResponseCode(APIVersion,method)
        #self.send_response(code)

        self.send_response(code)

        errorHeaders,chunkedResponse = UtilsPEP.obtainErrorResponseHeaders(APIVersion)
        for key in errorHeaders:
            self.send_header(key, errorHeaders[key])
        self.send_headers()
        self.end_headers() 
        #data = json.dumps(message).encode()
        data = json.dumps(messageBody).encode()
        if(chunkedResponse):
            self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
        else:
            self.wfile.write(data)

        self.close_connection
        logging.error("Response issued (do_HandleError): Communication is closed.")

    def do_GET(self):

        target_chunkedResponse=False
        if (self.path=="" or self.path=="/"):
            #To CI/CD.
            self.send_response(200)
            self.send_headers()
            self.end_headers()
            self.close_connection
        else:
            try:
                loggingPEPRequest(self)

                milli_secTotalReq=int(round(time.time() * 1000))

                headers,content_length = obtainRequestHeaders(self.headers)

                try:
                    #To find only admittable headers from request previously configured in config.cfg file.
                    value_index = headers.index("Error")
                except:
                    value_index = -1

                testSupported = UtilsPEP.validateNotSupportedMethodPath(APIVersion,self.command,self.path)

                if (value_index != -1):
                    logging.error("Error: " + str(headers["Error"]))
                    SimpleHTTPRequestHandler.do_HandleError(self,self.command,400,"Bad Request","Error obtaining headers.")

                else:

                    if (testSupported == False):
                        logging.error("Error: " + str(headers["Error"]))
                        SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method/path.")

                    else:

                        validation = True

                        validation = validationToken(headers,self.command,self.path)

                        if (validation == False):
                            SimpleHTTPRequestHandler.do_HandleError(self,self.command,401,"Unauthorized","The token is missing or invalid.")
                        else:
                            # We are sending this to the CB
                            result = CBConnection(self.command, self.path, headers, None)

                            errorConnection = False
                            try:
                                if(result==-1):
                                    errorConnection = True
                            except:
                                errorConnection = False

                            if(errorConnection):
                                SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
                            else:

                                # We send back the response to the client
                                self.send_response(result.code)

                                headersResponse, target_chunkedResponse = obtainResponseHeaders(result.headers)

                                #logging.info(" ******* Sending Headers back to client ******* ")
                                for key in headersResponse:
                                    if (PEPPROXY_CORS_ENABLED == 1 and 
                                       key in ["Access-Control-Allow-Origin","Access-Control-Allow-Methods","Access-Control-Allow-Headers",
                                       "Access-Control-Expose-Headers","Cache-Control"]):
                                        continue
                                    else:
                                        self.send_header(key, headersResponse[key])
                                    #self.send_header(key, headersResponse[key])

                                self.send_headers()
                                self.end_headers()

                                #logging.info("Sending the Body back to client")

                                # Link to resolve Transfer-Encoding chunked cases
                                # https://docs.amazonaws.cn/en_us/polly/latest/dg/example-Python-server-code.html

                                while True:
                                    data = result.read(chunk_size)
                            
                                    if (target_chunkedResponse):
                                        self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
                                    else:
                                        self.wfile.write(data)

                                    if data is None or len(data) == 0:
                                        break

                                if (target_chunkedResponse):
                                    self.wfile.flush()

                            self.close_connection
                            logging.info("Response issued: Communication is closed.")

                            milli_secTotalReq2=int(round(time.time() * 1000))

                            if(logginKPI.upper()=="Y".upper()):
                                logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

            except Exception as e:
                logging.error(str(e))
                SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
        

    def do_POST(self):
        target_chunkedResponse=False

        try:

            loggingPEPRequest(self)

            milli_secTotalReq=int(round(time.time() * 1000))

            headers,content_length = obtainRequestHeaders(self.headers)

            try:
                #To find only admittable headers from request previously configured in config.cfg file.
                value_index = headers.index("Error")
            except:
                value_index = -1

            testSupported = UtilsPEP.validateNotSupportedMethodPath(APIVersion,self.command,self.path)

            if (value_index != -1):
                logging.error("Error: " + str(headers["Error"]))
                SimpleHTTPRequestHandler.do_HandleError(self,self.command,400,"Bad Request","Error obtaining headers.")

            else:

                if (testSupported == False):
                    logging.error("Error: " + str(headers["Error"]))
                    SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method/path.")

                else:

                    #logging.info (" ********* OBTAIN BODY ********* ")
                    # We get the body
                    if (content_length>0):
                        #logging.info ("-------- self.rfile.read(content_length) -------")
                        post_body   = self.rfile.read(content_length)
                    else:
                        #logging.info ("-------- Lanzo self.rfile.read() -------")
                        post_body = b''

                    #logging.info(post_body)

                    validation = True
                    #NODE VERIFIER REQUEST has nothing to validate
                    if (self.path.upper() != node_verifier_post_verifycredential.upper() and self.path.upper() != node_verifier_post_verifyjwtcontent.upper()):
                        validation = validationToken(headers,self.command,self.path,post_body)

                    if (validation == False):
                        SimpleHTTPRequestHandler.do_HandleError(self,self.command,401,"Unauthorized","The token is missing or invalid.")

                    else:

                        if (self.path.upper() != node_verifier_post_verifycredential.upper() and self.path.upper() != node_verifier_post_verifyjwtcontent.upper()):
                            if (self.path.upper() != target2_thingdescription.upper()):
                                # We are sending this to the CB
                                result = CBConnection(self.command, self.path, headers, post_body)
                            else:
                                # We are sending this to the CB
                                result = CB2Connection(self.command, self.path, headers, post_body)
                        else:
                            #NODE VERIFIER REQUEST
                            #headers["device"] = pep_device
                            result = Node_Verifier_Connection(self.command, self.path,headers, post_body)

                        errorConnection = False
                        try:
                            if(result==-1):
                                errorConnection = True
                        except:
                            errorConnection = False

                        if(errorConnection):
                            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
                        else:

                            # We send back the response to the client
                            self.send_response(result.code)

                            headersResponse, target_chunkedResponse = obtainResponseHeaders(result.headers)

                            #logging.info(" ******* Sending Headers back to client ******* ")
                            for key in headersResponse:
                                if (PEPPROXY_CORS_ENABLED == 1 and 
                                   key in ["Access-Control-Allow-Origin","Access-Control-Allow-Methods","Access-Control-Allow-Headers",
                                   "Access-Control-Expose-Headers","Cache-Control"]):
                                    continue
                                else:
                                    self.send_header(key, headersResponse[key])
                                #self.send_header(key, headersResponse[key])

                            self.send_headers()
                            self.end_headers()

                            #logging.info("Sending the Body back to client")

                            # Link to resolve Transfer-Encoding chunked cases
                            # https://docs.amazonaws.cn/en_us/polly/latest/dg/example-Python-server-code.html

                            while True:
                                data = result.read(chunk_size)
                        
                                if (target_chunkedResponse):
                                    self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
                                else:
                                    self.wfile.write(data)
                                
                                if data is None or len(data) == 0:
                                    break
                
                            if (target_chunkedResponse):
                                self.wfile.flush()

                        self.close_connection
                        logging.info("Response issued: Communication is closed.")

                        milli_secTotalReq2=int(round(time.time() * 1000))

                        if(logginKPI.upper()=="Y".upper()):
                            logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

        except Exception as e:
            logging.error(str(e))
            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")

    def do_DELETE(self):
        target_chunkedResponse=False

        try:
            
            loggingPEPRequest(self)

            milli_secTotalReq=int(round(time.time() * 1000))

            headers,content_length = obtainRequestHeaders(self.headers)

            try:
                #To find only admittable headers from request previously configured in config.cfg file.
                value_index = headers.index("Error")
            except:
                value_index = -1

            testSupported = UtilsPEP.validateNotSupportedMethodPath(APIVersion,self.command,self.path)

            if (value_index != -1):
                logging.error("Error: " + str(headers["Error"]))
                SimpleHTTPRequestHandler.do_HandleError(self,self.command,400,"Bad Request","Error obtaining headers.")

            else:

                if (testSupported == False):
                    logging.error("Error: " + str(headers["Error"]))
                    SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method/path.")

                else:

                    validation = validationToken(headers,self.command,self.path)

                    if (validation == False):
                        SimpleHTTPRequestHandler.do_HandleError(self,self.command,401,"Unauthorized","The token is missing or invalid.")

                    else:
                    
                        # We are sending this to the CB
                        result = CBConnection(self.command, self.path, headers, None)

                        errorCBConnection = False
                        try:
                            if(result==-1):
                                errorCBConnection = True
                        except:
                            errorCBConnection = False

                        if(errorCBConnection):
                            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
                        else:        
                            # We send back the response to the client
                            self.send_response(result.code)

                            headersResponse, target_chunkedResponse = obtainResponseHeaders(result.headers)

                            #logging.info(" ******* Sending Headers back to client ******* ")
                            for key in headersResponse:
                                if (PEPPROXY_CORS_ENABLED == 1 and 
                                   key in ["Access-Control-Allow-Origin","Access-Control-Allow-Methods","Access-Control-Allow-Headers",
                                   "Access-Control-Expose-Headers","Cache-Control"]):
                                    continue
                                else:
                                    self.send_header(key, headersResponse[key])
                                #self.send_header(key, headersResponse[key])

                            self.send_headers()
                            self.end_headers()

                            #logging.info("Sending the Body back to client")

                            # Link to resolve Transfer-Encoding chunked cases
                            # https://docs.amazonaws.cn/en_us/polly/latest/dg/example-Python-server-code.html

                            while True:
                                data = result.read(chunk_size)
                        
                                if (target_chunkedResponse):
                                    self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
                                else:
                                    self.wfile.write(data)
                                
                                if data is None or len(data) == 0:
                                    break
                    
                            if (target_chunkedResponse):
                                self.wfile.flush()

                        self.close_connection
                        logging.info("Response issued: Communication is closed.")

                        milli_secTotalReq2=int(round(time.time() * 1000))

                        if(logginKPI.upper()=="Y".upper()):
                            logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

        except Exception as e:
            logging.error(str(e))
            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")

    def do_PATCH(self):
        target_chunkedResponse=False

        try:

            loggingPEPRequest(self)

            milli_secTotalReq=int(round(time.time() * 1000))

            headers,content_length = obtainRequestHeaders(self.headers)

            try:
                #To find only admittable headers from request previously configured in config.cfg file.
                value_index = headers.index("Error")
            except:
                value_index = -1

            testSupported = UtilsPEP.validateNotSupportedMethodPath(APIVersion,self.command,self.path)


            if (value_index != -1):
                logging.error("Error: " + str(headers["Error"]))
                SimpleHTTPRequestHandler.do_HandleError(self,self.command,400,"Bad Request","Error obtaining headers.")

            else:

                if (testSupported == False):
                    logging.error("Error: " + str(headers["Error"]))
                    SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method/path.")

                else:

                    #logging.info (" ********* OBTAIN BODY ********* ")
                    # We get the body
                    if (content_length>0):
                        #logging.info ("-------- self.rfile.read(content_length) -------")
                        patch_body   = self.rfile.read(content_length)
                    else:
                        #logging.info ("-------- Lanzo self.rfile.read() -------")
                        patch_body   = self.rfile.read()

                    #logging.info(patch_body)

                    validation = validationToken(headers,self.command,self.path,patch_body)

                    if (validation == False):
                        SimpleHTTPRequestHandler.do_HandleError(self,self.command,401,"Unauthorized","The token is missing or invalid.")

                    else:
                        # We are sending this to the CB
                        result = CBConnection(self.command, self.path,headers, patch_body)

                        errorCBConnection = False
                        try:
                            if(result==-1):
                                errorCBConnection = True
                        except:
                            errorCBConnection = False

                        if(errorCBConnection):
                            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
                        else:        
                            # We send back the response to the client
                            self.send_response(result.code)

                            headersResponse, target_chunkedResponse = obtainResponseHeaders(result.headers)

                            #logging.info(" ******* Sending Headers back to client ******* ")
                            for key in headersResponse:
                                if (PEPPROXY_CORS_ENABLED == 1 and 
                                   key in ["Access-Control-Allow-Origin","Access-Control-Allow-Methods","Access-Control-Allow-Headers",
                                   "Access-Control-Expose-Headers","Cache-Control"]):
                                    continue
                                else:
                                    self.send_header(key, headersResponse[key])
                                #self.send_header(key, headersResponse[key])

                            self.send_headers()
                            self.end_headers()

                            #logging.info("Sending the Body back to client")

                            # Link to resolve Transfer-Encoding chunked cases
                            # https://docs.amazonaws.cn/en_us/polly/latest/dg/example-Python-server-code.html

                            while True:
                                data = result.read(chunk_size)

                                if (target_chunkedResponse):
                                    self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
                                else:
                                    self.wfile.write(data)
                                
                                if data is None or len(data) == 0:
                                    break
                    
                            if (target_chunkedResponse):
                                self.wfile.flush()

                        self.close_connection
                        logging.info("Response issued: Communication is closed.")

                        milli_secTotalReq2=int(round(time.time() * 1000))

                        if(logginKPI.upper()=="Y".upper()):
                            logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

        except Exception as e:
            logging.error(str(e))
            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")            
            

    ##Actually not suppported
    def do_PUT(self):
        #SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method.")

        target_chunkedResponse=False

        try:

            loggingPEPRequest(self)

            milli_secTotalReq=int(round(time.time() * 1000))

            headers,content_length = obtainRequestHeaders(self.headers)

            try:
                #To find only admittable headers from request previously configured in config.cfg file.
                value_index = headers.index("Error")
            except:
                value_index = -1

            testSupported = UtilsPEP.validateNotSupportedMethodPath(APIVersion,self.command,self.path)


            if (value_index != -1):
                logging.error("Error: " + str(headers["Error"]))
                SimpleHTTPRequestHandler.do_HandleError(self,self.command,400,"Bad Request","Error obtaining headers.")

            else:

                if (testSupported == False):
                    logging.error("Error: " + str(headers["Error"]))
                    SimpleHTTPRequestHandler.do_HandleError(self,self.command,501,"Not Implemented","No supported method/path.")

                else:

                    put_body = None

                    try:

                        #logging.info (" ********* OBTAIN BODY ********* ")
                        # We get the body
                        if (content_length>0):
                            #logging.info ("-------- self.rfile.read(content_length) -------")
                            put_body   = self.rfile.read(content_length)
                        else:
                            #logging.info ("-------- Lanzo self.rfile.read() -------")
                            put_body   = self.rfile.read()
                        
                        #logging.info(put_body)

                    except:
                        put_body = None

                    validation = validationToken(headers,self.command,self.path,put_body)

                    if (validation == False):
                        SimpleHTTPRequestHandler.do_HandleError(self,self.command,401,"Unauthorized","The token is missing or invalid.")

                    else:
                        # We are sending this to the CB
                        result = CBConnection(self.command, self.path,headers, put_body)

                        errorCBConnection = False
                        try:
                            if(result==-1):
                                errorCBConnection = True
                        except:
                            errorCBConnection = False

                        if(errorCBConnection):
                            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")
                        else:        
                            # We send back the response to the client
                            self.send_response(result.code)

                            headersResponse, target_chunkedResponse = obtainResponseHeaders(result.headers)

                            #logging.info(" ******* Sending Headers back to client ******* ")
                            for key in headersResponse:
                                if (PEPPROXY_CORS_ENABLED == 1 and 
                                   key in ["Access-Control-Allow-Origin","Access-Control-Allow-Methods","Access-Control-Allow-Headers",
                                   "Access-Control-Expose-Headers","Cache-Control"]):
                                    continue
                                else:
                                    self.send_header(key, headersResponse[key])
                                #self.send_header(key, headersResponse[key])
                                
                            self.send_headers()
                            self.end_headers()

                            #logging.info("Sending the Body back to client")

                            # Link to resolve Transfer-Encoding chunked cases
                            # https://docs.amazonaws.cn/en_us/polly/latest/dg/example-Python-server-code.html

                            while True:
                                data = result.read(chunk_size)
                        
                                if (target_chunkedResponse):
                                    self.wfile.write(b"%X\r\n%s\r\n" % (len(data), data))
                                else:
                                    self.wfile.write(data)
                                
                                if data is None or len(data) == 0:
                                    break
                    
                            if (target_chunkedResponse):
                                self.wfile.flush()

                        self.close_connection
                        logging.info("Response issued: Communication is closed.")

                        milli_secTotalReq2=int(round(time.time() * 1000))

                        if(logginKPI.upper()=="Y".upper()):
                            logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

        except Exception as e:
            logging.error(str(e))
            SimpleHTTPRequestHandler.do_HandleError(self,self.command,500,"Internal Server Error","GENERAL")


#class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
#    """Handle requests in a separate thread."""
#    pass

# Multiple thread: https://www.youtube.com/watch?v=A--u3WseDhQ
class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""
    daemon_threads = True


# Multiple thread: https://www.youtube.com/watch?v=A--u3WseDhQ
def server_on_port(port):
    httpd = ThreadedHTTPServer( (pep_host, port), SimpleHTTPRequestHandler)

    if (pep_protocol.upper() == "https".upper()) :
        #httpd.socket = ssl.wrap_socket (httpd.socket,
        #    keyfile="./certs/server-priv-rsa.pem",
        #    certfile="./certs/server-public-cert.pem",
        #    server_side = True,
        #    ssl_version=ssl.PROTOCOL_TLS
        #)
        context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        context.load_cert_chain(certfile="./certs/server-public-cert.pem", keyfile="./certs/server-priv-rsa.pem")
        httpd.socket = context.wrap_socket(httpd.socket, server_side=True)

    try:
        httpd.serve_forever()
    except Exception as e:
        logging.error(str(e))
        pass

    httpd.server_close()

if __name__ == '__main__':

    # Multiple thread: https://www.youtube.com/watch?v=A--u3WseDhQ
    #Thread(target=server_on_port,args=[9090]).start()
    Thread(target=server_on_port,args=[pep_port]).start()
    #server_on_port(pep_port)

#    httpd = ThreadedHTTPServer( (pep_host, pep_port), SimpleHTTPRequestHandler )
#
#    if (pep_protocol.upper() == "https".upper()) :
#        httpd.socket = ssl.wrap_socket (httpd.socket,
#            keyfile="./certs/server-priv-rsa.pem",
#            certfile="./certs/server-public-cert.pem",
#            server_side = True,
#            ssl_version=ssl.PROTOCOL_TLS
#        )
#
#    try:
#        httpd.serve_forever()
#    except Exception as e:
#        logging.error(str(e))
#        pass
#
#    httpd.server_close()

#httpd = HTTPServer( (pep_host, pep_port), SimpleHTTPRequestHandler )
#
#httpd.socket = ssl.wrap_socket (httpd.socket,
#        keyfile="certs/server-priv-rsa.pem",
#        certfile='certs/server-public-cert.pem',
#        server_side = True)
#
#httpd.serve_forever()
