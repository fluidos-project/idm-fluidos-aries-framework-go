from .XACMLAttributeElement import XACMLAttributeElement
from xacmleditor import (ElementoActionMatch, 
                         ElementoAttributeValue, 
                         ElementoActionAttributeDesignator,
                         ElementoAction)
from OdinS_xacml_util import OdinS_xacml_util

class Action(XACMLAttributeElement):
    '''Elemento Action heredado de XACMLAttributeElement'''

    sortingParameter = "action"

    def __init__(self, name=None, ID=None, xacml_DataType=None, odins_xacml_util=None, elementoAction=None) -> None:
        '''Inicializa el objeto, con dos posibles constructores: uno con name, ID y xacml_DataType y otro con
        odins_xacml_util y elementoAction.'''
        # Constructor para inicializar el objeto con nombre, ID y xacml_DataType
        if name is not None and ID is not None and xacml_DataType is not None and odins_xacml_util is None and elementoAction is None:
            super().__init__(name, ID, self.sortingParameter if name != "" and ID != "" else None, xacml_DataType)
        # Constructor para inicializar el objeto a partir de un elemento XML de tipo Action
        elif odins_xacml_util is not None and elementoAction is not None and name is None and ID is None and xacml_DataType is None:
            super().__init__(None, None, self.sortingParameter)
            # Se obtienen los nodos de los elementos y se modifica el objeto Action para insertar el nombre y el XACMLID
            elementoActionMatch = odins_xacml_util.getChild(elementoAction, ElementoActionMatch.TIPO_ACTIONMATCH)
            elementoAttributeValue = odins_xacml_util.getChild(elementoActionMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
            elementoAttributeValueDesignator = odins_xacml_util.getChild(elementoAttributeValue, ElementoActionAttributeDesignator.TIPO_ACTIONATTRIBUTEDESIGNATOR)
            self.setName(elementoAttributeValue.getContenido())
            self.setXACMLID(elementoAttributeValueDesignator.getID())

    def getOdinS_XACML(self) -> OdinS_xacml_util:
        '''Devuelve un objeto OdinS_xacml_util que representa la acción en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoAction = odins_xacml_util.createPrincipal(ElementoAction.TIPO_ACTION)
        elementoActionMatch = odins_xacml_util.createChild(elementoAction, ElementoActionMatch.TIPO_ACTIONMATCH)
        elementoActionMatch.getAtributos()['MatchId'] = "urn:oasis:names:tc:xacml:1.0:function:string-equal"

        elementoAttributeValue = odins_xacml_util.createChild(elementoActionMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
        elementoAttributeValue.getAtributos()['DataType'] = "http://www.w3.org/2001/XMLSchema#string"
        elementoAttributeValue.setContenido(self.getName())

        elementoActionAttributeDesignator = odins_xacml_util.createChild(elementoActionMatch, ElementoActionAttributeDesignator.TIPO_ACTIONATTRIBUTEDESIGNATOR)
        eAADAtts = elementoActionAttributeDesignator.getAtributos()
        eAADAtts['AttributeId'] = self.getXACMLID()
        eAADAtts['DataType'] = "http://www.w3.org/2001/XMLSchema#string"

        return odins_xacml_util