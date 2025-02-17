from .ElementoXACML import ElementoXACML

class ElementoObligations(ElementoXACML):
    TIPO_OBLIGATIONS = "Obligations"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_OBLIGATIONS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True
    
    def getAllowedChild(self) -> list:
        return ["Obligation"]
    
    def getAllObligatory(self) -> list:
        return ["Obligation"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string