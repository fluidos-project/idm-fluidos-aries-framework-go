from flask import Flask, request
from odins.xacml.webpap import DBConnector, Action, Subject, Resource, Policy, Rule
import json
import logging
import traceback
import os
from waitress import serve
from datetime import datetime
from flask_cors import CORS

# Configuración inicial del serivdor API
app = Flask(__name__)
CORS(app)

dbConnector = DBConnector()
dbConnector.init('WebContent/WEB-INF/config.xml')
dbConnector.DBConnector("DISK_CONNECTOR")

# Configuración del logger
logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')
logging.getLogger().setLevel(logging.INFO)

# Opening JSON file
f = open('PAPConfigData/SubjectIdTypes.json')
# returns JSON object as a dictionary
configuration = json.load(f)
# Closing file
f.close()

HASH_Integration = 0
try:
    HASH_Integration = int(os.getenv('HASH_Integration'))
except Exception as e:
    logging.error(e)
    HASH_Integration = 0

try:
    HASH_Protocol = str(os.getenv('HASH_Protocol'))
except Exception as e:
    logging.error(e)
    HASH_Protocol = "https"

if (str(HASH_Protocol).upper() == "None".upper()) :
    HASH_Protocol = "https"

try:
    HASH_Host = str(os.getenv('HASH_Host'))
except Exception as e:
    logging.error(e)
    HASH_Host = "localhost"

if (str(HASH_Host).upper() == "None".upper()) :
    HASH_Host = "localhost"

HASH_Port = 8080
try:
    HASH_Port = int(os.getenv('HASH_Port'))
except Exception as e:
    logging.error(e)
    HASH_Port = 8080

HASH_BaseURL = str(HASH_Protocol) + "://" + str(HASH_Host) + ":" + str(HASH_Port)

try:
    HASH_Get_Resource = str(os.getenv('HASH_Get_Resource'))
except Exception as e:
    logging.error(e)
    HASH_Get_Resource = "/chain/events?entityid=urn:ngsi-ld:domain:{{domain}}"

if (str(HASH_Get_Resource).upper() == "None".upper()) :
    HASH_Get_Resource = "/chain/events?entityid=urn:ngsi-ld:domain:{{domain}}"

try:
    HASH_Post_Resource = str(os.getenv('HASH_Post_Resource'))
except Exception as e:
    logging.error(e)
    HASH_Post_Resource = "/chain/publish"

if (str(HASH_Post_Resource).upper() == "None".upper()) :
    HASH_Post_Resource = "/chain/publish"

try:
    HASH_Patch_Resource = str(os.getenv('HASH_Patch_Resource'))
except Exception as e:
    logging.error(e)
    HASH_Patch_Resource = "/chain/publish"

if (str(HASH_Patch_Resource).upper() == "None".upper()) :
    HASH_Patch_Resource = "/chain/publish"

@app.route("/")
def defaultRoute():
    '''Punto de entrada'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")
    return "PAP API", 200

@app.get("/obtaindomains")
def obtainDomains():

    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:
        read_accepted_values = os.environ.get('READ_ACCEPTED_VALUES').split(':')
        read = os.environ.get('READ')

        if read not in read_accepted_values:
            raise Exception(f"Env variable error: Read method not supported, accepted values: {read_accepted_values}")
        
        domains = dbConnector.getDomains(read)
        data = {
            "domains" : domains
        }
        jsonData = json.dumps(data)
        return jsonData, 200, {'Content-Type' : 'application/json'}
    
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\"," 
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}

@app.get("/obtainattributes")
def obtainAttributes():
    '''obtainattributes endpoint. Sólo se pueden enviar peticiones GET'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:

        header = request.headers.get('domain')
        read = os.environ.get('READ')
        read_accepted_values = os.environ.get('READ_ACCEPTED_VALUES').split(':')

        match read:
            case 'file':
                atts = dbConnector.getXACMLAttributes(header)
            case 'DLT':
                atts = dbConnector.getJSONAttributes(header)
            case _:
                raise Exception(f"Env variable error: Read method not supported, accepted values: {read_accepted_values}") 
        
        '''Se parsean los atributos obtenidos: para cada clave del diccionario de un atributo, se suprime
        el _XACMLAttributeElement__ generado, y se genera un diccionario con type y value para cada valor.'''
        attsJSON = json.dumps(atts, 
                              default=lambda o: 
                              {k.replace('_XACMLAttributeElement__', ''): {"type": "Property", "value": v} 
                               for k, v in o.__dict__.items()}, 
                               indent=4)
        logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")
        return attsJSON, 200, {'Content-Type': 'application/json'}
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\"," 
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}

