from xml.sax.handler import ContentHandler
from xml import sax
from .MiContentHandler import MiContentHandler
from .SAXErrorHandler import SAXErrorHandler
import logging
import xml.dom.minidom as minidom
from .ConversorDOM import ConversorDOM
from anytree import Node
from xml.sax.xmlreader import InputSource
import io

class AnalizadorSAX(ContentHandler):
    __contentHandler = None
    __erorHandler = None

    logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')

    def __init__(self):
        self.__contentHandler = MiContentHandler()
        self.__erorHandler = SAXErrorHandler()

    def analizar(self, uri : str):
        '''Se parsea el contenido de un URI a través de MiContentHandler'''
        try:
            parser = sax.make_parser()
            parser.setContentHandler(self.__contentHandler)
            parser.setErrorHandler(self.__erorHandler)

            with open(uri, 'rb') as f:
                parser.parse(f)
            
            if self.__contentHandler.getWarnings() != "":
                self.__erorHandler.setErrores(self.__contentHandler.getWarnings() + 
                                            self.__erorHandler.getErrores())
        except Exception as e:
            logging.error(e)

        return self.__contentHandler.getDatos()

    def analizarFromString(self, toParse : str):
        '''Se parsea el contenido de un string a través de MiContentHandler'''
        try:
            sax.parseString(toParse, self.__contentHandler, self.__erorHandler)
            if self.__contentHandler.getWarnings() != "":
                self.__erorHandler.setErrores(self.__contentHandler.getWarnings() + 
                                            self.__erorHandler.getErrores())
        except Exception as e:
            logging.error(e)

        return self.__contentHandler.getDatos()


    def getErrorHandler(self) -> str:
        '''Devuelve los errores obtenidos por el handler'''
        return self.__erorHandler.getErrores()
    
    def procesaSalvar(self, node : Node, uri : str):
        try:
            doc = ConversorDOM.convierte(node)
            minidomDoc = minidom.Document(doc)
            with open(uri, "w") as f:
                f.write(minidomDoc.toprettyxml())
        except Exception as e:
            logging.error(e)

    def procesaValidar(self, node : Node, os):
        try:
            doc = ConversorDOM.convierte(node)
            minidomDoc = minidom.Document(doc)
            with open(os, "w") as f:
                f.write(minidomDoc.toprettyxml())
        except Exception as e:
            logging.error(e)