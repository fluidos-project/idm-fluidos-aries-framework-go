#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

import json
from datetime import datetime
from pytz import timezone
#import os
from subprocess import Popen, PIPE

def processUri(uri):
    try:
        if (str(uri).upper().startswith("/ngsi-ld/v1".upper()) == False 
            and str(uri).upper() != "/scorpio/v1/info/".upper() 
            and str(uri).upper() != "/version".upper()
            ):
            uri = "/ngsi-ld/v1"+uri

        return uri
    except:
        return uri

def validateMethodPath(method,path):

    return True


def processBody(method,uri,body,sPAE,rPAE,noEncryptedKeys):

    bodyBackUp = body

    try:

        #Determine if method / uri is actually comtemplated by the process and run the corresponding process
        #depending de body.
        state = False

        bodyType = ""

        if (type(body) is dict):

            bodyType = "dict"

            #To detect JSON body format type.
            if ("id" in body or (("value" not in body) and ("object" not in body)) ):
                bodyType = "dict"
            else:
                if  (("id" not in body) and ("type" in body) and
                     (("value" in body) or ("object" in body))
                    ):
                    bodyType = "dict-attr"
                else:
                    print("Error processBody: Unsupported dict format: " + str(body))
                    return body, state
        else:
            if (type(body) is list):
                bodyType = "list"
            else:
                print("Error processBody: Unsupported type(body): " + str(type(body)))
                return body, state

        if ((method.upper()=="POST".upper() or method.upper()=="PATCH".upper()) 
                #and str(uri).upper().startswith("/ngsi-ld/v1/entities".upper())
                ):

            if (bodyType == "dict"):
                body,state = processCypher(body,sPAE,rPAE,noEncryptedKeys)
            else:
                if (bodyType == "dict-attr"):
                    body,state = processCypherAttr(body,sPAE,rPAE,noEncryptedKeys)
                else:
                    #list
                    body,state = processCypherList(body,sPAE,rPAE,noEncryptedKeys)

        ##POST - https://{HOST}:{PORT}/ngsi-ld/v1/entities
        #if (method.upper()=="POST".upper() and
        #    (uri.upper()=="/ngsi-ld/v1/entities/".upper()
        #    #or (str(uri).upper().startswith("/ngsi-ld/v1/entities".upper()) and str(uri).upper().endswith("/attrs".upper()))
        #    )):
        #
        #    if (bodyType == "dict"):
        #        body,state = processCypher(body,sPAE,rPAE,noEncryptedKeys)
        #    else:
        #        if (bodyType == "dict-attr"):
        #            body,state = processCypherAttr(body,sPAE,rPAE,noEncryptedKeys)
        #        else:
        #            body,state = processCypherList(body,sPAE,rPAE,noEncryptedKeys)
        #else:
        #
        #    #PATCH - https://{HOST}:{PORT}/ngsi-ld/v1/entities/{entityID}/attrs
        #    if (method.upper()=="PATCH".upper() and
        #        str(uri).upper().startswith("/ngsi-ld/v1/entities".upper())
        #        #sand str(uri).upper().endswith("/attrs".upper())
        #        ):
        #
        #        if (bodyType == "dict"):
        #            body,state = processCypher(body,sPAE,rPAE,noEncryptedKeys)
        #        else:
        #            if (bodyType == "dict-attr"):
        #                body,state = processCypherAttr(body,sPAE,rPAE,noEncryptedKeys)
        #            else:
        #                body,state = processCypherList(body,sPAE,rPAE,noEncryptedKeys)

        return body, state

    except:
        return bodyBackUp, False


