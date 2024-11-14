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
import os
from socketserver import ThreadingMixIn
import threading

import time
from datetime import datetime
import hashlib
import requests

import xmltodict

import re

from hashlib import sha256

#Obtain configuracion from config.cfg file.
cfg = configparser.ConfigParser()  
cfg.read(["./config.cfg"])  
host = cfg.get("GENERAL", "host")
port = int(cfg.get("GENERAL", "port"))
chunk_size=int(cfg.get("GENERAL", "chunk_size"))

PDP_Protocol = ""
try:
    PDP_Protocol = str(os.getenv('PDP_Protocol'))

    if (str(PDP_Protocol).upper() == "None".upper()) :
        PDP_Protocol = "http"
except Exception as e:
    logging.error(e)
    PDP_Protocol = "http"

XACML_Default_Domain = str(os.getenv('XACML_Default_Domain'))

XACML_Location_Type = ""
try:
    XACML_Location_Type = str(os.getenv('XACML_Location_Type'))

    if (str(XACML_Location_Type).upper() == "None".upper()) :
        XACML_Location_Type = "file"
except Exception as e:
    logging.error(e)
    XACML_Location_Type = "file"


XACML_API_Protocol = ""
try:
    XACML_API_Protocol = str(os.getenv('XACML_API_Protocol'))

    if (str(XACML_API_Protocol).upper() == "None".upper()) :
        XACML_API_Protocol = "http"
except Exception as e:
    logging.error(e)
    XACML_API_Protocol = "http"

XACML_API_Host = str(os.getenv('XACML_API_Host'))
XACML_API_Port = int(os.getenv('XACML_API_Port'))
XACML_API_Get_Resource = str(os.getenv('XACML_API_Get_Resource'))

HASH_Validation = ""
try:
    HASH_Validation = str(os.getenv('HASH_Validation'))

    if (str(HASH_Validation).upper() == "None".upper()) :
        HASH_Validation = "0"
except Exception as e:
    logging.error(e)
    HASH_Validation = "0"


HASH_Protocol = ""
try:
    HASH_Protocol = str(os.getenv('HASH_Protocol'))

    if (str(HASH_Protocol).upper() == "None".upper()) :
        HASH_Protocol = "http"
except Exception as e:
    logging.error(e)
    HASH_Protocol = "http"

HASH_Host = str(os.getenv('HASH_Host'))
HASH_Port = int(os.getenv('HASH_Port'))
HASH_Get_Resource = str(os.getenv('HASH_Get_Resource'))

PDP_CORS_ENABLED = int(os.getenv('PDP_CORS_ENABLED'))

# XACML_DomainList: contain a list of domain objects with local information obtained from the current XACML policies files. 
# It helps to reduce the number of access (open file, API requests) because if the timestamp is the same we have nothing to do.
XACML_DomainList = []

tabLoggingString = "\t\t\t\t\t\t\t"

logginKPI = cfg.get("GENERAL", "logginKPI")

gcontext = ssl.SSLContext()

# Opening JSON file
f = open('SubjectIdTypes.json')
# returns JSON object as a dictionary
configuration = json.load(f)
# Closing file
f.close()

# Obtain the index of a domain in the domain list (XACML_DomainList)
def obtainDomainIndex(domain):
    global XACML_DomainList

    try:

        for i in range(len(XACML_DomainList)):
            if XACML_DomainList[i]['domain'] == domain:
                return i
        return -1
    except:
        return -1

# Add a new element for the domain in the domain list (XACML_DomainList)
def initDomain(domain):
    global XACML_DomainList

    try:

        XACML_DomainList.append( { "domain": domain, "getmtime": 0, "XACMLJSON": {}, "XACMLString": "", "digest": "" } ) 

        return True

    except:
        return False

# If the timestamp is different, this function refresh the domain element in the domain list (XACML_DomainList). 
def refreshDomain(domain, getmtime, XACMLJSON, XACMLString, digest):
    global XACML_DomainList

    try:

        indexDomain = obtainDomainIndex(domain)

        if (indexDomain > -1):
            XACML_DomainList[indexDomain]["getmtime"] = getmtime
            XACML_DomainList[indexDomain]["XACMLJSON"] = XACMLJSON
            XACML_DomainList[indexDomain]["XACMLString"] = XACMLString
            XACML_DomainList[indexDomain]["digest"] = digest
        else:
            XACML_DomainList.append( { "domain": domain, "getmtime": getmtime, "XACMLJSON": XACMLJSON, "XACMLString": XACMLString, "digest": digest } ) 

        return True

    except:
        return False

# Obtain the header of the PDP request, the more inportant are:
# - content_length : usefull to read the body of the request.
# - content_type : to process the body depending of the case.
# - domain : to determine the XACML Policies that can must be used to issued the verdict.
def obtainRequestHeaders(RequestHeaders):

    headers = dict()

    content_length = 0
    content_type = ""
    domain = ""

    try:
        # We get the headers
        
        #logging.info (" ********* HEADERS BEFORE obtainRequestHeaders ********* ")
        #logging.info (RequestHeaders)
        
        for key in RequestHeaders:
            #logging.info("Included: " + str(key) + ":" + str(RequestHeaders[key]))

            headers[key] = RequestHeaders[key]

            if(key.upper()=="Content-Length".upper()):
                content_length = int(RequestHeaders[key])
            
            if(key.upper()=="Content-Type".upper()):
                content_type = str(RequestHeaders[key])

            if(key.upper()=="domain".upper()):
                domain = str(RequestHeaders[key])


    except Exception as e:
        logging.error(e)

        headers["Error"] = str(e)

    #logging.info (" ********* HEADERS AFTER obtainRequestHeaders ********* ")
    #logging.info (headers)

    return headers, content_length, content_type, domain

