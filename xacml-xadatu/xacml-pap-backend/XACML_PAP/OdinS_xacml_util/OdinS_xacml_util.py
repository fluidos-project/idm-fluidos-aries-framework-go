from anytree import Node, RenderTree
import xml.etree.ElementTree as ET
import xml.dom.minidom as minidom
from xacmleditor import (ElementoXACML, 
                         InsertarOrdenadoElemento, 
                         AnalizadorSAX, 
                         ElementoXACMLFactoryImpl, 
                         ConversorDOM, 
                         ElementoPolicy,
                         ElementoSubjects,
                         ElementoResources,
                         ElementoActions,
                         ElementoSubject,
                         ElementoResource,
                         ElementoAction,
                         ElementoRule)

class OdinS_xacml_util:
    __raiz = None # El nodo raíz de la estructura de árbol que se está construyendo
    __miDTM = Node("root") # Nodo raíz alternativo
    __treeNode = {} # Un diccionario que mapea los elementos ElementoXACML

    def __init__(self, node=None, xacml=None) -> None:
        '''Inicializa los objetos OdinS_xacml_util'''
        if node is None and xacml is None:
            self.__raiz = Node("Policy Document")
        elif node is not None and xacml is None:
            dom = minidom.getDOMImplementation()
            domDOC = dom.createDocument(None, None, None)
            domRoot = domDOC.importNode(node, True)
            self.createTree(domRoot.toxml)
        elif node is None and xacml is not None:
            self.createTree(xacml)

    def createTree(self, xacml : str) -> None:
        '''Crea el árbol a partir de un documento XACML.'''
        asax = AnalizadorSAX()
        DTM = asax.analizarFromString(xacml)
        self.__raiz = DTM.name.root
        for node in RenderTree(self.__raiz):
            if isinstance(node.node.name, ElementoXACML):
                self.__treeNode[node.node.name] = node.node

    def insertChild(self, father : ElementoXACML, child) -> None:
        '''Inserta un nodo hijo en el nodo padre especificado.'''
        nodoPadre = self.__treeNode[father]
        self.insertaNodos(nodoPadre, child.__treeNode[child.getPrincipal()])

    def insertaNodos(self, father : Node, child : Node) -> None:
        '''Inserta varios nodos hijos en el nodo padre especificado.'''
        e = child.name
        hijoNode = self.insertaNodo(e, father)
        e2 = hijoNode.name
        self.__treeNode[e2] = hijoNode
        hijos = child.children
        for hijo in hijos:
            self.insertaNodos(hijoNode, hijo)

    def createPrincipal(self, type : str):
        '''Crea un nodo principal especificado por el tipo.'''
        policyElement = self.crearNodos(type, self.__raiz)
        e = policyElement.name
        self.__treeNode[e] = policyElement
        return e

    def getPrincipal(self):
        '''Devuelve el nodo principal'''
        return self.__raiz.children[0].name

    def createChild(self, father : ElementoXACML, type : str):
        '''Crea un nodo hijo de un nodo padre especificado por el tipo.'''
        padre = self.__treeNode[father]
        hijo = self.crearNodos(type, padre)
        e = hijo.name
        self.__treeNode[e] = hijo
        return e

    def getChildren(self, father, tipo=None) -> list:
        '''Devuelve una lista de los hijos de un nodo padre. 
        Si se especifica un tipo, solo devuelve los hijos de ese tipo.'''
        childList = []
        padre = self.__treeNode[father]
        if tipo is None:
            hijos = padre.children
            for hijo in hijos:
                childList.append(hijo.name)
        else:
            hijos = self.getChildrenElements(padre, tipo)
            for hijo in hijos:
                childList.append(hijo.name)
        return childList

    def getChild(self, father, tipo):
        '''Devuelve el nodo hijo del padre del tipo especificado.'''
        padre = self.__treeNode[father]
        hijo = self.getChildElement(padre, tipo)
        if hijo is not None:
            return hijo.name
        return None

    def removeChild(self, father, child) -> None:
        '''Quita el hijo especificado del padre'''
        hijo = self.__treeNode[hijo]
        hijo.parent = None

    def crearNodos(self, s : str, nodo : Node) -> Node:
        '''Crea un nuevo nodo'''
        eXACML = ElementoXACMLFactoryImpl.obtenerElementoXACML(s, {})
        eXACML.setVacio(True)
        nodoXACML = Node(eXACML, parent=nodo)
        return nodoXACML

    def insertaNodo(self, eXACML : ElementoXACML, nodo : Node) -> Node:
        '''Inserta un nodo nuevo'''
        return Node(eXACML, parent=nodo)

    def getChildElement(self, ruleElement : Node, tipo : str):
        '''Devuelve el hijo del mismo tipo que el especificado. Por cómo funciona la librería
        anytree, si el hijo es un ElementoSubjects y el tipo no es ese, puede que el tipo buscado sea
        de ElementoResources o de ElementoActions. Si es así, el nodo avanza hasta encontrarlo. Si por el
        camino se encuentra una regla distinta, o el mismo tipo, el bucle para. Si es del mismo tipo,
        se devuelve. Sino, el bucle se detiene.'''
        hijos = ruleElement.children
        for hijo in hijos:
            if (hijo.name.getTipo() == ElementoSubjects.TIPO_SUBJECTS
                and tipo != ElementoSubjects.TIPO_SUBJECTS):
                if tipo == ElementoResources.TIPO_RESOURCES:
                    while True:
                        if (hijo.name.getTipo() == tipo
                            or hijo.name.getTipo() == ElementoRule.TIPO_RULE):
                            break
                        else:
                            hijo = hijo.children[0]
                elif tipo == ElementoActions.TIPO_ACTIONS:
                    while True:
                        if (hijo.name.getTipo() == tipo
                            or hijo.name.getTipo() == ElementoRule.TIPO_RULE):
                            break
                        else:
                            hijo = hijo.children[0]
            if hijo.name.getTipo() == tipo:
                return hijo
        return None

    def getChildrenElements(self, policyElement : Node, tipo : str) -> list:
        '''Obtiene todos los elementos hijos de la regla. Es posible que al iterar entre
        elementos se encuentren políticas distintas. De ser así, el bucle se detiene. Si
        se encuentra una nueva regla a la que se está buscando, también se detiene el bucle.
        De lo contrario, si el tipo coincide, se añade el elemento a la lista. Finalmente devuelve
        la lista de reglas.'''
        rules = []
        for hijo in RenderTree(policyElement):
            if (hijo.node.name.getTipo() == ElementoPolicy.TIPO_POLICY 
                and policyElement.name.getTipo() == ElementoPolicy.TIPO_POLICY):
                if hijo.node.name.getID() != policyElement.name.getID():
                    break
            elif (hijo.node.name.getTipo() == ElementoRule.TIPO_RULE
                and (tipo == ElementoAction.TIPO_ACTION or
                     tipo == ElementoResource.TIPO_RESOURCE or
                     tipo == ElementoSubject.TIPO_SUBJECT)):
                break
            elif hijo.node.name.getTipo() == tipo:
                rules.append(hijo.node)
        return rules


    def getDocument(self):
        '''Devuelve el documento generado con la raíz actual'''
        return ConversorDOM.convierte(ConversorDOM(), self.__raiz)

    def __str__(self):
        doc = self.getDocument()
        return doc.toprettyxml(indent="  ", encoding="UTF-8").decode()