#This process consider ONLY a JSON format.
def processCypher(body,sPAE,rPAE,noEncryptedKeys):

    #Here is an example of JSON structure:
    '''
    {
        "@context":[
                "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
                {
                    "Vehicle": "http://example.org/vehicle/Vehicle",
                    "brandName": "http://example.org/vehicle/brandName",
                    "speed": "http://example.org/vehicle/speed",
                    "color": "http://example.org/vehicle/color"
                }
        ],
        "id":"urn:ngsi-ld:Vehicle:TESTJUAN10",
        "type":"Vehicle",
        "brandName":{
            "type":"Property",
            "value":"Mercedes",
            "encrypt_cpabe":{
                "type":"Property",
                "value":"att1 att2 2of2"
            }
        },
        "speed":{
            "type":"Property",
            "value":80,
            "encrypt_cpabe2":{
                "type":"Property",
                "value":"admin"
            }
        },
        "color":{
            "type":"Property",
            "value":"Red",
            "encrypt_cpabe":{
                "type":"Property",
                "value":"att4 att5 2of2"
            }
        }  
    }
    '''

    bodyBackUp = body

    try:

        encryptAttributes,state = obtainAttributesToCipher(body,sPAE,rPAE,noEncryptedKeys)
        if(state == False):
            return bodyBackUp, False

        body,state  = cipherBodyAttributes(body,encryptAttributes)
        if(state == False):
            return bodyBackUp, False

        return body, True

    except:
        return bodyBackUp, False

#This process consider ONLY a JSON format.
def processCypherAttr(body,sPAE,rPAE,noEncryptedKeys):

    #Here is an example of JSON structure:
    '''
    {
        "type":"Property",
        "value":"Mercedes",
        "encrypt_cpabe":{
            "type":"Property",
            "value":"att1 att2 2of2"
        }
    }
    '''

    bodyBackUp = body

    try:

        if  ("encrypt_cpabe" in body):
            if  ("type" in body["encrypt_cpabe"] and "value" in body["encrypt_cpabe"]):
                if  (body["encrypt_cpabe"]["type"]=="Property"):
                    if(str(body["encrypt_cpabe"]["value"]).strip()!=""):

                        jsonBody = {"attribute": body}

                        body,state  = cipherBodyAttributes(jsonBody,["attribute"])
                        if(state == False):
                            return bodyBackUp, False
                        else:
                            body = body["attribute"]


        return body, True

    except:
        return bodyBackUp, False

#This process consider ONLY a JSON format.
def processCypherList(body,sPAE,rPAE,noEncryptedKeys):

    #Here is an example of JSON structure:
    '''
    [
        {
            "id":"urn:ngsi-ld:Vehicle:TESTJUAN10",
            "type":"Vehicle",
            "brandName":{
                "type":"Property",
                "value":"Mercedes",
                "encrypt_cpabe":{
                    "type":"Property",
                    "value":"att1 att2 2of2"
                }
            },
            "color":{
                "type":"Property",
                "value":"Red",
                "encrypt_cpabe":{
                    "type":"Property",
                    "value":"att4 att5 2of2"
                }
            }
        },
        {
            "id":"urn:ngsi-ld:Vehicle:TESTJUAN11",
            "type":"Vehicle",
            "brandName":{
                "type":"Property",
                "value":"Mercedes",
                "encrypt_cpabe":{
                    "type":"Property",
                    "value":"att1 att2 2of2"
                }
            },
            "speed":{
                "type":"Property",
                "value":80,
                "encrypt_cpabe2":{
                    "type":"Property",
                    "value":"admin"
                }
            }
        }
    ]
    '''

    bodyBackUp = body

    try:
        #Hay que hacer un for e ir concatenando en un array de resultado, cuando hubiera un fallo en cualquiera de los elementos devuelve error
        bodyResult = []

        for m in range(len(body)):
            #body[m]
            encryptAttributes,state = obtainAttributesToCipher(body[m],sPAE,rPAE,noEncryptedKeys)
            if(state == False):
                return bodyBackUp, False

            bodyResultElem,state  = cipherBodyAttributes(body[m],encryptAttributes)
            if(state == False):
                return bodyBackUp, False
            else:
                bodyResult.append(bodyResultElem)

        return bodyResult, True

    except:
        return bodyBackUp, False