@app.patch("/saveattributes")
def saveAttributes():
    '''saveattributes endpoint. Sólo se pueden enviar peticiones PATCH'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:

        atts = []
        data = request.get_json()
        # Se leen todos los objetos de la lista y se reconstruyen para guardar los atributos

        for obj in data:
            elem = None
            name = obj['name']['value']
            xacml_id = obj['xacml_id']['value']
            xacml_DataType = obj['xacml_DataType']['value']
            match obj['sortedValue']['value']:
                case 'action':
                    elem = Action(name, xacml_id, xacml_DataType)
                case 'resource':
                    elem = Resource(name, xacml_id, xacml_DataType)
                case 'subject':
                    elem = Subject(name, xacml_id, xacml_DataType)
            atts.append(elem)

        header = request.headers.get('domain')
        read = os.environ.get('READ')
        read_accepted_values = os.environ.get('READ_ACCEPTED_VALUES').split(':')

        match read:
            case 'file':
                dbConnector.storeXACMLAttributes(atts, header)
            case 'DLT':
                dbConnector.storeXACMLAttributes(atts, header) # Using files as backups
                dbConnector.storeJSONAttributes(atts, header)
            case _:
                raise Exception(f"Env variable error: Read method not supported, accepted values: {read_accepted_values}") 

        logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")
        return "{\"response\": \"XACMLAttributes saved successfully\"}", 200, {'Content-Type': 'application/json'}
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\"," 
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}

@app.get("/obtainpolicies")
def obtainPolicies():
    '''obtainpolicies endpoint. Sólo se pueden enviar peticiones GET'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:
        # Carga las políticas y las parsea a JSON a través de la función policyToJSON
        header = request.headers.get('domain')

        read_accepted_values = os.environ.get('READ_ACCEPTED_VALUES').split(':')
        read = os.environ.get('READ')
               
        if read not in read_accepted_values:
            raise Exception(f"Env variable error: Read method not supported, accepted values: {read_accepted_values}")

        policies = dbConnector.loadPolicies(header)
        
        policiesJSON = json.dumps(policies, default=policyToJSON, indent=4)
        logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")
        return policiesJSON, 200, {'Content-Type': 'application/json'}
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\"," 
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}