# HTTP/HTTPS Server.
class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):

    # Handle OPTIONS requests.
    def do_OPTIONS(self):
        self.send_response(200)
        self.send_headers()
        self.end_headers()

    # To include CORS hedars if it is configured.
    def send_headers (self):
        if (PDP_CORS_ENABLED == 1):
            #self.send_header('Access-Control-Allow-Credentials', 'true')
            self.send_header('Access-Control-Allow-Origin', '*')
            #self.send_header('Access-Control-Allow-Methods', '*')
            self.send_header('Access-Control-Allow-Methods', 'GET,POST,OPTIONS')
            self.send_header('Access-Control-Allow-Headers', '*')
            self.send_header('Access-Control-Expose-Headers', '*')
            #self.send_header('Cache-Control','no-store,no-cache,must-revalidate')
            self.send_header('Cache-Control','no-cache,private,no-store,must-revalidate,max-stale=0,post-check=0,pre-check=0')
    
    # To send Error response.
    def do_HandleError(self):
        self.send_response(500)
        self.send_headers()
        self.end_headers() 
        data = json.dumps("Internal server error").encode()

        self.wfile.write(data)

        self.close_connection
    
    # Handle POST requests.
    def do_POST(self):
        
        try:

            logging.info("POST '" + str(self.path) + "' Request (" + self.address_string() +  ")... ")

            milli_secTotalReq=int(round(time.time() * 1000))

            # Validate the path of the request.
            if (processUri("POST", self.path) == False):

                logging.error("Error: Invalid path: '" + str(self.path) + "'")
                self.do_HandleError()
            else:

                headers,contentLength,contentType, domain = obtainRequestHeaders(self.headers)

                # If domain header is present it determines el domain to consider,
                # if it is not the case the default domain is used.
                if (domain == ""):

                    domain = XACML_Default_Domain

                #logging.info("domain: " + domain)

                #logging.info(self.headers)
                
                #logging.info (" ********* OBTAIN BODY ********* ")
                # We get the body
                if (contentLength>0):
                    post_body   = self.rfile.read(contentLength)
                else:
                    #logging.info ("-------- Lanzo self.rfile.read() -------")
                    post_body   = self.rfile.read()

                # Obtain the index of the domain in the domain list (XACML_DomainList)
                indexDomain = obtainDomainIndex(domain)

                #logging.info("indexDomain: " + str(indexDomain))

                newDomain = False

                # If the domanin does not included in the domain list (XACML_DomainList), it will be included.
                if (indexDomain == -1 ):

                    if (initDomain(domain)):
                        indexDomain = obtainDomainIndex(domain)
                        if (indexDomain == -1 ):
                            logging.error("The domain '" + domain + "' could not be initialized (1).")
                        else:
                            newDomain = True
                    else:
                        logging.error("The domain '" + domain + "' could not be initialized (2).")

                test = 0
                currentGetmtime = 0

                if newDomain:
                    logging.info("Detected new domain '" + domain + "'.")
                    # Detects if local information of the domain must be refreshed, for it the timestamp of the last change of the XACML Policies will be used.
                    test,currentGetmtime = detectXACMLChanges(domain,indexDomain) 

                # If test == 1 it means that local information of the domain must be refreshed.
                if (test == 1):
                    try:
                        logging.info("XACML Policies (" +  domain + ") changed: " + str(datetime.fromtimestamp(currentGetmtime)))
                    except:
                        logging.info("XACML Policies (" +  domain + ") changed: " + str(currentGetmtime))

                    # Refresh the local information of the domain.
                    if (recalculateDomain(domain,indexDomain,currentGetmtime) == False):
                        
                        result = "Can not refresh the domain variables."
                         
                        logging.error(result)

                        self.send_response(500)
                        # Send headers
                        self.send_header('Content-Type','text/plain; charset=utf-8')
                        self.send_header('Content-Length', str(len(result)))
                        self.send_headers()
                        self.end_headers()
                        self.wfile.write(result.encode('utf8'))
                        self.close_connection

                        return
                else:

                    # If test == -1 it means that the service has no access to XACML Policies last change and it os considered as an error.
                    if (test == -1):

                        result = "Can not validate if the XACML policies changed."
                        
                        logging.error(result)

                        self.send_response(500)
                        # Send headers
                        self.send_header('Content-Type','text/plain; charset=utf-8')
                        self.send_header('Content-Length', str(len(result)))
                        self.send_headers()
                        self.end_headers()
                        self.wfile.write(result.encode('utf8'))
                        self.close_connection

                        return

                result = ""
                resourceValue = ""

                # Optional: from XACML Policies we can obtain a HASH that could be stored, for instance in a DLT. This code compares if the HASH of the XACML Policies and 
                # the HASH stored (accesible through API) is the same.
                if (HASH_Validation == "1"):

                    if (isSameHash(indexDomain) == False):

                        result = "HASH is not accesible or does not correspond with the stored one."

                        self.send_response(500)
                        # Send headers
                        self.send_header('Content-Type','text/plain; charset=utf-8')
                        self.send_header('Content-Length', str(len(result)))
                        self.send_headers()
                        self.end_headers()
                        self.wfile.write(result.encode('utf8'))
                        self.close_connection

                        return  

                # Consume verdict request
                if (str(self.path).upper() == "/pdp/verdict".upper() or str(self.path).upper() == "/XACMLServletPDP/".upper()):
                    #"/XACMLServletPDP/" will be removed once was replaced by "/pdp/verdict" in Capability Manager component.

                    # Handling "application/json" bodies.
                    if (contentType.upper() == "application/json".upper()):

                        #bodyJSON = json.loads(post_body.decode('utf8').replace("'", '"'))
                        bodyJSON = json.loads(post_body.decode('utf8'))

                        sendVerdict = False

                        messageLog = ""

                        # For to obtain the elements of array of the request. 
                        for i in range(len(bodyJSON)):

                            bodyValue = bodyJSON[i]["body"]

                            xpars = xmltodict.parse(bodyValue)
                            jsonStringValue = json.dumps(xpars)

                            #jsonValue = json.loads(jsonStringValue.replace("'", '"'))
                            jsonValue = json.loads(jsonStringValue)

                            # From the "body" field of each element of the array, we have a string with the next format 
                            
                            '''
                            {"Request": {
                                "@xmlns": "urn:oasis:names:tc:xacml:2.0:context:schema:os", 
                                "Subject": {
                                    "@SubjectCategory": "urn:oasis:names:tc:xacml:1.0:subject-category:access-subject", 
                                    "Attribute": {
                                        "@AttributeId": "{{subjectTypeValue}}", 
                                        "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                        "AttributeValue": "{{subjectValue}}"
                                    }
                                }, 
                                "Resource": {
                                    "Attribute": {
                                        "@AttributeId": "urn:oasis:names:tc:xacml:1.0:resource:resource-id", 
                                        "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                        "AttributeValue": "{{resourceValue}}"
                                    }
                                }, 
                                "Action": {
                                    "Attribute": {
                                        "@AttributeId": "urn:oasis:names:tc:xacml:1.0:action:action-id", 
                                        "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                        "AttributeValue": "{{actionValue}}"
                                    }
                                }, 
                                "Environment": null
                                }
                            }
                            '''

                            actionValue = jsonValue["Request"]["Action"]["Attribute"]["AttributeValue"]
                            subjectTypeValue = jsonValue["Request"]["Subject"]["Attribute"]["@AttributeId"]
                            subjectValue = jsonValue["Request"]["Subject"]["Attribute"]["AttributeValue"]
                            resourceValue = jsonValue["Request"]["Resource"]["Attribute"]["AttributeValue"]
                            
                            #logging.info("actionValue: " + actionValue)
                            #logging.info("subjectTypeValue: " + subjectTypeValue)
                            #logging.info("subjectValue: " + subjectValue)
                            #logging.info("resourceValue: " + resourceValue)

                            messageLog = "{ 'Action': '" + actionValue + "', " \
                                        "'SubjectType': '" + subjectTypeValue + "', " \
                                        "'Subject': '" + subjectValue + "', " \
                                        "'Resource': '" + resourceValue + "' }"

                            # To obtain the verdict.
                            result = obtainVerdict(indexDomain, actionValue, subjectTypeValue, subjectValue, resourceValue)                        

                            if ("<Decision>NotApplicable</Decision>" not in result):

                                if ("<Decision>Permit</Decision>" in result):
                                    logging.info("POST Response: VERDICT: Permit.")
                                else:
                                    logging.info(messageLog)
                                    logging.info("POST Response: VERDICT: Deny.")
                                
                                sendVerdict = True

                                self.send_response(200)
                                # Send headers
                                self.send_header('Content-Type','text/plain; charset=utf-8')
                                self.send_header('SubjectType',str(subjectTypeValue))
                                self.send_header('Subject',str(subjectValue))
                                self.send_header('Content-Length', str(len(result)))
                                self.send_headers()
                                self.end_headers()
                                self.wfile.write(result.encode('utf8'))
                                self.close_connection

                                break

                        if (sendVerdict == False):

                            logging.info(messageLog)

                            result = "<Response>\n" \
                            "  <Result ResourceID=\"" + resourceValue + "\">\n" \
                            "    <Decision>NotApplicable</Decision>\n" \
                            "    <Status>\n" \
                            "      <StatusCode Value=\"urn:oasis:names:tc:xacml:1.0:status:ok\"/>\n" \
                            "    </Status>\n" \
                            "  </Result>\n" \
                            "</Response>"

                            logging.info("POST Response: VERDICT: NotApplicable.")

                            self.send_response(200)
                            # Send headers
                            self.send_header('Content-Type','text/plain; charset=utf-8')
                            self.send_header('Content-Length', str(len(result)))
                            self.send_headers()
                            self.end_headers()
                            self.wfile.write(result.encode('utf8'))
                            self.close_connection

                    else :

                        xpars = xmltodict.parse(post_body)
                        jsonStringValue = json.dumps(xpars)

                        #jsonValue = json.loads(jsonStringValue.replace("'", '"'))
                        jsonValue = json.loads(jsonStringValue)

                        # From the body request we have a string with the next format 
                        '''
                        {"Request": {
                            "@xmlns": "urn:oasis:names:tc:xacml:2.0:context:schema:os", 
                            "Subject": {
                                "@SubjectCategory": "urn:oasis:names:tc:xacml:1.0:subject-category:access-subject", 
                                "Attribute": {
                                    "@AttributeId": "{{subjectTypeValue}}", 
                                    "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                    "AttributeValue": "{{subjectValue}}"
                                }
                            }, 
                            "Resource": {
                                "Attribute": {
                                    "@AttributeId": "urn:oasis:names:tc:xacml:1.0:resource:resource-id", 
                                    "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                    "AttributeValue": "{{resourceValue}}"
                                }
                            }, 
                            "Action": {
                                "Attribute": {
                                    "@AttributeId": "urn:oasis:names:tc:xacml:1.0:action:action-id", 
                                    "@DataType": "http://www.w3.org/2001/XMLSchema#string", 
                                    "AttributeValue": "{{actionValue}}"
                                }
                            }, 
                            "Environment": null
                            }
                        }
                        '''

                        actionValue = jsonValue["Request"]["Action"]["Attribute"]["AttributeValue"]
                        subjectTypeValue = jsonValue["Request"]["Subject"]["Attribute"]["@AttributeId"]
                        subjectValue = jsonValue["Request"]["Subject"]["Attribute"]["AttributeValue"]
                        resourceValue = jsonValue["Request"]["Resource"]["Attribute"]["AttributeValue"]
                            
                        #logging.info("actionValue: " + actionValue)
                        #logging.info("subjectTypeValue: " + subjectTypeValue)
                        #logging.info("subjectValue: " + subjectValue)
                        #logging.info("resourceValue: " + resourceValue)

                        messageLog = "{ 'Action': '" + actionValue + "', " \
                                            "'SubjectType': '" + subjectTypeValue + "', " \
                                            "'Subject': '" + subjectValue + "', " \
                                            "'Resource': '" + resourceValue + "' }"

                        # To obtain the verdict.
                        result = obtainVerdict(indexDomain, actionValue, subjectTypeValue, subjectValue, resourceValue)

                        if ("<Decision>Permit</Decision>" in result):
                            decision = "Permit"
                            logging.info("POST Response: VERDICT: Permit.")
                        else:

                            logging.info(messageLog)

                            if ("<Decision>Deny</Decision>" in result):
                                decision = "Deny"
                                logging.info("POST Response: VERDICT: Deny.")
                            else:
                                decision = "NotApplicable"
                                logging.info("POST Response: VERDICT: NotApplicable.")

                        # Obtain actual date in ISO format
                        date = datetime.now().isoformat()
                        id_string = f"{actionValue}{subjectValue}{date}{decision}{resourceValue}"
                        unique_id = hashlib.md5(id_string.encode()).hexdigest()

                        subjectValue = jsonValue["Request"]["Subject"]["Attribute"]["AttributeValue"]
                        resourceValue = jsonValue["Request"]["Resource"]["Attribute"]["AttributeValue"]
                        actionValue = jsonValue["Request"]["Action"]["Attribute"]["AttributeValue"]

                        json_data = {
                            "Action": actionValue,
                            "Subject": subjectValue,
                            "Resource": resourceValue,
                            "Timestamp": date,
                            "Id": unique_id,
                            "Decision": decision
                        }

                        json_paylaod = json.dumps(json_data)

                        try:
                            response = requests.post('http://172.16.10.118:3002/xadatu/auth/register', data=json_paylaod, headers={'Content-Type': 'application/json'})
                            if response.status_code == 200:
                                logging.info("Access request registered in blockchain")
                            else:
                                logging.error(f"There was an error sending the POST request: {response.status_code} {response}")
                        except requests.exceptions.RequestException as e:
                            logging.error(f"Connection error: {e}")

                        self.send_response(200)
                        # Send headers
                        self.send_header('Content-Type','text/plain; charset=utf-8')
                        self.send_header('SubjectType',str(subjectTypeValue))
                        self.send_header('Subject',str(subjectValue))
                        self.send_header('Content-Length', str(len(result)))
                        self.send_headers()
                        self.end_headers()
                        self.wfile.write(result.encode('utf8'))
                        self.close_connection
                
                #Consume attribute list request (/pdp/attrlist)
                else:

                    if (contentType.upper() == "application/json".upper()):

                        #bodyJSON = json.loads(post_body.decode('utf8').replace("'", '"'))
                        bodyJSON = json.loads(post_body.decode('utf8'))

                        validateFields = True
                        try:
                            resourceValue = bodyJSON["re"]
                            actionValue = bodyJSON["ac"]
                            
                        except:
                            validateFields = False

                        if ( validateFields == False):

                            logging.error("Body must contain JSON with 're' and 'ac' fields.")

                            self.send_response(400)
                            # Send headers
                            self.send_header('Content-Type','application/json; charset=utf-8')
                            self.send_headers()
                            self.end_headers()
                            self.wfile.write(json.dumps({"error": "Body must contain JSON with 're' and 'ac' fields."}).encode())
                            self.close_connection
                        
                        else:

                            result,attrList = obtainAttributeList(indexDomain, actionValue, resourceValue);

                            if (result == -1):
                                self.do_HandleError()
                            else: 
                                if (len(attrList) == 0):					
                                
                                    self.send_response(404)
                                    # Send headers
                                    self.send_header('Content-Type','application/json; charset=utf-8')
                                    self.send_headers()
                                    self.end_headers()
                                    self.wfile.write(json.dumps({"error": "Attribute list not found for 're' and 'ac' defined."}).encode())
                                    self.close_connection

                                else:

                                    self.send_response(200)
                                    # Send headers
                                    self.send_header('Content-Type','application/json; charset=utf-8')
                                    self.send_headers()
                                    self.end_headers()
                                    self.wfile.write(json.dumps(attrList).encode())
                                    self.close_connection

                    else:

                        logging.error("Content-Type header must be 'application/json'.")

                        self.send_response(400)
                        # Send headers
                        self.send_header('Content-Type','application/json; charset=utf-8')
                        self.send_headers()
                        self.end_headers()
                        self.wfile.write(json.dumps({"error": "Content-Type header must be 'application/json'."}).encode())
                        self.close_connection

            logging.info("Response issued: Communication is closed.")
            
            milli_secTotalReq2=int(round(time.time() * 1000))

            if(logginKPI.upper()=="Y".upper()):
                logging.info("Total(ms) request: " + str(milli_secTotalReq2 - milli_secTotalReq))

        except Exception as e:
            logging.error(str(e))
            
            self.do_HandleError()
    
    # Handle GET requests.
    def do_GET(self):

        # Validate the path of the request.
        if (processUri("GET", self.path)):
            self.send_response(200)
            self.send_headers()
            self.end_headers()
            self.close_connection
        else:
            self.do_HandleError()