#This process consider ONLY a JSON format.
def obtainAttributesToCipher(body,sPAE,rPAE,noEncryptedKeys):

    encryptAttributes = []

    try:

        #FIND attributes must be encrypted and append into the array list encryptAttributes.
        for key in body:

            validateCriteria = True

            if(key.lower() not in noEncryptedKeys):

                for i in range(len(rPAE)):

                    if(len(rPAE[i])>0):

                        validateCriteria = True

                        for j in range(len(rPAE[i])):

                            #Obtain criteria.
                            #Uses separator to create a path array into the attribute.

                            pos0=rPAE[i][j][0].split(sPAE)
                            #Uses separator to create a path array into the attribute.
                            pos1=rPAE[i][j][1]

                            bodyAux=body[key]
                            k=0
                            lenPos0 = len(pos0)

                            while(k<lenPos0):

                                findKeyLevel = False

                                for key2 in bodyAux:

                                    if (pos0[k]!="" and pos0[k]==key2) :

                                        #Critera key - OK
                                        findKeyLevel = True

                                        if(k<lenPos0-1): 
                                            newbody = bodyAux[key2]
                                        else:
                                            if(str(pos1)!="") and str(pos1)!=str(bodyAux[key2]):
                                                validateCriteria = False
                                            break
 
                                if(findKeyLevel==False):
                                    validateCriteria = False
                                    break
                                else:
                                    if (k<lenPos0-1):
                                        bodyAux = newbody
                                    k = k + 1
                                              
                            #Critera fails
                            if(validateCriteria == False):
                                break
                        #Critera fails
                        if(validateCriteria == True):
                            break
            else:
                validateCriteria = False

            #print("processCreateEntityBody - resultado ( " + key + ")")
            #print(validateCriteria)

            if(validateCriteria):
                encryptAttributes.append(key)



        return encryptAttributes, True

    except:
        return encryptAttributes, False

def getstatusoutput(command):
    process = Popen(command, stdout=PIPE,stderr=PIPE)
    out, err = process.communicate()

    #print("out")
    #print(out)
    #print("err")
    #print(err)

    return (process.returncode, out)

#This process consider ONLY a JSON format.
def cipherBodyAttributes(body,encryptAttributes):

    bodyBackUp = body

    try:

        #print("body - BEFORE ENCRYPT")
        #print(body)
        
        #Encrypt values of attributes of array list encryptAttributes.
        for m in range(len(encryptAttributes)):
            
            try:
                #Verify, encriptation type.
                if(body[encryptAttributes[m]]["encrypt_cpabe"]["type"]=="Property"): #CPABE ENCRYPTATION
                    if(str(body[encryptAttributes[m]]["encrypt_cpabe"]["value"]).strip()!=""): #CPABE ENCRYPTATION

                        #Cipher attribute value
                        codeValue, outValue = getstatusoutput(["java","-jar","cpabe_cipher.jar",
                        str(body[encryptAttributes[m]]["encrypt_cpabe"]["value"]),
                        str(body[encryptAttributes[m]]["value"])])

                        if(codeValue == 0):
                            #Assign cipher attribute value
                            body[encryptAttributes[m]]["value"] =   outValue.decode('utf8')

                            #Change attribute value metadata
                            body[encryptAttributes[m]]["encrypt_cpabe"]["value"] = "encrypt_cpabe"
            
            except Exception as e:
                print(e)
                return bodyBackUp, False

        #print("body - AFTER ENCRYPT")
        #print(body)

        return body, True

    except:
        return bodyBackUp, False

def errorHeaders(method=None,message=None):

    headers = dict()

    '''
    #GET - headersError
    if(method.upper()=="GET"):
    else:
        #POST - headersError
        if(method.upper()=="POST"):
        else: #PATCH - headersError
            if(method.upper()=="PATCH"):
            else:#PUT - headersError                    
                if(method.upper()=="PUT"): 
    '''
    
    headers['X-Content-Type-Options'] = 'nosniff'
    headers['X-XSS-Protection'] = '1; mode=block'
    headers['Cache-Control'] = 'no-cache, no-store, max-age=0, must-revalidate'
    headers['Pragma'] = 'no-cache'
    headers['Expires'] = '0'
    headers['X-Frame-Options'] = 'DENY'
    headers['Content-Type'] = 'application/ld+json;charset=UTF-8'
    headers['Transfer-Encoding'] = 'chunked'
    ##headers['Connection'] = 'close'close

    #Second value is True because API send Transfer-Encoding=chunked header response.
    return  headers, True


def errorBody(method,code,title,details):
#    # Current time in UTC
#    now_utc = datetime.now(timezone('UTC'))
#    date_time = now_utc.strftime("%Y-%m-%dT%H:%M:%S.") + now_utc.strftime("%f")[:-3] + now_utc.strftime("%z")
#    return {'timestamp':str(date_time),'status': 500,'error':'Internal Server Error','message':'GENERAL'}
    return {"code": code, "error": title, "details": details}

#def errorCode(method):
#
#    return 500