from .NamedObject import NamedObject
from OdinS_xacml_util import OdinS_xacml_util
from xacmleditor import (ElementoTarget, 
                         ElementoRule, 
                         ElementoResource,
                         ElementoAction,
                         ElementoSubject, 
                         ElementoSubjects, 
                         ElementoActions,
                         ElementoResources)
from .Resource import Resource
from .Subject import Subject
from .Action import Action
class Rule(NamedObject):
    __RuleId = None # Identificador de la regla
    __resources = [] # Lista que guarda cada elemento Resource de la regla
    __subjects = [] # Lista que guarda cada elemento Subject de la regla
    __Effect = None # Efecto de la regla
    __actions = [] # Lista que guarda cada elemento Action de la regla

    # Efectos de la regla
    RULE_PERMIT = "Permit"
    RULE_DENY = "Deny"

    # Efectos de las acciones
    ACTION_READ = "Read"
    ACTION_WRITE = "Write"

    def __init__(self, name=None, odins_xacml_util=None, elementoRule=None) -> None:
        self.__resources = [] # Se generan otras listas para que no se compartan entre elementos
        self.__subjects = []
        self.__actions = []
        # Constructor para generar una nueva regla
        if name is not None and odins_xacml_util is None and elementoRule is None:
            self.__Effect = self.RULE_PERMIT
            self.__RuleId = name
        # Constructor para crear una regla a partir de un elemento XACML
        elif name is None and odins_xacml_util is not None and elementoRule is not None:
            elementoPolicyAtts = elementoRule.getAtributos()
            self.__RuleId = elementoPolicyAtts['RuleId']
            self.__Effect = elementoPolicyAtts['Effect']

            elementoTarget = odins_xacml_util.getChild(elementoRule, ElementoTarget.TIPO_TARGET)

            # Se obtienen los Subjects de la regla
            elementoSubjects = odins_xacml_util.getChild(elementoTarget, ElementoSubjects.TIPO_SUBJECTS)
            if elementoSubjects is None:
                self.__subjects.append(Subject())
            else:
                hijos = odins_xacml_util.getChildren(elementoSubjects, ElementoSubject.TIPO_SUBJECT)
                for hijo in hijos:
                    subject = Subject(odins_xacml_util=odins_xacml_util, elementoSubject=hijo)
                    self.__subjects.append(subject)

            # Se obtienen los elementos Resource de la regla
            elementoResources = odins_xacml_util.getChild(elementoTarget, ElementoResources.TIPO_RESOURCES)
            if elementoResources is None:
                self.__resources.append(Resource())
            else:
                hijos = odins_xacml_util.getChildren(elementoResources, ElementoResource.TIPO_RESOURCE)
                for hijo in hijos:
                    resource = Resource(odins_xacml_util=odins_xacml_util, elementoResource=hijo)
                    self.__resources.append(resource)

            # Se obtienen los elementos Action de la regla
            elementoActions = odins_xacml_util.getChild(elementoTarget, ElementoActions.TIPO_ACTIONS)
            if elementoActions is None:
                self.__actions.append(Action())
            else:
                hijos = odins_xacml_util.getChildren(elementoActions, ElementoAction.TIPO_ACTION)
                for hijo in hijos:
                    action = Action(odins_xacml_util=odins_xacml_util, elementoAction=hijo)
                    self.__actions.append(action)
    
    def getName(self) -> str:
        '''Devuelve el nombre de la regla'''
        return self.__RuleId
    
    def getKeyForMap(self) -> str:
        '''Devuelve la clave del diccionario en el que se guarda la regla'''
        return self.getName()
    
    def setName(self, name) -> None:
        '''Establece un nuevo nombre para la regla'''
        self.__RuleId = name

    def getResources(self) -> list:
        '''Devuelve la lista de objetos Resource'''
        return self.__resources
    
    def getResourceMap(self) -> dict:
        '''Devuelve un diccionario con los objetos Resource'''
        resources = {}
        for resource in self.getResources():
            resources[resource.getKeyForMap()] = resource
        return resources

    def setResources(self, resources : list) -> None:
        '''Establece una nueva lista de elementos Resource'''
        self.__resources = resources

    def getSubjects(self) -> list:
        '''Devuelve la lista de elementos Subject'''
        return self.__subjects
    
    def getSubjectsMap(self) -> dict:
        '''Devuelve un diccionario con los objetos Subject'''
        subjects = {}
        for subject in self.getSubjects():
            subjects[subject.getKeyForMap()] = subject
        return subjects

    def setSubjects(self, subjects : list) -> None:
        '''Establece una nueva lista de elementos Subject'''
        self.__subjects = subjects

    def getPolicy(self) -> str:
        '''Devuelve el efecto de la regla'''
        return self.__Effect
    
    def setPolicy(self, policy : str) -> None:
        '''Establece un nuevo efecto para la regla'''
        self.__Effect = policy

    def getActions(self) -> list:
        return self.__actions
    
    def getActionsMap(self) -> dict:
        '''Devuelve un diccionario con los objetos Action'''
        actions = {}
        for action in self.getActions():
            actions[action.getKeyForMap()] = action
        return actions

    def setActions(self, actions : list) -> None:
        '''Establece una nueva lista de objetos Action'''
        self.__actions = actions

    def getOdinS_XACML(self) -> OdinS_xacml_util:
        '''Devuelve un objeto OdinS_xacml_util que representa la regla en XML'''
        odins_xacml_util = OdinS_xacml_util()
        elementoRule = odins_xacml_util.createPrincipal(ElementoRule.TIPO_RULE)
        elementoRuleAtts = elementoRule.getAtributos()
        elementoRuleAtts['Effect'] = self.__Effect
        elementoRuleAtts['RuleId'] = self.__RuleId

        elementoTarget = odins_xacml_util.createChild(elementoRule, ElementoTarget.TIPO_TARGET)
        #if not Subject() in self.__subjects and len(self.__subjects) != 0:
        if not Subject() in self.__subjects:
            elementoSubjects = odins_xacml_util.createChild(elementoTarget, ElementoSubjects.TIPO_SUBJECTS)
            for subject in self.__subjects:
                odins_xacml_util.insertChild(elementoSubjects, subject.getOdinS_XACML())
        
        #if not Resource() in self.__resources and len(self.__resources) != 0:
        if not Resource() in self.__resources:
            elementoResources = odins_xacml_util.createChild(elementoTarget, ElementoResources.TIPO_RESOURCES)
            for resource in self.__resources:
                odins_xacml_util.insertChild(elementoResources, resource.getOdinS_XACML())

        #if not Action() in self.__actions and len(self.__actions) != 0:
        if not Action() in self.__actions:
            elementoActions = odins_xacml_util.createChild(elementoTarget, ElementoActions.TIPO_ACTIONS)
            for action in self.__actions:
                odins_xacml_util.insertChild(elementoActions, action.getOdinS_XACML())

        return odins_xacml_util

    def __str__(self) -> str:
        '''Devuelve una cadena que representa la regla'''
        return (f'{type(self).__name__}: [ name: {self.__RuleId} ]')