# Validate if path request is valid/supported.
def processUri(method, uri):
    
    try:
        if (method.upper() == "POST".upper() and (str(uri).upper() == "/pdp/verdict".upper() or str(uri).upper() == "/pdp/attrlist".upper() or 
            str(uri).upper() == "/XACMLServletPDP/".upper())):
            #"/XACMLServletPDP/" will be removed once was replaced by "/pdp/verdict" in Capability Manager component.
            return True
        else:
            if (method.upper() == "GET".upper() and str(uri).upper() == "/pdp/test".upper()):
                return True

        return False

    except:
        return False

# To obtain the verdict it is based in local information of XACML Policies of the domain. (Field XACMLJSON of XACML_DomainList)
def obtainVerdict(indexDomain, actionValue, subjectTypeValue, subjectValue, resourceValue):
    global XACML_DomainList

    XACMLJSON = XACML_DomainList[indexDomain]["XACMLJSON"]

    result = ""

    try:
        output_dict = [x for x in configuration if str(subjectTypeValue).lower() == str(x["subjectIdCapMan"]).lower()]

        if (len(output_dict) == 1):
            subjectTypeValue = output_dict[0]["subjectIdPAP"]
        else:
            logging.info("Subject Type mapping could not be performed.")
            result =    "<Response>\n" \
                    "  <Result ResourceID=\"" + resourceValue + "\">\n" \
                    "    <Decision>NotApplicable</Decision>\n" \
                    "    <Status>\n" \
                    "      <StatusCode Value=\"urn:oasis:names:tc:xacml:1.0:status:ok\"/>\n" \
                    "    </Status>\n" \
                    "  </Result>\n" \
                    "</Response>"
            return result

        verdict = "NotApplicable"

        findAction = False
        findSubject = False
        findPolicyTriplet = False

        if ("Policy" not in XACMLJSON["PolicySet"]):
            logging.info("There are not defined policies.")
            result =    "<Response>\n" \
                    "  <Result ResourceID=\"" + resourceValue + "\">\n" \
                    "    <Decision>" + verdict + "</Decision>\n" \
                    "    <Status>\n" \
                    "      <StatusCode Value=\"urn:oasis:names:tc:xacml:1.0:status:ok\"/>\n" \
                    "    </Status>\n" \
                    "  </Result>\n" \
                    "</Response>"
            return result

        #logging.info(XACMLJSON)
    
        for i in range(len(XACMLJSON["PolicySet"]["Policy"])):

            if ("Rule" not in XACMLJSON["PolicySet"]["Policy"][i]):
                continue
            
            for j in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"])):

                findAction = False
                findSubject = False
                findPolicyTriplet = False

                if (type(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]) == type(None)):
                    continue

                # Obtain if the Subject is present (findSubject).
                if ("Subjects" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                    #logging.info("subjectTypeValue: " + subjectTypeValue)
                    #logging.info("subjectValue: " + subjectValue)

                    if (type(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]) != type(None)):
                        for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"])):

                            subjectTypeRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"][k]["SubjectMatch"]["SubjectAttributeDesignator"]["@AttributeId"]
                            subjectRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"][k]["SubjectMatch"]["AttributeValue"]["#text"]

                            #logging.info("subjectTypeRule: " + subjectTypeRule)
                            #logging.info("subjectRule: " + subjectRule)

                            if (subjectTypeValue == subjectTypeRule and subjectValue == subjectRule):
                                findSubject = True
                                break
                    else:
                        # If Subjects field is not present, it understands that the Subject is found.
                        findSubject = True

                else:
                    # If Subjects field is not present, it understands that the Subject is found.
                    findSubject = True

                if (findSubject):
               
                    # Obtain if the Action is present (findAction).
                    if ("Actions" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                        #logging.info("actionValue: " + actionValue)
                        if (type(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]) != type(None)):
                            for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"])):

                                actionRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"][k]["ActionMatch"]["AttributeValue"]["#text"]
                                #logging.info("actionRule: " + actionRule)

                                if (actionValue == actionRule):
                                    findAction = True
                                    break
                        else:
                            # If Actions field is not present, it understands that the Action is found.
                            findAction = True

                    else:
                        # If Actions field is not present, it understands that the Action is found.
                        findAction = True

                    if (findAction):

                        # Obtain if the Resource is present and in this sense the triplet is present (findPolicyTriplet).
                        if ("Resources" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):
                            #logging.info("resourceValue: " + resourceValue)
                       
                            if (type(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]) != type(None)):
                                for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"])):

                                    resourceRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"][k]["ResourceMatch"]["AttributeValue"]["#text"]
                                    #logging.info("resourceRule: " + resourceRule)

                                    if (resourceValue == resourceRule):
                                        verdict = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["@Effect"]
                                        findPolicyTriplet = True                                    
                                        break

                                    else:

                                        #java.util.regex.Pattern pattern1 = java.util.regex.Pattern.compile("^" + resource + "$");
                                        #java.util.regex.Matcher matcher1 = pattern1.matcher( resourceValue );

                                        pattern1 = re.compile("^" + resourceRule + "$")
                                        matcher1 = pattern1.match( resourceValue )

                                        if (matcher1):
                                            verdict = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["@Effect"]
                                            findPolicyTriplet = True
                                            break

                                        '''
                                        else:

                                            #java.util.regex.Pattern pattern2 = java.util.regex.Pattern.compile("^" + java.util.regex.Pattern.quote(resource) + "$");
                                            #java.util.regex.Matcher matcher2 = pattern2.matcher( resourceValue );

                                        '''
                            else:
                                # If Resources field is not present, it understands that the Resource is found.
                                verdict = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["@Effect"]
                                findPolicyTriplet = True
                                break
                            
                            if (findPolicyTriplet):
                                break
                                                                    
                        else:
                            # If Resources field is not present, it understands that the Resource is found.
                            verdict = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["@Effect"]
                            findPolicyTriplet = True
                            break
                
                #if (findPolicyTriplet):
                #    break
            
            if (findPolicyTriplet):
                break

        #if (findPolicyTriplet == False):
        #	verdict = "NotApplicable";
                    
        result =    "<Response>\n" \
                    "  <Result ResourceID=\"" + resourceValue + "\">\n" \
                    "    <Decision>" + verdict + "</Decision>\n" \
                    "    <Status>\n" \
                    "      <StatusCode Value=\"urn:oasis:names:tc:xacml:1.0:status:ok\"/>\n" \
                    "    </Status>\n" \
                    "  </Result>\n" \
                    "</Response>"
    
    except Exception as e:
        logging.error(e)
        result =    "<Response>\n" \
                    "  <Result ResourceID=\"" + resourceValue + "\">\n" \
                    "    <Decision>NotApplicable</Decision>\n" \
                    "    <Status>\n" \
                    "      <StatusCode Value=\"urn:oasis:names:tc:xacml:1.0:status:ok\"/>\n" \
                    "    </Status>\n" \
                    "  </Result>\n" \
                    "</Response>"

    return result

