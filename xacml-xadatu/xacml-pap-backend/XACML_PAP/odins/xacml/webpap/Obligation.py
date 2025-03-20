from enum import Enum
from OdinS_xacml_util import OdinS_xacml_util
from xacmleditor import ElementoObligation

class Fulfill(Enum):
    Permit = 1
    Deny = 2

class Obligation:
    __id = None # ID de la obligación
    __fulFillOn = None # Efecto de la obligación, de la clase Fulfill

    def __init__(self, id=None, odins_xacml_util=None, elementoObligation=None) -> None:
        '''Constructor del objeto'''
        # Constructor para crear una nueva obligación
        if id is not None and odins_xacml_util is None and elementoObligation is None:
            self.__id = id
            self.__fulFillOn = Fulfill.Permit
        # Constructor para crear una nueva obligación a partir de un elemento XML
        elif id is None and odins_xacml_util is not None and elementoObligation is not None:
            eObligationAtts = elementoObligation.getAtributos()
            self.__id = eObligationAtts['ObligationId']
            sFulfillOn = eObligationAtts['FulfillOn']
            if sFulfillOn == "Permit":
                self.__fulFillOn = Fulfill.Permit
            elif sFulfillOn == "Deny":
                self.__fulFillOn = Fulfill.Deny

    def getId(self) -> str:
        '''Devuelve el ID de la obligación'''
        return self.__id
    
    def setId(self, id : str) -> None:
        '''Establece un nuevo ID a la obligación'''
        self.__id = id
    
    def getFulFillOn(self) -> Fulfill:
        '''Devuelve el efecto de la política'''
        return self.__fulFillOn
    
    def setFulFillOn(self, FulFillOn : Fulfill) -> None:
        '''Establece un nuevo efecto para la política'''
        self.__fulFillOn = FulFillOn

    def getOdinS_XACML(self) -> OdinS_xacml_util:
        '''Devuelve un objeto odins_xacml_util que representa la obligación de la política en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoObligation = odins_xacml_util.createPrincipal(ElementoObligation.TIPO_OBLIGATION)
        elementoObligation.getAtributos()['ObligationId'] = self.__id
        if self.__fulFillOn == Fulfill.Permit:
            elementoObligation.getAtributos()['FulfillOn'] = "Permit"
        elif self.__fulFillOn == Fulfill.Deny:
            elementoObligation.getAtributos()['FulfillOn'] = "Deny"

        return odins_xacml_util