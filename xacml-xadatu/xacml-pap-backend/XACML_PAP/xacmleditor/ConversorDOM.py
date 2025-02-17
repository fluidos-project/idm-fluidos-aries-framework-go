import logging
import xml.etree.ElementTree as ET
from xml.dom.minidom import Document
from anytree import Node

class ConversorDOM:
    logging.basicConfig(format='%(filename)s:%(lineno)d - %(levelname)s - %(message)s')

    def __init__(self):
        pass

    def convierte(self, root : Node):
        '''Se convierte el contenido del nodo a un documento'''
        doc = None
        try:
            if isinstance(root.name, str):
                doc = Document()
            else:
                return doc

            hijos = root.children
            
            for hijo in hijos:
                doc.appendChild(self.__procesaHijo(hijo, doc))

            return doc
        except Exception as e:
            logging.error(e)
        
        return doc

    def __procesaHijo(self, node : Node, doc : Document):
        '''Se procesa el contenido del nodo y se añade al documento'''
        e = node.name
        ret = doc.createElement(e.getTipo())
        
        hijos = node.children

        for hijo in hijos:
            ret.appendChild(self.__procesaHijo(hijo, doc))

        atributos = e.getAtributos()
        claves = atributos.keys()
        
        for clave in claves:
            ret.setAttribute(clave, atributos[clave])

        texto = e.getContenido()
        if texto != "":
            txtNode = doc.createTextNode(texto)
            ret.appendChild(txtNode)

        return ret
        