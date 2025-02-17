from .XACMLAttributeElement import XACMLAttributeElement
from xacmleditor import (ElementoSubjectMatch,
                         ElementoAttributeValue,
                         ElementoSubjectAttributeDesignator,
                         ElementoSubject)
from OdinS_xacml_util import OdinS_xacml_util

class Subject(XACMLAttributeElement):
    '''Elemento Subject heredado de XACMLAttributeElement'''
    sortingParameter = "subject"

    def __init__(self, name=None, ID=None, xacml_DataType=None, odins_xacml_util=None, elementoSubject=None) -> None:
        '''Inicializa el objeto, con dos posibles constructores: uno con name, ID y xacml_DataType y otro con
        odins_xacml_util y elementoAction.'''
        # Constructor para inicializar el objeto con nombre, ID y xacml_DataType
        if name is not None and ID is not None and xacml_DataType is not None and odins_xacml_util is None and elementoSubject is None:
            super().__init__(name, ID, self.sortingParameter if name != "" and ID != "" else None, xacml_DataType)
        # Constructor para inicializar el objeto a partir de un elemento XML de tipo Subject
        elif odins_xacml_util is not None and elementoSubject is not None and name is None and ID is None and xacml_DataType is None:
            super().__init__(None, None, self.sortingParameter)
            elementoSubjectMatch = odins_xacml_util.getChild(elementoSubject, ElementoSubjectMatch.TIPO_SUBJECTMATCH)
            elementoAttributeValue = odins_xacml_util.getChild(elementoSubjectMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
            elementoSubjectValueDesignator = odins_xacml_util.getChild(elementoAttributeValue, ElementoSubjectAttributeDesignator.TIPO_SUBJECTATTRIBUTEDESIGNATOR)
            self.setName(elementoAttributeValue.getContenido())
            self.setXACMLID(elementoSubjectValueDesignator.getID())

    def getOdinS_XACML(self) -> OdinS_xacml_util:
        '''Devuelve un objeto OdinS_xacml_util que representa el elemento Subject en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoSubject = odins_xacml_util.createPrincipal(ElementoSubject.TIPO_SUBJECT)
        elementoSubjectMatch = odins_xacml_util.createChild(elementoSubject, ElementoSubjectMatch.TIPO_SUBJECTMATCH)
        elementoSubjectMatch.getAtributos()['MatchId'] = "urn:oasis:names:tc:xacml:1.0:function:string-equal"

        elementoAttributeValue = odins_xacml_util.createChild(elementoSubjectMatch, ElementoAttributeValue.TIPO_ATTRIBUTEVALUE)
        elementoAttributeValue.getAtributos()['DataType'] = "http://www.w3.org/2001/XMLSchema#string"
        elementoAttributeValue.setContenido(self.getName())

        elementoSubjectAttributeDesignator = odins_xacml_util.createChild(elementoSubjectMatch, ElementoSubjectAttributeDesignator.TIPO_SUBJECTATTRIBUTEDESIGNATOR)
        eAADAtts = elementoSubjectAttributeDesignator.getAtributos()
        eAADAtts['AttributeId'] = self.getXACMLID()
        eAADAtts['DataType'] = "http://www.w3.org/2001/XMLSchema#string"

        return odins_xacml_util