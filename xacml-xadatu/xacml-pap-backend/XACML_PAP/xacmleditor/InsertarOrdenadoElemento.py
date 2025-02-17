from anytree import Node
from xacmleditor import ElementoXACML

class InsertarOrdenadoElemento:
    '''Actualmente no usado, puesto que ya se insertan ordenados los elementos. 
    Aún así, existe su implementación, por si se quiere usar más adelante.'''
    _nodoPadre = None
    _nodoHijo = None
    _elem = None
    elementosXACML = []

    def __init__(self, nodoPadre : Node, nodoHijo : Node):
        self._nodoPadre = nodoPadre
        self._nodoHijo = nodoHijo
        elem = nodoPadre.name

    def buscarPosicion(self, size : int, elemNum : ElementoXACML):
        pos = -1
        for i in range(size):
            if pos == -1:
                if self._elem.getPosicion(self.elementosXACML[i]) > self._elem.getPosicion(elemNum):
                    pos = i

        return pos
    
    def moverHaciaDerecha(self, desde : int, hasta : int):
        for i in range(hasta, desde, -1):
            self.elementosXACML.insert(i, self.elementosXACML[i-1])

    def ordenarInsercion(self):
        pos = None
        elemNum = None
        hijos = self._nodoPadre.children
        i = 0

        for hijo in hijos:
            self.elementosXACML.insert(i, hijo.name)
            i += 1
        
        self.elementosXACML.insert(i, self._nodoHijo.name)

        for i in range(1, len(self.elementosXACML)):
            elemNum = self.elementosXACML[i]
            pos = self.buscarPosicion(i, elemNum)
            if pos != -1:
                self.moverHaciaDerecha(pos, i)
                self.elementosXACML.insert(pos, elemNum)

    def getPosicion(self):
        for i in range(len(self.elementosXACML)):
            if self.elementosXACML[i] == self._nodoHijo.name:
                return i

        return len(self.elementosXACML)