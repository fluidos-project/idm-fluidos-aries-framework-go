from .ElementoXACML import ElementoXACML

class ElementoDescription(ElementoXACML):
    TIPO_DESCRIPTION = "Description"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_DESCRIPTION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string