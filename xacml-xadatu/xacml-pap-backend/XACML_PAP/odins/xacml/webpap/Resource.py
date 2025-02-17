from .XACMLAttributeElement import XACMLAttributeElement
from xacmleditor import (ElementoResourceMatch,
                         ElementoAttributeValue,
                         ElementoResource,
                         ElementoResourceAttributeDesignator)
from OdinS_xacml_util import OdinS_xacml_util

class Resource(XACMLAttributeElement):
    '''Elemento Resource heredado de XACMLAttributeElement'''
    sortingParameter = "resource"

    def __init__(self, name=None, ID=None, xacml_DataType=None, odins_xacml_util=None, elementoResource=None):
        '''Inicializa el objeto, con dos posibles constructores: uno con name, ID y xacml_DataType y otro con
        odins_xacml_util y elementoResource.'''
        # Constructor para inicializar el objeto con nombre, ID y xacml_DataType
        if name is not None and ID is not None and xacml_DataType is not None and odins_xacml_util is None and elementoResource is None:
            super().__init__(name, ID, self.sortingParameter if name != "" and ID != "" else None, xacml_DataType)
        # Constructor para inicializar el objeto a partir de un elemento XML de tipo Resource
        elif odins_xacml_util is not None and elementoResource is not None and name is None and ID is None and xacml_DataType is None:
            super().__init__(None, None, self.sortingParameter)
            elementoResourceMatch = odins_xacml_util.getChild(elementoResource, ElementoResourceMatch.TIPO_RESOURCEMATCH)
            elementoAttributeValue = odins_xacml_util.getChild(elementoResourceMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
            elementoResourceValueDesignator = odins_xacml_util.getChild(elementoAttributeValue, ElementoResourceAttributeDesignator.TIPO_RESOURCEATTRIBUTEDESIGNATOR)
            self.setName(elementoAttributeValue.getContenido())
            self.setXACMLID(elementoResourceValueDesignator.getID())

    def getOdinS_XACML(self):
        '''Devuelve un objeto OdinS_xacml_util que representa el recurso en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoResource = odins_xacml_util.createPrincipal(ElementoResource.TIPO_RESOURCE)
        elementoResourceMatch = odins_xacml_util.createChild(elementoResource, ElementoResourceMatch.TIPO_RESOURCEMATCH)
        elementoResourceMatch.getAtributos()['MatchId'] = "urn:oasis:names:tc:xacml:1.0:function:string-equal"

        elementoAttributeValue = odins_xacml_util.createChild(elementoResourceMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
        elementoAttributeValue.getAtributos()['DataType'] = "http://www.w3.org/2001/XMLSchema#string"
        elementoAttributeValue.setContenido(self.getName())

        elementoResourceAttributeDesignator = odins_xacml_util.createChild(elementoResourceMatch, ElementoResourceAttributeDesignator.TIPO_RESOURCEATTRIBUTEDESIGNATOR)
        eAADAtts = elementoResourceAttributeDesignator.getAtributos()
        eAADAtts['AttributeId'] = self.getXACMLID()
        eAADAtts['DataType'] = "http://www.w3.org/2001/XMLSchema#string"

        return odins_xacml_util