# To obtain the list of attribute for the presentation of credentials based in local information of XACML Policies of the domain. (Field XACMLJSON of XACML_DomainList)
def obtainAttributeList(indexDomain, actionValue, resourceValue):
    global XACML_DomainList

    XACMLJSON = XACML_DomainList[indexDomain]["XACMLJSON"]

    result = 0

    attrList = []

    try: 

        findAction = False
        findResource = False

        if ("Policy" not in XACMLJSON["PolicySet"]):
            logging.info("There are not defined policies.")
            return result, attrList

        for i in range(len(XACMLJSON["PolicySet"]["Policy"])):

            if ("Rule" not in XACMLJSON["PolicySet"]["Policy"][i]):
                continue
            
            for j in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"])):

                findAction = False
                findResource = False

                if (type(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]) == type(None)):
                    continue

                # Obtain if the Action is present (findAction).
                if ("Actions" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                    #logging.info("actionValue: " + actionValue)
                       
                    for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"])):

                        actionRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"][k]["ActionMatch"]["AttributeValue"]["#text"]
                        #logging.info("actionRule: " + actionRule)

                        if (actionValue == actionRule):
                            findAction = True
                            break

                else:
                    # If Actions field is not present, it understands that the Action is found.
                    findAction = True

                if (findAction):

                    # Obtain if the Resource is present
                    if ("Resources" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):
                        #logging.info("resourceValue: " + resourceValue)
                       
                        for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"])):
                                
                            resourceRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"][k]["ResourceMatch"]["AttributeValue"]["#text"]
                            #logging.info("resourceRule: " + resourceRule)
                                
                            if (resourceValue == resourceRule):
                                findResource = True                                    
                                break

                            else:

                                #java.util.regex.Pattern pattern1 = java.util.regex.Pattern.compile("^" + resource + "$");
                                #java.util.regex.Matcher matcher1 = pattern1.matcher( resourceValue );
                                    
                                pattern1 = re.compile("^" + resourceRule + "$")
                                matcher1 = pattern1.match( resourceValue )

                                if (matcher1):
                                    findResource = True                                    
                                    break

                                '''
                                else:

                                    #java.util.regex.Pattern pattern2 = java.util.regex.Pattern.compile("^" + java.util.regex.Pattern.quote(resource) + "$");
                                    #java.util.regex.Matcher matcher2 = pattern2.matcher( resourceValue );

                                '''
                                                                                                     
                    else:
                        # If Resources field is not present, it understands that the Resource is found.
                        findResource = True
                    
                    if (findResource):

                        # Obtain if the Subject is present (findSubject).
                        if ("Subjects" in XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                             for k in range(len(XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"])):

                                subjectTypeRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"][k]["SubjectMatch"]["SubjectAttributeDesignator"]["@AttributeId"]
                                #subjectRule = XACMLJSON["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"][k]["SubjectMatch"]["AttributeValue"]["#text"]

                                elem =  subjectTypeRule.split(':')

                                #logging.info(elem[len(elem)-1])

                                if (elem[len(elem)-1] not in attrList):
                                    attrList.append(elem[len(elem)-1])                                
                    
    except Exception as e:
        logging.error(e)
    
        result = -1
        attrList = []

    return result, attrList


# Detects if local information of the domain must be refreshed, for it the timestamp of the last change of the XACML Policies will be used.
def detectXACMLChanges(domain,indexDomain):

    global XACML_DomainList, XACML_Location_Type, XACML_API_Protocol, XACML_API_Host, XACML_API_Port, XACML_API_Get_Resource    

    result = 0
    currentGetmtime = 0  

    try:

        # Obtain the timestamp from the date of update of the file.
        if (XACML_Location_Type.upper() == "file".upper()):

            #currentGetmtime = os.path.getmtime(r"continue-a.xml")
            currentGetmtime = os.path.getmtime(r"./Policies/" + str(domain) + ".xml")

            if currentGetmtime != XACML_DomainList[indexDomain]["getmtime"]:
                result = 1
        else:

            # Obtain the timestamp from a field of an API response.
            if (XACML_Location_Type.upper() == "api".upper()):

                currentGetmtime, xacmldata = obtainXACMLInfo(indexDomain)

                if (currentGetmtime != "-1"):

                    if currentGetmtime != XACML_DomainList[indexDomain]["getmtime"]:
                        logging.info("Datetime (" +  XACML_DomainList[indexDomain]["domain"] + ") does not correspond with the stored one:")
                        logging.info("\t - Local Datetime: " + str(XACML_DomainList[indexDomain]["getmtime"]))
                        logging.info("\t - Remote Datetime: " + str(currentGetmtime))

                        result = 1
                else:
                    result = -1

            else:
                result = -1


    except Exception as e:
        logging.error(e)
        currentGetmtime = 0
        result = -1

    # Returns the result of the comparation and the timestamp value obtained.
    return result, currentGetmtime

# Refresh the local information of the domain.
def recalculateDomain(domain,indexDomain,currentGetmtime):

    global XACML_DomainList, XACML_Location_Type
    
    result = False

    try:

        # Obtain the XACML Policies string from the corresponding file or through an API.        
        if (XACML_Location_Type.upper() == "file".upper()):
            #Read de policies.
            #f = open("continue-a.xml", "r")
            f = open("./Policies/" + str(domain)+".xml", "r")
            contentXACML = f.read()
            f.close()
        else:
            if (XACML_Location_Type.upper() == "api".upper()):

                # API Request to obtain XACML Policies string.
                getmtime, contentXACMLaux = obtainXACMLInfo(indexDomain)

                if (contentXACMLaux == ""):
                    return False
                else:

                    test = False
                    
                    contentXACML = ""

                    # To replace '\"' by '"' in the XACML Policies string obtained (scape " character)
                    for index in range(len(contentXACMLaux)):

                        if (test):
                            if (ord(contentXACMLaux[index]) == 34):
                                contentXACML = contentXACML + contentXACMLaux[index]
                            else:
                                contentXACML = contentXACML + contentXACMLaux[index-1] + contentXACMLaux[index]
                            test = False

                        else:

                            if (ord(contentXACMLaux[index]) == 92):
                                test = True
                            else:
                                contentXACML = contentXACML + contentXACMLaux[index]

        xpars = xmltodict.parse(contentXACML)
        XACMLString = json.dumps(xpars)

        #Calculate HASH
        XACMLdigest = sha256(XACMLString.encode('utf-8')).hexdigest()
        #logging.info(XACMLdigest)

        #jsonXACMLValue = json.loads(XACMLString.replace("'", '"'))
        jsonXACMLValue = json.loads(XACMLString)

        #logging.info(jsonXACMLValue)

        if ("Policy" not in jsonXACMLValue["PolicySet"]):

            XACMLJSON = jsonXACMLValue

            #logging.info(XACMLJSON)

            # Refresh the local domain information
            if (refreshDomain(domain, currentGetmtime, XACMLJSON, XACMLString, XACMLdigest)):
                return True
            else:
                return False

        if (type(jsonXACMLValue["PolicySet"]["Policy"]) != type([])):

            if (type(jsonXACMLValue["PolicySet"]["Policy"]) == type({})):

                jsonXACMLValue["PolicySet"]["Policy"] = [ jsonXACMLValue["PolicySet"]["Policy"] ]

            else:

                XACMLJSON = jsonXACMLValue

                #logging.info(XACMLJSON)

                # Refresh the local domain information
                if (refreshDomain(domain, currentGetmtime, XACMLJSON, XACMLString, XACMLdigest)):
                    return True
                else:
                    return False

        for i in range(len(jsonXACMLValue["PolicySet"]["Policy"])): 

            if ("Rule" not in jsonXACMLValue["PolicySet"]["Policy"][i]):
                continue

            if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"]) != type([])):

                if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"]) == type({})):

                    jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"] = [ jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"] ]
                
                else:
                    jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"] = []

            for j in range(len(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"])): 

                if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]) == type(None)):
                    continue 
                
                if ("Subjects" in jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                    if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]) != type(None)):
                        if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"]) != type([])):

                            if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"]) == type({})):

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"] = [ jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"] ]

                            else:

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Subjects"]["Subject"] = []

                if ("Actions" in jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):
                
                    if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]) != type(None)):
                        if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"]) != type([])):

                            if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"]) == type({})):

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"] = [ jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"] ]

                            else:

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Actions"]["Action"] = []

                if ("Resources" in jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]):

                    if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]) != type(None)):
                        if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"]) != type([])):

                            if (type(jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"]) == type({})):

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"] = [ jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"] ]

                            else:

                                jsonXACMLValue["PolicySet"]["Policy"][i]["Rule"][j]["Target"]["Resources"]["Resource"] = []

        #logging.info(jsonXACMLValue)
        
        XACMLJSON = jsonXACMLValue

        # Refresh the local domain information
        if (refreshDomain(domain, currentGetmtime, XACMLJSON, XACMLString, XACMLdigest)):
            return True
        else:
            return False

    except Exception as e:
        logging.error(e)
        result = False

    return result

