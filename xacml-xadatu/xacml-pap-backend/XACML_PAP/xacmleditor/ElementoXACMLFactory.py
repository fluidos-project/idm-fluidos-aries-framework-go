from abc import ABC, abstractmethod

class ElementoXACMLFactory(ABC):

    @abstractmethod
    def obtenerElementoXACML(tipo, atributos):
        pass