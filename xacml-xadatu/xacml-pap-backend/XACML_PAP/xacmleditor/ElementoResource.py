from .ElementoXACML import ElementoXACML

class ElementoResource(ElementoXACML):
    TIPO_RESOURCE = "Resource"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RESOURCE)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["ResourceMatch"]
    
    def getAllObligatory(self) -> list:
        return ["ResourceMatch"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string