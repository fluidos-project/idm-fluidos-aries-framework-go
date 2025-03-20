from .DBManager import DBManager
from pyexistdb.db import ExistDB
import logging

class ExistDBManager(DBManager):
    '''Clase heredada de DBManager. Se guarda la información en una base de datos eXist'''
    __existDBURL = None
    __existDBUser = None
    __existDBPassword = None
    __policiesURI = None
    __xacmlAttsURI = None
    __baseURI = None

    DRIVER = "org.exist.xmldb.DatabaseImpl"
    ALL_ROLE_NAME = "ALL-ROLES"
    ALL_PERMIT_ATTRIBUTES = "ALL-PERMIT-ATTRIBUTES"
    NULL_ATTRIBUTE = "NULL_ATTRIBUTE"

    logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')

    def __init__(self, properties=None, existDBURL=None, existDBUser=None, existDBPassword=None) -> None:
        '''Se inicializa el objeto ExistDBManager'''
        if properties is not None:
            pass
        else:
            self.__existDBURL = existDBURL
            self.__existDBUser = existDBUser
            self.__existDBPassword = existDBPassword
            self.__baseURI = existDBURL + "pap/"
            self.__policiesURI = self.__baseURI + "policies/"
            self.__xacmlAttsURI = self.__baseURI + "xacmlAtt/"


    def neededConfigParameters(self) -> list:
        '''Devuelve los parámetros necesarios para el objeto ExistDBManager'''
        return ['EXIST_DB_URL', 'EXIST_DB_USER', 'EXIST_DB_PASSWORD', 'POLICIES_URI']

    def retrievePolicySet(self) -> None:
        '''Operación no soportada.'''
        raise NotImplementedError("Not supported yet.")

    def storePolicySet(self, policySet : str) -> None:
        '''Operación no soportada.'''
        raise NotImplementedError("Not supported yet.")

    def storeXACMLAttributes(self, attributes : list) -> None:
        '''Operación no soportada.'''
        raise NotImplementedError("Not supported yet.")

    def getXACMLAtribbutes(self) -> None:
        '''Operación no soportada.'''
        raise NotImplementedError("Not supported yet.")
    
    def __createExistCollection(self, URI : str, collection : str) -> bool:
        '''Crea una colección en la base de datos eXist. Devuelve True si ha sido creada, False si no.'''
        try:
            eXDB = ExistDB(self.__existDBURL, self.__existDBUser, self.__existDBPassword)
            col = eXDB.hasCollection(URI + collection)
            if col is None:
                eXDB.createCollection(URI + collection)
            return True
        except Exception as ex:
            logging.critical(ex)

        return False

    def createBBDD(self) -> bool:
        '''Crea una base de datos. Devuelve True si ha sido creada, False si no.'''
        if (self.__createExistCollection(self.__existDBURL, "pap") and self.__createExistCollection(self.__baseURI, "policies") and
            self.__createExistCollection(self.__baseURI, "xacmlAtts")):
            return True
        else:
            return False


    def checkExist(self) -> bool:
        '''Comprueba si existe la base de datos. Devuelve True si existe, False si no.'''
        try:
            eXDB = ExistDB(self.__existDBURL, self.__existDBUser, self.__existDBPassword)
            col = eXDB.hasCollection(self.__policiesURI)
            if col is None:
                return False
            return True
        except Exception as e:
            return False