# API Request to obtain XACML Policies string.
def obtainXACMLInfo(indexDomain):
    global XACML_DomainList, XACML_API_Get_Resource, XACML_API_Protocol, XACML_API_Host, XACML_API_Port
    
    try:

        resource = XACML_API_Get_Resource.replace("{{domain}}",XACML_DomainList[indexDomain]["domain"])

        #Send request to obtain HASH value.
        responseCode,data = connectionAPI(XACML_API_Protocol, XACML_API_Host, XACML_API_Port, "GET", resource, {}, None)

        errorConnectionAPI = False
        try:
            if(responseCode==-1):
                errorConnectionAPI = True
        except:
            errorConnectionAPI = False

        if (errorConnectionAPI==False):

            # Expected response code.
            if (responseCode == 200):

                # Expected response body.
                '''
                {
                    "id":"urn:ngsi-ld:xacml:{{domain}}",
                    "type":"xacml",
                    "version":{
                        "type": "Property",
                        "value": "{{version}}"
                    },
                    "xacml":{
                        "type": "Property",
                        "value": "{{XACML_Value}}"
                    },
                    "timestamp": {
                        "type": "Property",
                        "value": "{{datetime}}" #Example: "2023-02-16 22:34:00" (UTC)
                    }
                }
                '''

                #bodyJSON = json.loads(data.decode('utf8').replace("'", '"'))
                bodyJSON = json.loads(data.decode('utf8'))

                try:
                    
                    return bodyJSON["timestamp"]["value"], bodyJSON["xacml"]["value"]
                
                except Exception as e:
                    logging.error(e)
                    logging.error("Can not obtain the XACML timestamp (" + str(responseCode) + ") - No JSON data.")
                    return "-1", ""

            else:
                logging.error("Can not obtain the XACML timestamp (" + str(responseCode) + ")")
                return "-1", ""

        else:
            logging.error("Can not connect to obtain the XACML timestamp.")
            return "-1", ""

    except Exception as e:
        logging.error(e)
        return "-1", ""

