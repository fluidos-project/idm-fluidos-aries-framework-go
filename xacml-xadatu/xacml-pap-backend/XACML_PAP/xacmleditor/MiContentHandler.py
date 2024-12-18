from xml.sax import ContentHandler, SAXException
from xml.sax.xmlreader import Locator, AttributesImpl
from .ElementoXACML import ElementoXACML
from .ElementoXACMLFactory import ElementoXACMLFactory
from .ElementoXACMLFactoryImpl import ElementoXACMLFactoryImpl
from anytree import Node, RenderTree
from .ElementoAttributeValue import ElementoAttributeValue

class MiContentHandler(ContentHandler):
    profundidad = 0
    elemento = None
    __warnings = ""
    resto = False
    __datos = None
    __elementoActual = None

    def startPrefixMapping(self, prefix : str, uri : str):
        pass

    def getDatos(self):
        return self.__datos
    
    def setDocumentLocator(self, locator : Locator):
        pass

    def startDocument(self):
        self.__elementoActual = Node("Policy Document")

    def endDocument(self):
        self.__datos = Node(self.__elementoActual)

    def processingInstruction(self, target, data):
        pass

    def endPrefixMapping(self, prefix):
        pass

    def startElement(self, name, atts):
        '''Empieza a generarse un elemento con sus atributos. Se obtiene los nombres
        de cada atributo y se añade al árbol.'''
        atributos = {}
        qNames = atts.getQNames()
        qNamesIterator = iter(qNames)

        for i in range(atts.getLength()):
            attName = next(qNamesIterator)
            atributos[attName] = atts.get(attName)

        if not self.resto:
            elem = ElementoXACMLFactoryImpl.obtenerElementoXACML(name, atributos)
            nodoActual = Node(elem, parent=self.__elementoActual)
            if isinstance(elem, ElementoAttributeValue):
                self.resto = True
            self.__elementoActual = nodoActual
            self.profundidad += 1
            self.elemento = name
        else:
              self.__elementoActual.name.addContenido(f'<{name}>')

    def endElement(self, name : str):
        '''Marca la finalización de la creación de un elemento.'''
        self.profundidad -= 1
        e = self.__elementoActual.name # .name se usa para obtener el contenido del nodo
        if isinstance(e, ElementoAttributeValue) and e.getTipo() == name:
            self.resto = False
        elif self.resto is True:
            e.addContenido(f"</{name}>")

        if self.resto is False:
            if e is None:
                delNode = self.__elementoActual
                self.__elementoActual = self.__elementoActual.parent
                delNode.parent = None
                del delNode
                match name:
                    case "AnySubject" | "AnyResource" | "AnyAction":
                        self.__warnings += f"<{name}> has been removed from XACML 2.0\n"
                    case _:
                        self.__warnings += f"<{name}> is not recognized.\n"
            else:
                self.__warnings += f"<{name} is not recognized.\n"
        else:
            if self.__elementoActual.is_leaf():
                e.setVacio(True)
            self.__elementoActual = self.__elementoActual.parent

    def characters(self, content : str):
        if self.__elementoActual is not None:
            e = self.__elementoActual.name
            if isinstance(e, ElementoXACML):
                content = content.replace("\t", "")
                content = content.strip()
                e.addContenido(content)

    def ignorableWhitespace(self, ch, start, end):
        pass

    def skippedEntity(self, name : str):
        pass

    def getWarnings(self) -> str:
        return self.__warnings
