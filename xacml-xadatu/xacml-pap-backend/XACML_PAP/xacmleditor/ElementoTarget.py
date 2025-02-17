from .ElementoXACML import ElementoXACML

class ElementoTarget(ElementoXACML):
    TIPO_TARGET = "Target"

    __allowedChild = ["Subjects", "Resources", "Actions", "Environments"]
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_TARGET)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return self.__allowedChild
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string
    
    def getPosicion(self, elementoXACML : ElementoXACML) -> int:
        if elementoXACML.getTipo() in self.__allowedChild:
            return self.__allowedChild.index(elementoXACML.getTipo())
        else:
            return super().getPosicion(elementoXACML)
        
    def getMaxNumChild(self, elementoXACML : ElementoXACML):
        for i in range(len(self.__allowedChild)):
            if elementoXACML.getTipo() == self.__allowedChild[i]:
                return 1
        return super().getMaxNumChild(elementoXACML)