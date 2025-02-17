from bs4 import BeautifulSoup
import os
import xml.etree.ElementTree as ET
import xml.dom.minidom as minidom
from .dbmanager.DiskDBManager import DiskDBManager
from .dbmanager.ExistDBManager import ExistDBManager
import logging
from OdinS_xacml_util import OdinS_xacml_util
from xacmleditor import ElementoPolicySet, ElementoPolicy, ElementoTarget
from .Policy import Policy

class DBConnector:
    '''Clase entrante del PAP. Una vez inicializada, se pueden obtener 
    y guardar tanto los atributos como las políticas del XACML.'''
    __connector = None
    __diskConnectorConfigPath = None
    __existConnectorConfigPath = None
    __diskConnectorProperties = None
    __existConnectorProperties = None

    DISK_CONNECTOR = "DISK_CONNECTOR"
    EXIST_CONNECTOR = "EXIST_CONNECTOR"

    logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')

    def init(self, configFilePath : str) -> None:
        '''Punto de entrada de DBConnector. Ejecuta la función loadConfigFile, 
        con el configFilePath que se le pasa por parámetro'''
        self._loadConfigFile(configFilePath)

    def _loadConfigFile(self, configFilePath : str) -> None:
        '''Carga el archivo de configuración y lo parsea
        para poder obtener las propiedades necesarias.'''
        try:
            with open(configFilePath, 'r') as f:
                bsData = BeautifulSoup(f, 'xml')

            '''Cogemos cada una de las propiedades y las parseamos'''
            root = bsData.find('properties')

            for linea in root.find_all():
                self._setProperty(linea.get('key'), linea.text.strip())

            diskConnectorExists = os.path.isfile(self.__diskConnectorConfigPath)

            if diskConnectorExists:
                diskConnectorDict = {}
                with open(self.__diskConnectorConfigPath, 'r') as dcFile:
                    dcData = BeautifulSoup(dcFile, 'xml')
                rootDC = dcData.find('properties')
                for linea in rootDC.find_all():
                    diskConnectorDict[linea.get('key')] = linea.text.strip()
                self.__diskConnectorProperties = diskConnectorDict
            else:
                logging.warning("disk connector config file does not exist: " + self.__diskConnectorConfigPath)

            existConnectorExists = os.path.isfile(self.__existConnectorConfigPath)

            if existConnectorExists:
                existConnectorDict = {}
                with open(self.__existConnectorConfigPath, 'r') as ecFile:
                    ecData = BeautifulSoup(ecFile, 'xml')
                rootEC = ecData.find('properties')
                for linea in rootEC.find_all():
                    existConnectorDict[linea.get('key')] = linea.text.strip()
                self.__diskConnectorProperties = existConnectorDict
            else:
                logging.warning("eXist connector config file does not exist: " + self.__existConnectorConfigPath)

        except IOError as ex:
            logging.critical(ex)

    def saveConfig(self, connector : str, properties : dict) -> None:
        '''Guardamos las propiedades para el conector indicado en XML.'''
        try:
            if connector == self.DISK_CONNECTOR:
                self.__diskConnectorProperties = properties
                root = ET.Element('properties')
                for key, value in properties.items():
                    entry = ET.SubElement(root, 'entry', key=key)
                    entry.text = value
                xmlString = minidom.parseString(ET.tostring(root)).toprettyxml(indent="", encoding="UTF-8")
                with open(self.__diskConnectorConfigPath, "w") as xml:
                    xml.write(xmlString)
            elif connector == self.EXIST_CONNECTOR:
                self.__diskConnectorProperties = properties
                root = ET.Element('properties')
                for key, value in properties.items():
                    entry = ET.SubElement(root, 'entry', key=key)
                    entry.text = value
                xmlString = minidom.parseString(ET.tostring(root)).toprettyxml(indent="", encoding="UTF-8")
                with open(self.__existConnectorConfigPath, "w") as xml:
                    xml.write(xmlString)
        except IOError as ex:
            logging.critical(ex)

            self.__connector = connector
            self.__connector = self.DBConnector(self.__connector)

    def _setProperty(self, current : str, property : str) -> None:
        '''Asigna la propiedad deseada al conector actual'''
        if current == 'CONNECTOR':
            self.__connector = property
        elif current == 'DISK_CONNECTOR_CONFIG_PATH':
            self.__diskConnectorConfigPath = property
        elif current == 'EXIST_CONNECTOR_CONFIG_PATH':
            self.__existConnectorConfigPath = property

    def getCurrentDiskParameter(self) -> dict:
        '''Devuelve las propiedades de DiskDBManager'''
        return self.__diskConnectorProperties

    def getCurrentExistParameter(self) -> dict:
        '''Devuelve las propiedades de ExistDBManager'''
        return self.__existConnectorProperties

    def getDiskParameters(self) -> list:
        '''Devuelve los parámetros necesarios de DiskDBManager'''
        return DiskDBManager.neededConfigParameters(self)

    def getExistParameters(self) -> list:
        '''Devuelve los parámetros necesarios de ExistDBManager'''
        return ExistDBManager.neededConfigParameters(self)

    def getCurrentDBConnectorName(self) -> str:
        '''Devuelve el nombre del conector actual'''
        return self.__connector

    ALL_SELECTED = '[ALL-SELECTED]'
    ruleMap = {}
    __dbConnector = None
    __dbManager = None

    def getDBConnector(self):
        '''Devuelve el objeto DBConnector actual'''
        if self.__dbConnector is None:
            self.__dbConnector = self.DBConnector(self.__connector)
        return self.__dbConnector

    def DBConnector(self, connector : str):
        '''Constructor del objeto DBConnector'''
        if connector == self.DISK_CONNECTOR:
            self.__dbManager = DiskDBManager(self.__diskConnectorProperties)
        elif connector == self.EXIST_CONNECTOR:
            self.__dbManager = ExistDBManager(properties=self.__diskConnectorProperties)
        return self # Devuelve "self", es decir, el objeto en sí

    def loadPolicies(self, header : str) -> list:
        '''Devuelve las políticas actuales en una lista'''
        policies = []

        read = os.environ.get('READ')
        match read:
            case 'file':
                policySet = self.__dbManager.retrievePolicySet(header)
            case 'DLT':
                policySet = self.__dbManager.retrievePolicySetJSON(header)

        if policySet == "" or policySet is None:
            return policies
        
        try:
            odins_xacml_util = OdinS_xacml_util(xacml=policySet)
            elementoPrincipal = odins_xacml_util.getPrincipal()
            if isinstance(elementoPrincipal, ElementoPolicySet):
                elementoPolicySet = elementoPrincipal
                hijos = odins_xacml_util.getChildren(elementoPolicySet, ElementoPolicy.TIPO_POLICY)
                for hijo in hijos:
                    policy = Policy(odins_xacml_util=odins_xacml_util, elementoPolicy=hijo)
                    policies.append(policy)
        except Exception as e:
            logging.critical(e)

        return policies

    def importPolicy(self, data : str) -> Policy | None:
        '''Importa la política deseada y la devuelve. Si falla, devuelve None'''
        try:
            odins_xacml_util = OdinS_xacml_util(data)
            ePrincipal = odins_xacml_util.getPrincipal()
            if isinstance(ePrincipal, ElementoPolicy):
                elementoPolicy = ePrincipal
                policy = Policy(odins_xacml_util=odins_xacml_util, elementoPolicy=elementoPolicy)
                return policy
        except Exception as e:
            logging.critical(e)

        return None

    def storeXACMLAttributes(self, attributes : list, header : str) -> None:
        '''Se guardan los atributos de la lista en el DBManager actual'''
        self.__dbManager.storeXACMLAttributes(attributes, header)

    def storeJSONAttributes(self, attributes : list, header : str) -> None:
        '''Se guardan los atributos de la lista en formato JSON en el DBManager actual'''
        self.__dbManager.storeJSONAttributes(attributes, header)

    def getXACMLAttributes(self, header : str) -> list:
       '''Devuelve los atributos del DBManager actual'''
       return self.__dbManager.getXACMLAtribbutes(header)

    def getJSONAttributes(self, header : str) -> list:
        '''Devuelve los atributos del DBManager actual en formato JSON'''
        return self.__dbManager.getJSONAttributes(header)

    def storePolicies(self, policies : list, header : str) -> None:
        '''Se guardan las políticas de la lista en el DBManager actual, con
        el formato correspondiente'''
        odins_xacml_util = OdinS_xacml_util()
        elementoPolicySet = odins_xacml_util.createPrincipal(ElementoPolicySet.TIPO_POLICYSET)
        elementoPolicySet.getAtributos()["xmlns"] = "urn:oasis:names:tc:xacml:2.0:policy:schema:os"
        elementoPolicySet.getAtributos()['PolicyCombiningAlgId'] = "urn:oasis:names:tc:xacml:1.0:policy-combining-algorithm:first-applicable"
        elementoPolicySet.getAtributos()["PolicySetId"] = "POLICY_SET"
        target = odins_xacml_util.createChild(elementoPolicySet, ElementoTarget.TIPO_TARGET)

        for policy in policies:
            odins_xacml_util.insertChild(elementoPolicySet, policy.getOdinS_XACML())
        

        read = os.environ.get('READ')
        match read:
            case 'file':
                self.__dbManager.storePolicySet(str(odins_xacml_util), header)
            case 'DLT':
                self.__dbManager.storePolicySet(str(odins_xacml_util), header)  # Using files as backups
                self.__dbManager.storePolicySetJSON(str(odins_xacml_util), header)

        return str(odins_xacml_util)
    def getDomains(self, read : str) -> list:
        ''' Devuelve una lista de los dominios'''
        return self.__dbManager.getDomains(read)