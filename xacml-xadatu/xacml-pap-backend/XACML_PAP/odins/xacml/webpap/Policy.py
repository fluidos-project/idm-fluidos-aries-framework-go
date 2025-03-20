from .NamedObject import NamedObject
from OdinS_xacml_util import OdinS_xacml_util
from xacmleditor import (ElementoPolicy, 
                         ElementoRule,
                        ElementoTarget,
                        ElementoObligations,
                        ElementoObligation)
from .Rule import Rule
from .Obligation import Obligation

class Policy(NamedObject):
    __Rules = {} # Diccionario que almacena las reglas de la política
    __ruleList = [] # Lista que almacena las reglas de la política (para mantener el orden)
    __RuleCombiningAlgId = None # Identificador del algoritmo de combinación de reglas
    __PolicyId = None # Identificador de la política
    __obligations = [] # Lista que almacena las obligaciones de la política

    def __init__(self, policyId=None, odins_xacml_util=None, elementoPolicy=None) -> None:
        '''Constructor del objeto'''
        self.__Rules = {}
        self.__ruleList = []
        self.__obligations = []
        # Constructor para crear una nueva política
        if policyId is not None and odins_xacml_util is None and elementoPolicy is None:
            self.__PolicyId = policyId
            self.__RuleCombiningAlgId = "urn:oasis:names:tc:xacml:1.0:rule-combining-algorithm:first-applicable"
        # Constructor para crear una política a partir de un elemento XML
        elif policyId is None and odins_xacml_util is not None and elementoPolicy is not None:
            elementoPolicyAtts = elementoPolicy.getAtributos()
            self.__RuleCombiningAlgId = elementoPolicyAtts['RuleCombiningAlgId']
            self.__PolicyId = elementoPolicyAtts['PolicyId']

            # Obtenemos las reglas de la política
            elementoRules = odins_xacml_util.getChildren(elementoPolicy, ElementoRule.TIPO_RULE)
            for eRule in elementoRules:
                rule = Rule(odins_xacml_util=odins_xacml_util, elementoRule=eRule)
                self.__Rules[rule.getKeyForMap()] = rule
                self.__ruleList.append(rule)

            # Obtenemos las obligaciones de la política
            elemObligations = odins_xacml_util.getChild(elementoPolicy, ElementoObligations.TIPO_OBLIGATIONS)
            if elemObligations is not None:
                elementosObligation = odins_xacml_util.getChildren(elemObligations, ElementoObligation.TIPO_OBLIGATION)
                for elem in elementosObligation:
                    obligation = Obligation(odins_xacml_util=odins_xacml_util, elementoObligation=elem)
                    self.__obligations.append(obligation)

    def getName(self) -> str:
        '''Devuelve el identificador de la política'''
        return self.getPolicyId()
    
    def getKeyForMap(self) -> None:
        '''Devuelve la clave usada para almacenar la política'''
        return None
    
    def getPolicyId(self) -> str:
        '''Devuelve el ID de la política'''
        return self.__PolicyId
    
    def setPolicyId(self, policyId : str) -> None:
        '''Establece el ID de la política'''
        self.__PolicyId = policyId

    def getCombiningAlg(self) -> str:
        '''Devuelve el identificador del algoritmo de combinación de regla'''
        return self.__RuleCombiningAlgId
    
    def setCombiningAlg(self, combiningAlg : str) -> None:
        '''Establece el identificador de algoritmo de combinación de regla'''
        self.__RuleCombiningAlgId = combiningAlg

    def getRules(self) -> list:
        '''Devuelve la lista de reglas'''
        return self.__ruleList
    
    def getRulesMap(self) -> dict:
        '''Devuelve el diccionario de reglas'''
        return self.__Rules
    
    def setRules(self, ruleList=None, ruleDict=None) -> None:
        '''Guarda las reglas. Si se pone en ruleList, lo guarda en la lista.
        Si se pone en ruleDict, se pone en el diccionario de reglas. No se puede
        guardar en ambos sitios a la vez'''
        if ruleList is not None and ruleDict is None:
            self.__ruleList = ruleList
        elif ruleList is None and ruleDict is not None:
            self.__Rules = ruleDict
    
    def getObligations(self) -> list:
        '''Devuelve las obligaciones de la política'''
        return self.__obligations
    
    def setObligations(self, obligations : list) -> None:
        '''Guarda las obligaciones de la política'''
        self.__obligations = obligations

    def getRuleByName(self, name : str) -> Rule:
        '''Devuelve el objeto Rule por el nombre de la misma'''
        for rule in self.__ruleList:
            if rule.getName() == name:
                return rule
    
    def deleteRuleByName(self, name : str) -> None:
        '''Borra una regla de la lista de reglas por el nombre de la misma'''
        for rule in self.__ruleList:
            if rule.getName() == name:
                self.__ruleList.remove(name)

    def getOdinS_XACML(self) -> OdinS_xacml_util:
        '''Devuelve un objeto odins_xacml_util que representa la política en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoPolicy = odins_xacml_util.createPrincipal(ElementoPolicy.TIPO_POLICY)
        elementoPolicyAtts = elementoPolicy.getAtributos()
        elementoPolicyAtts['PolicyId'] = self.__PolicyId
        elementoPolicyAtts['RuleCombiningAlgId'] = self.__RuleCombiningAlgId

        odins_xacml_util.createChild(elementoPolicy, ElementoTarget.TIPO_TARGET)
        rules = self.getRules()
        for rule in rules:
            odins_xacml_util.insertChild(elementoPolicy, rule.getOdinS_XACML())

        if len(self.__obligations) != 0:
            elementoObligations = odins_xacml_util.createChild(elementoPolicy, ElementoObligations.TIPO_OBLIGATIONS)
            for obligation in self.__obligations:
                odins_xacml_util.insertChild(elementoObligations, obligation.getOdinS_XACML())

        return odins_xacml_util

    def __str__(self):
        '''Devuelve una cadena que representa la política'''
        return (f'{type(self).__name__}: [ {self.getName()} ]')