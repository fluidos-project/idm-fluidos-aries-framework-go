from abc import ABC, abstractmethod

class NamedObject(ABC):
    '''Clase abstracta, usada para Policy y Rule.'''

    @abstractmethod
    def getName():
        pass
    
    @abstractmethod
    def getKeyForMap():
        pass