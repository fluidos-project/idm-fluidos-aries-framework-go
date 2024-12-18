from abc import ABC, abstractmethod


class DBManager(ABC):
    '''Clase abstracta de DBManager. Usada en los módulos DiskDBManager y ExistDBManager'''
    @abstractmethod
    def retrievePolicySet(self):
        pass

    @abstractmethod
    def storePolicySet(self, policySet):
        pass

    @abstractmethod
    def storeXACMLAttributes(self, attributes):
        pass

    @abstractmethod
    def getXACMLAtribbutes(self):
        pass