# Compares if the HASH of the XACML Policies and the HASH stored (accesible through API) is the same.
def isSameHash(indexDomain):
    global XACML_DomainList, HASH_Get_Resource, HASH_Protocol, HASH_Host, HASH_Port
    
    try:

        resource = HASH_Get_Resource.replace("{{domain}}",XACML_DomainList[indexDomain]["domain"])

        #Send request to obtain HASH value.
        responseCode,data = connectionAPI(HASH_Protocol, HASH_Host, HASH_Port, "GET", resource, {}, None)

        errorConnectionAPI = False
        try:
            if(responseCode==-1):
                errorConnectionAPI = True
        except:
            errorConnectionAPI = False

        if (errorConnectionAPI==False):

            # Expected response code.
            if (responseCode == 201 or responseCode == 200):

                # Expected response body.
                '''
                {
                    "id":"urn:ngsi-ld:domain:{{domain}}",
                    "type":"domain",
                    "version":{
                        "type": "Property",
                        "value": "{{version}}"
                    },
                    "digest":{
                        "type": "Property",
                        "value": "{{HASH_Value}}"
                    },
                    "timestamp": {
                        "type": "Property",
                        "value": "{{datetime}}" #Example: "2023-02-16 22:34:00" (UTC)
                    }
                }
                '''

                #bodyJSON = json.loads(data.decode('utf8').replace("'", '"'))
                bodyJSON = json.loads(data.decode('utf8'))

                try:

                    if (str(XACML_DomainList[indexDomain]["digest"]) == str(bodyJSON["digest"]["value"])):
                        #logging.info("Is the same HASH --> True")
                        return True
                    else:
                        logging.error("HASH (" +  XACML_DomainList[indexDomain]["domain"] + ") does not correspond with the stored one:")
                        logging.error("\t - Local HASH: " + str(XACML_DomainList[indexDomain]["digest"]))
                        logging.error("\t - Remote HASH: " + str(bodyJSON["digest"]["value"]))
                        return False
                
                except Exception as e:
                    logging.error(e)
                    logging.error("Can not obtain the HASH (" + str(responseCode) + ") - No JSON data.")
                    return False
            else:
                logging.error("Can not obtain the HASH (" + str(responseCode) + ")")
                return False

        else:
            logging.error("Can not connect to obtain the HASH.")
            return False

    except Exception as e:
        logging.error(e)
        return False