@app.patch("/savepolicies")
def savePolicies():
    '''savepolicies endpoint. Sólo se pueden enviar peticiones PATCH'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:
        policies = []
        data = request.get_json()
        '''Se leen todos los objetos de la lista y se van reconstruyendo para 
        obtener las políticas y reglas correspondientes, tanto si existen como si no'''
        for obj in data:
            policy = Policy(policyId=obj.get('PolicyId', {}).get('value'))
            rules = obj.get('Rules', [])
            ruleDict = {}
            ruleList = []
            for r in rules:
                ruleID = r.get('RuleId', {}).get('value')
                rule = Rule(name=ruleID)
                effect = r.get('Effect', {}).get('value')
                if effect and effect != 'Permit':
                    rule.setPolicy(effect)
                subjects = []
                resources = []
                actions = []
                for subject in r.get('Subjects', []):
                    name = subject.get('AttributeValue')
                    id = subject.get('AttributeId')
                    subject = Subject(name, id, "#string")
                    subjects.append(subject)
                for resource in r.get('Resources', []):
                    name = resource.get('AttributeValue')
                    id = resource.get('AttributeId')
                    resource = Resource(name, id, "#string")
                    resources.append(resource)
                for action in r.get('Actions', []):
                    name = action.get('AttributeValue')
                    id = action.get('AttributeId')
                    action = Action(name, id, "#string")
                    actions.append(action)
                rule.setSubjects(subjects)
                rule.setResources(resources)
                rule.setActions(actions)
                ruleDict[ruleID] = rule
                ruleList.append(rule)
            policy.setRules(ruleDict=ruleDict)
            policy.setRules(ruleList=ruleList)
            policies.append(policy)

        header = request.headers.get('domain')
        xacml_policies = dbConnector.storePolicies(policies, header)

        logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")

        if (HASH_Integration == 1):
            responseStoreHASH = storeHASH(header, xacml_policies)

        return "{\"response\": \"Policies saved successfully\"}", 200, {'Content-Type': 'application/json'}
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\"," 
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}

def policyToJSON(policy):
    """
    Convierte el objeto Policy a un diccionario JSON.
    """
    # Se crea un diccionario vacío para ir añadiendo las propiedades de la política
    policyDict = {}
    try:
        # Se comprueba si la política tiene PolicyId y RuleCombiningAlgId y se añaden al diccionario
        if hasattr(policy, 'getPolicyId') and policy.getPolicyId() is not None:
            policyDict['PolicyId'] = {"type": "Property", "value": policy.getPolicyId()}
        if hasattr(policy, 'getCombiningAlg') and policy.getCombiningAlg() is not None:
            policyDict['RuleCombiningAlgId'] = {"type": "Property", "value": policy.getCombiningAlg()}

        # Se comprueba si la política tiene reglas y se añaden al diccionario
        if hasattr(policy, 'getRules') and policy.getRules():
            rules = []
            for rule in policy.getRules():
                if rule is not None:
                    # Se crea un diccionario vacío para ir añadiendo las propiedades de la regla
                    rule_dict = {}
                    if hasattr(rule, 'getName') and rule.getName() is not None:
                        rule_dict['RuleId'] = {"type": "Property", "value": rule.getName()}
                    if hasattr(rule, 'getPolicy') and rule.getPolicy() is not None:
                        rule_dict['Effect'] = {"type": "Property", "value": rule.getPolicy()}

                    # Se comprueba si la regla tiene subjects, resources y actions y se añaden al diccionario
                    for attr_name, attr_method in [('Subjects', 'getSubjects'), ('Resources', 'getResources'), ('Actions', 'getActions')]:
                        if hasattr(rule, attr_method):
                            attr_list = getattr(rule, attr_method)()
                            if len(attr_list)>0:
                                if any([attr.getName() is not None and attr.getXACMLID() is not None for attr in attr_list]):
                                    rule_dict[attr_name] = [
                                        {"AttributeValue": attr.getName(), "AttributeId": attr.getXACMLID()}
                                        for attr in attr_list if attr.getName() is not None and attr.getXACMLID() is not None
                                    ]
                            else:
                                rule_dict[attr_name] = []
                            
                    # Se añade el diccionario de la regla a la lista de reglas
                    rules.append(rule_dict)

            # Se añade la lista de reglas al diccionario de la política
            policyDict['Rules'] = rules
    except Exception as e:
        logging.info("\033[31mError parsing policies to JSON\033[0m")

    # Se devuelve el diccionario JSON generado
    return policyDict

@app.get("/obtainSubjectIdTypes")
def obtainSubjectIdTypes():
    '''obtainSubjectIdTypes endpoint. Sólo se pueden enviar peticiones GET'''
    time = datetime.now().strftime("%d/%b/%Y %H:%M:%S")
    try:
        logging.info(f"\033[32m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 200\033[0m")
        return configuration, 200, {'Content-Type': 'application/json'}
    except Exception as e:
        # Se hace un traceback de dónde ha ocurrido el error y se devuelve por JSON
        tb = traceback.extract_tb(e.__traceback__)
        file, line_number, function, line_text = tb[-1]
        logging.error(f"\033[31m{request.remote_addr} - [{time}] \"{request.method} {request.path} {request.environ.get('SERVER_PROTOCOL')}\" 500\033[0m")
        return ("{\"response\": \"Error: " + str(e) + "\","
                "\"file\": \"" + file + "\","
                "\"line\": \"" + str(line_number) + "\","
                "\"function\": \"" + function + "\","
                "\"text\": \"" + line_text + "\"}"), 500, {'Content-Type': 'application/json'}


import requests
import xmltodict
from hashlib import sha256

def storeHASH(domain : str, xacml_policies : str):

    response = ""
    try:
        '''Guarda el HASH basado en las políticas actuales'''

        xpars = xmltodict.parse(xacml_policies)

        xacml_policies = json.dumps(xpars)

        #Calculate HASH
        XACMLdigest = sha256(xacml_policies.encode('utf-8')).hexdigest()

        if domain is None:
            domain = os.environ.get('DEFAULT_POLICY_HEADER')

        response = requests.get( f'{HASH_BaseURL}' + f'{os.environ.get("HASH_Get_Resource")}' + f'{domain}?attrs=version')

        if (response.status_code == 404):
            version = 1
        else:
            data = response.json()
            version = int(data.get("version").get("value"))
            version += 1

        if (version == 1):

            payload = {
                    "id": f"urn:ngsi-ld:domain:{domain}",
                    "type":"domain",
                    "version":{
                        "type": "Property",
                        "value": str(version)
                    },
                    "digest":{
                        "type": "Property",
                        "value": str(XACMLdigest)
                    },
                    "timestamp": {
                        "type": "Property",
                        "value": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                    }
            }

            response = requests.post(url=f'{HASH_BaseURL}' + f'{os.environ.get("HASH_Post_Resource")}',
                headers={"content-type" : "application/json"}, data= json.dumps(payload))
        else:
            payload = {
                "version":{
                    "type": "Property",
                    "value": str(version)
                },
                "digest":{
                    "type": "Property",
                    "value": str(XACMLdigest)
                },
                "timestamp": {
                    "type": "Property",
                    "value": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                }
            }
            response = requests.patch(url=f'{HASH_BaseURL}' + f'{os.environ.get("HASH_Patch_Resource")}' + f'{domain}/attrs'
                , headers={"content-type" : "application/json"}, data= json.dumps(payload))

    except Exception as e:
        logging.info("\033[31mError storing HASH\033[0m")

        logging.info(str(e))

    return response


if __name__ == '__main__':
    logging.info("\033[32mPAP API started\033[0m")
    serve(app, host="0.0.0.0", port=8080)
