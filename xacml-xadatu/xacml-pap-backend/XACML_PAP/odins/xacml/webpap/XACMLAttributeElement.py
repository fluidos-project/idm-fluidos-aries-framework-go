from collections import OrderedDict

class XACMLAttributeElement():
    '''Clase base de módulos como Action, Resource o Subject'''
    __name = None # Nombre del atributo
    __xacml_id = None # ID del Atributo
    __sortedValue = None # SortedValue del atributo
    __xacml_DataType = None # DataType del atributo

    def __init__(self, name="", xacmlID="", sortedValue=None, xacml_DataType=None) -> None:
        '''Constructor del objeto'''
        self.__name = name
        self.__xacml_id = xacmlID
        self.__sortedValue = sortedValue
        self.__xacml_DataType = xacml_DataType


    def getName(self) -> str:
        '''Devuelve el nombre del atributo'''
        return self.__name
    
    def getXACMLID(self) -> str:
        '''Devuelve el ID del atributo XACML'''
        return self.__xacml_id
    
    def getSortedValue(self) -> str:
        '''Devuelve el sortedValue del atributo XACML'''
        return self.__sortedValue
    
    def getDataType(self) -> str:
        '''Devuelve el tipo de dato del atributo XACML'''
        return self.__xacml_DataType
    
    def getKeyForMap(self) -> str:
        '''Devuelve la clave del diccionario del atributo'''
        return (self.getName() + "<()>" + self.getXACMLID())
    
    def setName(self, name) -> None:
        '''Establece un nuevo nombre para el atributo'''
        self.__name = name

    def setXACMLID(self, xacmlID : str) -> None:
        '''Establece un nuevo ID para el atributo XACML'''
        self.__xacml_id = xacmlID

    def setSortedValue(self, sortedValue : str) -> None:
        '''Establece un nuevo sortedValue para el atributo XACML'''
        self.__sortedValue = sortedValue

    def setDataType(self, dataType : str) -> None:
        '''Establece un nuevo tipo de datos para el atributo XACML'''
        self.__xacml_DataType = dataType

    def hashCode(self) -> int:
        '''Devuelve el hashCode del atributo XACML'''
        return hash(self.__name)

    def getListableMessage(self):
        oDict = OrderedDict([(self.__xacml_id, self.__name)])
        return oDict[0]
    
    def __str__(self) -> str:
        return (f'{type(self).__name__}: [ name: {self.__name}, xacml_id: {self.__xacml_id},' 
                f' sortedValue: {self.__sortedValue}, xacml_DataType: {self.__xacml_DataType} ]')