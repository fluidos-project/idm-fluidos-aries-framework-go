from .DBManager import DBManager
import os
from odins.xacml.webpap.XACMLAttributeElement import XACMLAttributeElement
import logging
import json
import requests
import re
import xml.etree.ElementTree as ET
from odins.xacml.webpap.Action import Action
from odins.xacml.webpap.Resource import Resource
from odins.xacml.webpap.Subject import Subject
from datetime import datetime
import xml.dom.minidom as minidom


class DiskDBManager(DBManager):
    '''Clase heredada de DBManager. Se guardan los archivos en el disco duro'''
    __policySetPath = None
    __xacmlAttributesPath = None
    xacmlAttributes = []

    
    logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')

    def __init__(self, properties : dict) -> None:
        '''Se inicializa el objeto DiskDBManager'''
        self.__policySetPath = properties['POLICY_SET_PATH']
        self.__xacmlAttributesPath = properties["XACMLATTS_PATH"]

    def neededConfigParameters(self) -> list:
        '''Devuelve los parámetros necesarios para el objeto DiskDBManager'''
        return ['POLICY_SET_PATH', 'XACMLATTS_PATH']

    def retrievePolicySet(self, header : str) -> str:
        '''Devuelve el XML de las políticas actuales'''

        if header is None:
            header = "continue-a"

        print(f'Retrieving the PolicySet from {self.__policySetPath} + {header}.xml')
        try:
            with open(self.__policySetPath + f"{header}.xml", 'r') as f:
                policySet = f.read()
            return policySet
        except IOError as ex:
            logging.critical(ex)
        return ''

    def retrievePolicySetJSON(self, header : str) -> str:

        if header is None:
            header = os.environ.get('DEFAULT_POLICY_HEADER')

        #response = requests.get( f'{os.environ.get("DLT_POLICIES")}' + f'{header}?options=sysAttrs')
        response = requests.get( f'{os.environ.get("DLT_POLICIES")}' + f'{header}')
        if (response.status_code != 404):    
            data = response.json()

            policySet = data.get('xacml').get('value').replace('\\n', '\n').replace('\\"', '"')

            return policySet
        else:
            return ''  

    def storePolicySet(self, policySet : str, header : str) -> None:
        '''Guarda las políticas actuales'''

        if header is None:
            header = "continue-a"

        try:
            with open(self.__policySetPath + f"{header}.xml", 'w') as f:
                f.write(policySet)
        except IOError as ex:
            logging.critical(ex)
   
    def storePolicySetJSON(self, policySet : str, header : str) -> None:
        '''Guarda las políticas actuales en formato JSON'''

        if header is None:
            header = os.environ.get('DEFAULT_POLICY_HEADER')

        response = requests.get( f'{os.environ.get("DLT_POLICIES")}' + f'{header}?attrs=version')
        
        if (response.status_code == 404):
            version = 1
        else:
            data = response.json()
            version = int(data.get("version").get("value"))
            version += 1

        policySet = str(policySet)
        policySet = policySet.split("\n")
        policySet = "".join(policySet)

        if (version == 1):

            payload = {
                "id": f"urn:ngsi-ld:xacml:{header}",
                "type" : "xacml",
                "version" : {
                    "type" : "Property" , 
                    "value" : str(version),
                },
                "xacml" : {
                    "type" : "Property", 
                    "value" : policySet
                },
                "timestamp" : {
                    "type" : "Property",
                    "value" : datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                }
            }

            requests.post(url=f'{os.environ.get("DLT_POST")}',
                headers={"content-type" : "application/json"}, data= json.dumps(payload))
        else:
            payload = {
                "version" : {
                    "type" : "Property" , 
                    "value" : str(version),
                },
                "xacml" : {
                    "type" : "Property", 
                    "value" : policySet
                },
                "timestamp" : {
                    "type" : "Property",
                    "value" : datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                }
            }
            response = requests.patch(url= f'{os.environ.get("DLT_POLICIES")}' + f'{header}/attrs'
                , headers={"content-type" : "application/json"}, data= json.dumps(payload))


    def storeXACMLAttributes(self, attributes : list, header : str) -> None:
        '''Guarda los atributos actuales en la ruta de los atributos XACML'''
        root = ET.Element('attributes')
        root.set("storingDate", datetime.now().strftime("%Y-%m-%dT%H:%M:%S.%fZ"))
        attsElem = ET.SubElement(root, "attributes")
        for attribute in attributes:
            entry = ET.SubElement(attsElem, 'attribute')
            entry.set("name", attribute.getName())
            entry.set("xacml_id", attribute.getXACMLID())
            entry.set("sortedValue", attribute.getSortedValue())
            entry.set("xacml_DataType", attribute.getDataType())
        xmlString = minidom.parseString(ET.tostring(root)).toprettyxml(indent="  ", encoding="UTF-8")

        if header is None:
            header = os.environ.get('DEFAULT_ATTRIBUTES_HEADER')

        with open(self.__xacmlAttributesPath + f"{header}.xml", "wb") as xml:
            xml.write(xmlString)

    def storeJSONAttributes(self, attributes : list, header : str) -> None:
        root = ET.Element('attributes')
        root.set("storingDate", datetime.now().strftime("%Y-%m-%dT%H:%M:%S.%fZ"))
        attsElem = ET.SubElement(root, "attributes")
        for attribute in attributes:
            entry = ET.SubElement(attsElem, 'attribute')
            entry.set("name", attribute.getName())
            entry.set("xacml_id", attribute.getXACMLID())
            entry.set("sortedValue", attribute.getSortedValue())
            entry.set("xacml_DataType", attribute.getDataType())
        xmlString = minidom.parseString(ET.tostring(root)).toprettyxml(indent="  ", encoding="UTF-8") 

        if header is None:
            header = os.environ.get('DEFAULT_ATTRIBUTES_HEADER')

        response = requests.get(f'{os.environ.get("DLT_ATTRIBUTES")}' + f'{header}?attrs=version')

        if (response.status_code == 404):
            version = 1
        else:
            data = response.json()
            version = int(data.get("version").get("value"))
            version += 1

        if (version == 1):

            payload = {
                "id": f"urn:ngsi-ld:attributes:{header}",
                "type" : "attributes",
                "version" : {
                    "type" : "Property" , 
                    "value" : str(version),
                },
                "attributes" : {
                    "type" : "Property", 
                    "value" : xmlString.decode('utf-8')
                },
                "timestamp" : {
                    "type" : "Property",
                    "value" : datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                }
            }

            requests.post(url=f'{os.environ.get("DLT_POST")}',
                headers={"content-type" : "application/json"}, data= json.dumps(payload))
        else:
            payload = {
                "version" : {
                    "type" : "Property" , 
                    "value" : str(version),
                },
                "attributes" : {
                    "type" : "Property", 
                    "value" : xmlString.decode('utf-8')
                },
                "timestamp" : {
                    "type" : "Property",
                    "value" : datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                }
            }
            response = requests.patch(url= f'{os.environ.get("DLT_ATTRIBUTES")}' + f'{header}/attrs'
                , headers={"content-type" : "application/json"}, data= json.dumps(payload))   

    def getXACMLAtribbutes(self, header : str) -> list:
        '''Devuelve los atributos XACML actuales'''

        if header is None:
            header = os.environ.get('DEFAULT_ATTRIBUTES_HEADER')

        xacmlAttributes = []
        parsedXACMLAtts = self.parseXACMLAttributes(self.__xacmlAttributesPath + f"{header}.xml")
        xacmlAttributes.extend(parsedXACMLAtts)
        return xacmlAttributes

        '''xacmlAttributes = []
        try:
            for filename in os.listdir(self.__xacmlAttributesPath):
                if filename.endswith(".xml"):
                    print("Searching for Attributes in the file: " + self.__xacmlAttributesPath + filename)
                    parsedXACMLAtts = self.parseXACMLAttributes(self.__xacmlAttributesPath + filename)
                    xacmlAttributes.extend(parsedXACMLAtts)
        except Exception as ex:
            logging.critical(ex)
            raise Exception(ex)
        return xacmlAttributes'''

    def getJSONAttributes(self, header : str) -> list:

        jsonAttributes = []

        if header is None:
            header = os.environ.get('DEFAULT_ATTRIBUTES_HEADER')

        #response = requests.get(f'{os.environ.get("DLT_ATTRIBUTES")}' + f'{header}?options=sysAttrs')
        response = requests.get(f'{os.environ.get("DLT_ATTRIBUTES")}' + f'{header}')

        if response.status_code != 404:
            parsedJSONAtts = self.parseJSONAttributes(response.json())
            jsonAttributes.extend(parsedJSONAtts)

        return jsonAttributes

    def parseXACMLAttributes(self, resourceFile : str) -> list:
        '''Parsea los atributos actuales y devuelve una lista 
        con los objetos creados'''
        attributes = []
        
        try:
            XMLFile = ET.parse(resourceFile)
            root = XMLFile.getroot()
            attributesList = root.find("attributes")
        except Exception as e:
            return attributes

        for attribute in attributesList:
            name = attribute.get("name")
            xacmlID = attribute.get("xacml_id")
            sortedValue = attribute.get("sortedValue")
            xacml_DataType = attribute.get("xacml_DataType")

            if name is not None and xacmlID is not None and sortedValue is not None:
                if sortedValue == "resource":
                    attributes.append(Resource(name, xacmlID, xacml_DataType))
                elif sortedValue == "action":
                    attributes.append(Action(name, xacmlID, xacml_DataType))
                elif sortedValue == "subject":
                    attributes.append(Subject(name, xacmlID, xacml_DataType))       
        return attributes

    def parseJSONAttributes(self, data : str) -> list:
        attributes = []
  
        attributesValue = data.get('attributes').get('value').replace('\\n', '\n').replace('\\"', '"')

        if attributesValue:
            try:
                root = ET.fromstring(attributesValue)
                attributesList = root.find('attributes')
                
                for attribute in attributesList:
                    name = attribute.get("name")
                    xacmlID = attribute.get("xacml_id")
                    sortedValue = attribute.get("sortedValue")
                    xacml_DataType = attribute.get("xacml_DataType")

                    if name is not None and xacmlID is not None and sortedValue is not None:
                        if sortedValue == "resource":
                            attributes.append(Resource(name, xacmlID, xacml_DataType))

                        elif sortedValue == "action":
                            attributes.append(Action(name, xacmlID, xacml_DataType))

                        elif sortedValue == "subject":
                            attributes.append(Subject(name, xacmlID, xacml_DataType))
            except ET.ParseError as e:
                logging.error(f"Error parsing JSON: {str(e)}")

        return attributes        

    def getDomains(self, read : str) -> list:

        domains = []
        match read:
            case 'file':
                try:
                    for filename in os.listdir(self.__policySetPath):
                        if filename.endswith(".xml"):
                            domains.append(filename.replace(".xml", ""))                   
                    return domains        
                except Exception as ex:
                    logging.critical(ex)
                    raise Exception(ex)
                
            case 'DLT':
                response = requests.get(f"{os.environ.get('DLT_GET_DOMAINS')}")
                data = response.text
                IDs = re.findall('urn:ngsi-ld:xacml:[0-9a-zA-Z_-]+', data)
                domains = [ID.replace('urn:ngsi-ld:xacml:', '') for ID in IDs]
                return domains
        