# Standard funtion to request via API.
def connectionAPI(protocol, host, port, method, uri, headers = {}, body = None):

    try:

        #logging.info("")
        #logging.info("")
        #logging.info("********* connectionAPI *********")

        # send some data
        
        if(protocol.upper() == "http".upper() or protocol.upper() == "https".upper()):

            if(protocol.upper() == "http".upper()):
                conn = http.client.HTTPConnection(host, port)
            else:                    
                conn = http.client.HTTPSConnection(host,port,
                                                context=gcontext)
            #logging.info("API REQUEST:\n" +
            #             tabLoggingString + "- Host: " + host + "\n" + 
            #             tabLoggingString + "- Port: " + str(port) + "\n" + 
            #             tabLoggingString + "- Method: " + method + "\n" + 
            #             tabLoggingString + "- URI: " + uri + "\n" + 
            #             tabLoggingString + "- Headers: " + str(headers) + "\n" + 
            #             tabLoggingString + "- Body: " + str(body))

            conn.request(method, uri, body, headers)

            response = conn.getresponse()

            #logging.info("connectionAPI - RESPONSE")
            #logging.info(" SUCCESS : API response - code: "      + str(response.code))
            #logging.info("Headers: ")
            #logging.info(response.headers)

            #logging.info(response)

            code = response.code
            reason = response.reason
            data = response.read()
            conn.close()

            return code, data

            #logging.info("********* connectionAPI - END *********")

        else: 
            logging.error("connectionAPI - Protocol not valid: " + protocol)
            return -1, {}

    except Exception as e:
        logging.error(e)
        return -1, {}


logPath="./"
fileName="out"

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(threadName)-12.12s] [%(levelname)-5.5s]  %(message)s",
    handlers=[
        logging.FileHandler("{0}/{1}.log".format(logPath, fileName)),
        logging.StreamHandler(sys.stdout)
    ])

class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""
    pass

#if __name__ == '__main__':
#
#    httpd = ThreadedHTTPServer( (host, port), SimpleHTTPRequestHandler )
#
#    if (PDP_Protocol.upper() == "https".upper()) :
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



# this is the deamon thread to refresh local data information.
def polUpdate():
    global XACML_DomainList

    time.sleep(60)
    while True:

        try:
        
            for i in range(len(XACML_DomainList)):

                try:
                    if str(XACML_DomainList[i]["getmtime"]) == "0":
                        continue

                    # Detects if local information of the domain must be refreshed, for it the timestamp of the last change of the XACML Policies will be used.
                    test,currentGetmtime = detectXACMLChanges(str(XACML_DomainList[i]["domain"]),i) 

                    # If test == 1 it means that local information of the domain must be refreshed.
                    if (test == 1):
                        try:
                            logging.info("XACML Policies (" +  str(XACML_DomainList[i]["domain"]) + ") changed: " + str(datetime.fromtimestamp(currentGetmtime)))
                        except:
                            logging.info("XACML Policies (" +  str(XACML_DomainList[i]["domain"]) + ") changed: " + str(currentGetmtime))

                        # Refresh the local information of the domain.
                        if (recalculateDomain(str(XACML_DomainList[i]["domain"]),i,currentGetmtime) == False):

                            logging.error("ERROR:  Can not refresh the domain ("+  str(XACML_DomainList[i]["domain"]) + ") variables.")

                        else:

                            if (test == -1):

                                logging.error("ERROR:  Can not validate if the XACML policies of the domain ("+  str(XACML_DomainList[i]["domain"]) + ") changed.")

                            else:

                                logging.info("Refreshed XACML Policies (" +  str(XACML_DomainList[i]["domain"]) + ")")

                except Exception as e:
                    logging.error(e)

            time.sleep(60)

        except Exception as e:
                    logging.error(e)

# here´s the http server starter

def run():
  
    httpd = ThreadedHTTPServer( (host, port), SimpleHTTPRequestHandler )

    if (PDP_Protocol.upper() == "https".upper()) :
        httpd.socket = ssl.wrap_socket (httpd.socket,
            keyfile="./certs/server-priv-rsa.pem",
            certfile="./certs/server-public-cert.pem",
            server_side = True,
            ssl_version=ssl.PROTOCOL_TLS
        )

    try:
        httpd.serve_forever()
    except Exception as e:
        logging.error(str(e))
        pass

    httpd.server_close()



if __name__ == '__main__':

    # init the thread as deamon
    d = threading.Thread(target=polUpdate, name='Daemon')
    d.setDaemon(True)
    d.start()

    # runs the HTTP server
    run()

